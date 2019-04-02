package main

import "fmt"

func main() {

	done:= make(chan int)

	go func() {
		fmt.Println("input...")
		done <- 1
	}()
	for {
		select {
		case <-done:
			fmt.Println("break")
			return
		default:
			fmt.Println("123")
		}
	}

}
