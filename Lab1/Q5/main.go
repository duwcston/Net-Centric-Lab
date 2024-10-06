package main

import (
	"fmt"
)

func main() {
	fmt.Println("Q5: Matching Brackets")
	typeOfBracket := "{{[]}}"
	fmt.Println(typeOfBracket)
	if isCorrect(typeOfBracket) {
		fmt.Println("The brackets are matched and nested correct.")
	} else {
		fmt.Println("The brackets are matched and nested incorrect.")
	}
	fmt.Println("====================================")
}
