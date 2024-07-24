package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	width  = 80
	height = 24
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?")

func main() {
	rand.Seed(time.Now().UnixNano())
	
	fmt.Print("\033[2J\033[?25l")
	
	drops := make([]int, width)
	
	for {
		fmt.Print("\033[H")
		
		for i := 0; i < width; i++ {
			if drops[i] == 0 && rand.Float32() < 0.1 {
				drops[i] = 1
			}
			
			if drops[i] > 0 {
				color := rand.Intn(6) + 31
				fmt.Printf("\033[%dm%c\033[0m", color, chars[rand.Intn(len(chars))])
				
				drops[i]++
				if drops[i] > height {
					drops[i] = 0
				}
			} else {
				fmt.Print(" ")
			}
		}
		
		time.Sleep(50 * time.Millisecond)
	}
}
