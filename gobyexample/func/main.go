package main

import "fmt"

func main() {
	fmt.Println(plus(1, '2'))
	fmt.Println(returnManyValue(1, '2'))
	fmt.Println(sum(1, 2, 3, 4, 5, 6, 7, 8, 9))
	args := []int{1, 3, 5, 7, 9}
	fmt.Println(sum(args...)) //分片传变参需要加...
}

/*
 * 带参函数
 */
func plus(a int, b int) interface{} {
	return float64(a + b)
}

/*
 * 多返回值
 */
func returnManyValue(a int64, b int64) (interface{}, int64) {
	return a, a + b
}

/*
 * 动态参数
 */
func sum(args ...int) int {
	var sum = 0
	for _, item := range args {
		sum += item
	}
	return sum
}
