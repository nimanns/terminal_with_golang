package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"math"
	"os"
	"time"
)

func main() {
	input_file := flag.String("input", "input.png", "Input PNG file")
	output_file := flag.String("output", "output.png", "Output PNG file")
	mode := flag.String("mode", "edge_detect", "Image processing mode")
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
		case "gs":
			convert_to_grayscale(input_image, output_image, width, height)
		case "edge_detect":
			edge_detection(input_image, output_image, width, height)
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


func convert_to_grayscale(input image.Image, output *image.RGBA, width, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
			output.Set(x, y, color.RGBA{
				R: gray,
				G: gray,
				B: gray,
				A: uint8(a >> 8),
			})
		}
	}
}

func edge_detection(input image.Image, output *image.RGBA, width, height int) {
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			gx := sobel_x(input, x, y)
			gy := sobel_y(input, x, y)
			magnitude := uint8(math.Sqrt(float64(gx*gx + gy*gy)))
			output.Set(x, y, color.RGBA{magnitude, magnitude, magnitude, 255})
		}
	}
}

func sobel_x(img image.Image, x, y int) int {
	return -int(brightness(img.At(x-1, y-1))) - 2*int(brightness(img.At(x-1, y))) - int(brightness(img.At(x-1, y+1))) +
		int(brightness(img.At(x+1, y-1))) + 2*int(brightness(img.At(x+1, y))) + int(brightness(img.At(x+1, y+1)))
}

func sobel_y(img image.Image, x, y int) int {
	return -int(brightness(img.At(x-1, y-1))) - 2*int(brightness(img.At(x, y-1))) - int(brightness(img.At(x+1, y-1))) +
		int(brightness(img.At(x-1, y+1))) + 2*int(brightness(img.At(x, y+1))) + int(brightness(img.At(x+1, y+1)))
}

func brightness(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	return uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
}
