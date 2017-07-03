package main

import (
	"fmt"
	"log"
)

// To run:
// go run callback.go

func iterateEvenNumbers(max int, cb func(n int) error) error {
	if max < 0 {
		return fmt.Errorf("'max' is %d, must be >= 0", max)
	}
	for i := 2; i <= max; i += 2 {
		err := cb(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func printEvenNumbers(max int) {
	err := iterateEvenNumbers(max, func(n int) error {
		fmt.Printf("%d\n", n)
		return nil
	})
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

func main() {
	printEvenNumbers(8)
	printEvenNumbers(-1)
}
