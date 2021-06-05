package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672"
	}

	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		panic("Error trying to open a connection: " + err.Error())
	}

	channel, err := connection.Channel()
	if err != nil {
		panic("Error trying to open a channel: " + err.Error())
	}

	//PUBLISHER
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		panic("Error trying to declare exchange: " + err.Error())
	}

	//SUBSCRIBER
	_, err = channel.QueueDeclare("my_queue", true, false, false, false, nil)
	if err != nil {
		panic("Error trying to declare the queue: " + err.Error())
	}

	err = channel.QueueBind("my_queue", "#", "events", false, nil)
	if err != nil {
		panic("Error trying to bind the queue: " + err.Error())
	}

	msgs, err := channel.Consume("my_queue", "", false, false, false, false, nil)
	if err != nil {
		panic("error while consuming the queue: " + err.Error())
	}

	//PUBLISH 10 MESSAGES
	go func() {
		message := amqp.Publishing{
			Body: []byte("Hello World"),
		}
		for i := 0; i < 10; i++  {
			err = channel.Publish("events", "events.test", false, false, message)
			if err != nil {
				panic("Error trying to publish an event: " + err.Error())
			}
		}
	}()

	//WAIT FOR MESSAGES
	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(false)
	}
}
