package main

import (
	"net/http"
	"io"
	"fmt"
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
	fmt.Printf(getWebsite("http://example.org"))
}
