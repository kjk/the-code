package main

// To run:
// go run main.go

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

var (
	flgVersion bool
)

func parseCmdLineFlags() {
	flag.BoolVar(&flgVersion, "version", false, "if true, print version and exit")
	flag.Parse()
	if flgVersion {
		fmt.Printf("Build on %s from sha1 %s\n", buildTime, sha1ver)
		os.Exit(0)
	}
}

func servePlainText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

// /app/debug
func handleDebug(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("url: %s %s", r.Method, r.RequestURI)
	a := []string{s}

	a = append(a, "Headers:")
	for k, v := range r.Header {
		if len(v) == 0 {
			a = append(a, k)
		} else if len(v) == 1 {
			s = fmt.Sprintf("  %s: %v", k, v[0])
			a = append(a, s)
		} else {
			a = append(a, "  "+k+":")
			for _, v2 := range v {
				a = append(a, "    "+v2)
			}
		}
	}

	a = append(a, "")
	a = append(a, fmt.Sprintf("ver: https://github.com/kjk/go-cookbook/commit/%s", sha1ver))
	a = append(a, fmt.Sprintf("built on: %s", buildTime))

	s = strings.Join(a, "\n")
	servePlainText(w, s)
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}

	mux.HandleFunc("/app/debug", handleDebug)

	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func startHTTPServer() {
	httpAddr := "127.0.0.1:4040"
	httpSrv := makeHTTPServer()
	httpSrv.Addr = httpAddr
	fmt.Printf("Visit http://%s/app/debug\n", httpAddr)
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListendAndServe() failed with %s\n", err)
	}
}

func main() {
	parseCmdLineFlags()

	startHTTPServer()
}
