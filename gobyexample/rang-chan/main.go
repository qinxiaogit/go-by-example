package main

import "fmt"

func main() {
	//一个非空的通道也是可以关闭的，但是通道中剩下的值仍然可以被接收到
	queue:=make(chan string,3)
	queue<-"one1"
	queue<-"one2"
	queue<-"one3"
	close(queue)
	for index:=range(queue)  {
		fmt.Println(index)
	}
}
