package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/klauspost/compress/zstd"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getDataToCompress2() []byte {
	d, err := os.ReadFile("index.js")
	panicIfErr(err)
	return d
}

func zstdCompress(d []byte, concurrency int) []byte {
	var buf bytes.Buffer
	opts := []zstd.EOption{zstd.WithEncoderLevel(zstd.SpeedBetterCompression)}
	if concurrency > 0 {
		opts = append(opts, zstd.WithEncoderConcurrency(concurrency))
	}
	w, err := zstd.NewWriter(&buf, opts...)
	panicIfErr(err)
	_, err = w.Write(d)
	panicIfErr(err)
	err = w.Close()
	panicIfErr(err)
	return buf.Bytes()
}

func main() {
	fmt.Printf("Test of zstd compression\n")
	d := getDataToCompress2()
	n := runtime.GOMAXPROCS(0)
	for i := 0; i <= n; i++ {
		startTime := time.Now()
		dc := zstdCompress(d, i)
		dur := time.Since(startTime)
		fmt.Printf("Concurrency %d: %d -> %d in %s\n", i, len(d), len(dc), dur)
	}
}
