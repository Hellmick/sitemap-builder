package main

import (
	"log"
	"errors"
	"net/http"
	"io"
	"fmt"
	"flag"
	"github.com/Hellmick/sitemap-builder/linkparser"
	"strings"
	"encoding/xml"
	"os"
)

type Sitemap struct {
	Domain	string
	Urls	[]string
}

type Url struct {
	Loc		string	`xml:"url>loc"`
	Lastmod		string	`xml:"url>lastmod"`
	Changefreq	string	`xml:"url>changefreq"`
	Priority	string	`xml:"url>priority"`
}

type Urlset struct {
	Urls	[]Url	`xml:"url"`
	Xmlns	string	`xml:"xmlns,attr"`
}

func generateXml(sitemap *Sitemap) ([]byte, error) {
	urlset := Urlset{
		Urls: []Url{},
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}
	for _, sitemap_url := range sitemap.Urls {
		url := Url{
			Loc: sitemap_url,
			Lastmod: "",
			Changefreq: "",
			Priority: "",
		}
		urlset.Urls = append(urlset.Urls, url)
	}
	out, err := xml.MarshalIndent(urlset, " ", " ")
	xml := []byte(xml.Header + string(out))
	
	return xml, err
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

func getWebsite(url string, domain string) (string, error) {
	if retreiveDomain(url) == "" {
		url = domain + url
	}

	if !strings.Contains(url, domain) {
		return "", nil 
	}

	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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

func breadthFirstSearch(url string, sitemap *Sitemap, depth int) (*Sitemap, error) {
	if depth < 1 {
		err := errors.New("Algorythm depth must be greater thatn 0")
		return sitemap, err
	}

	visited := make(map[string]bool)
	queue := []string{}
	queue = append(queue, url)
	for i := 0; i < depth; i++ {
		for _, link := range queue {
			html, err := getWebsite(link, sitemap.Domain)
			if err != nil || html == "" || visited["link"] == true{
				continue
			}
			if visited["link"] == false {
				visited["link"] = true 
			}
			links := linkparser.FindLinks(html)
			filteredUrls := filterLinks(sitemap, links)
			queue = popFirst(queue)
			
			for _, url := range filteredUrls {
				sitemap.Urls = append(sitemap.Urls, url)
				queue = append(queue, url)
			}
		}
	}

	return sitemap, nil
}

func writeToFile(filename string, xml []byte) error {
	err := os.WriteFile(filename, xml, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	url := flag.String("u", "", "URL to build the sitemap for")
	depth := flag.Int("d", 3, "Depth of the BFS algorythm")
	filename := flag.String("f", "sitemap.xml", "Output file name")

	flag.Parse()

	if !strings.Contains(*url, "http://") && !strings.Contains(*url, "https://") {
		fmt.Printf("URL must contain the protocol!\n")
		return
	}

	if *url == "" {
		fmt.Printf("Usage: sb -u http://example.com -d depth (default 3)\n")
		return
	}

	fmt.Printf("SitemapBuilder starts its work.\n")

	sitemap, err := breadthFirstSearch(*url, createSitemap(retreiveDomain(*url), []string{}), *depth)
	if err != nil {
		log.Fatal(err)
	}

	sitemapXml, err := generateXml(sitemap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(xml.Header + string(sitemapXml) + "\n")
	if filename != nil {
		err := writeToFile(*filename, sitemapXml)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Print("Done.\n")
}
