package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

// To run:
// go run main.go

func main() {
	idXid := xid.New()
	idKsuid := ksuid.New()
	t := time.Unix(1000000, 0)
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	idUlid := ulid.MustNew(ulid.Timestamp(t), entropy)

	fmt.Printf("id generated with github.com/rs/xid:          %s\n", idXid.String())
	fmt.Printf("id generated with github.com/segmentio/ksuid: %s\n", idKsuid.String())
	fmt.Printf("id generated with github.com/oklog/ulid:      %s\n", idUlid.String())
}
