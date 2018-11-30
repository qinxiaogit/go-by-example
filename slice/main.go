package main

import "fmt"

func main(){
	s:= make([]string,3)
	fmt.Print(s)
	s[0] = "hello world"
	s[1] = "hello world"
	s[2] = "hello world"
	s = append(s, "xiaoming")

	fmt.Println(s[:])
	fmt.Println(make([]int,10,10))
//	c :=s
//	fmt.Println(len(s))
//	twoD := make([][][] string,10)
//	fmt.Println(twoD)
//	for i:=0;i<10 ;i++  {
		//twoD[0][i] = append(twoD[0][i],"hello")
	//}
	//fmt.Println(twoD)
}
