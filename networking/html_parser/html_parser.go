// basic html parser, reads a made-up format of text from a txt file and generates html elements based on the text
// the idea is to learn go by making my own mini framework
// place .txt file in the root of this directory

package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
)

func parseFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var elements []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.SplitN(line, "->", 2)
        if len(parts) == 2 {
            tag := strings.TrimSpace(parts[0])
            content := strings.TrimSpace(parts[1])
            element := fmt.Sprintf("<%s>%s</%s>", tag, content, tag)
            elements = append(elements, element)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return elements, nil
}

func generateHTML(elements []string) string {
    body := strings.Join(elements, "\n    ")
    return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Generated Page</title>
</head>
<body>
    %s
</body>
</html>
`, body)
}

func main() {
    inputFile := "main.txt"
    outputFile := "output.html"

    elements, err := parseFile(inputFile)
    if err != nil {
        log.Fatalf("Error parsing file: %v", err)
    }

    html := generateHTML(elements)

    err = ioutil.WriteFile(outputFile, []byte(html), 0644)
    if err != nil {
        log.Fatalf("Error writing HTML file: %v", err)
    }

    fmt.Printf("HTML file generated: %s\n", outputFile)
}
