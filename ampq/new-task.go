package main

import (
	"os"
	"log"
	"github.com/streadway/amqp"
	"strings"
	)

func main(){
	conn,err := amqp.Dial("amqp:localhost:5672")
	if err!= nil{
		log.Fatalf("%s-%s",err,"Failed to connect to RabbitMQ")
	}
	defer conn.Close()
	ch,err:= conn.Channel()
	if err!= nil{
		log.Fatalf("%s-%s",err,"Failed to channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err!=nil{
		log.Fatalf("%s-%s",err,"Failed to declare a queue")
	}
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	//failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}


func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
