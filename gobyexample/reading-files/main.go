package main

import (
	"fmt"
	"io/ioutil"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data2, err := ioutil.ReadFile("/tmp/test.file")
	checkErr(err)
	fmt.Println(string(data2))
}
