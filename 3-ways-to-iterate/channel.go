package main

import (
	"fmt"
	"log"
)

// To run:
// go run channel.go

// IntWithError combines an integer value and an error
type IntWithError struct {
	Value int
	Err   error
}

func generateEvenNumbers(max int) chan IntWithError {
	ch := make(chan IntWithError)
	go func() {
		if max < 0 {
			ch <- IntWithError{
				Value: 0,
				Err:   fmt.Errorf("'max' is %d and should be >= 0", max),
			}
		}
		for i := 2; i <= max; i += 2 {
			ch <- IntWithError{
				Value: i,
				Err:   nil,
			}
		}
		close(ch)
	}()
	return ch
}

func printEvenNumbers(max int) {
	for val := range generateEvenNumbers(max) {
		if val.Err != nil {
			log.Fatalf("Error: %s\n", val.Err)
		}
		fmt.Printf("%d\n", val.Value)
	}
}

func main() {
	printEvenNumbers(8)
	printEvenNumbers(-1)
}
