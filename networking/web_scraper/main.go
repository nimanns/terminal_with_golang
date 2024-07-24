package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://example.com"
	links, err := scrape_links(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d links on %s:\n", len(links), url)
	for i, link := range links {
		fmt.Printf("%d. %s\n", i+1, link)
	}
}

func scrape_links(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch the page: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			href = strings.TrimSpace(href)
			if href != "" && !strings.HasPrefix(href, "#") {
				links = append(links, href)
			}
		}
	})

	return links, nil
}
