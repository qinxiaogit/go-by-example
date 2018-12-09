package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	var init int64=0
	for i:=0; i<2;i++  {
		go func() {
			for   {
				//init++
				atomic.AddInt64(&init,1)
				fmt.Println("init：",init)
				runtime.Gosched()
			}
		}()
	}
	time.Sleep(time.Second)

	fmt.Println(atomic.LoadInt64(&init))
}
