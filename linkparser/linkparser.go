package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"golang.org/x/net/html"
)

type Link struct {
	Url  string
	Text string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var links []Link

func retreiveText(node *html.Node) string {
	var text string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			text += child.Data
		} else {
			text += retreiveText(child)
		}
	}
	return text
}

func processNode(node *html.Node) (Link, error) {
	var link Link
	for _, attribute := range node.Attr {
		if attribute.Key == "href" {
			link.Url = attribute.Val
		}
		link.Text = strings.TrimSpace(retreiveText(node))
	}

	return link, nil
}

func findLinks(node *html.Node) {

	if node.Type == html.ElementNode && node.Data == "a" {
		current_link, err := processNode(node)
		check(err)
		links = append(links, current_link)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findLinks(child)
	}

}

func main() {
	filename := flag.String("f", "file.html", "filename of the parsed file")
	flag.Parse()

	dat, err := os.ReadFile(*filename)
	check(err)

	reader := bytes.NewReader(dat)
	node_tree, err := html.Parse(reader)
	check(err)

	findLinks(node_tree)
	for _, link := range links {
		fmt.Printf("Text: %s, URL: %s\n", link.Text, link.Url)
	}

}
