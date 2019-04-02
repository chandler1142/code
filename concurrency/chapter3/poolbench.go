package main

import "fmt"

func main() {

	stringStream := make(chan string)
	go func() {
		stringStream <- "hello world"
	}()
	fmt.Println(<-stringStream)

}

