package main

import (
	"net/http"
	"io"
	"fmt"
	"flag"
	"github.com/Hellmick/sitemap-builder/linkparser"
	"strings"
)

type Sitemap struct {
	Domain string
	Urls []string
}

func createSitemap(domain string, urls []string) *Sitemap{
	sitemap := new(Sitemap)
	if domain != ""{
		sitemap.Domain = domain
	}
	if len(urls) > 0 {
		sitemap.Urls = urls
	}

	return sitemap
}

func getWebsite(url string, domain string) string {
	if retreiveDomain(url) == "" {
		url = domain + url
	}

	if !strings.Contains(url, domain) {
		return ""
	}

	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "https://" + url
	}

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

func retreiveDomain(url string) string {
	url = strings.Replace(url, "https://www", "", 1)
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "http://www", "", 1)
	url = strings.Split(url, "/")[0]

	return url
}

func contains(stringSlice []string, theString string) bool {
	for _, stringFromSlice := range stringSlice {
		if theString == stringFromSlice {
			return true
		}
	}
	return false
}

func filterLinks(sitemap *Sitemap, links []linkparser.Link) []string {
	urls := []string{}
	for _, link := range links {
		if !strings.Contains(link.Url, sitemap.Domain) {
			continue
		}
		if !contains(sitemap.Urls, link.Url) && link.Url != "" {
			urls = append(urls, link.Url)
		}
	}
	return urls
}

func popFirst(slice []string) []string{
	return slice[1:]
}

func breadthFirstSearch(url string, sitemap *Sitemap, depth int) *Sitemap {
	queue := []string{}
	queue = append(queue, url)
	for i := 0; i < depth; i++ {
		fmt.Sprintf("Iteration #%v", i+1)
		for _, link := range queue {
			html := getWebsite(link, sitemap.Domain)
			if html == "" {
				continue
			}
			links := linkparser.FindLinks(html)
			filteredUrls := filterLinks(sitemap, links)
			queue = popFirst(queue)
			for _, url := range filteredUrls{
				sitemap.Urls = append(sitemap.Urls, url)
				queue = append(queue, url)
			}
		}
	}

	return sitemap
}

func main() {
	fmt.Printf("Start\n")
	url := flag.String("u", "", "URL to build the sitemap for")
	depth := flag.Int("d", 3, "Depth of the BFS algorythm")
	flag.Parse()

	if !strings.Contains(*url, "http://") && !strings.Contains(*url, "https://") {
		fmt.Printf("URL must contain the protocol!")
		return
	}

	if *url == "" {
		fmt.Printf("Usage: sb -u http://example.com -d depth (default 3)")
		return
	}

	sitemap := createSitemap(retreiveDomain(*url), []string{})
	sitemap = breadthFirstSearch(*url, sitemap, *depth)
	

	for _, url := range sitemap.Urls {
		fmt.Printf("Url:%s\n", url)
	}
	

	fmt.Print("End\n")
}
