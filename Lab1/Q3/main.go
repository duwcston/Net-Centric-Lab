package main

import (
	"fmt"
)

func main() {

	fmt.Println("Q3: Luhn")
	if validating("4539 3195 0343 6467") {
		fmt.Println("The number is valid.")
	} else {
		fmt.Println("The number is invalid.")
	}
	fmt.Println("====================================")
}
