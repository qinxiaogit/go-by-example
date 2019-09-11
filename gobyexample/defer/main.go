package main

import (
	"fmt"
	"os"
)

func main() {
	file := CreateFile("/tmp/test.file")
	defer CloseFile(file)
	WirteFile(file)

	for i := 0; i < 10; i++ {
		defer test(i)
	}
}
func test(i int) {
	fmt.Println("---------test--------", i)
}

func CreateFile(fileName string) *os.File {

	fmt.Println("createing...")
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	return f
}

func WirteFile(file *os.File) {
	fmt.Println("writing")
	fmt.Fprintln(file, "data")
}

func CloseFile(file *os.File) {
	fmt.Println("close file")
	file.Close()
}
