package main

import (
	"fmt"
	"time"
)

func worker(id int,jobs <-chan int,results chan <-int){
	for j:= range jobs{

		fmt.Println("worker ",id ,"processing job",j)
		time.Sleep(time.Second)
		results<-j*2

	}

}

func main() {
	jobs:=make(chan int,100)
	results:=make(chan int,100)
	for i:=0;i<1000000 ;i++  {
		go worker(i,jobs,results)
	}

	for j:=1;j<=999999 ; j++ {
		jobs<-j
	}

	close(jobs)

	for w:=1;w<999999;w++ {
		<-results
	}
	}
