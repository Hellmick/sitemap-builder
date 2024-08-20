package linkparser

import (
	"testing"
)

func TestFindLinks(t *testing.T) {	
	
	want := make(map[string]string)
	want["https://loremipsum.org"] = "Lorem Ipsum"
	htmlText := `<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8"/>
	</head>
	<body>
		<h1>Lorem Ipsum</h1>
		<a href="https://loremipsum.org">Lorem Ipsum</a>
	</body>
	</html>`
	
	links, err := FindLinks(htmlText)
	if err != nil {
		t.Fatal(err)
	}

	for _, link := range links {
		if want[link.Url] != link.Text {
			t.Fatal("The link is not correct")
		}
	}
}
