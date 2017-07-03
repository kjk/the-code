package main

import (
	"fmt"
	"log"
)

// To run:
// go run channel.go

// IntWithError combines an integer value and an error
type IntWithError struct {
	Int int
	Err error
}

func generateEvenNumbers(max int) chan IntWithError {
	ch := make(chan IntWithError)
	go func() {
		defer close(ch)
		if max < 0 {
			ch <- IntWithError{
				Err: fmt.Errorf("'max' is %d and should be >= 0", max),
			}
			return
		}

		for i := 2; i <= max; i += 2 {
			ch <- IntWithError{
				Int: i,
			}
		}
	}()
	return ch
}

func printEvenNumbers(max int) {
	for val := range generateEvenNumbers(max) {
		if val.Err != nil {
			log.Fatalf("Error: %s\n", val.Err)
		}
		fmt.Printf("%d\n", val.Int)
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
