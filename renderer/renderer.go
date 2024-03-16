package renderer

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/libeks/go-scene-renderer/scenes"
	"github.com/schollz/progressbar"
	"golang.org/x/sync/semaphore"
)

const (
	frameConcurrency  = 10   // should depend on video preset. Too many and you'll operate close to full memory, slowing rendering down.
	generateVideoPNGs = true // set to false to debug ffmpeg settings without recreating image files (files have to exist in .tmp/)
)

var (
	cleanUpFrameCache = false
)

// Renderer does two things - tracks progress of per-frame goroutines, and updates
// a progress bar based on the number of image rows that have been rendered so far
type Renderer struct {
	lineChannel chan int // each line completion is sent on lineChannel
	fileChannel chan int // each file completion is sent on fileChannel
	doneChannel chan int // doneChannel sends a message when all frames are rendered
}

func newRenderer() Renderer {
	return Renderer{
		lineChannel: make(chan int, 10),
		fileChannel: make(chan int, 10),
		doneChannel: make(chan int, 1),
	}
}

func RenderVideo(scene scenes.DynamicScene, vp VideoPreset, outFile string) error {
	start := time.Now()
	// clean up frames in temp directory before starting
	tmpDirectory := ".tmp"
	fileWildcardPattern := filepath.Join(".", tmpDirectory, "frame_*.png")
	outFileFormat := filepath.Join(tmpDirectory, "frame_%03d.png")
	if generateVideoPNGs {
		if err := cleanUpTempFiles(fileWildcardPattern); err != nil {
			return err
		}
		if err := createSubdirectories(outFileFormat); err != nil {
			return err
		}
		fmt.Printf("Rendering frames...\n")
		r := newRenderer()
		var sem = semaphore.NewWeighted(int64(frameConcurrency))
		go r.progressbar(vp.nFrameCount, vp.nFrameCount*vp.height) // start progressbar before launching goroutines to not deadlock
		for i := range vp.nFrameCount {
			if err := sem.Acquire(context.Background(), 1); err != nil {
				return err
			}
			go func() {
				outFile := fmt.Sprintf(outFileFormat, i)
				f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				t := float64(i) / float64(vp.nFrameCount-1) // range [0.0, 1.0]
				frameObj := scene.GetFrame(t)
				frame := r.getImage(frameObj, vp.ImagePreset)
				err = png.Encode(f, frame)
				if err != nil {
					panic(err)
				}
				sem.Release(1)
				r.fileChannel <- 1

			}()
		}
		r.wait() // block until completion

		fmt.Printf("PNG frame generation took %s\n", time.Since(start))
	}

	// render video file from png frame images in .tmp/
	// encoder := "yuv444p"
	encoder := "yuv420p"
	format := "libx265"
	// format := "libx264"
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
		"-crf", "15",
		// "-x264-params", "ref=4",
		// "-preset", "slow",
		// "-x265-params", "lossless=1",
		// "-b:v", "10000k",
		"-tag:v", "hvc1",
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

func RenderPNG(scene scenes.Frame, im ImagePreset, outfile string) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var frame image.Image
	r := newRenderer()
	go r.progressbar(1, im.height) // block until completion
	go func() {
		frame = r.getImage(scene, im)
		r.fileChannel <- 1
	}()
	r.wait()
	return png.Encode(f, frame)
}

func (r Renderer) progressbar(nFiles, nLines int) {
	fileProgress := 0
	lineProgress := 0
	bar := progressbar.New(nLines)
	for {
		select {
		case <-r.lineChannel:
			lineProgress += 1
			bar.Add(1)
		case <-r.fileChannel:
			fileProgress += 1
			if fileProgress == nFiles {
				r.doneChannel <- 1
				return
			}
		}
	}
}
func (r Renderer) wait() {
	<-r.doneChannel
}

func (r Renderer) getImage(scene scenes.Frame, ip ImagePreset) image.Image {
	img := image.NewRGBA(
		image.Rect(
			0, 0, ip.width, ip.height,
		),
	)
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
			img.Set(x, ip.height-y, pixelColor)
		}
		r.lineChannel <- 1
	}
	return img
}
