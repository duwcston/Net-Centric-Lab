package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	// Add the number of goroutines to wait for
	wg.Add(3)
	go performTask("Task 1", &wg)
	go performTask("Task 2", &wg)
	go performTask("Task 3", &wg)
	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All tasks completed.")
}
func performTask(taskName string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Starting %s\n", taskName)
	time.Sleep(2 * time.Second) // Simulate some work
	fmt.Printf("Completed %s\n", taskName)
}

// This go program use WaitGroup to synchroinze goroutines
// The main goroutine will init 3 goroutines to perform some tasks
// The main goroutine will wait for all the goroutines to complete
// The output will be the start and end of each task
// The output will be in order of the task number
// The output will be "All tasks completed." at the end
