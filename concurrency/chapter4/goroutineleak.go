package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func main() {

	doWork := func(strings <-chan string) <-chan interface{} {
		complete := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(complete)
			for s := range strings {
				fmt.Println(s)
			}
		}()
		return complete
	}

	doWork(nil)
	fmt.Println("Done.")
	time.Sleep(1 * time.Second)
	fmt.Println("goroutine nums: " + strconv.Itoa(runtime.NumGoroutine()))

}
