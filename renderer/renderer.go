package renderer

import (
	"context"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/scenes"
	"github.com/schollz/progressbar"
	"golang.org/x/sync/semaphore"
)

const (
	frameConcurrency       = 10    // should depend on video preset. Too many and you'll operate close to full memory, slowing rendering down.
	generateVideoPNGs      = false // set to false to debug ffmpeg settings without recreating image files (files have to exist in .tmp/)
	minWindowWidth         = 10
	minWindowCount         = 1
	wireframeTriangleDepth = false
	applyWireframe         = false // draw wireframes on top of rendered objects
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

func RenderVideo(scene scenes.DynamicScene, vp VideoPreset, outFile string, wireframe bool) error {
	start := time.Now()
	// clean up frames in temp directory before starting
	tmpDirectory := ".tmp"
	fileWildcardPattern := filepath.Join(".", tmpDirectory, "frame_*.png")
	outFileFormat := filepath.Join(tmpDirectory, "frame_%03d.png")
	if generateVideoPNGs {
		fmt.Printf("Preparing setup...\n")
		if err := cleanUpTempFiles(fileWildcardPattern); err != nil {
			return err
		}
		if err := createSubdirectories(outFileFormat); err != nil {
			return err
		}
		r := newRenderer()
		var sem = semaphore.NewWeighted(int64(frameConcurrency))
		go r.progressbar(vp.nFrameCount, vp.nFrameCount*vp.height) // start progressbar before launching goroutines to not deadlock

		fmt.Printf("Rendering frames...\n")
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
				var frame *Image
				if wireframe {
					if wireframeTriangleDepth {
						frame = r.getTriangleDepthImage(frameObj, vp.ImagePreset)
					} else {
						frame = r.getWireframeImage(frameObj, vp.ImagePreset)
					}

				} else {
					frame = r.getWindowedImage(frameObj, vp.ImagePreset)
					if applyWireframe {
						frame = r.applyWireframeToImage(frame, frameObj, vp.ImagePreset)
					}
				}
				err = png.Encode(f, frame.GetImage())
				if err != nil {
					panic(err)
				}
				sem.Release(1)
				r.fileChannel <- 1

			}()
		}
		r.wait() // block until completion

		fmt.Printf("\nPNG frame generation took %s\n", time.Since(start))
	}
	fmt.Printf("Encoding with ffmpeg...\n")
	// render video file from png frame images in .tmp/
	encoder := "yuv420p"
	format := "libx265"
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-framerate", fmt.Sprintf("%d", vp.frameRate),
		"-i", outFileFormat,
		"-c:v", format,
		"-pix_fmt", encoder,
		"-profile:v", "main",
		"-level", "3.1",
		"-preset", "medium",
		// "-vf", "lutyuv=u=128:v=128",
		// "-b:v", "2600k",
		"-crf", "15",
		// "-x265-params", "lossless=1",
		"-tag:v", "hvc1",
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

func RenderPNG(scene scenes.StaticScene, im ImagePreset, outfile string, wireframe bool) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var frame *Image
	r := newRenderer()
	go r.progressbar(1, im.width) // block until completion
	go func() {
		if wireframe {
			if wireframeTriangleDepth {
				frame = r.getTriangleDepthImage(scene, im)
			} else {
				frame = r.getWireframeImage(scene, im)
			}

		} else {
			frame = r.getWindowedImage(scene, im)
			if applyWireframe {
				frame = r.applyWireframeToImage(frame, scene, im)
			}
		}
		png.Encode(f, frame.GetImage())
		r.fileChannel <- 1
	}()
	r.wait()
	return nil
}

func (r Renderer) progressbar(nFiles, nLines int) {
	fileProgress := 0
	lineProgress := 0
	bar := progressbar.New(nLines)
	for {
		select {
		case prog := <-r.lineChannel:
			lineProgress += prog
			bar.Add(prog)
		case prog := <-r.fileChannel:
			fileProgress += prog
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

func (r Renderer) getWindowedImage(scene scenes.StaticScene, ip ImagePreset) *Image {
	img := NewImage(ip)
	windows := subdivideSceneIntoWindows(scene, ip)
	pixelCount := 0
	for _, window := range windows {
		for x := window.xMin; x < window.xMax; x++ {
			for y := window.yMin; y < window.yMax; y++ {
				xR, yR := getImageSpace(x, ip.width), getImageSpace(y, ip.height)
				var pixelColor colors.Color
				if ip.interpolateN > 1 {
					// if len(window.triangles) > 0 && ip.interpolateN > 1 {

					samples := make([]colors.Color, ip.interpolateN)
					for i := range ip.interpolateN {
						dx, dy := getPixelWiggle(ip.width), getPixelWiggle(ip.height)
						samples[i] = window.GetColor(xR+rand.Float64()*dx, yR+rand.Float64()*dy)
					}
					pixelColor = colors.Average(samples)
				} else {
					pixelColor = window.GetColor(xR, yR)
				}

				img.Set(x, y, pixelColor)
				pixelCount += 1
				if pixelCount%ip.height == 0 {
					r.lineChannel <- 1
				}
			}
		}
	}
	return img
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func toImageDimension(d float64, pixelCount int) *int {
	// if d < -1.0 || d > 1.0 {
	// 	return nil
	// }
	v := int((d/2 + 0.5) * float64(pixelCount))
	return &v
}

func toImagePixel(p geometry.Pixel, width, height int) *RasterPixel {
	x := toImageDimension(p.X, width)
	y := toImageDimension(p.Y, height)
	// if x == nil || y == nil {
	// 	return nil
	// }
	return &RasterPixel{
		X: *x,
		Y: *y,
	}
}

func (r Renderer) applyWireframeToImage(img *Image, scene scenes.StaticScene, ip ImagePreset) *Image {
	triangles, _ := scene.Flatten()
	for _, tri := range triangles {
		for _, line := range tri.GetWireframe() {
			sceneA, aDepth := line.A.ToPixel()
			sceneB, bDepth := line.B.ToPixel()
			if sceneA == nil || sceneB == nil {
				fmt.Printf("Skipping line %s since it may be behind the screen\n", line)
				continue
			}
			pixA := toImagePixel(*sceneA, ip.width, ip.height)
			pixB := toImagePixel(*sceneB, ip.width, ip.height)
			if pixA == nil || pixB == nil {
				fmt.Printf("Skipping line %s since one or both pixels are outside of screen\n", line)
				continue
			}
			greenBlack := colors.SimpleGradient{
				colors.Green,
				colors.Black,
			}
			ratio := 8.0
			colorA := greenBlack.Interpolate(2*maths.Sigmoid(aDepth/ratio) - 1)
			colorB := greenBlack.Interpolate(2*maths.Sigmoid(bDepth/ratio) - 1)
			img.RenderLine(NewRasterLine(
				*pixA,
				*pixB,
			), colors.SimpleGradient{colorA, colorB})
		}
		bbox := tri.GetBoundingBox()
		pixA := toImagePixel(bbox.TopLeft, ip.width, ip.height)
		pixB := toImagePixel(bbox.BottomRight, ip.width, ip.height)
		if pixA == nil || pixB == nil {
			continue
		}
		img.RenderLine(NewRasterLine(*pixA, RasterPixel{pixA.X, pixB.Y}), colors.SimpleGradient{colors.Red, colors.Red})
		img.RenderLine(NewRasterLine(*pixA, RasterPixel{pixB.X, pixA.Y}), colors.SimpleGradient{colors.Red, colors.Red})
		img.RenderLine(NewRasterLine(*pixB, RasterPixel{pixA.X, pixB.Y}), colors.SimpleGradient{colors.Red, colors.Red})
		img.RenderLine(NewRasterLine(*pixB, RasterPixel{pixB.X, pixA.Y}), colors.SimpleGradient{colors.Red, colors.Red})
	}
	return img
}

func (r Renderer) getWireframeImage(scene scenes.StaticScene, ip ImagePreset) *Image {
	img := NewImage(ip)
	// set to black bakcground
	img.Fill(colors.Black)
	r.applyWireframeToImage(img, scene, ip)
	r.lineChannel <- ip.height
	return img
}

func (r Renderer) getTriangleDepthImage(scene scenes.StaticScene, ip ImagePreset) *Image {
	img := NewImage(ip)
	// set to black bakcground
	pixelColor := colors.Black
	for x := 0; x < ip.width; x++ {
		for y := 0; y < ip.height; y++ {
			img.Set(x, y, pixelColor)
		}
		r.lineChannel <- 1
	}
	windows := subdivideSceneIntoWindows(scene, ip)
	gradient := colors.SimpleGradient{colors.Black, colors.Red}
	for _, window := range windows {
		for x := window.xMin; x < window.xMax; x++ {
			for y := window.yMin; y < window.yMax; y++ {
				nTriangles := len(window.triangles)
				if x == window.xMin || y == window.yMin {
					pixelColor = colors.Green
				} else {
					pixelColor = gradient.Interpolate(float64(nTriangles) / 1.0)
				}
				img.Set(x, y, pixelColor)
			}
		}
	}
	return img
}
