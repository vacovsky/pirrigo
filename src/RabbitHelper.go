package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func RabbitConnect() {
	conn, err := amqp.Dial(RMQCONNSTRING)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

}
func RabbisSend() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "hello"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func RebbitReceive() {
	conn, err := amqp.Dial(RMQCONNSTRING)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	autoAck := false

	msgs, err := ch.Consume(q.Name, "", autoAck, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msgs { // the d stands for Delivery
		log.Printf(string(d.Body[:])) // or whatever you want to do with the message
		d.Ack(false)
	}
	failOnError(err, "Failed to declare a queue")
}
