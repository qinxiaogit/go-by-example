package main

import "fmt"

type test struct {
	width int
	height int
}

func (r *test)circ() int {
	return (r.height+r.width)*2
}

func (r *test)area() int{
	return r.height*r.width
}

func main() {
	init_test:=test{1,2}
	fmt.Println(init_test.circ())

	fmt.Println(init_test.area())
	rp:=&init_test

	fmt.Println(rp.circ())

	fmt.Println(rp.area())
}
