package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:',.<>/?`~")
var rows, cols int
var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	rows = 40
	cols = 100
}

func call_clear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("unsupported platform")
	}
}

func random_char() rune {
	return chars[rand.Intn(len(chars))]
}

func random_color() string {
	return fmt.Sprintf("\033[3%dm", rand.Intn(7)+1)
}

func random_bg_color() string {
	return fmt.Sprintf("\033[4%dm", rand.Intn(7)+1)
}

func random_bold() string {
	if rand.Intn(2) == 0 {
		return "\033[1m"
	}
	return ""
}

func random_underline() string {
	if rand.Intn(2) == 0 {
		return "\033[4m"
	}
	return ""
}

func reset_color() string {
	return "\033[0m"
}

func visual_madness() {
	for {
		call_clear()
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				fmt.Printf("%s%s%s%s%c", random_color(), random_bg_color(), random_bold(), random_underline(), random_char())
			}
			fmt.Print(reset_color() + "\n")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func moving_parts() {
	for {
		call_clear()
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				if rand.Intn(100) < 5 {
					fmt.Printf("%s%s%s%s%c", random_color(), random_bg_color(), random_bold(), random_underline(), random_char())
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Print(reset_color() + "\n")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func draw_circle(x, y, r int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			distance := math.Sqrt(float64((i-x)*(i-x) + (j-y)*(j-y)))
			if int(distance) == r {
				fmt.Printf("%s%s%s%s%c", random_color(), random_bg_color(), random_bold(), random_underline(), random_char())
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print(reset_color() + "\n")
	}
}

func animate_circles() {
	x, y, r := rows/2, cols/2, 10
	dx, dy := 1, 1
	for {
		call_clear()
		draw_circle(x, y, r)
		x += dx
		y += dy
		if x-r <= 0 || x+r >= rows {
			dx = -dx
		}
		if y-r <= 0 || y+r >= cols {
			dy = -dy
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	go visual_madness()
	go moving_parts()
	animate_circles()
}

