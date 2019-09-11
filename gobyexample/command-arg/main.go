package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	argsWithProg := os.Args
	argsWithOutProg := os.Args[1:]

	//arg := os.Args[3]
	fmt.Println(argsWithProg)
	fmt.Println(argsWithOutProg)
	//fmt.Println(arg)

	wordPtr := flag.String("world", "hello", "golang")

	numPtr := flag.Int("numb", 10, "an int")

	var svar string

	flag.StringVar(&svar, "svar", "bar", "a string val")
	flag.Parse()
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numPtr)
	//fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())

}
