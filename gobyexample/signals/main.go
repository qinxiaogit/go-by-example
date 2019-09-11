package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sign := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sign
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exit")
}
