package main

import (
	"fmt"
	"time"
)

func worker(do chan bool) {
	fmt.Println("working...")
	time.Sleep(time.Second)
	fmt.Println("doing...")
	do <- true
}

func main() {
	done := make(chan bool)
	go worker(done)

	fmt.Println(<-done)
}
