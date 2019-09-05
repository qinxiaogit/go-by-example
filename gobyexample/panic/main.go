package main

import (
	"os"
)

func main() {


	//panic("this is a problem")

	_,err:= os.Open("/tmp/test0.file")
	if err!=nil{
		panic(err)
	}
}
