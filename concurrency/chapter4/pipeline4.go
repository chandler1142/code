package main

import (
	"fmt"
	"math/rand"
)

func main() {

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Int()
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
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

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}

}
