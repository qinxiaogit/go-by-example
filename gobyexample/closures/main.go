package main

import "fmt"

func intSeq() func() int{
	i:=0
	return func() int {
		i+=1
		return i
	}
}
func main(){
	nextInt := intSeq()
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	nextIntTwo := intSeq()
	fmt.Println(nextIntTwo())
	//递归
	fmt.Println(fact(10))
}

/**
 * 递归
 */
func fact(n int) int{
	if(n<=1){
		return n
	}
	return n*fact(n-1)
}