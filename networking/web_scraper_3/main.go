package main

import (
	"fmt"
	"io"
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
