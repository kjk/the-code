package main

// https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html
// To run:
// go run main.go

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chilts/sid"
	guuid "github.com/google/uuid"
	"github.com/kjk/betterguid"
	"github.com/lithammer/shortuuid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
)

func genShortUUID() {
	id := shortuuid.New()
	fmt.Printf("github.com/lithammer/shortuuid: %s\n", id)
}

func genUUID() {
	id := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", id.String())
}

func genXid() {
	id := xid.New()
	fmt.Printf("github.com/rs/xid:              %s\n", id.String())
}

func genKsuid() {
	id := ksuid.New()
	fmt.Printf("github.com/segmentio/ksuid:     %s\n", id.String())
}

func genBetterGUID() {
	id := betterguid.New()
	fmt.Printf("github.com/kjk/betterguid:      %s\n", id)
}

func genUlid() {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	fmt.Printf("github.com/oklog/ulid:          %s\n", id.String())
}

func genSonyflake() {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	// Note: this is base16, could shorten by encoding as base62 string
	fmt.Printf("github.com/sony/sonyflake:      %x\n", id)
}

func genSid() {
	id := sid.Id()
	fmt.Printf("github.com/chilts/sid:          %s\n", id)
}

func genUUIDv4() {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("uuid.NewV4() failed with %s\n", err)
	}
	fmt.Printf("github.com/satori/go.uuid:      %s\n", id)
}

func main() {
	genXid()
	genKsuid()
	genBetterGUID()
	genUlid()
	genSonyflake()
	genSid()
	genShortUUID()
	genUUIDv4()
	genUUID()
}
