package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kjk/betterguid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

// To run:
// go run main.go

func main() {
	{
		id := xid.New()
		fmt.Printf("id generated with github.com/rs/xid:           %s\n", id.String())
	}

	{
		id := ksuid.New()
		fmt.Printf("id generated with github.com/segmentio/ksuid:  %s\n", id.String())

	}

	{
		id := betterguid.New()
		fmt.Printf("id generated with github.com/kjk/betterguid:   %s\n", id)
	}

	{
		t := time.Now().UTC()
		entropy := rand.New(rand.NewSource(t.UnixNano()))
		id := ulid.MustNew(ulid.Timestamp(t), entropy)
		fmt.Printf("id generated with github.com/oklog/ulid:       %s\n", id.String())

	}
}
