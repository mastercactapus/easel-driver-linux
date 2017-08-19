package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const baseURL = "https://easel.inventables.com/downloads"

type versionDesc struct {
	Version string
	URL     string
	Check   string
}

var rx = regexp.MustCompile(`EaselDriver-\d+\.\d+\.\d+.pkg$`)

func getVersionDesc(vers, href string) *versionDesc {
	u, err := url.Parse(href)
	if err != nil {
		// not what we're looking for
		return nil
	}
	if vers == "" && !rx.MatchString(u.Path) {
		return nil
	} else if vers != "" && !strings.HasSuffix(u.Path, "EaselDriver-"+vers+".pkg") {
		return nil
	}

	file := path.Base(u.Path)

	return &versionDesc{
		Version: strings.TrimPrefix(strings.TrimSuffix(file, ".pkg"), "EaselDriver-"),
		URL:     href,
		Check:   u.RawQuery,
	}
}

func main() {
	var vers string
	if len(os.Args) > 1 {
		vers = os.Args[1]
		log.Println("looking for version:", vers)
	} else {
		log.Println("looking for latest version")
	}
	resp, err := http.Get(baseURL)
	if err != nil {
		log.Fatalln("GET", baseURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalln("GET", baseURL, resp.Status)
	}

	z := html.NewTokenizer(resp.Body)

	var v *versionDesc
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
					v = getVersionDesc(vers, a.Val)
					if v != nil {
						break scan
					}
					break
				}
			}
		}
	}

	if v == nil {
		log.Fatalln("failed to find current version")
	}

	log.Println("fetching version:", v.Version)

	resp, err = http.Get(v.URL)
	if err != nil {
		log.Fatalln("failed to fetch:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalln("failed to fetch:", resp.Status)
	}
	file := "EaselDriver-" + v.Version + ".pkg"
	fd, err := os.Create(file)
	if err != nil {
		log.Fatalln("failed to create file:", err)
	}
	defer fd.Close()
	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		log.Fatalln("failed to download package:", err)
	}
	fmt.Println(file)
}
