package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
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

var fanIn = func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	rand := func() interface{} {
		return rand.Intn(500000000)
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n ", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Println("working goroutine: " + strconv.Itoa(runtime.NumGoroutine()))

		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))

}
