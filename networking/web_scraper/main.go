package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type page_data struct {
	url         string
	title       string
	links       []string
	images      []string
	word_count  int
}

func main() {
	start_urls := []string{
		"https://example.com",
		"https://golang.org",
		"https://github.com",
	}

	results := scrape_multiple(start_urls, 5)

	save_results(results, "scrape_results.txt")

	fmt.Printf("Scraped %d pages. Results saved to scrape_results.txt\n", len(results))
}

func scrape_multiple(urls []string, max_concurrent int) []page_data {
	var wg sync.WaitGroup
	sem := make(chan struct{}, max_concurrent)
	results_chan := make(chan page_data, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			data, err := scrape_page(url)
			if err != nil {
				log.Printf("Error scraping %s: %v", url, err)
				return
			}
			results_chan <- data
		}(url)
	}

	go func() {
		wg.Wait()
		close(results_chan)
	}()

	var results []page_data
	for data := range results_chan {
		results = append(results, data)
	}

	return results
}

func scrape_page(url string) (page_data, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return page_data{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return page_data{}, fmt.Errorf("failed to fetch the page: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return page_data{}, err
	}

	data := page_data{
		url:   url,
		title: doc.Find("title").Text(),
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			href = strings.TrimSpace(href)
			if href != "" && !strings.HasPrefix(href, "#") {
				data.links = append(data.links, href)
			}
		}
	})

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			src = strings.TrimSpace(src)
			if src != "" {
				data.images = append(data.images, src)
			}
		}
	})

	data.word_count = len(strings.Fields(doc.Text()))

	return data, nil
}

func save_results(results []page_data, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, data := range results {
		fmt.Fprintf(file, "URL: %s\n", data.url)
		fmt.Fprintf(file, "Title: %s\n", data.title)
		fmt.Fprintf(file, "Word Count: %d\n", data.word_count)
		fmt.Fprintf(file, "Links (%d):\n", len(data.links))
		for _, link := range data.links {
			fmt.Fprintf(file, "  - %s\n", link)
		}
		fmt.Fprintf(file, "Images (%d):\n", len(data.images))
		for _, img := range data.images {
			fmt.Fprintf(file, "  - %s\n", img)
		}
		fmt.Fprintln(file, strings.Repeat("-", 50))
	}

	return nil
}
