package main

import (
	"fmt"
	"log"
)

// To run:
// go run next.go

type EvenNumberIterator struct {
	max  int
	currValue int
	err  error
}

func NewEvenNumberIterator(max int) *EvenNumberIterator {
	var err error
	if max < 0 {
		err = fmt.Errorf("'max' is %d, should be >= 0", max)
	}
	return &EvenNumberIterator{
		max:  max,
		currValue: 0,
		err:  err,
	}
}

func (i *EvenNumberIterator) Next() bool {
	if i.err != nil {
		return false
	}
	i.currValue += 2
	return i.currValue <= i.max
}

func (i *EvenNumberIterator) Value() int {
	if i.err != nil {
		panic(i.err.Error())
	}
	return i.currValue
}

func (i *EvenNumberIterator) Err() error {
	return i.err
}

func printEvenNumbers(max int) {
	iter := NewEvenNumberIterator(max)
	for iter.Next() {
		n := iter.Value()
		fmt.Printf("n: %d\n", n)
	}
	if iter.Err() != nil {
		log.Fatalf("error: %s\n", iter.Err())
	}
}

func main() {
	printEvenNumbers(7)
}