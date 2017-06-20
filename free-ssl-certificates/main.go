package main

// To run:
// go run main.go

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	htmlIndex    = `<html><body>Welcome!</body></html>`
	inProduction = false
	httpPort     = "127.0.0.1:8080"
)

func handeIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlIndex)
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handeIndex)

	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}

func main() {
	var httpsSrv, httpSrv *http.Server
	if inProduction {
		dataDir := "."
		hostPolicy := func(ctx context.Context, host string) error {
			// Note: change to your real domain
			allowedHost := "www.mydomain.com"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		httpsSrv = makeHTTPServer()
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}
		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			fmt.Printf("Starting HTTPS server on %s\n", httpsSrv.Addr)
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.LstendAndServeTLS() failed with %s", err)
			}
		}()
	}

	httpSrv = makeHTTPServer()
	httpSrv.Addr = httpPort
	fmt.Printf("Starting HTTP server on %s\n", httpPort)
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}
}
