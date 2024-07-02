package main

import (
	"net/http"
	"io"
	"fmt"
	"flag"
	"github.com/Hellmick/sitemap-builder/linkparser"
)

func getWebsite(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func main() {
	fmt.Printf("Start\n")
	url := flag.String("d", "", "Domain to build the sitemap for")
	flag.Parse()
	html := getWebsite(*url)
	
	links := linkparser.FindLinks(html)
	for _, link := range links {
		fmt.Println("Url:", link.Url, "Text:", link.Text, "\n")
	}
}
