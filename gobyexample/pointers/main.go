package main

import "fmt"

func zeroVal(ivar int) {
	ivar++
}
func zerPtr(ivar *int) {
	*ivar++
}
func main() {
	i := 1
	zeroVal(i)
	fmt.Println(i)
	zerPtr(&i)
	fmt.Println(i)
}
