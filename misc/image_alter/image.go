package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	input_file := flag.String("input", "input.png", "Input PNG file")
	output_file := flag.String("output", "output.png", "Output PNG file")
	mode := flag.String("mode", "invert", "Image processing mode")
	seed := flag.Int64("seed", time.Now().UnixNano(), "Random seed for shuffling")
	flag.Parse()

	input, err := os.Open(*input_file)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer input.Close()

	input_image, err := png.Decode(input)
	if err != nil {
		fmt.Printf("Error decoding input image: %v\n", err)
		return
	}

	bounds := input_image.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	output_image := image.NewRGBA(bounds)

	switch *mode {
		case "shuffle":
			shuffle_pixels(input_image, output_image, width, height, *seed)
		case "invert":
			invert_colors(input_image, output_image, width, height)
		default:
			fmt.Println("Invalid mode selected. Using default shuffle mode.")
			shuffle_pixels(input_image, output_image, width, height, *seed)
	}

	output, err := os.Create(*output_file)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer output.Close()

	err = png.Encode(output, output_image)

	if err != nil {
		fmt.Printf("Error encoding output image: %v\n", err)
		return
	}

	fmt.Println("Image processing completed successfully.")
}

func shuffle_pixels(input image.Image, output *image.RGBA, width, height int, seed int64) {
	pixels := make([]color.Color, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels[y*width+x] = input.At(x, y)
		}
	}

	rand.Seed(seed)
	rand.Shuffle(len(pixels), func(i, j int) {
		pixels[i], pixels[j] = pixels[j], pixels[i]
	})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			output.Set(x, y, pixels[y*width+x])
		}
	}
}


func invert_colors(input image.Image, output *image.RGBA, width, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			output.Set(x, y, color.RGBA{
				R: uint8(255 - r>>8),
				G: uint8(255 - g>>8),
				B: uint8(255 - b>>8),
				A: uint8(a >> 8),
			})
		}
	}
}
