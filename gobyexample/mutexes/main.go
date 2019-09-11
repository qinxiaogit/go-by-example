package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var state = make(map[int]int)
	mutex := sync.Mutex{}
	var ops int64 = 0
	for i := 0; i < 100; i++ {
		go func() {
			for {
				total := 0
				key := rand.Intn(100)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddInt64(&ops, 1)
				runtime.Gosched()
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(100)
				val := rand.Intn(100)

				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddInt64(&ops, 1)
				runtime.Gosched()
			}
		}()
	}
	time.Sleep(time.Second * 10)
	opsFinsh := atomic.LoadInt64(&ops)
	fmt.Println("finsh：", opsFinsh)
	mutex.Lock()
	fmt.Println("state：", state)
	mutex.Unlock()
}
