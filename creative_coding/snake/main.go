package main

import (
    "fmt"
    "github.com/nsf/termbox-go"
    "math/rand"
    "time"
)

type point struct {
    x, y int
}

type snake struct {
    body      []point
    direction point
}

var (
    snake_game  snake
    food        point
    width       = 20
    height      = 20
    game_over   bool
    score       int
    initial_tick = 100 * time.Millisecond
)

func init_game() {
    rand.Seed(time.Now().UnixNano())
    snake_game = snake{
        body: []point{
            {x: width / 2, y: height / 2},
        },
        direction: point{x: 1, y: 0},
    }
    place_food()
    score = 0
    game_over = false
}

func place_food() {
    food = point{
        x: rand.Intn(width),
        y: rand.Intn(height),
    }
}

func move_snake() {
    head := snake_game.body[0]
    new_head := point{
        x: head.x + snake_game.direction.x,
        y: head.y + snake_game.direction.y,
    }

    if new_head.x < 0 || new_head.x >= width || new_head.y < 0 || new_head.y >= height {
        game_over = true
        return
    }

    for _, segment := range snake_game.body {
        if segment == new_head {
            game_over = true
            return
        }
    }

    snake_game.body = append([]point{new_head}, snake_game.body...)
    if new_head == food {
        place_food()
        score++
    } else {
        snake_game.body = snake_game.body[:len(snake_game.body)-1]
    }
}

func draw() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    for _, segment := range snake_game.body {
        char := 'O'
        if segment == snake_game.body[0] {
            char = 'H'
        }
        termbox.SetCell(segment.x, segment.y, char, termbox.ColorGreen, termbox.ColorDefault)
    }

    termbox.SetCell(food.x, food.y, '@', termbox.ColorRed, termbox.ColorDefault)

    for i, c := range fmt.Sprintf("Score: %d", score) {
        termbox.SetCell(i, height, c, termbox.ColorWhite, termbox.ColorDefault)
    }

    termbox.Flush()
}

func poll_event() termbox.Event {
    return termbox.PollEvent()
}

func main() {
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    init_game()
    ticker := time.NewTicker(initial_tick)
    defer ticker.Stop()

    event_queue := make(chan termbox.Event)
    go func() {
        for {
            event_queue <- poll_event()
        }
    }()

    for {
        if game_over {
            termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
            msg := "Game Over! Press 'R' to restart or 'Esc' to quit."
            for i, c := range msg {
                termbox.SetCell((width-len(msg))/2+i, height/2, c, termbox.ColorRed, termbox.ColorDefault)
            }
            termbox.Flush()

            ev := poll_event()
            if ev.Type == termbox.EventKey {
                if ev.Ch == 'r' || ev.Ch == 'R' {
                    init_game()
                    ticker = time.NewTicker(initial_tick)
                } else if ev.Key == termbox.KeyEsc {
                    break
                }
            }
        } else {
            select {
            case ev := <-event_queue:
                if ev.Type == termbox.EventKey {
                    switch ev.Key {
                    case termbox.KeyArrowUp:
                        if snake_game.direction.y == 0 {
                            snake_game.direction = point{x: 0, y: -1}
                        }
                    case termbox.KeyArrowDown:
                        if snake_game.direction.y == 0 {
                            snake_game.direction = point{x: 0, y: 1}
                        }
                    case termbox.KeyArrowLeft:
                        if snake_game.direction.x == 0 {
                            snake_game.direction = point{x: -1, y: 0}
                        }
                    case termbox.KeyArrowRight:
                        if snake_game.direction.x == 0 {
                            snake_game.direction = point{x: 1, y: 0}
                        }
                    case termbox.KeyEsc:
                        game_over = true
                    }
                }
            case <-ticker.C:
                move_snake()
                draw()
                if score != 0 && score%5 == 0 {
                    ticker.Stop()
                    ticker = time.NewTicker(initial_tick - time.Duration(score*5)*time.Millisecond)
                }
            }
        }
    }

    termbox.Close()
    println("Thanks for playing!")
}

