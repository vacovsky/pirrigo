package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func RabbitConnect() {
	conn, ERR = amqp.Dial(RMQCONNSTRING)
	failOnError(ERR, "Failed to connect to RabbitMQ")

}
func RabbitSend(queueName string, body string) {
	RabbitConnect()
	defer conn.Close()
	fmt.Println("Sending", body, "to", queueName)
	ch, ERR := conn.Channel()

	q, ERR := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(ERR, "Failed to declare a queue")

	ERR = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(ERR, "Failed to publish a message")
}

func RabbitReceive(queueName string) {
	for {
		conn, ERR := amqp.Dial(RMQCONNSTRING)
		failOnError(ERR, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, ERR := conn.Channel()
		failOnError(ERR, "Failed to open a channel")
		defer ch.Close()

		q, ERR := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		autoAck := true

		msgs, ERR := ch.Consume(q.Name, "", autoAck, false, false, false, nil)
		if ERR != nil {
			panic(ERR.Error())
		}

		for d := range msgs { // the d stands for Delivery
			//			log.Printf(string(d.Body[:])) // or whatever you want to do with the message
			//			d.Ack(false)
			//			if queueName == SETTINGS.RabbitStopQueue {
			//				KILL = strconv.ParseBool(d.Body["kill"])
			//			}
			fmt.Println(string(d.Body[:]))
		}
		failOnError(ERR, "Failed to declare a queue")
	}
}

//{
//                            'sid': task[1],
//                            'schedule_id': task[0],
//                            'duration': task[2]
//                        }
