package main

import (
	"log"
	"github.com/streadway/amqp"
	//"encoding/binary"
	//"bytes"
)
type aaa struct {
	test int64
	name string
}

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

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err!= nil{
		log.Fatalf("%s:%s",err,"Failed to declare a queue")
	}
	//var bin_buf  bytes.Buffer
	//body := aaa{test:1000,name:"hello world"}
	body := "hello world"
	//binary.Write(&bin_buf, binary.BigEndian, body)

	for i:=0;i<1<<1;i++  {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:[]byte(body)   ,
			})
	}
	log.Printf(" [x] Sent %s", body)
}