package main

import (
	"fmt"
)

func main() {
	message:=make(chan string)
	sig:=make(chan bool)
	select {
	case msg:=<-message:
		fmt.Println("receivedMgs",msg)
	default:
		fmt.Println("no receive")

	}

	msg:="hi"
	select {
	case message<-msg:
		fmt.Println("sent message",msg)
	default:
		fmt.Println("no message")
	}

	select {
	case msg:=<-message:
		fmt.Println("msg->",msg)
	case sig:=<-sig:
		fmt.Println(sig)
	default:
		fmt.Println("no active")
	}
}
