package main

import (
	"fmt"
	"strings"
	s "strings"
)

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

func Any(vs []string, callback func(string) bool) bool {
	for _, v := range vs {
		if callback(v) {
			return true
		}
	}
	return false
}

func All(vs []string, callback func(string) bool) bool {

	for _, v := range vs {
		if !callback(v) {
			return false
		}
	}
	return true
	//var tempStr[]string
	//for _,v:=range vs{
	//	tempStr = append(tempStr, callback(v))
	//}
	//return tempStr
}

func Fiter(vs []string, callback func(string) bool) []string {

	var tmpStr []string
	for _, v := range vs {
		if callback(v) {
			tmpStr = append(tmpStr, v)
		}
	}
	return tmpStr
}

func Map(vs []string, callback func(string) string) []string {

	var tmpStr = make([]string, len(vs))
	for index, v := range vs {
		tmpStr[index] = callback(v)
	}
	return tmpStr
}

/************************** string function *****************************/

var p = fmt.Println

func main() {
	var strs = []string{"peach", "apple", "pear", "plum"}

	fmt.Println(Index(strs, "pear"))
	fmt.Println(Include(strs, "pear"))
	fmt.Println(All(strs, func(s string) bool {
		return strings.HasPrefix(s, "p")
	}))

	fmt.Println(Any(strs, func(s string) bool {
		return true
	}))

	fmt.Println(Fiter(strs, func(s string) bool {
		return strings.Contains(s, "e")
	}))
	fmt.Println(Map(strs, strings.ToUpper))

	/************************** string function *****************************/
	p("Contains:  ", s.Contains("test", "es"))
	p("Count:     ", s.Count("test", "t"))
	p("HasPrefix: ", s.HasPrefix("test", "te"))
	p("HasSuffix: ", s.HasSuffix("test", "st"))
	p("Index:     ", s.Index("test", "e"))
	p("Join:      ", s.Join([]string{"a", "b"}, "-"))
	p("Repeat:    ", s.Repeat("a", 5))
	p("Replace:   ", s.Replace("foo", "o", "0", -1))
	p("Replace:   ", s.Replace("foo", "o", "0", 1))
	p("Split:     ", s.Split("a-b-c-d-e", "-"))
	p("ToLower:   ", s.ToLower("TEST"))
	p("ToUpper:   ", s.ToUpper("test"))
	p()
	p("char", "hello"[1])
	p("len", len("len"))

}
