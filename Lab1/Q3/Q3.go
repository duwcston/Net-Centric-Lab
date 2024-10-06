package main

import (
	"fmt"
	"strconv"
	"strings"
)

func doubling(number string) int {
	product, err := strconv.Atoi(number)
	if err != nil {
		// ... handle error
		panic(err)
	}
	if product*2 > 9 {
		product = product*2 - 9
	} else {
		product *= 2
	}
	return product
}

func validating(number string) bool {
	fmt.Println("Number: " + number)
	number = strings.ReplaceAll(number, " ", "")
	if len(number) <= 1 {
		return false
	}
	total := 0
	for i := len(number) - 2; i >= 0; i -= 2 {
		total += doubling(string(number[i]))
	}
	for i := len(number) - 1; i >= 0; i -= 2 {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			// ... handle error
			panic(err)
		}
		total += digit
	}
	fmt.Printf("Check sum: %d\n", total)

	return total%10 == 0
}
