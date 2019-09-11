package main

import (
	"fmt"
	"math"
)

func main() {

	a := 10
	if a/2 > 3 {
		fmt.Println(math.Pow(float64(a), 2))
	} else {
		fmt.Println(math.Sqrt(float64(a)))
	}
	if num := 12; num > 13 {
		fmt.Println(num)
	} else if num < 10 {
		fmt.Println(num / 10)
	} else {
		fmt.Println(num * 10)
	}
	//fmt.Println(num)
}
