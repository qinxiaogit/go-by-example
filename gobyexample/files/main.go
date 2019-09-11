package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//dat,err:= ioutil.ReadFile("./main.go")
	//check(err)
	//fmt.Printf("%s",dat)

	f, err := os.OpenFile("./files/temp", os.O_APPEND, os.FileMode(1))
	check(err)
	//fmt.Println(f)
	b1 := make([]byte, 5)
	f.Seek(6, 0)
	n1, err := f.Read(b1)
	check(err)
	fmt.Println(n1, string(b1))
	//f.Write([]byte{115, 111, 109, 101, 10})
	f.WriteString("大家好")
	f.Close()
}
