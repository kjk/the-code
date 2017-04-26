package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func copyAndCapture(w io.Writer, r io.Reader) []byte {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			break
		}
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			os.Stdout.Write(d)
		}
	}
	return out
}

func main() {
	cmd := exec.Command("ls", "-lah")
	var stdout, stderr []byte
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()

	go func() {
		stdout = copyAndCapture(os.Stdout, stdoutIn)
	}()

	go func() {
		stderr = copyAndCapture(os.Stderr, stderrIn)
	}()

	err := cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
