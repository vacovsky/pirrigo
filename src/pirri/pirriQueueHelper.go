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
			fmt.Println(ERR, "Waiting 15 seconds and attempting to connect to RabbitMQ again.")
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
	if ERR != nil {
		fmt.Println(ERR, "Unable to open channel.")
		return
	}

	q, ERR := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if ERR != nil {
		fmt.Println(ERR, "Failed to declare a queue")
		return
	}

	ERR = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if ERR != nil {
		fmt.Println(ERR, "Failed to publish message")
	}

}

func rabbitReceive(queueName string) {
	rabbitConnect()
	defer conn.Close()

	for {
		fmt.Println(ERR, "Failed to connect to RabbitMQ")

		ch, ERR := conn.Channel()
		if ERR != nil {
			fmt.Println(ERR, "Failed to open a channel")
		}
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

		time.Sleep(500 * time.Millisecond)

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

func listenForTasks() {
	defer WG.Done()
	for {
		ORQMutex.Lock()
		q := OfflineRunQueue
		ORQMutex.Unlock()

		var task *Task
		if len(q) > 0 {
			fmt.Println("Task found.")
			ORQMutex.Lock()
			task, OfflineRunQueue = OfflineRunQueue[len(OfflineRunQueue)-1],
				OfflineRunQueue[:len(OfflineRunQueue)-1]
			ORQMutex.Unlock()
			task.execute()
		}
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}
