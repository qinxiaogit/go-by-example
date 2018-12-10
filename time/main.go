package main

import (
	"fmt"
	"time"
)

func main(){
	now :=time.Now()
	fmt.Println(now)

	then:= time.Date(1991,10,20,12,12,10,12,time.UTC)
	fmt.Println(then)
	fmt.Println(then.Year())
	fmt.Println(then.Month())
	fmt.Println(then.Day())
	fmt.Println(then.Location())
	fmt.Println(then.Before(now))
	fmt.Println(then.After(now))
	fmt.Println(then.Equal(now))
	fmt.Println(then.Sub(now))
	fmt.Println(then.Add(-then.Sub(now)))
	/******************  timestamp  *************************/

	fmt.Println(now.Unix())

	fmt.Println(now.Format(time.RFC3339))

	t1,e:=time.Parse(time.RFC3339,"2012-11-01T22:08:41+00:00")
	fmt.Println(t1,e)
}
