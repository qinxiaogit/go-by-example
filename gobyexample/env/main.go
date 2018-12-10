package main

import (
	"fmt"
	"os"
	"strings"
)

func main(){
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0])
	}
}