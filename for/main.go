package main

import "fmt"

func main(){
	i:=1
	for ;i<20;i++{
		fmt.Println(i)
	}
	for {
		if(i<2){
			break;
		}
		println("i:",i)
		i--;
	}
	for i<200 {
		fmt.Println(i)
		i++
	}
}
