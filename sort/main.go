package main

import (
	"fmt"
	"sort"
)

type ByLength[]string

func (r ByLength)Len()int{
	return len(r)
}
func (r ByLength)Less(start int,end int)bool{
	return len(r[start])<len(r[end])
}
func (r ByLength)Swap(start,end int){
	r[start],r[end] = r[end],r[start]
}


func main() {
	strs:=[]string{"1","hello ","world","c","c++","php","python","golang"}

	//sort.Strings(strs)

	fmt.Println(strs)

	ints:= []int{1,4,9,0,129,198,2}

	sort.Ints(ints)
	fmt.Println(ints)
	s:=sort.IntsAreSorted(ints)
	fmt.Println(s)
	/****        sort      by      function       ***/

	sort.Sort(ByLength(strs))
	fmt.Println(strs)

}
