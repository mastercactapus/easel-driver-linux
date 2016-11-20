package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mastercactapus/gocpio"
	"github.com/mastercactapus/goxar"
)

type versionDesc struct {
	Version string
	URL     string
	Check   string
}

func getPayload(xarData []byte) (io.ReadCloser, error) {
	b := bytes.NewReader(xarData)
	r, err := xar.NewReader(b, int64(b.Len()))
	if err != nil {
		return nil, err
	}

	filename := "IrisLib-0.3.1.pkg/Payload"
	for _, f := range r.File {
		if f.Name != filename {
			continue
		}
		return f.Open()

	}

	return nil, fmt.Errorf("could not find file in xar: %s", filename)
}

func main() {

	var input struct {
		Version versionDesc
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln("read stdin:", err)
	}

	err = json.Unmarshal(data, &input)
	if err != nil {
		log.Fatalln("decode input:", err, "\n\nFrom:\n", string(data))
	}

	log.Println("Downloading:", input.Version.URL)
	resp, err := http.Get(input.Version.URL)
	if err != nil {
		log.Fatalln("GET", input.Version.URL, err)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("download:", err)
	}
	resp.Body.Close()

	root := os.Args[1]
	log.Println("Extracting to:", root)
	r, err := getPayload(data)
	if err != nil {
		log.Fatalln("processing xar:", err)
	}
	defer r.Close()

	gr, err := gzip.NewReader(r)
	if err != nil {
		log.Fatalln("gunzip payload:", err)
	}

	cr := cpio.NewReader(gr)

	for {
		hdr, err := cr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("extract cpio payload:", err)
		}
		filename := filepath.Join(root, hdr.Name)
		if hdr.FileInfo().IsDir() {
			err = os.MkdirAll(filename, hdr.FileInfo().Mode())
			if err != nil {
				log.Fatalln("mkdir,", filename, ":", err)
			}
			continue
		}
		fd, err := os.Create(filename)
		if err != nil {
			log.Fatalln(filename, ":", err)
		}
		_, err = io.Copy(fd, cr)
		if err != nil {
			log.Fatalln("extract", filename, err)
		}
		fd.Close()
	}
}
