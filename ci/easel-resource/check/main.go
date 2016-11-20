package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

type versionDesc struct {
	Version string
	URL     string
	Check   string
}

const baseURL = "https://easel.inventables.com/downloads"

func getVersionDesc(href string) *versionDesc {
	u, err := url.Parse(href)
	if err != nil {
		// not what we're looking for
		return nil
	}
	if !strings.HasSuffix(u.Path, ".pkg") {
		return nil
	}

	file := path.Base(u.Path)
	if !strings.HasPrefix(file, "EaselDriver-") {
		return nil
	}

	return &versionDesc{
		Version: strings.TrimPrefix(strings.TrimSuffix(file, ".pkg"), "EaselDriver-"),
		URL:     href,
		Check:   u.RawQuery,
	}
}

func main() {
	var input struct {
		Source struct {
			URI string
		}
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

	if input.Source.URI == "" {
		input.Source.URI = baseURL
	}

	resp, err := http.Get(input.Source.URI)
	if err != nil {
		log.Fatalln("GET", input.Source.URI, err)
	}
	if resp.StatusCode != 200 {
		log.Fatalln("GET", input.Source.URI, resp.Status)
	}

	versions := make([]versionDesc, 0, 50)

	z := html.NewTokenizer(resp.Body)
scan:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break scan
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key != "href" {
						continue
					}
					vs := getVersionDesc(a.Val)
					if vs != nil {
						versions = append(versions, *vs)
					}
					break
				}
			}
		}
	}

	for i, v := range versions {
		if v != input.Version {
			continue
		}

		versions = versions[:i+1]

		break
	}

	err = json.NewEncoder(os.Stdout).Encode(versions)
	if err != nil {
		log.Fatalln("encode:", err)
	}
}
