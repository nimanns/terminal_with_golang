package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go <text_to_loop>")
		return
	}

	text := os.Args[1]
	width, _, _ := term.GetSize(0)

	if width <= 0 {
		width = 80 // fallback width if unable to determine terminal width
	}

	// Truncate text if it's longer than the terminal width
	if len(text) > width {
		text = text[:width]
	}

	padding := strings.Repeat(" ", width-len(text))

	for {
		for i := 0; i <= len(text); i++ {
			output := text[i:] + text[:i] + padding
			fmt.Print("\r" + output[:width])
			time.Sleep(100 * time.Millisecond)
		}
	}
}
