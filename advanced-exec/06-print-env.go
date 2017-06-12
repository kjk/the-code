package main

import (
	"fmt"
	"os"
)

const (
	envName = "MY_TEST_ENV_VARIABLE"
)

func main() {
	envVal := os.Getenv(envName)
	if envVal == "" {
		fmt.Printf("Env variable '%s' is not set\n", envName)
	} else {
		fmt.Printf("Env variable '%s' is set to '%s'\n", envName, envVal)
	}
}
