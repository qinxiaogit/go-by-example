package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Printf("%v", os.Args)
	currentMsc := time.Now().Nanosecond()
	for index, vo := range os.Args {
		fmt.Printf("%d-%s\n", index, vo)
	}
	lastMsc := time.Now().Nanosecond()
	fmt.Println(lastMsc, currentMsc, lastMsc-currentMsc)
}
