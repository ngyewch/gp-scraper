package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

type DataHolder struct {
	Apps []AppData `json:"apps"`
}

type AppData struct {
	PackageName string `json:"packageName"`
	Name string `json:"name"`
	Developer string `json:"developer"`
	Categories []string `json:"categories"`
	Tags []string `json:"tags"`
	Version string `json:"version"`
}

func doScrape(urlString string) AppData {
	fmt.Println(urlString)

	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	appData := AppData{}
	appData.PackageName = u.Query().Get("id")

	resp, err := http.Get(urlString)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	genreMatcher := func(n *html.Node) bool {
		return scrape.Attr(n, "itemprop") == "genre"
	}
	iconMatcher := func(n *html.Node) bool {
		return scrape.Attr(n, "itemprop") == "image"
	}
	softwareVersionMatcher := func(n *html.Node) bool {
		return scrape.Attr(n, "itemprop") == "softwareVersion"
	}

	name, ok := scrape.Find(root, scrape.ByClass("id-app-title"))
	if ok {
		appData.Name = scrape.Text(name)
	}
	genre, ok := scrape.Find(root, genreMatcher)
	if ok {
		appData.Categories = append(appData.Categories, scrape.Text(genre))
	}
	icon, ok := scrape.Find(root, iconMatcher)
	if ok {
		iconSrc := scrape.Attr(icon, "src")
		iconUrl, err := url.Parse(iconSrc)
		if err != nil {
			panic(err)
		}
		if iconUrl.Scheme == "" {
			iconSrc = "https:" + iconSrc
		}

		resp, err = http.Get(iconSrc)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		outputFile, err := os.Create("output/" + appData.PackageName + ".png")
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		_, err = io.Copy(outputFile, resp.Body)
		if err != nil {
			panic(err)
		}
	}
	version, ok := scrape.Find(root, softwareVersionMatcher)
	if ok {
		appData.Version = strings.TrimSpace(scrape.Text(version))
	}

	return appData
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dataHolder := DataHolder{}

	err = os.MkdirAll("output", 0777)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dataHolder.Apps = append(dataHolder.Apps, doScrape(scanner.Text()))
	}

	outputFile, err := os.Create("output/data.json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	jsonEncoder := json.NewEncoder(outputFile)
	jsonEncoder.Encode(dataHolder)
}
