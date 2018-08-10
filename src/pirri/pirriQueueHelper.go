package pirri

import (
	"time"

	"github.com/vacovsky/pirrigo/src/logging"
	"github.com/vacovsky/pirrigo/src/settings"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var conn *amqp.Connection

func rabbitConnect() {
	set := settings.Service()
	log := logging.Service()
	conn, err := amqp.Dial(set.RabbitMQ.ConnectionString)
	if err != nil {
		for conn == nil {
			log.LogError("Unable to connect to RabbitMQ.  Trying again in 15 seconds.",
				zap.String("AMQPConnectionString", set.RabbitMQ.ConnectionString),
				zap.String("error", err.Error()),
			)
			time.Sleep(time.Duration(15) * time.Second)
			conn, err = amqp.Dial(set.RabbitMQ.ConnectionString)
			if err != nil {
				log.LogError("Unable to connect to RabbitMQ.  Fatal?  Probably..",
					zap.String("AMQPConnectionString", set.RabbitMQ.ConnectionString),
					zap.String("error", err.Error()))
			}
		}
	}
}

func rabbitSend(queueName string, body string) {

	log := logging.Service()
	rabbitConnect()
	defer conn.Close()
	log.LogEvent(`Sending message to queue`,
		zap.String("messageBody", body),
		zap.String("queueName", queueName))

	ch, err := conn.Channel()
	if err != nil {
		log.LogError("Unable to open AMQP channel for sending message.",
			zap.String("messageBody", body),
			zap.String("error", err.Error()))
		return
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.LogError("Failed to declare a queue.",
			zap.String("error", err.Error()),
			zap.String("queueName", queueName),
		)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.LogError("Failed publish message to the queue.",
			zap.String("error", err.Error()),
			zap.String("queueName", q.Name),
			zap.String("contentType", "text/plain"),
			zap.String("queueName", q.Name),
			zap.String("messageBody", body),
		)
	}

}

func RabbitReceive(queueName string) {
	rabbitConnect()
	log := logging.Service()
	defer conn.Close()

	for {
		ch, err := conn.Channel()
		if err != nil {
			log.LogError("Failed to open a channel for receiving on RabbitMQ",
				zap.String("error", err.Error()),
				zap.String("queueName", queueName))
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		autoAck := true

		msgs, err := ch.Consume(q.Name, "", autoAck, false, false, false, nil)
		if err != nil {
			log.LogError("Failed to declare a queue.", zap.String("error", err.Error()))
		}

		time.Sleep(500 * time.Millisecond)

		for d := range msgs {
			messageHandler(queueName, d.Body)
		}
	}
}

func messageHandler(queueName string, message []byte) {
	log := logging.Service()
	set := settings.Service()
	if queueName == set.RabbitMQ.TaskQueue {
		log.LogEvent(
			`Sending message to RabbitMQ Server`,
			zap.String("queueName", queueName),
			zap.String("messageBody", string(message)))
		reactToGpioMessage(message)
	}
}

func ListenForTasks() {
	defer WG.Done()
	for {
		ORQMutex.Lock()
		q := OfflineRunQueue
		ORQMutex.Unlock()

		var task *Task
		if len(q) > 0 {
			ORQMutex.Lock()
			task, OfflineRunQueue = OfflineRunQueue[len(OfflineRunQueue)-1],
				OfflineRunQueue[:len(OfflineRunQueue)-1]
			ORQMutex.Unlock()
			task.execute()
		}
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}
