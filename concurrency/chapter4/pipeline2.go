package main

import (
	"fmt"
)

func main() {

	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				//将元素添加入intStream中
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int, ) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- multiplier * i:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int, ) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- additive + i:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	intStream := generator(done, 1, 2, 3, 4, )

	outStream := multiply(done, add(done, intStream, 3), 4)

	go func() {
		//time.Sleep(1*time.Second)
		fmt.Println("break pipeline")
		done <- 1
	}()

	for i := range outStream {
		fmt.Println(i)
	}

}
