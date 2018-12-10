package main

import "fmt"

func main()  {
	var a[5] int
	fmt.Print(a)
	b :=[5] int64{1,2,3,4}
	fmt.Println(b)
	var c[10][20] int
	for i:=0;i<10;i++{
		for  j:=0; j<20; j++ {
			c[i][j] = i*j
		}
		fmt.Println(c[i])
	}

}
