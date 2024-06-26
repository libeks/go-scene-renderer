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
	frameConcurrency       = 10   // should depend on video preset. Too many and you'll operate close to full memory, slowing rendering down.
	generateVideoPNGs      = true // set to false to debug ffmpeg settings without recreating image files (files have to exist in .tmp/)
	minWindowWidth         = 3
	minWindowCount         = 1
	wireframeTriangleDepth = false
	applyWireframe         = false // draw wireframes on top of rendered objects
	render_h265            = false // if false, will render with h264
)

var (
	cleanUpFrameCache = false
)

type fileReport struct {
	frameID int
	// triangleChecks  int
	// trianglesInView int
}

type chunkReport struct {
	pixels            int
	triangleChecks    int
	trianglesInWindow int
}

// Renderer does two things - tracks progress of per-frame goroutines, and updates
// a progress bar based on the number of image rows that have been rendered so far
type Renderer struct {
	lineChannel chan chunkReport // each line completion is sent on lineChannel
	fileChannel chan fileReport  // each file completion is sent on fileChannel
	doneChannel chan struct{}    // doneChannel sends a message when all frames are rendered
	offsets     []Offset
}

func newRenderer(ip ImagePreset) Renderer {
	offsets := make([]Offset, ip.interpolateN)
	dx, dy := getPixelWiggle(ip.width), getPixelWiggle(ip.height)
	for i := range ip.interpolateN {
		offsets[i] = Offset{rand.Float64() * dx, rand.Float64() * dy}
	}
	return Renderer{
		lineChannel: make(chan chunkReport, 10),
		fileChannel: make(chan fileReport, 10),
		doneChannel: make(chan struct{}, 1),
		offsets:     offsets,
	}
}

func RenderVideo(scene scenes.DynamicScene, vp VideoPreset, outFile string, wireframe bool, triDepth bool) error {
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
		r := newRenderer(vp.ImagePreset)
		var sem = semaphore.NewWeighted(int64(frameConcurrency))
		go r.progressbar(vp.nFrameCount, vp.nFrameCount*vp.width*vp.height) // start progressbar before launching goroutines to not deadlock

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

				} else if triDepth {
					frame = r.getTriangleDepthImage(frameObj, vp.ImagePreset)
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
				r.fileChannel <- fileReport{
					frameID: i,
				}

			}()
		}
		r.wait() // block until completion

		fmt.Printf("\nPNG frame generation took %s\n", time.Since(start))
	}
	fmt.Printf("Encoding with ffmpeg...\n")
	// render video file from png frame images in .tmp/
	params := []string{
		"-y",
		"-framerate", fmt.Sprintf("%.2f", vp.frameRate),
		"-i", outFileFormat,
	}
	if render_h265 {
		params = append(params,
			"-c:v", "libx265",
			// "-pix_fmt", "yu v420p",
			"-pix_fmt", "yuv420p10le",
			// "-profile:v", "main",
			"-level", "3.1",
			"-preset", "medium",
			"-crf", "14",
			"-tag:v", "hvc1",
		)
	} else {
		params = append(params,
			"-c:v", "libx264",
			"-pix_fmt", "yuv420p",
			"-profile:v", "main",
			"-level", "3.1",
			"-preset", "medium",
			"-crf", "14",
		)
	}
	params = append(params, outFile)

	// fmt.Printf("params %v\n", params)
	cmd := exec.Command(
		"ffmpeg",
		params...)
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

func RenderPNG(scene scenes.StaticScene, im ImagePreset, outfile string, wireframe bool, triDepth bool) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var frame *Image
	r := newRenderer(im)
	go r.progressbar(1, im.width) // block until completion
	go func() {
		if wireframe {
			if wireframeTriangleDepth {
				frame = r.getTriangleDepthImage(scene, im)
			} else {
				frame = r.getWireframeImage(scene, im)
			}
		} else if triDepth {
			frame = r.getTriangleDepthImage(scene, im)
		} else {
			frame = r.getWindowedImage(scene, im)
			if applyWireframe {
				frame = r.applyWireframeToImage(frame, scene, im)
			}
		}
		png.Encode(f, frame.GetImage())
		r.fileChannel <- fileReport{
			frameID: 0,
		}
	}()
	r.wait()
	return nil
}

func (r Renderer) progressbar(nFiles, nPixels int) {
	fileProgress := 0
	pixelProgress := 0
	bar := progressbar.New(nPixels)
	for {
		select {
		case windowProg := <-r.lineChannel:
			pixelProgress += windowProg.pixels
			bar.Add(windowProg.pixels)
		case <-r.fileChannel:
			fileProgress += 1
			if fileProgress == nFiles {
				r.doneChannel <- struct{}{}
				return
			}
		}
	}
}

func (r Renderer) wait() {
	<-r.doneChannel
	fmt.Println("")
}

type Offset struct {
	dx float64
	dy float64
}

func (r Renderer) getWindowedImage(scene scenes.StaticScene, ip ImagePreset) *Image {
	img := NewImage(ip)
	windows := subdivideSceneIntoWindows(scene, ip)
	var imageTriangles, imageChecks int
	for _, window := range windows {
		var nTriangles, windowChecks int
		for x := window.xMin; x < window.xMax; x++ {
			for y := window.yMin; y < window.yMax; y++ {
				xR, yR := getImageSpace(x, ip.width), getImageSpace(y, ip.height)
				var pixelColor colors.Color
				if ip.interpolateN > 1 {
					samples := make([]colors.Color, ip.interpolateN)
					for i, offset := range r.offsets {
						var nChecks int
						samples[i], nTriangles, nChecks = window.GetColor(xR+offset.dx, yR+offset.dy)
						windowChecks += nChecks
					}
					pixelColor = colors.Average(samples)
				} else {
					var nChecks int
					pixelColor, nTriangles, nChecks = window.GetColor(xR, yR)
					windowChecks += nChecks
				}

				img.Set(x, y, pixelColor)
			}
		}
		windowPixels := (window.yMax - window.yMin) * (window.xMax - window.xMin)
		r.lineChannel <- chunkReport{
			pixels:            windowPixels,
			triangleChecks:    windowChecks,
			trianglesInWindow: nTriangles,
		}
		imageTriangles += windowPixels * nTriangles
		imageChecks += windowChecks
	}
	imagePixels := ip.width * ip.height
	fmt.Printf("Image had %d pixels, %.3f triangles per pixel, %.3f checks per pixel\n",
		imagePixels,
		float64(imageTriangles)/float64(imagePixels),
		float64(imageChecks)/float64(imagePixels),
	)
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
			pixA := toImagePixel(line.A, ip.width, ip.height)
			pixB := toImagePixel(line.B, ip.width, ip.height)
			if pixA == nil || pixB == nil {
				fmt.Printf("Skipping line %v since one or both pixels are outside of screen\n", line)
				continue
			}
			greenBlack := colors.SimpleGradient{
				Start: colors.Green,
				End:   colors.Green,
			}
			ratio := 8.0
			colorA := greenBlack.Interpolate(2*maths.Sigmoid(line.ADepth/ratio) - 1)
			colorB := greenBlack.Interpolate(2*maths.Sigmoid(line.BDepth/ratio) - 1)
			img.RenderLine(NewRasterLine(
				*pixA,
				*pixB,
			), colors.SimpleGradient{Start: colorA, End: colorB})
		}
		bbox := tri.GetBoundingBox()
		pixA := toImagePixel(bbox.TopLeft, ip.width, ip.height)
		pixB := toImagePixel(bbox.BottomRight, ip.width, ip.height)
		if pixA == nil || pixB == nil {
			continue
		}
		gradient := colors.SimpleGradient{Start: colors.Red, End: colors.Red}
		img.RenderLine(NewRasterLine(*pixA, RasterPixel{pixA.X, pixB.Y}), gradient)
		img.RenderLine(NewRasterLine(*pixA, RasterPixel{pixB.X, pixA.Y}), gradient)
		img.RenderLine(NewRasterLine(*pixB, RasterPixel{pixA.X, pixB.Y}), gradient)
		img.RenderLine(NewRasterLine(*pixB, RasterPixel{pixB.X, pixA.Y}), gradient)
	}
	return img
}

func (r Renderer) getWireframeImage(scene scenes.StaticScene, ip ImagePreset) *Image {
	img := NewImage(ip)
	// set to black bakcground
	img.Fill(colors.Black)
	r.applyWireframeToImage(img, scene, ip)
	r.lineChannel <- chunkReport{
		pixels: ip.height * ip.width,
	}
	return img
}

func (r Renderer) getTriangleDepthImage(scene scenes.StaticScene, ip ImagePreset) *Image {
	img := NewImage(ip)
	// set to black bakcground
	img.Fill(colors.Black)
	r.lineChannel <- chunkReport{
		pixels: ip.height * ip.width,
	}
	drawBorders := false
	windows := subdivideSceneIntoWindows(scene, ip)
	gradient := colors.LinearGradient{Points: []colors.Color{colors.Red, colors.Green, colors.White}}
	var pixelColor colors.Color
	for _, window := range windows {
		for x := window.xMin; x < window.xMax; x++ {
			for y := window.yMin; y < window.yMax; y++ {
				nTriangles := len(window.triangles)
				if nTriangles == 0 {
					pixelColor = colors.Black
				} else {
					if drawBorders && (x == window.xMin || y == window.yMin) {
						pixelColor = colors.Green
					} else {
						pixelColor = gradient.Interpolate(float64(nTriangles) / 30.0)
					}
				}
				img.Set(x, y, pixelColor)
			}
		}
	}
	return img
}
