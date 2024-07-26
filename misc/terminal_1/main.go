package main

import (
	"fmt"
	"time"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
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
	rows = 20
	cols = 50
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

func reset_color() string {
	return "\033[0m"
}

func visual_madness() {
	for {
		call_clear()
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				fmt.Printf("%s%s%c", random_color(), random_bg_color(), random_char())
			}
			fmt.Print(reset_color() + "\n")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	visual_madness()
}