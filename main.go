package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/renderer"
	"github.com/libeks/go-scene-renderer/scenes"
)

const (
	PNG_FORMAT = "png"
	MP4_FORMAT = "mp4"
)

var (
	// defaultVideoPreset = renderer.VideoPresetTest
	defaultVideoPreset = renderer.VideoPresetHiDef
	// defaultImagePreset = renderer.ImagePresetTest
	defaultImagePreset = renderer.ImagePresetHiDef

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
	// scene = scenes.DummySpinningCube(scenes.Uniform{color.Black})
	scene = scenes.DummySpinningCube(scenes.SineWave{
		XYRatio:      0.1,
		SigmoidRatio: 2,
		SinCycles:    3,
		Gradient:     color.Grayscale,
	})
	// scene = scenes.DummyTriangle()

	// scene = scenes.HorizGradient{
	// 	Gradient: gradient,
	// }

	// scene = scenes.SineWave{
	// 	XYRatio:      0.1,
	// 	SigmoidRatio: 2,
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
	case PNG_FORMAT:
		t := 0.5
		err := renderer.RenderPNG(scene.GetFrame(t), defaultImagePreset, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	case MP4_FORMAT:
		err := renderer.RenderVideo(scene, defaultVideoPreset, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	default:
		log.Fatalf("Unknown format %s", format)
	}
}
