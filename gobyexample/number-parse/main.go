package main

import (
	"fmt"
	"strconv"
)

func main() {

	f, _ := strconv.ParseFloat("3.14", 64)
	fmt.Println(f)
	i, _ := strconv.ParseInt("1024", 0, 64)
	fmt.Println(i)
	d, _ := strconv.ParseInt("0x1c8", 0, 64)
	fmt.Println(d)
	u, _ := strconv.ParseUint("100", 0, 64)
	fmt.Println(u)

	k, _ := strconv.Atoi("132")
	fmt.Println(k)
	a, _ := strconv.Atoi("wat")
	fmt.Println(a)
}
