package main

import (
	"context"
	"fmt"
	"image"
	go_color "image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/scenes"
	"golang.org/x/sync/errgroup"
)

const (
	// frameSpacing = 7
	// nFrameCount  = 100
	// width        = 1000
	// height       = 1000

	// interpolateN = 1

	cleanUpFrameCache = false

	GIF_FORMAT = "gif"
	PNG_FORMAT = "png"
	MP4_FORMAT = "mp4"
)

type ImagePreset struct {
	width        int
	height       int
	interpolateN int
}

type VideoPreset struct {
	ImagePreset
	nFrameCount int
	frameRate   int
}

var (
	videoPresetTest = VideoPreset{
		ImagePreset: ImagePreset{
			width:        200,
			height:       200,
			interpolateN: 1,
		},
		nFrameCount: 30,
		frameRate:   15,
	}
	videoPresetHiDef = VideoPreset{
		ImagePreset: ImagePreset{
			width:        1000,
			height:       1000,
			interpolateN: 1,
		},
		nFrameCount: 100,
		frameRate:   30,
	}
	// defaultVideoPreset = videoPresetTest
	defaultVideoPreset = videoPresetHiDef
	// defaultImagePreset =
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
	scene = scenes.DummySpinningCube()
	// scene = scenes.DummyTriangle()

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
	// case GIF_FORMAT:
	// 	err := renderGIF(scene, width, height, nFrameCount, outFile)
	// 	if err != nil {
	// 		fmt.Printf("Failure %s\n", err)
	// 	}
	case PNG_FORMAT:
		t := 0.5
		err := renderPNG(scene.GetFrame(t), ImagePreset{
			width:        200,
			height:       200,
			interpolateN: 1,
		}, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	case MP4_FORMAT:
		err := renderVideo(scene, defaultVideoPreset, outFile)
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

func renderVideo(scene scenes.DynamicScene, vp VideoPreset, outFile string) error {
	start := time.Now()
	ctx := context.Background()
	// clean up frames in temp directory before starting
	tmpDirectory := ".tmp"
	fileWildcardPattern := filepath.Join(".", tmpDirectory, "frame_*.png")
	if err := cleanUpTempFiles(fileWildcardPattern); err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	// Write output files to temporary directory
	outFileFormat := filepath.Join(tmpDirectory, "frame_%03d.png")
	if err := createSubdirectories(outFileFormat); err != nil {
		return err
	}
	for i := range vp.nFrameCount {
		g.Go(func() error {
			outFile := fmt.Sprintf(outFileFormat, i)
			f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer f.Close()
			fmt.Printf("frame %d\n", i)
			t := float64(i) / float64(vp.nFrameCount-1) // range [0.0, 1.0]
			frameObj := scene.GetFrame(t)
			frame := getImage(frameObj, vp.ImagePreset)
			fmt.Printf("finished frame %d\n", i)
			return png.Encode(f, frame)
		})

	}
	if err := g.Wait(); err != nil {
		return err
	}
	fmt.Printf("PNG frame generation took %s\n", time.Since(start))
	fmt.Printf("Finished rendering PNG frames\n")
	// encoder := "yuv444p"
	encoder := "yuv420p"
	// format := "libx265"
	format := "libx264"
	cmd := exec.Command(
		"ffmpeg", "-y",
		// "-f", "lavfi",
		"-framerate", fmt.Sprintf("%d", vp.frameRate),
		"-i", outFileFormat,
		"-c:v", format,
		"-pix_fmt", encoder,
		"-profile:v", "main",
		"-level", "3.1",
		"-preset", "medium",
		"-crf", "23",
		"-x264-params", "ref=4",
		// "-preset", "slow",
		// "-x265-params", "lossless=1",
		"-b:v", "5000k",
		// "-i", "anullsrc=channel_layout=stereo:sample_rate=44100",
		// "-c:a", "aac",
		outFile)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(stdout))
		return err
	}
	fmt.Print(string(stdout))
	if cleanUpFrameCache {
		return cleanUpTempFiles(fileWildcardPattern)
	}
	return nil
}

func cleanUpTempFiles(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, f := range files {
		// fmt.Printf("About to remove %s\n", f)
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func createSubdirectories(outFileFormat string) error {
	return os.MkdirAll(filepath.Dir(outFileFormat), os.ModePerm)
}

// func getGIFFrame(scene scenes.GIFScene, width, height int, t float64) *image.Paletted {
// 	grid := getPixelGrid(scene, width, height, t)
// 	palette := generateFramePalette(scene, grid, t)
// 	now := time.Now()
// 	img := image.NewPaletted(
// 		image.Rect(
// 			0, 0, width, height,
// 		),
// 		// palette.WebSafe,
// 		color.ToInterfaceSlice(palette),
// 	)

// 	for pixel, color := range grid {
// 		img.Set(pixel.X, pixel.Y, color)
// 	}
// 	fmt.Printf("Palette setting took %s\n", time.Since(now))
// 	return img
// }

func generateFramePalette(scene scenes.GIFScene, pixels map[Pixel]color.Color, t float64) []color.Color {
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

func getPixelGrid(scene scenes.Frame, width, height int) map[Pixel]color.Color {
	start := time.Now()
	grid := map[Pixel]color.Color{}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			xR, yR := getImageSpace(x, width), getImageSpace(y, height)
			grid[Pixel{x, y}] = scene.GetColor(xR, yR)
		}
	}
	fmt.Printf("Pixel generation took %s\n", time.Since(start))
	return grid
}

func getImage(scene scenes.Frame, ip ImagePreset) image.Image {
	grid := map[Pixel]go_color.Color{}
	for x := 0; x < ip.width; x++ {
		for y := 0; y < ip.height; y++ {
			xR, yR := getImageSpace(x, ip.width), getImageSpace(y, ip.height)
			var pixelColor color.Color
			if ip.interpolateN > 1 {

				samples := make([]color.Color, ip.interpolateN)
				for i := range ip.interpolateN {
					dx, dy := getPixelWiggle(ip.width), getPixelWiggle(ip.height)
					samples[i] = scene.GetColor(xR+rand.Float64()*dx, yR+rand.Float64()*dy)
				}
				pixelColor = color.Average(samples)
			} else {
				pixelColor = scene.GetColor(xR, yR)
			}

			// insert pixels with flipped y- coord, so y would be -1 at the bottom, +1 at the top of the image
			grid[Pixel{x, ip.height - y}] = pixelColor
		}
	}
	img := image.NewRGBA(
		image.Rect(
			0, 0, ip.width, ip.height,
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

func renderPNG(scene scenes.Frame, im ImagePreset, outfile string) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	frame := getImage(scene, im)
	return png.Encode(f, frame)
}

// func renderGIF(scene scenes.GIFScene, width, height, nFrames int, outfile string) error {
// 	start := time.Now()
// 	defer func() {
// 		fmt.Printf("GIF generation took %s\n", time.Since(start))
// 	}()
// 	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	frames := make([]*image.Paletted, nFrames)
// 	var wg sync.WaitGroup
// 	for i := range nFrames {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			fmt.Printf("frame %d\n", i)
// 			t := float64(i) / float64(nFrames-1) // range [0.0, 1.0]
// 			frames[i] = getGIFFrame(scene, width, height, t)
// 			fmt.Printf("finished frame %d\n", i)
// 		}()
// 	}
// 	wg.Wait()
// 	delays := getUniformDelays(nFrames, frameSpacing)
// 	animation := gif.GIF{
// 		Image:     frames,
// 		Delay:     delays, //[]int in 10ms
// 		LoopCount: 0,
// 	}
// 	return gif.EncodeAll(f, &animation)

// }
