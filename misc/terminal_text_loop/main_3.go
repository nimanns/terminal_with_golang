package main

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/term"
)

type circle struct {
	radius float64
	color  int
}

func main() {
	width, height, err := term.GetSize(0)
	if err != nil {
		width, height = 80, 24 
	}

	centerX, centerY := width/2, height/2
	circles := make([]circle, 0)
	maxRadius := math.Sqrt(float64(width*width + height*height)) / 2

	clear_screen()
	hide_cursor()

	for {
		buffer := make([][]rune, height)
		for i := range buffer {
			buffer[i] = make([]rune, width)
			for j := range buffer[i] {
				buffer[i][j] = ' '
			}
		}

		if len(circles) == 0 || int(circles[len(circles)-1].radius)%10 == 0 {
			circles = append(circles, circle{radius: 0, color: 31 + len(circles)%7})
		}

		for i := range circles {
			circles[i].radius += 0.5
			if circles[i].radius > maxRadius {
				circles = circles[1:] 
				continue
			}

			drawCircle(buffer, centerX, centerY, circles[i].radius)
		}

		move_cursor_to_top()

		for y, row := range buffer {
			for x, ch := range row {
				color := getColorAtPosition(circles, centerX, centerY, x, y)
				if color != 0 {
					fmt.Printf("\x1b[%dm%c\x1b[0m", color, ch)
				} else {
					fmt.Print(string(ch))
				}
			}
			fmt.Println()
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func drawCircle(buffer [][]rune, centerX, centerY int, radius float64) {
	for theta := 0.0; theta < 2*math.Pi; theta += 0.01 {
		x := int(float64(centerX) + radius*math.Cos(theta))
		y := int(float64(centerY) + radius*math.Sin(theta)/2) 

		if x >= 0 && x < len(buffer[0]) && y >= 0 && y < len(buffer) {
			buffer[y][x] = '●'
		}
	}
}

func getColorAtPosition(circles []circle, centerX, centerY, x, y int) int {
	for i := len(circles) - 1; i >= 0; i-- {
		dx := float64(x - centerX)
		dy := float64(y - centerY) * 2 
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance <= circles[i].radius {
			return circles[i].color
		}
	}
	return 0
}

func clear_screen() {
	fmt.Print("\x1b[2J")
}

func move_cursor_to_top() {
	fmt.Print("\x1b[H")
}

func hide_cursor() {
	fmt.Print("\x1b[?25l")
}
