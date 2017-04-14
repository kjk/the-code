package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// add file to .zip archive
func addZipFileMust(zw *zip.Writer, path, zipName string) {
	fmt.Printf("Adding file '%s' as '%s' to .zip file\n", path, zipName)
	fi, err := os.Stat(path)
	panicIfErr(err)
	fih, err := zip.FileInfoHeader(fi)
	panicIfErr(err)
	fih.Name = zipName
	fih.Method = zip.Deflate
	d, err := ioutil.ReadFile(path)
	panicIfErr(err)
	fw, err := zw.CreateHeader(fih)
	panicIfErr(err)
	_, err = fw.Write(d)
	panicIfErr(err)
	// fw is just a io.Writer so we can't Close() it. It's not necessary as
	// it's implicitly closed by the next Create(), CreateHeader()
	// or Close() call on zip.Writer
}

func addZipDirMust(zw *zip.Writer, dir string) {
	filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		panicIfErr(err)
		if fi.IsDir() || !fi.Mode().IsRegular() {
			return nil
		}
		addZipFileMust(zw, path, path)
		return nil
	})
}

func main() {
	f, err := os.Create("assets.zip")
	panicIfErr(err)
	defer f.Close()

	zw := zip.NewWriter(f)
	addZipDirMust(zw, "assets")
	err = zw.Close()
	panicIfErr(err)
	fmt.Printf("Created assets.zip with bundled resources.\n")
}
