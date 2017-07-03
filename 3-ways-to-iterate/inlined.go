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
	for i := 2; i < max; i += 2 {
		fmt.Printf("%d\n", i)
	}
}

func main() {
	printEvenNumbers(7)
}
