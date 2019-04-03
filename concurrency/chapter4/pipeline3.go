package main

import (
	"fmt"
)

func main() {

	repeat := func(done <-chan interface{}, values ...int) <-chan int {
		valueStream := make(chan int)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan int, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			index := 0
			for v := range valueStream {
				select {
				case <-done:
					return
				case takeStream <- v:
				}
				if index >= num {
					break
				}
				index++
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}

}
