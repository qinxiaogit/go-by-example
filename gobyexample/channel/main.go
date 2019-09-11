package main

import "fmt"

func main() {
	message := make(chan string)

	go func() {
		message <- "ping hello world"
	}()
	msg := <-message

	fmt.Println(msg)
}
