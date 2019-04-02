package main

import (
	"fmt"
	"sync"
)

func main() {

	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	i := 0
	subscribe := func(cond *sync.Cond, fn func()) {
		//cond.L.Lock()
		//defer cond.L.Unlock()
		//cond.Wait()
		fn()
		if i == 2 {
			cond.Signal()
		}
	}

	go subscribe(cond, func() {
		fmt.Println("this is goroutine 1")
		i = i + 1
	})

	go subscribe(cond, func() {
		fmt.Println("this is goroutine 2")
		i = i + 1
	})

	for i != 2 {
		cond.Wait()
	}
	//cond.Broadcast()
	cond.L.Unlock()

}
