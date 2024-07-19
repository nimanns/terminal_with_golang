package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/term"
)

type square struct {
	x, y, dx, dy int
	color        int
}

func main() {
	width, height, err := term.GetSize(0)
	if err != nil {
		width, height = 80, 24 // fallback size
	}

	num_squares := 5
	squares := make([]square, num_squares)
	directions := []int{-1, 1}

	for i := range squares {
		squares[i] = square{
			x:     rand.Intn(width),
			y:     rand.Intn(height),
			dx:    directions[rand.Intn(2)],
			dy:    directions[rand.Intn(2)],
			color: 31 + rand.Intn(7), // ANSI colors from 31 to 37
		}
	}

	clear_screen()
	hide_cursor()

	buffer := make([][]rune, height)
	for i := range buffer {
		buffer[i] = make([]rune, width)
	}

	for {
		for y := range buffer {
			for x := range buffer[y] {
				buffer[y][x] = ' '
			}
		}

		for i := range squares {
			squares[i].x += squares[i].dx
			squares[i].y += squares[i].dy

			if squares[i].x <= 0 || squares[i].x >= width-1 {
				squares[i].dx *= -1
			}
			if squares[i].y <= 0 || squares[i].y >= height-1 {
				squares[i].dy *= -1
			}

			squares[i].x = (squares[i].x + width) % width
			squares[i].y = (squares[i].y + height) % height

			buffer[squares[i].y][squares[i].x] = '■'
		}

		move_cursor_to_top()

		for y, row := range buffer {
			for x, ch := range row {
				for _, s := range squares {
					if s.x == x && s.y == y {
						fmt.Printf("\x1b[%dm%c\x1b[0m", s.color, ch)
						goto next
					}
				}
				fmt.Print(string(ch))
			next:
			}
			fmt.Println()
		}

		time.Sleep(50 * time.Millisecond)
	}
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
