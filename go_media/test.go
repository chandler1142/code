package main

import "fmt"

func fun1(a *int) {
	fmt.Println(*a)
	*a--
	*a += 3
}

func main() {
	a := 3
	fun1(&a)
	fmt.Println(a)
}
