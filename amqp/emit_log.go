package main

import (
	"log"
	"os"
	"github.com/streadway/amqp"
	"github.com/qinxiaogit/go-by-example/amqp/tools"
)


func main(){
	conn,err := amqp.Dial("amqp:localhost:5672");
	defer conn.Close();
	if err!=nil{
		log.Fatalf("%s:%s",err,"Failed to connect to RabbitMQ")
	}

	ch ,err := conn.Channel()
	if err!= nil{
		log.Fatalf("%s:%s",err,"Failed to open a channel")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"log",//name
		"fanout",//type
		true,//durable
		false,//auto-deleted
		false,//internal
		false,//no wait
		nil,// arguments
		)

	body := tools.bodyFrom(os.Args)
	err = ch.Publish("logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body)},
	)

	tools.failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

