package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	maxStudents    = 30
	totalStudents  = 100
	maxReadingTime = 4
)

// Student represents a student in the library
type Student struct {
	id   int
	time int
}

func main() {
	// Capture the start time
	startTime := time.Now()

	// Create a channel to store students
	students := make(chan Student, maxStudents)
	// Create a wait group to wait for all students to finish
	var wg sync.WaitGroup

	// Create a map to track assigned IDs
	assignedIDs := make(map[int]bool)

	// Create a goroutine for each student
	for i := 0; i < totalStudents; i++ {
		time.Sleep(time.Second)
		duration := rand.Intn(maxReadingTime) + 1

		// Generate a unique random ID
		var id int
		for {
			id = rand.Intn(totalStudents) + 1
			if !assignedIDs[id] {
				assignedIDs[id] = true
				break
			}
		}

		student := Student{
			id:   id,
			time: duration,
		}

		// Add the student to the wait group
		wg.Add(1)
		go func(s Student) {
			defer wg.Done()
			fmt.Print("Student in the library: ", len(students), "\n")
			// Check if the channel is full
			select {
			case students <- s:
				openingTime := time.Since(startTime).Seconds() - 1
				fmt.Printf("Time %.0f: Student %d starts reading at the library\n", openingTime, s.id)
				time.Sleep(time.Second * time.Duration(s.time))
				openingTime = time.Since(startTime).Seconds() - 1
				fmt.Printf("Time %.0f: Student %d leaves the library with %d hours\n", openingTime, s.id, s.time)
				<-students
			default:
				openingTime := time.Since(startTime).Seconds() - 1
				fmt.Printf("Time %.0f: Student %d is waiting to enter the library\n", openingTime, s.id)
				students <- s
				openingTime = time.Since(startTime).Seconds() - 1
				fmt.Printf("Time %.0f: Student %d starts reading at the library\n", openingTime, s.id)
				time.Sleep(time.Second * time.Duration(s.time))
				openingTime = time.Since(startTime).Seconds() - 1
				fmt.Printf("Time %.0f: Student %d leaves the library with %d hours\n", openingTime, s.id, s.time)
				<-students
			}
			fmt.Print("Student in the library: ", len(students), "\n")
		}(student)
	}

	wg.Wait()
	go func() {
		close(students)
	}()

	// wg.Wait()
	openingTime := time.Since(startTime).Seconds()
	fmt.Printf("Time %.0f: No more students. Lets call it a day\n", openingTime)
}
