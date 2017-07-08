package main

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
// To run:
// go run 05-change-environment.go

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	envName = "MY_TEST_ENV_VARIABLE"
)

func runEnvironTest(envValue string) {
	cmd := exec.Command("go", "run", "05-print-env-helper.go")
	if envValue != "" {
		newEnv := append(os.Environ(), fmt.Sprintf("%s=%s", envName, envValue))
		cmd.Env = newEnv
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("%s", out)
}

func main() {
	runEnvironTest("")
	runEnvironTest("test value")
}
