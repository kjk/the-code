package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"

	"github.com/kjk/u"
)

var (
	visitedURLS map[string]*URLData
)

// URLData represents result of http GET on a given url
type URLData struct {
	url         string // original url
	size        int64
	statusCode  int // 200, 404 etc.
	header      http.Header
	bodySha1Hex string
}

func usage() {
	fmt.Printf("Usage:\nlinkchecker <url>\n")
}

func makeFullURL(base string, href string) *url.URL {
	baseURL, err := url.Parse(base)
	if err != nil {
		fmt.Printf("Failed to parse '%s', error: %s\n", baseURL, err)
		return nil
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		fmt.Printf("Failed to parse '%s', error: %s\n", href, err)
		return nil
	}
	fullURL := baseURL.ResolveReference(hrefURL)
	return fullURL
}

func main() {
	args := os.Args
	if len(args) != 2 {
		usage()
		os.Exit(1)
	}
	startURL := os.Args[1]
	parsedStartURL, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("'%s' is not a valid url. Error: %s\n", startURL, err)
	}
	fmt.Printf("Checking links for '%s'\n", startURL)
	fmt.Printf("parsed: %+v\n", parsedStartURL)
	visitedURLS = make(map[string]*URLData)
	toVisit := []string{startURL}
	startHostname := parsedStartURL.Hostname()
	for len(toVisit) > 0 {
		uri := toVisit[0]
		toVisit = toVisit[1:]
		// don't visit the same url twice
		if _, exists := visitedURLS[uri]; exists {
			continue
		}

		parsed, err := url.Parse(uri)
		if err != nil {
			fmt.Printf("Failed to parse '%s', error: %s\n", uri, err)
			continue
		}

		rsp, err := http.Get(uri)
		if err != nil {
			fmt.Printf("http.Get('%s') failed with %s\n", uri, err)
			continue
		}
		defer rsp.Body.Close()
		d, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Printf("failed reading body of '%s', error: %s\n", uri, err)
			continue
		}
		sha1Hex := u.Sha1HexOfBytes(d)
		data := &URLData{
			url:         uri,
			statusCode:  rsp.StatusCode,
			header:      rsp.Header,
			bodySha1Hex: sha1Hex,
			size:        int64(len(d)),
		}
		visitedURLS[uri] = data
		// TODO: handle 301, 302
		fmt.Printf("Downloaded %s, size: %d, code: %d\n", uri, len(d), rsp.StatusCode)

		if parsed.Hostname() != startHostname {
			fmt.Printf("Not following '%s' because external link\n", uri)
			continue
		}

		hrefs := extractURLs(d)
		fmt.Printf("Extracted %d urls\n", len(hrefs))
		for _, href := range hrefs {
			fullURL := makeFullURL(uri, href)
			if nil == fullURL {
				fmt.Printf("Failed to make full url out of '%s' and '%s'\n", uri, href)
				continue
			}
			//fmt.Printf("'%s' => '%s'\n", href, fullURL)
			toVisit = append(toVisit, fullURL.String())
		}
	}
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) string {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

func extractURLs(d []byte) []string {
	r := bytes.NewBuffer(d)
	z := html.NewTokenizer(r)
	var res []string
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			t := z.Token()
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			uri := getHref(t)
			if uri != "" {
				res = append(res, uri)
			}
		}
	}
	return res
}
