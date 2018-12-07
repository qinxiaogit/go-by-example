package main

import (
	"fmt"
	"time"
)

func main(){
	i:=2
	switch i {
	case 3:
		fmt.Println("hello:",i)
		break
	default:
		fmt.Println("default")
		break

	}
	switch time.Now().Weekday() {
	case time.Sunday,time.Friday:
			fmt.Println("hello ")
			break
	case time.Monday:
		fmt.Println("world")
		break

	}
	t := time.Now()
	switch  {
	case t.Hour()<12:
		fmt.Println("before:",t.Hour())
		break
	default:
		fmt.Println("after:",t.Hour())

	}
}
