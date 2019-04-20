package test

import (
	"testing"
)

var repeat = func(done <-chan interface{}, values ...string) <-chan string {
	valueStream := make(chan string)
	go func() {
		defer close(valueStream)
		for _, v := range values {
			select {
			case <-done:
				return
			case valueStream <- v:
			}
		}
	}()
	return valueStream
}

var take = func(done <-chan interface{}, valueStream <-chan string, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := num; i > 0 || i == -1; {
			if i != -1 {
				i--
			}
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

var toString = func(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()
	return stringStream
}

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range toString(done, take(done, repeat(done, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {

	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {

	}

}
