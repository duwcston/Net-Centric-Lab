package main

import (
	"fmt"
	"math/rand"
)

var dna = []string{"A", "T", "C", "G"}

func randomDNA(size int) []string {
	DNA := make([]string, size)

	for i := range size {
		DNA[i] = dna[rand.Intn(4)]
	}
	return DNA 
}

func hammingDistance(size int) int {
	var DNA1 = randomDNA(size)
	var DNA2 = randomDNA(size)
	
	if len(DNA1) != len(DNA2) {
		return -1
	}

	distance := 0
	for i := range DNA1 {
		if DNA1[i] != DNA2[i] {
			distance++
		}
	}

	fmt.Println("DNA1:", DNA1)
	fmt.Println("DNA2:", DNA2)
	fmt.Printf("Hamming distance: ")
	return distance
}
