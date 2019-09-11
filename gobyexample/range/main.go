package main

import "fmt"

func main() {
	//rg :=make([]string,10,10)
	rg := []int{1, 3, 5, 7}
	sum := 0
	for _, num := range rg {
		sum += num
	}
	fmt.Println(rg, sum)

	strstr := make(map[string]string)
	strstr["hello"] = "world"
	strstr["banana"] = "香蕉"
	for key, value := range strstr {
		fmt.Println(key, value)
	}
	kvs := map[string]string{"a": "apple", "b": "banana"}

	fmt.Println(kvs, "\n", strstr, kvs["a"])
	for _, k := range kvs["a"] {
		fmt.Printf("%x\n", k)
	}
}
