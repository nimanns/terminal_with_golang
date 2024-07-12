package main

import (
	"image/color"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
	circleRadius = 20
)

type Game struct {
	x, y    float64
	dx, dy  float64
}

func (g *Game) Update() error {
	g.x += g.dx
	g.y += g.dy

	if g.x <= circleRadius || g.x >= screenWidth-circleRadius {
		g.dx = -g.dx
	}
	if g.y <= circleRadius || g.y >= screenHeight-circleRadius {
		g.dy = -g.dy
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(g.x), float32(g.y), circleRadius, color.RGBA{255, 0, 0, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		x:  screenWidth / 2,
		y:  screenHeight / 2,
		dx: 4,
		dy: 3,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bouncing Circle")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
