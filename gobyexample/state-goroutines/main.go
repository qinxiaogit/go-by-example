package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type readOps struct {
	key  int
	resp chan int
}

type wirteOps struct {
	key  int
	val  int
	resp chan bool
}

func main() {
	var ops int64

	reads := make(chan *readOps)
	wirtes := make(chan *wirteOps)

	go func() {

		var states = make(map[int]int)
		for {
			select {

			case wirte := <-wirtes:

				states[wirte.key] = wirte.val
				wirte.resp <- true
			case read := <-reads:
				read.resp <- states[read.key]

			}
		}

	}()

	for r := 100; r > 0; r-- {
		go func() {

			for {
				read := &readOps{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- read
				<-read.resp
				atomic.AddInt64(&ops, 1)

			}

		}()
	}

	for w := 0; w < 10; w++ {
		go func() {

			for {
				wirte := &wirteOps{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool),
				}
				wirtes <- wirte
				<-wirte.resp
				atomic.AddInt64(&ops, 1)
			}

		}()
	}

	time.Sleep(time.Second)

	osFinsh := atomic.LoadInt64(&ops)
	fmt.Println(osFinsh)

}
