package main

import (
	"fmt"
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

 	//doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
	//	terminated := make(chan interface{})
	//	go func() {
	//		defer fmt.Println("doWork exited")
	//		defer close(terminated)
	//		for {
	//			select {
	//			case s := <-strings:
	//				fmt.Println(s)
	//			case t := <-done:
	//				fmt.Println(t)
	//				return
	//			}
	//		}
	//	}()
	//	return terminated
	//}
	//
	//done := make(chan interface{})
	//input := make(chan string)
	//terminated := doWork(done, input)
	//
	//go func() {
	//	input <- "hello"
	//	input <- "world"
	//}()
	//go func() {
	//	time.Sleep(1 * time.Second)
	//	fmt.Println("Canceling doWork goroutine...")
	//	//close(done)
	//	done <- 1
	//}()
	//
	//<-terminated
	//fmt.Println("Done.")

}
