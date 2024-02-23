package main

import (
	"fmt"
	"log"
	"os"

	"image"
	"image/color"
	"image/gif"

	"github.com/libeks/go-scene-renderer/scenes"
)

const (
	frameSpacing = 10
	nFrameCount  = 30
	width        = 500
	height       = 500
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		log.Fatal("Insufficient arguments, expect <output.gif>.")
	}

	outfile := ""
	if len(argsWithoutProg) > 0 {
		outfile = argsWithoutProg[0]
	}
	scene := scenes.SineWaveWBump{
		Frame:        scenes.PictureFrame{width, height},
		XYRatio:      20,
		SigmoidRatio: 3,
		SinCycles:    3,
	}
	err := renderScene(scene, width, height, nFrameCount, outfile)
	if err != nil {
		fmt.Printf("Failure %s\n", err)
	}
}

type Pixel struct {
	X int
	Y int
}

func getFrame(scene scenes.Scene, width, height int, t float64) *image.Paletted {
	colorMap := map[color.Color]struct{}{}
	grid := map[Pixel]color.Color{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixelColor := scene.GetPixel(x, y, t)
			grid[Pixel{x, y}] = pixelColor
			colorMap[pixelColor] = struct{}{}
		}
	}
	fmt.Printf("%d colors in palette \n", len(colorMap))
	colors := make([]color.Color, 0, len(colorMap))
	for color, _ := range colorMap {
		colors = append(colors, color)
	}
	img := image.NewPaletted(
		image.Rect(
			0, 0, width, height,
		),
		// palette.WebSafe,
		colors,
	)
	for pixel, color := range grid {
		img.Set(pixel.X, pixel.Y, color)
	}

	return img
}

func getUniformDelays(nFrames, delay int) []int {
	delays := make([]int, nFrames)
	for i := range nFrames {
		delays[i] = delay
	}
	return delays
}

func renderScene(scene scenes.Scene, width, height, nFrames int, outfile string) error {
	fmt.Printf("file '%s'\n", outfile)
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	frames := make([]*image.Paletted, nFrames)
	for i := range nFrames {
		fmt.Printf("frame %d\n", i)
		t := float64(i) / float64(nFrames-1) // range [0.0, 1.0]
		frames[i] = getFrame(scene, width, height, t)
	}
	delays := getUniformDelays(nFrames, frameSpacing)
	animation := gif.GIF{
		Image:     frames,
		Delay:     delays, //[]int in 10ms
		LoopCount: 10,
	}
	return gif.EncodeAll(f, &animation)

}
