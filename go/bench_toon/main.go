package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alpkeskin/gotoon"
	"github.com/toon-format/toon-go"
)

type S3Config struct {
	ID                 int
	UserID             int
	Endpoint           string
	AccessKey          string
	SecretKey          string
	Bucket             string
	Region             string
	DirPrefix          string
	Name               string
	Description        string
	EmailNotifications string
	UploadPageURL      string
}

var testConfig = &S3Config{
	ID:                 1,
	UserID:             88,
	Endpoint:           "https://s3.amazonaws.com",
	AccessKey:          "asdf;lkj;lkasdfj",
	SecretKey:          "asldfk;lasjdf;alsjdf;lj;alsdf;lj;asdf;",
	Bucket:             "mybucket",
	Region:             "us-east-1",
	DirPrefix:          "uploads/",
	Name:               "My S3 Config",
	Description:        "A test S3 configuration",
	EmailNotifications: "daily",
	UploadPageURL:      "https://example.com/upload",
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if false {
		d, err := gotoon.Encode(testConfig)
		panicIfErr(err)
		fmt.Printf("%s", string(d))
	}
	if false {
		d, err := toon.Marshal(testConfig)
		panicIfErr(err)
		fmt.Printf("%s", string(d))
	}
	cmd := exec.Command("go", "test", "-bench=.")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running `%s` benchmarks: %v\n", cmd.String(), err)
	}
}
