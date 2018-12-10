package main

import "fmt"

type person struct {
	name string
	num int
}

func main() {
	my_person := person{"xiao",1}
	fmt.Println(my_person,my_person.name)
}
