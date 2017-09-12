package main

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
)

func read_gzip_file(path string) ([]byte, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := gzip.NewReader(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	log.Printf("%s: %d bytes uncompressed", path, len(buf))
	return buf, nil
}

func main() {
	_, err := read_gzip_file("macbeth.html.gz")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
