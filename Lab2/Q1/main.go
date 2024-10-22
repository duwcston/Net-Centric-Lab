package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// func count(word string, result map[rune]int) {
// 	startTime := time.Now()
// 	counts := make(map[rune]int)
// 	for _, r := range word {
// 		counts[r]++
// 	}
// 	for key, value := range counts {
// 		result[key] += value
// 	}
// 	fmt.Println("Time taken to count characters:", time.Since(startTime))
// }

func countGo(word string, resultChan chan<- map[rune]int, wg *sync.WaitGroup) {
	defer wg.Done()
	counts := make(map[rune]int)
	for _, r := range word {
		counts[r]++
	}
	resultChan <- counts
}

func main() {
	var wg sync.WaitGroup
	data, err := os.ReadFile("file.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println("Data read from file:", string(data))
	str := strings.ToLower(string(data))
	words := strings.Fields(str) // Use Fields to split by any whitespace

	wordChan := make(chan map[rune]int, len(words)) // Buffered channel

	wg.Add(len(words))

	startTime := time.Now()
	for _, word := range words {
		go countGo(word, wordChan, &wg)
	}

	go func() {
		wg.Wait()
		close(wordChan)
	}()

	result := make(map[rune]int)
	blankCount := strings.Count(string(data), " ")
	for wordCounts := range wordChan {
		for key, value := range wordCounts {
			result[key] += value
		}
	}

	if blankCount > 0 {
		result[' '] = blankCount
	}

	for key, value := range result {
		if key == ' ' {
			fmt.Printf("<blank>: %d\n", value)
		} else {
			fmt.Printf("%c: %d\n", key, value)
		}
	}

	fmt.Println("Time taken to count characters using Go concurrency:", time.Since(startTime))
}
