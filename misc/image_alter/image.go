package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

func main() {
	input_file, _ := os.Open("input.png")
	defer input_file.Close()

	input_image, _ := png.Decode(input_file)
	bounds := input_image.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	output_image := image.NewRGBA(bounds)

	pixels := make([]color.Color, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels[y*width+x] = input_image.At(x, y)
		}
	}

	rand.Shuffle(len(pixels), func(i, j int) {
		pixels[i], pixels[j] = pixels[j], pixels[i]
	})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			output_image.Set(x, y, pixels[y*width+x])
		}
	}

	output_file, _ := os.Create("output.png")
	defer output_file.Close()

	png.Encode(output_file, output_image)
}
