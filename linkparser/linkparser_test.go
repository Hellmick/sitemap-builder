package linkparser

import (
	"testing"
	"regexp"
)

func TestRetreiveText(t *testing T) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil
		panic(err)
	}

	
}

func TestFindLinks(t *testing T) {
	want := []Link{}
	want.append()
}
