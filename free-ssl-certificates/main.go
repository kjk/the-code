package main

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
	inProduction = true
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
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.LstendAndServeTLS() failed with %s", err)
			}
		}()
	}

	httpSrv = makeHTTPServer()
	httpSrv.Addr = ":80"
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}
}

/*

## Why HTTPS?

Having HTTPS for your website is important.

HTTPS protocol uses SSL to encrypt the traffic between browser and server.

If the browser sends confidential information to the server (e.g. username/password) you must encrypt it. Otherwise any random person sitting next to your user in a cafe might be sniffing wifi traffic and stealing their credentials.

If you care about SEO, Google ranks HTTPS websites higher than HTTP websites.

## Getting someone else to give you certificate.

Before we learn how to support HTTPS directly in your Go server, let's talk about simpler options.

You can use a third-party service like Cloudflare as a layer on top of your service. One of the features they offer (even in the free tier) is HTTPS proxy.

You only need to provide HTTP endpoint. You change DNS of your domain to point to CloudFlare's IP address and configure HTTPS redirect in their web interface by providing IP address of your server.

Browser talks to CloudFlare, which takes care of provisioning SSL certificate and proxies the traffic to your server.

AWS, Google Cloud and other hosting providers also provide this for servers hosted on their infrastructures.

## Directly supporting HTTPS

Not so long ago if you wanted a certificate, you had to pay few hundred dollars a year for a single domain.

Let's Encrypt is a non-profit organization that provides certificates for free and has HTTP API for obtaining certificates, which allows automating the process.

Before Let's Encrypt you would buy a certificate, which is just a bunch of bytes. You would typically save the data to a file on the server and configure the server with it.

With Let's Encrypt you can use their API to ask for the certificate at startup. Thankfully, all the hard work of implementing it has already bee done by others. There are a couple of Go libraries that implement Let's Encrypt support. I've been using golang.org/x/crypto/acme/autocert for several months now. It's been stable and has additional benefit of coming from members of Go project.

Here's how to start an HTTPS web server that uses free SSL certificates from Let's Encrypt:

```go
```

There are some important things to note.

1\. The standard port for HTTPS is 443

2\. You can run only HTTP, only HTTPS or both.

3\. If the server doesn't have certificate, it'll use HTTP API to ask Let's Encrypt servers for it.

Those requests are throttled to 20 per week to avoid over-loading Let's Encrypt servers.

It's therefore important to cache the certificate somewhere. In our example we cache them on disk, using autocert.DirCache() cache.

Cache is an interface so you could implemnt your own storage e.g. in a SQL database or Redis.

4\. You must have DNS set up correctly. The way protocol works, if you ask for a certificate for "www.mydomain.com", you need to provide HTTP callback on that exact doamin. If that DNS name doesn't resolve to the IP address of your server, you'll not get your certificate.

That makes local testing of HTTPS support hard.

5\. You might be wondering: what is this HostPolicy business?

As I mentioned, Let's Certificate throttles certificate provisioning so you need to ensure the server won't ask for certificates for domains you don't care about. Autocert docs [explain this well](https://...)

Our example assumes most common case: a server that only responds to a single domain. You can easily customize the logic.

## How free certificates from Let's Encrypt came to be

Arguably due to a design mistake, SSL protocol not only encrypts but also proves site's identity to the browser. It provided accountability so that we can trace the ownership of google.com and see that it is indeed owned by Google, Inc in US, and not Ivan The Hacker in Moscow.

We implement that accountability by trusting a very small number of companies (Certificate Authorities) to issue certificates that prove the identity of the website owner.

A certificate is just a bunch of bytes constructed with clever cryptographic methods.

Website owner then configures his web server with that certificate.

Browsers check the certificate and if they trust it's legit, they use the crypto bits in certificate to secure the connection. If they don't trust the certificate, they warn the user that something fishy is going on and block the website.

When you apply for a certificate, Certificate Authority has to verify your identity by asking you to submit necessary papwerwork and ensuring it's valid.

Verifying identity requires labor. Keeping certificates safe requires labor. It's reasonable that Certificate Authorities charge for the service of issuing certifcates.

Unfortunately, the nature of trust is such that browsers can only trust a small number of companies so the authority to issue certificates is strictly managed. We don't want any random company to become a rogue certificate authority and start issuing certificates for google.com domains to spammers.

Economically a market controlled by small number of companies tends to have a cartel that keeps prices high due to lack of competition.

That's exactly what happened in SSL certificates market. You can have a low-end server for $60/year and a certificate alone could cost 5x as much.

That was a problem because the cost of SSL certificates was a significant barrier to adopting encryption.

A few companies decided to poll their resources and solve that.

They funded Let's Encrypt organization, which became a Certificate Authority, wrote necessary software and is running the servers that do the work of issuing certificates via API.
*/
