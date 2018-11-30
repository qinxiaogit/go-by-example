package main

import "fmt"

func main(){
	m:=make(map[string]string)
	m["string"] = "world"
	name,value := m["string"]
	fmt.Print(m)
	fmt.Print(name,value)
	m2 :=make(map[string] map[string] string)
	m2["hello"] = m
	fmt.Println(m2)
	//m2:=make(map(make([]string,20) int))
}
