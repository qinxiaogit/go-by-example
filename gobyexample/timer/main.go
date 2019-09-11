package main

import (
	"fmt"
	"time"
)

func main() {
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C
	fmt.Println("time1 1 expired ")

	timer2 := time.NewTimer(time.Second)

	go func() {
		<-timer2.C
		fmt.Println("time2 2 expired")
	}()

	//stop:=timer2.Stop()
	if timer2.Stop() {
		fmt.Println("timer2 is stop")
	}

}
