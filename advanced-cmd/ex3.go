package main

// to execute: go run ex3.go

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "-lah")
	var stdout, stderr []byte
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()

	go func() {
		buf := make([]byte, 1024, 1024)
		for {
			n, err := stdoutIn.Read(buf)
			if err != nil {
				break
			}
			if n > 0 {
				d := buf[:n]
				stdout = append(stdout, d...)
				os.Stdout.Write(d)
			}
		}
	}()

	go func() {
		buf := make([]byte, 1024, 1024)
		for {
			n, err := stderrIn.Read(buf)
			if err != nil {
				break
			}
			if n > 0 {
				d := buf[:n]
				stderr = append(stderr, d...)
				os.Stderr.Write(d)
			}
		}
	}()

	err := cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
