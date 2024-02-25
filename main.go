package main

import (
	"fmt"
	"image"
	go_color "image/color"
	"image/gif"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/scenes"
)

const (
	frameSpacing = 7
	nFrameCount  = 20
	width        = 1000
	height       = 1000

	interpolateN = 9

	GIF_FORMAT = "gif"
	PNG_FORMAT = "png"
)

var (
	// gradient = color.LinearGradient{
	// 	Points: []color.Color{
	// 		color.Hex("#6CB4F5"),
	// 		color.Hex("#EBF56C"),
	// 		color.Hex("#F5736C"),
	// 	},
	// }
	// gradient = color.LinearGradient{
	// 	Points: []color.Color{
	// 		color.Hex("#F590C1"), // pink
	// 		color.Hex("#000"),
	// 		color.Hex("#90E8F5"), // light blue
	// 		color.Hex("#000"),
	// 		color.Hex("#F590C1"), // pink
	// 	},
	// }

	// gradient = color.LinearGradient{
	// 	Points: []color.Color{
	// 		color.Hex("#FFF"), // black
	// 		color.Hex("#DDF522"),
	// 		color.Hex("#A0514C"),
	// 		color.Hex("#000"), // white
	// 	},
	// }
	// scene = scenes.SineWaveWCross{
	// 	XYRatio:      0.0001,
	// 	SigmoidRatio: 2.0,
	// 	SinCycles:    3,
	// 	TScale:       0.3,
	// 	// TOffset:      0.0,
	// 	// Gradient:     color.Grayscale,
	// 	Gradient: gradient,
	// }
	scene = scenes.DummyTriangle()

	// scene = scenes.HorizGradient{
	// 	Gradient: gradient,
	// }

	// scene = scenes.SineWave{
	// 	XYRatio:      0.1,
	// 	SigmoidRatio: 3,
	// 	SinCycles:    3,
	// 	Gradient:     color.Grayscale,
	// }
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 2 {
		log.Fatal("Insufficient arguments, expect <type> <output.gif>.")
	}

	format := argsWithoutProg[0]
	outFile, err := filepath.Abs(argsWithoutProg[1])
	if err != nil {
		log.Fatalf("Invalid file path %s", err)
	}
	switch format {
	case GIF_FORMAT:
		err := renderGIF(scene, width, height, nFrameCount, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	case PNG_FORMAT:
		t := 0.5
		err := renderPNG(scene, width, height, t, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	default:
		log.Fatalf("Unknown format %s", format)
	}
}

type Pixel struct {
	X int
	Y int
}

func getGIFFrame(scene scenes.Scene, width, height int, t float64) *image.Paletted {
	grid := getPixelGrid(scene, width, height, t)
	palette := generateFramePalette(scene, grid, t)
	now := time.Now()
	img := image.NewPaletted(
		image.Rect(
			0, 0, width, height,
		),
		// palette.WebSafe,
		color.ToInterfaceSlice(palette),
	)

	for pixel, color := range grid {
		img.Set(pixel.X, pixel.Y, color)
	}
	fmt.Printf("Palette setting took %s\n", time.Since(now))
	return img
}

func generateFramePalette(scene scenes.Scene, pixels map[Pixel]color.Color, t float64) []color.Color {
	start := time.Now()
	defer func() {
		fmt.Printf("Palette generation took %s\n", time.Since(start))
	}()
	if palette := scene.GetColorPalette(t); len(palette) > 0 {
		return palette
	}
	colorMap := map[color.Color]struct{}{}
	for _, c := range pixels {
		colorMap[c] = struct{}{}
	}
	palette := make([]color.Color, 0, len(colorMap))
	for color, _ := range colorMap {
		palette = append(palette, color)
	}
	fmt.Printf("%d colors in palette \n", len(colorMap))
	if len(colorMap) > 255 {
		var err error
		palette, err = computePalette(palette)
		if err != nil {
			panic(err)
		}
	}
	return palette
}

func getPixelGrid(scene scenes.Scene, width, height int, t float64) map[Pixel]color.Color {
	start := time.Now()
	grid := map[Pixel]color.Color{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			xR, yR := getImageSpace(x, width), getImageSpace(y, height)
			grid[Pixel{x, y}] = scene.GetColor(xR, yR, t)
		}
	}
	fmt.Printf("Pixel generation took %s\n", time.Since(start))
	return grid
}

func getImage(scene scenes.Scene, width, height int, t float64) image.Image {
	grid := map[Pixel]go_color.Color{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			xR, yR := getImageSpace(x, width), getImageSpace(y, height)
			samples := make([]color.Color, interpolateN)
			for i := range interpolateN {
				dx, dy := getPixelWiggle(width), getPixelWiggle(height)
				samples[i] = scene.GetColor(xR+rand.Float64()*dx, yR+rand.Float64()*dy, t)
			}
			pixelColor := color.Average(samples)
			// pixelColor := scene.GetColor(xR, yR, t)
			grid[Pixel{x, y}] = pixelColor
		}
	}
	img := image.NewRGBA(
		image.Rect(
			0, 0, width, height,
		),
	)
	for pixel, color := range grid {
		img.Set(pixel.X, pixel.Y, color)
	}
	return img
}

// convert coordinate from pixel space (0, pixels-1) to image space (-1.0, 1.0)
func getImageSpace(x, pixels int) float64 {
	return 2*float64(x)/float64(pixels) - 1.0
}

// get the width/height of a pixel in image space, in one dimension
func getPixelWiggle(pixels int) float64 {
	return 2.0 / float64(pixels)
}

func getUniformDelays(nFrames, delay int) []int {
	delays := make([]int, nFrames)
	for i := range nFrames {
		delays[i] = delay
	}
	return delays
}

func renderPNG(scene scenes.Scene, width, height int, t float64, outfile string) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	frame := getImage(scene, width, height, t)
	return png.Encode(f, frame)
}

func renderGIF(scene scenes.Scene, width, height, nFrames int, outfile string) error {
	start := time.Now()
	defer func() {
		fmt.Printf("GIF generation took %s\n", time.Since(start))
	}()
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	frames := make([]*image.Paletted, nFrames)
	var wg sync.WaitGroup
	for i := range nFrames {
		wg.Add(1)
		go func() {
			defer wg.Done()

			fmt.Printf("frame %d\n", i)
			t := float64(i) / float64(nFrames-1) // range [0.0, 1.0]
			frames[i] = getGIFFrame(scene, width, height, t)
			fmt.Printf("finished frame %d\n", i)
		}()
	}
	wg.Wait()
	delays := getUniformDelays(nFrames, frameSpacing)
	animation := gif.GIF{
		Image:     frames,
		Delay:     delays, //[]int in 10ms
		LoopCount: 0,
	}
	return gif.EncodeAll(f, &animation)

}
