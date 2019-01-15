package main

import "fmt"

func main()  {

	slicetest := make([]int,100)

	printSlice(slicetest)
}

func printSlice(x []int){
	fmt.Printf("len=%d cap=%d slice=%v\n",len(x),cap(x),x)
}