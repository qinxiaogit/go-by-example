package main

import (
	"log"
	"github.com/streadway/amqp"
)
func main() {
	conn, err := amqp.Dial("amqp:localhost:5672");
	defer conn.Close();
	if err != nil {
		log.Fatalf("%s:%s", err, "Failed to connect to RabbitMQ")
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s:%s", err, "Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("%s:%s", err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("%s:%s", err, "Failed to register a consumer")
	}

	forever := make(chan bool)
	go func() {
		for d:=range msgs{
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}