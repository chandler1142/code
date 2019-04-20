package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var primeFinder = func(done <-chan interface{}, valueStream <-chan int) <-chan interface{} {
	primeValueStream := make(chan interface{})
	go func() {
		for v := range valueStream {
			select {
			case <-done:
				return
			default:
				for i := 2; i < v; i++ {
					if v%i == 0 {
						break
					}
					if i == v-1 {
						primeValueStream <- v
					}
				}
			}
		}
	}()
	return primeValueStream
}

var repeatFn = func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
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

var toInt = func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case intStream <- v.(int):
			}
		}
	}()
	return intStream
}

var take = func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	resultStream := make(chan interface{})
	go func() {
		defer close(resultStream)
		index := 0
		for v := range valueStream {
			select {
			case <-done:
				return
			case resultStream <- v:
			}
			index ++
			if index >= num {
				break
			}
		}
	}()
	return resultStream
}

func main() {
	rand := func() interface{} {
		return rand.Intn(500000000)
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes: ")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Println("working goroutine: " + strconv.Itoa(runtime.NumGoroutine()))
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))

}
