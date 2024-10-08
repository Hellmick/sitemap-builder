package linkparser

import (
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

func FindLinks(html_text string) ([]Link, error) {
	reader := strings.NewReader(html_text)
	var links []Link
	var findLinks func(node *html.Node) error
	
	doc, err := html.Parse(reader)
	if err != nil {
		return links, err
	}

	findLinks = func(node *html.Node) error {
		
		//TODO: make it concurrent
		if node.Type == html.ElementNode && node.Data == "a" {
			current_link, err := processNode(node)
			if err != nil {
				return err
			}
			links = append(links, current_link)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
		return nil
	}
	
	err = findLinks(doc)
	return links, err

}



