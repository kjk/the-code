package main

import (
	"fmt"
	"log"
)

// To run:
// go run next.go

// EvenNumberIterator generates even number
type EvenNumberIterator struct {
	max       int
	currValue int
	err       error
}

// NewEvenNumberIterator creates new number iterator
func NewEvenNumberIterator(max int) *EvenNumberIterator {
	var err error
	if max < 0 {
		err = fmt.Errorf("'max' is %d, should be >= 0", max)
	}
	return &EvenNumberIterator{
		max:       max,
		currValue: 0,
		err:       err,
	}
}

// Next advances to next even number. Returns false on end of iteration.
func (i *EvenNumberIterator) Next() bool {
	if i.err != nil {
		return false
	}
	i.currValue += 2
	return i.currValue <= i.max
}

// Value returns current even number
func (i *EvenNumberIterator) Value() int {
	if i.err != nil || i.currValue > i.max {
		panic("Value is not valid after iterator finished")
	}
	return i.currValue
}

// Err returns iteration error.
func (i *EvenNumberIterator) Err() error {
	return i.err
}

func printEvenNumbers(max int) {
	iter := NewEvenNumberIterator(max)
	for iter.Next() {
		fmt.Printf("n: %d\n", iter.Value())
	}
	if iter.Err() != nil {
		log.Fatalf("error: %s\n", iter.Err())
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
