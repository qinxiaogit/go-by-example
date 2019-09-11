package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(time.Second)
		c1 <- "hello"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "world"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("recv:", msg1)
		case msg2 := <-c2:
			fmt.Println("recv2:", msg2)
		}
	}
}
