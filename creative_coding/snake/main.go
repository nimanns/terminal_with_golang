
package main

import (
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
    snake_game snake
    food       point
    width      = 20
    height     = 20
    game_over  bool
)

func init() {
    rand.Seed(time.Now().UnixNano())
    snake_game = snake{
        body: []point{
            {x: width / 2, y: height / 2},
        },
        direction: point{x: 1, y: 0},
    }
    place_food()
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
    } else {
        snake_game.body = snake_game.body[:len(snake_game.body)-1]
    }
}

func draw() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    for _, segment := range snake_game.body {
        termbox.SetCell(segment.x, segment.y, 'O', termbox.ColorGreen, termbox.ColorDefault)
    }

    termbox.SetCell(food.x, food.y, '@', termbox.ColorRed, termbox.ColorDefault)

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

    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    event_queue := make(chan termbox.Event)
    go func() {
        for {
            event_queue <- poll_event()
        }
    }()

    for !game_over {
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
        }
    }

    termbox.Close()
    println("Game Over!")
}

