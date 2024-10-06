package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Q2: Scrabble Score")
	var words string
	fmt.Print("Type a string: ")
	fmt.Scan(&words)
	// fmt.Scanf("%s", &words) // Use for multiple words
	fmt.Println(calculateScarbbleScore(strings.ToUpper(words)))
	fmt.Println("====================================")
}
