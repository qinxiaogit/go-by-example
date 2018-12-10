package main

import (
	"os"
)

func main() {


	//panic("this is a problem")

	_,err:= os.Create("/tmp/test.file")
	if err!=nil{
		panic(err)
	}
}
