package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/libeks/go-scene-renderer/renderer"
)

const (
	PNG_FORMAT = "png"
	MP4_FORMAT = "mp4"
	do_pprof   = false
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
	argsWithoutProg := flag.Args()
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
