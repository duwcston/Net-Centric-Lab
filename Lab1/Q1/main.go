package main

import (
	"fmt"
)

func main() {
	fmt.Println("Q1: Hamming")
	var size int
	fmt.Print("Input the DNA size: ")
	fmt.Scan(&size)
	for i := 1; i <= 1000; i++ {
		fmt.Println("Pair number", i)
		fmt.Println(hammingDistance(size))
		fmt.Println("====================================")
	}
}
