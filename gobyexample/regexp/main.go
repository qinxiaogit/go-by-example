package main

import (
	"fmt"
	"regexp"
)

func main() {

	match ,_:=regexp.MatchString("\\d*","18227755589")
	fmt.Println(match)
	r,_:=regexp.Compile("\\d*")
	fmt.Println(r.FindString("18227a755589"))

}
