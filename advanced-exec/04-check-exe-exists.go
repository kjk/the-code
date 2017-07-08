package main

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
// to run:
// go run 04-check-exe-exists.go

import (
	"fmt"
	"os/exec"
)

func checkExeExists(exe string) {
	path, err := exec.LookPath(exe)
	if err != nil {
		fmt.Printf("didn't find '%s' executable\n", exe)
		return
	}
	fmt.Printf("'%s' executable is '%s'\n", exe, path)
}

func main() {
	checkExeExists("ls")
	checkExeExists("ls2")
}
