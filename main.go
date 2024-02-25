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
	"sync"
	"time"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/scenes"
)

const (
	frameSpacing = 5
	nFrameCount  = 10
	width        = 500
	height       = 500

	interpolateN = 1

	GIF_FORMAT = "gif"
	PNG_FORMAT = "png"
)

var (
	scene = scenes.SineWaveWCross{
		Frame:        scenes.PictureFrame{width, height},
		XYRatio:      0.01,
		SigmoidRatio: 3,
		SinCycles:    3,
		// Gradient:     color.Grayscale,
		Gradient: color.Gradient{
			Start: color.Hex("#DDF522"),
			End:   color.Hex("#A0514C"),
		},
	}

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

	outFile := argsWithoutProg[1]
	format := argsWithoutProg[0]
	switch format {
	case GIF_FORMAT:
		err := renderGIF(scene, width, height, nFrameCount, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	case PNG_FORMAT:
		t := 1.0
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
	start := time.Now()
	colorMap := map[color.Color]struct{}{}
	grid := map[Pixel]go_color.Color{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			xR, yR := getImageSpace(x, width), getImageSpace(y, height)
			pixelColor := scene.GetColor(xR, yR, t)
			grid[Pixel{x, y}] = pixelColor
			colorMap[pixelColor] = struct{}{}
		}
	}
	fmt.Printf("%d colors in palette \n", len(colorMap))
	var palette []color.Color
	if palette = scene.GetColorPalette(t); len(palette) == 0 {
		palette = make([]color.Color, 0, len(colorMap))
		for color, _ := range colorMap {
			palette = append(palette, color)
		}
		if len(colorMap) > 255 {
			var err error
			palette, err = computePalette(palette)
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Printf("Palette generation took %s\n", time.Since(start))
	t2 := time.Now()
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
	fmt.Printf("Pixel setting took %s\n", time.Since(t2))
	return img
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
		LoopCount: 10,
	}
	return gif.EncodeAll(f, &animation)

}
