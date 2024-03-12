package renderer

import (
	"context"
	"fmt"
	"image"
	go_color "image/color"
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
	frameConcurrency = 10 // should depend on video preset. Too many and you'll run out of memory.
)

var (
	cleanUpFrameCache = false
)

type Renderer struct {
	lineChannel chan int
	fileChannel chan int
	doneChannel chan int
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
	if err := cleanUpTempFiles(fileWildcardPattern); err != nil {
		return err
	}

	outFileFormat := filepath.Join(tmpDirectory, "frame_%03d.png")
	if err := createSubdirectories(outFileFormat); err != nil {
		return err
	}
	fmt.Printf("Rendering frames...\n")
	r := newRenderer()
	var sem = semaphore.NewWeighted(int64(frameConcurrency))
	go r.progressbar(vp.nFrameCount, vp.nFrameCount*vp.height) // block until completion
	for i := range vp.nFrameCount {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			return err
		}
		go func() {
			// fmt.Printf("starting frame %d\n", i)
			outFile := fmt.Sprintf(outFileFormat, i)
			f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			t := float64(i) / float64(vp.nFrameCount-1) // range [0.0, 1.0]
			// fmt.Printf("getting frame %d\n", i)
			frameObj := scene.GetFrame(t)
			frame := r.getImage(frameObj, vp.ImagePreset)
			err = png.Encode(f, frame)
			if err != nil {
				panic(err)
			}
			// fmt.Printf("ending frame %d\n", i)
			sem.Release(1)
			r.fileChannel <- 1

		}()
	}
	r.wait() // block until completion

	// fmt.Printf("Got %d lines, %d files \n", lineProgress, fileProgress)
	// if err := g.Wait(); err != nil {
	// 	return err
	// }
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
			// fmt.Println("received line")
		case <-r.fileChannel:
			// fmt.Printf("File progress %d\n", fileProgress)
			fileProgress += 1
			if fileProgress == nFiles {
				// fmt.Printf("sending done signal\n")
				r.doneChannel <- 1
				// fmt.Printf("sent done signal\n")
				return
			}
		}
	}
}
func (r Renderer) wait() {
	// fmt.Printf("waiting\n")
	<-r.doneChannel
	// fmt.Printf("done waiting\n")
}

func (r Renderer) getImage(scene scenes.Frame, ip ImagePreset) image.Image {
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
		r.lineChannel <- 1
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
