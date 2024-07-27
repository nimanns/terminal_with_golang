package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	url := "https://example.com"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	element_map := make(map[string][]string)
	traverse_node(doc, element_map)

	output_file, err := os.Create("output.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer output_file.Close()

	write_html(output_file, element_map)
	
	word_cloud_file, err := os.Create("word_cloud.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer word_cloud_file.Close()

	generate_word_cloud(word_cloud_file, element_map)
}

func traverse_node(n *html.Node, element_map map[string][]string) {
	if n.Type == html.ElementNode {
		element_type := strings.ToLower(n.Data)
		element_value := get_node_content(n)
		element_map[element_type] = append(element_map[element_type], element_value)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse_node(c, element_map)
	}
}

func get_node_content(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var content string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		content += get_node_content(c)
	}
	return strings.TrimSpace(content)
}

func write_html(w io.Writer, element_map map[string][]string) {
	fmt.Fprintf(w, "<!DOCTYPE html>\n<html>\n<body>\n")
	for element_type, values := range element_map {
		fmt.Fprintf(w, "<h1>%s</h1>\n<ul>\n", element_type)
		for _, value := range values {
			fmt.Fprintf(w, "<li>%s</li>\n", value)
		}
		fmt.Fprintf(w, "</ul>\n")
	}
	fmt.Fprintf(w, "</body>\n</html>")
}

func generate_word_cloud(w io.Writer, element_map map[string][]string) {
	word_count := make(map[string]int)
	for _, values := range element_map {
		for _, value := range values {
			words := strings.Fields(strings.ToLower(value))
			for _, word := range words {
				if len(word) > 3 {
					word_count[word]++
				}
			}
		}
	}

	fmt.Fprintf(w, "<!DOCTYPE html>\n<html>\n<head>\n")
	fmt.Fprintf(w, "<style>\n.word_cloud {position: relative; width: 800px; height: 400px;}\n")
	fmt.Fprintf(w, ".word {position: absolute; font-family: Arial;}\n</style>\n")
	fmt.Fprintf(w, "</head>\n<body>\n<div class=\"word_cloud\">\n")

	for word, count := range word_count {
		font_size := 10 + count*2
		top := rand.Intn(350)
		left := rand.Intn(750)
		color := fmt.Sprintf("#%06x", rand.Intn(0xFFFFFF))
		fmt.Fprintf(w, "<span class=\"word\" style=\"font-size:%dpx;top:%dpx;left:%dpx;color:%s;\">%s</span>\n",
			font_size, top, left, color, word)
	}

	fmt.Fprintf(w, "</div>\n</body>\n</html>")
}

