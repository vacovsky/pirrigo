package main

import (
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var conn *amqp.Connection

func rabbitConnect() {
	conn, err := amqp.Dial(RMQCONNSTRING)
	if err != nil {
		for conn == nil {
			getLogger().LogError("Unable to connect to RabbitMQ.  Trying again in 15 seconds.",
				zap.String("AMQPConnectionString", RMQCONNSTRING),
				zap.String("error", err.Error()),
			)
			time.Sleep(time.Duration(15) * time.Second)
			conn, err = amqp.Dial(RMQCONNSTRING)
			if err != nil {
				getLogger().LogError("Unable to connect to RabbitMQ.  Fatal?  Probably..",
					zap.String("AMQPConnectionString", RMQCONNSTRING),
					zap.String("error", err.Error()))
			}
		}
	}
}

func rabbitSend(queueName string, body string) {
	rabbitConnect()
	defer conn.Close()
	getLogger().LogEvent(`Sending message to queue`,
		zap.String("messageBody", body),
		zap.String("queueName", queueName))

	ch, err := conn.Channel()
	if err != nil {
		getLogger().LogError("Unable to open AMQP channel for sending message.",
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
		getLogger().LogError("Failed to declare a queue.",
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
		getLogger().LogError("Failed publish message to the queue.",
			zap.String("error", err.Error()),
			zap.String("queueName", q.Name),
			zap.String("contentType", "text/plain"),
			zap.String("queueName", q.Name),
			zap.String("messageBody", body),
		)
	}

}

func rabbitReceive(queueName string) {
	rabbitConnect()
	defer conn.Close()

	for {
		ch, err := conn.Channel()
		if err != nil {
			getLogger().LogError("Failed to open a channel for receiving on RabbitMQ",
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
			getLogger().LogError("Failed to declare a queue.", zap.String("error", err.Error()))
		}

		time.Sleep(500 * time.Millisecond)

		for d := range msgs {
			messageHandler(queueName, d.Body)
		}
	}
}

func messageHandler(queueName string, message []byte) {
	if queueName == SETTINGS.RabbitMQ.TaskQueue {
		getLogger().LogEvent(
			`Sending message to RabbitMQ Server`,
			zap.String("queueName", queueName),
			zap.String("messageBody", string(message)))
		reactToGpioMessage(message)
	}
}

func listenForTasks() {
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
