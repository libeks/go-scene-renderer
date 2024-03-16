package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/renderer"
	"github.com/libeks/go-scene-renderer/scenes"
)

const (
	PNG_FORMAT = "png"
	MP4_FORMAT = "mp4"
	do_pprof   = false
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
	if do_pprof {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	var imageFlag = flag.String("image", "default", "image options, either <width>,<height>,<interpolate> or one of default/test/hidef")
	var videoFlag = flag.String("video", "default", "video options, either <width>,<height>,<interpolate>,<nframes>,<frameRate> or one of default/test/intermediate/hidef")

	flag.Parse()
	fmt.Printf("image flag: %s\n", *imageFlag)
	fmt.Printf("video flag: %s\n", *videoFlag)
	argsWithoutProg := flag.Args()
	fmt.Printf("Args: %+v\n", argsWithoutProg)
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
		imagePreset, err := renderer.ParseImagePreset(*imageFlag)
		if err != nil {
			log.Fatalf("%s", err)
		}
		err = renderer.RenderPNG(scene.GetFrame(t), imagePreset, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	case MP4_FORMAT:
		videoPreset, err := renderer.ParseVideoPreset(*videoFlag)
		if err != nil {
			log.Fatalf("%s", err)
		}
		err = renderer.RenderVideo(scene, videoPreset, outFile)
		if err != nil {
			fmt.Printf("Failure %s\n", err)
		}
	default:
		log.Fatalf("Unknown format %s", format)
	}

	if do_pprof {
		f1, err := os.Create("mem.pprof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f1.Close() // error handling omitted for example
		if err := pprof.WriteHeapProfile(f1); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
