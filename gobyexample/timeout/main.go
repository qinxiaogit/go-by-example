package main

import (
	"fmt"
	"time"
)

func main(){

	c1:=make(chan string,1)
	go func() {
		time.Sleep(time.Second)
		c1<-"hello world"
	}()

	select {
	case res:=<-c1:
		fmt.Println(res)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}

	c2:=make(chan string,1)
	go func() {
		time.Sleep(time.Second*2)
		c2<-"go lang"
	}()

	select {
	case res:=<-c2:
		fmt.Println(res)
	case <-time.After(time.Second*3):
		fmt.Println("timeout 2")
	}


}
