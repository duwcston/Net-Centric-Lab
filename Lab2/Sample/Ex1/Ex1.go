package main

import (
	"fmt"
	"strconv"
)

func main() {
	ch := make(chan string)
	for i := 0; i < 10; i++ {
		go func(i int) {
			for j := 0; j < 10; j++ {
				ch <- "Goroutine : " + strconv.Itoa(i)
			}
		}(i)
	}
	for k := 1; k <= 100; k++ {
		fmt.Println(k, <-ch)
	}
}

//This go program init 10 goroutines to write value to a shared channel
//The main goroutine will read value from the channel and print it out
//The output will be 100 lines of "Goroutine : x" where x is the goroutine id
//The output will be in order of the goroutine id
