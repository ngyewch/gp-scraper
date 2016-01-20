package main

import (
	"fmt"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

func main() {
	// request and parse the front page
	resp, err := http.Get("https://play.google.com/store/apps/details?id=com.alegrium.billionaire")
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	title, ok := scrape.Find(root, scrape.ByClass("id-app-title"))
	if ok {
		fmt.Printf("title=%s\n", scrape.Text(title))
	}
	icon, ok := scrape.Find(root, scrape.ByClass("cover-image"))
	if ok {
		fmt.Printf("icon=%s\n", scrape.Attr(icon, "src"))
	}
}
