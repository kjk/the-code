package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup

	httpSrv := makeHTTPServer()
	httpSrv.Addr = ":5000"

	go func() {
		wg.Add(1)
		fmt.Printf("Starting http server on %s\n", httpSrv.Addr)
		err := httpSrv.ListenAndServe()
		// mute error caused by Shutdown()
		if err == http.ErrServerClosed {
			err = nil
		}
		fatalIfErr(err)
		fmt.Printf("HTTP server shutdown gracefully\n")
		wg.Done()
	}()

	// Note: we could listen on https the same way

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt /* SIGINT */, syscall.SIGTERM)
	// wait for a signal
	sig := <-c
	fmt.Printf("Got signal %s\n", sig)
	// will cause ListenAndServe() stop with http.ErrServerClosed error
	httpSrv.Shutdown(nil)
	wg.Wait()
	fmt.Printf("did shutdown http servers\n")
	// add any other cleanup like closing database connection, closing log files etc.
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleMainPage)
	// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	s := `Welcome to graceful shotdown test`
	servePlainText(w, s)
}

// helper functions

func servePlainText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func fmtArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format := args[0].(string)
	if len(args) == 1 {
		return format
	}
	return fmt.Sprintf(format, args[1:]...)
}

func fatalIfErr(err error, args ...interface{}) {
	if err == nil {
		return
	}
	s := fmtArgs(args...)
	if s == "" {
		s = err.Error()
	}
	panic(s)
}
