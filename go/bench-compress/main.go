package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/dustin/go-humanize"
	"github.com/klauspost/compress/zstd"
)

type benchResult struct {
	name string
	data []byte
	dur  time.Duration
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func logf(s string, args ...any) {
	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}
	fmt.Print(s)
}

func push[S ~[]E, E any](s *S, els ...E) {
	*s = append(*s, els...)
}

func humanSize(n int) string {
	un := uint64(n)
	return humanize.Bytes(un)
}

func benchFileCompress(path string) {
	d, err := os.ReadFile(path)
	panicIfErr(err)

	var results []benchResult
	gzipCompress := func(d []byte) []byte {
		var buf bytes.Buffer
		w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		panicIfErr(err)
		_, err = w.Write(d)
		panicIfErr(err)
		err = w.Close()
		panicIfErr(err)
		return buf.Bytes()
	}

	zstdCompress := func(d []byte, level zstd.EncoderLevel) []byte {
		var buf bytes.Buffer
		w, err := zstd.NewWriter(&buf, zstd.WithEncoderLevel(level), zstd.WithEncoderConcurrency(1))
		panicIfErr(err)
		_, err = w.Write(d)
		panicIfErr(err)
		err = w.Close()
		panicIfErr(err)
		return buf.Bytes()
	}

	brCompress := func(d []byte, level int) []byte {
		var dst bytes.Buffer
		w := brotli.NewWriterLevel(&dst, level)
		_, err := w.Write(d)
		panicIfErr(err)
		err = w.Close()
		panicIfErr(err)
		return dst.Bytes()
	}

	var cd []byte
	logf("compressing with gzip\n")
	t := time.Now()
	cd = gzipCompress(d)
	push(&results, benchResult{"gzip", cd, time.Since(t)})

	logf("compressing with brotli: default (level 6)\n")
	t = time.Now()
	cd = brCompress(d, brotli.DefaultCompression)
	push(&results, benchResult{"brotli default", cd, time.Since(t)})

	logf("compressing with brotli: best (level 11)\n")
	t = time.Now()
	cd = brCompress(d, brotli.BestCompression)
	push(&results, benchResult{"brotli best", cd, time.Since(t)})

	logf("compressing with zstd level: better (3)\n")
	t = time.Now()
	cd = zstdCompress(d, zstd.SpeedBetterCompression)
	push(&results, benchResult{"zstd better", cd, time.Since(t)})

	logf("compressing with zstd level: best (4)\n")
	t = time.Now()
	cd = zstdCompress(d, zstd.SpeedBestCompression)
	push(&results, benchResult{"zstd best", cd, time.Since(t)})

	sort.Slice(results, func(i, j int) bool {
		return len(results[i].data) < len(results[j].data)
	})
	logf("\nBy size:\n")
	for _, r := range results {
		logf("%14s: %6d (%s) in %s\n", r.name, len(r.data), humanSize(len(r.data)), r.dur)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].dur < results[j].dur
	})
	logf("\nBy time:\n")
	for _, r := range results {
		logf("%14s: %6d (%s) in %s\n", r.name, len(r.data), humanSize(len(r.data)), r.dur)
	}
}

func main() {
	path := filepath.Join("..", "test_compress", "index.js")
	benchFileCompress(path)
}
