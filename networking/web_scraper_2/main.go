package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://picsum.photos/"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error_fetching_the_url:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("error_parsing_the_html:", err)
		return
	}

	var image_urls []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			image_urls = append(image_urls, src)
		}
	})

	output := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>scraped_images</title>
    <style>
        .image_grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 16px;
            padding: 16px;
        }
        .image_grid img {
            width: 100%;
            height: auto;
            object-fit: cover;
        }
    </style>
</head>
<body>
    <div class="image_grid">
`

	for _, img_url := range image_urls {
		output += fmt.Sprintf("        <img src=\"%s\" alt=\"scraped_image\">\n", img_url)
	}

	output += `
    </div>
</body>
</html>
`

	err = os.WriteFile("scraped_images.html", []byte(output), 0644)
	if err != nil {
		fmt.Println("error_saving_the_output_html:", err)
		return
	}

	fmt.Println("scraping_completed_output_saved_to_scraped_images_html")
}
