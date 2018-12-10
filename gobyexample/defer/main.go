package main

import (
	"fmt"
	"os"
)

func main() {
	file :=CreateFile("/tmp/test.file")
	defer CloseFile(file)
	WirteFile(file)


}

func CreateFile(fileName string) *os.File{

	fmt.Println("createing...")
	f,err:=os.Create(fileName)
	if err!=nil{
		panic(err)
	}
	return f
}

func WirteFile(file *os.File){
	fmt.Println("writing")
	fmt.Fprintln(file, "data")
}

func CloseFile(file *os.File){
	fmt.Println("close file")
	file.Close()
}
