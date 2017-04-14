package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	flgUseBundledAssets bool
	httpAddr            = ":6789"
	mimeTextPlain       = "text/plain; charset=utf-8"
	// loaded only once at startup. maps a file path of the resource to its data
	assetsFromZip map[string][]byte
)

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func mimeTypeByExtensionExt(name string, defaultMimeType string) string {
	ext := strings.ToLower(filepath.Ext(name))
	result := mime.TypeByExtension(ext)
	if result == "" {
		return defaultMimeType
	}
	return result
}

func loadAssetsFromZipMust() {
	zr, err := zip.OpenReader("assets.zip")
	panicIfErr(err)
	defer zr.Close()
	assetsFromZip = make(map[string][]byte)
	for _, f := range zr.File {
		// convert windows-style path to unix style
		name := strings.Replace(f.Name, "\\", "/", -1)
		rc, err := f.Open()
		panicIfErr(err)
		d, err := ioutil.ReadAll(rc)
		rc.Close()
		panicIfErr(err)
		assetsFromZip[name] = d
	}
}

func serveData(w http.ResponseWriter, r *http.Request, path string, data []byte) {
	contentType := mimeTypeByExtensionExt(path, mimeTextPlain)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func serveAssetFromFilesystem(w http.ResponseWriter, r *http.Request, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	serveData(w, r, path, data)
}

func serveAssetFromZip(w http.ResponseWriter, r *http.Request, path string) {
	data := assetsFromZip[path]
	if data == nil {
		http.NotFound(w, r)
		return
	}
	serveData(w, r, path, data)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	if uri == "/" {
		uri = "html/index.html"
	}

	assetPath := filepath.Join("assets", strings.ToLower(uri))
	fmt.Printf("handleIndex: from zip: %v, uri: %s, asset path: %s\n", flgUseBundledAssets, r.URL.Path, assetPath)
	if flgUseBundledAssets {
		serveAssetFromZip(w, r, assetPath)
	} else {
		serveAssetFromFilesystem(w, r, assetPath)
	}
}

func main() {
	flag.BoolVar(&flgUseBundledAssets, "use-bundled-assets", false, "if true, serve assets from assets.zip")
	flag.Parse()

	if flgUseBundledAssets {
		loadAssetsFromZipMust()
	}

	http.HandleFunc("/", handleIndex)

	fmt.Printf("Starting web server on port %s. Serving bundled assets: %v!.\nOpen http://localhost%s in your browser\n", httpAddr, flgUseBundledAssets, httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("http.ListenAndServe() failed with %s\n", err)
	}
}
