package main

import (
	"fmt"
	"log"
)

// To run:
// go run inlined.go

func printEvenNumbers(max int) {
	if max < 0 {
		log.Fatalf("'max' is %d, should be >= 0", max)
	}
	for i := 2; i <= max; i += 2 {
		fmt.Printf("%d\n", i)
	}
}

func main() {
	fmt.Printf("Even numbers up to 8:\n")
	printEvenNumbers(8)
	fmt.Printf("Even numbers up to 9:\n")
	printEvenNumbers(9)
	fmt.Printf("Error: even numbers up to -1:\n")
	printEvenNumbers(-1)
}
