package main

import (
	"fmt"
	"log"

	"github.com/kjk/common/filerotate"
)

func main() {
	didClose := func(path string, didRotate bool) {
		fmt.Printf("closed file '%s', didRotate: %v\n", path, didRotate)
		if didRotate {
			// here you can implement e.g. compressing a log file and
			// uploading them to s3 for long-term storage
		}
	}
	f, err := filerotate.NewDaily(".", "log.txt", didClose)
	if err != nil {
		log.Fatalf("filerotate.NewDaily() failed with '%s'\n", err)
	}
	_, err = f.Write([]byte("hello"))
	if err != nil {
		log.Fatalf("f.Write() failed with '%s'\n", err)
	}
	err = f.Close()
	if err != nil {
		log.Fatalf("f.Close() failed with '%s'\n", err)
	}
}
