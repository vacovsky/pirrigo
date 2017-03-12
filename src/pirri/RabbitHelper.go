package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func rabbitConnect() {
	conn, ERR = amqp.Dial(RMQCONNSTRING)
	if ERR != nil {
		for conn == nil {
			fmt.Println("Waiting 15 seconds and attempting to connect to RabbitMQ again.")
			time.Sleep(time.Duration(15) * time.Second)
			conn, ERR = amqp.Dial(RMQCONNSTRING)
		}
	}
}

func rabbitSend(queueName string, body string) {
	rabbitConnect()
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

func rabbitReceive(queueName string) {
	rabbitConnect()
	defer conn.Close()

	for {
		failOnError(ERR, "Failed to connect to RabbitMQ")
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
			fmt.Println(ERR)
		}

		for d := range msgs { // the d stands for Delivery
			fmt.Println(string(d.Body[:]))
			messageHandler(queueName, d.Body)
		}
	}
}

func messageHandler(queueName string, message []byte) {
	if queueName == SETTINGS.RabbitMQ.TaskQueue {
		fmt.Println(queueName, message)
		reactToGpioMessage(message)
	}
}
