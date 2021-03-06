package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
)

func main() {

	newRandStream := func(done chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done:= make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints: ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	done <- 1

	fmt.Println("goroutine nums: " + strconv.Itoa(runtime.NumGoroutine()))


}
