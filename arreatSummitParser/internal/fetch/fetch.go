package fetch

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func LoadCategory(category string, variant string) (string, error) {
	var url string
	if category == "circlets" {
		if variant != "normal" {
			return "", nil
		}
		url = fmt.Sprintf("https://classic.battle.net/diablo2exp/items/%s.shtml", category)
	} else {
		url = fmt.Sprintf("https://classic.battle.net/diablo2exp/items/%s/%s.shtml", variant, category)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return "", err
		}
	default:
		reader = resp.Body
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	html, err := doc.Html()
	if err != nil {
		return "", err
	}

	return html, nil
}
