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
	"github.com/libeks/go-scene-renderer/geometry"
	"github.com/libeks/go-scene-renderer/maths"
	"github.com/libeks/go-scene-renderer/objects"
	"github.com/libeks/go-scene-renderer/scenes"
	"github.com/schollz/progressbar"
	"golang.org/x/sync/semaphore"
)

const (
	frameConcurrency  = 10   // should depend on video preset. Too many and you'll operate close to full memory, slowing rendering down.
	generateVideoPNGs = true // set to false to debug ffmpeg settings without recreating image files (files have to exist in .tmp/)
	minWindowWidth    = 4
	minWindowCount    = 1
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
				var frame image.Image
				if wireframe {
					if combScene, ok := frameObj.(scenes.CombinedScene); ok {
						frame = r.getTriangleDepthImage(combScene, vp.ImagePreset)
					} else {
						frame = r.getWireframeImage(frameObj, vp.ImagePreset)
					}

				} else {
					if combScene, ok := frameObj.(scenes.CombinedScene); ok {
						frame = r.getWindowedImage(combScene, vp.ImagePreset)
					} else {
						frame = r.getImage(frameObj, vp.ImagePreset)
					}

				}
				err = png.Encode(f, frame)
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

func RenderPNG(scene scenes.Frame, im ImagePreset, outfile string, wireframe bool) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var frame image.Image
	r := newRenderer()
	go r.progressbar(1, im.width) // block until completion
	go func() {
		if wireframe {
			if combScene, ok := scene.(scenes.CombinedScene); ok {
				frame = r.getTriangleDepthImage(combScene, im)
			} else {
				frame = r.getWireframeImage(scene, im)
			}

		} else {
			if combScene, ok := scene.(scenes.CombinedScene); ok {
				frame = r.getWindowedImage(combScene, im)
			} else {
				frame = r.getImage(scene, im)
			}
		}
		png.Encode(f, frame)
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
			// fmt.Printf("Got %d on line channel\n", prog)
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

// Window specifies a renderable area along with the triangles within in
type Window struct {
	// coordinates are in image pixel space
	xMin      int                 // inclusive
	xMax      int                 // non-inclusive
	yMin      int                 // inclusive
	yMax      int                 // non-inclusive
	triangles []*objects.Triangle // list of triangles whose bounding box intersects the window
}

func (w Window) Width() int {
	return w.xMax - w.xMin
}

func (w Window) Height() int {
	return w.yMax - w.yMin
}

func (w Window) Bisect(ip ImagePreset) []Window {
	// window is too small to be divided any further
	if w.xMax-w.xMin < minWindowWidth {
		return nil
	}
	// window doesn't have enough triangles to be divided further
	if len(w.triangles) <= minWindowCount {
		return nil
	}
	xMid := (w.xMax-w.xMin)/2 + w.xMin
	yMid := (w.yMax-w.yMin)/2 + w.yMin
	midXImg, midYImg := getImageSpace(xMid, ip.width), getImageSpace(yMid, ip.height)
	tlW := []*objects.Triangle{}
	trW := []*objects.Triangle{}
	blW := []*objects.Triangle{}
	brW := []*objects.Triangle{}
	for _, tri := range w.triangles {
		bbox := tri.GetBoundingBox()
		if bbox.TopLeft.X <= midXImg && bbox.TopLeft.Y <= midYImg {
			tlW = append(tlW, tri)
		}
		if bbox.BottomRight.X >= midXImg && bbox.TopLeft.Y <= midYImg {
			trW = append(trW, tri)
		}
		if bbox.TopLeft.X <= midXImg && bbox.BottomRight.Y >= midYImg {
			blW = append(blW, tri)
		}
		if bbox.BottomRight.X >= midXImg && bbox.BottomRight.Y >= midYImg {
			brW = append(brW, tri)
		}
	}
	return []Window{
		{w.xMin, xMid, w.yMin, yMid, tlW},
		{xMid, w.xMax, w.yMin, yMid, trW},
		{w.xMin, xMid, yMid, w.yMax, blW},
		{xMid, w.xMax, yMid, w.yMax, brW},
	}
}

func (w Window) Subscene(scene scenes.CombinedScene) scenes.CombinedScene {
	objects := make([]objects.Object, len(w.triangles))
	for i, tri := range w.triangles {
		objects[i] = tri
	}
	return scenes.CombinedScene{
		Objects:    objects,
		Background: scene.Background,
	}
}

func subdivideSceneIntoWindows(scene scenes.Frame, ip ImagePreset) []Window {
	// start with one window for the whole image. Assume that all objects fall within the image
	windows := []Window{
		{0, ip.width, 0, ip.height, scene.Flatten()},
	}
	maxTriangles := 0
	totalWork := 0
	finalWindows := []Window{}
	for len(windows) > 0 {
		unprocessedWindows := append([]Window{}, windows...)
		windows = []Window{}
		for _, window := range unprocessedWindows {
			newOnes := window.Bisect(ip)
			if len(newOnes) == 0 {
				nTriangles := len(window.triangles)
				if nTriangles > maxTriangles {
					maxTriangles = nTriangles
				}
				totalWork += nTriangles * window.Width() * window.Height()
				finalWindows = append(finalWindows, window)
			} else {
				windows = append(windows, newOnes...)
			}
		}
	}
	// fmt.Printf(
	// 	"Image has %d windows, max # of triangles is %d, total work (pixels * triangles) is %d\n",
	// 	len(finalWindows), maxTriangles, totalWork,
	// )
	return finalWindows
}

func (r Renderer) getWindowedImage(scene scenes.CombinedScene, ip ImagePreset) image.Image {
	img := image.NewRGBA(
		image.Rect(
			0, 0, ip.width, ip.height,
		),
	)
	windows := subdivideSceneIntoWindows(scene, ip)
	pixelCount := 0
	for _, window := range windows {
		wScene := window.Subscene(scene)
		for x := window.xMin; x < window.xMax; x++ {
			for y := window.yMin; y < window.yMax; y++ {
				xR, yR := getImageSpace(x, ip.width), getImageSpace(y, ip.height)
				var pixelColor color.Color
				if len(window.triangles) > 0 && ip.interpolateN > 1 {

					samples := make([]color.Color, ip.interpolateN)
					for i := range ip.interpolateN {
						dx, dy := getPixelWiggle(ip.width), getPixelWiggle(ip.height)
						samples[i] = wScene.GetColor(xR+rand.Float64()*dx, yR+rand.Float64()*dy)
					}
					pixelColor = color.Average(samples)
				} else {
					pixelColor = wScene.GetColor(xR, yR)
				}

				// insert pixels with flipped y- coord, so y would be -1 at the bottom, +1 at the top of the image
				img.Set(x, ip.height-y-1, pixelColor)
				pixelCount += 1
				if pixelCount%ip.height == 0 {
					r.lineChannel <- 1
				}
			}
		}
	}
	return img
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
			img.Set(x, ip.height-y-1, pixelColor)
		}
		r.lineChannel <- 1
	}
	return img
}

type RasterLine struct {
	A RasterPixel
	B RasterPixel
}

type RasterPixel struct {
	X int
	Y int
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// adapted from https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func (r Renderer) renderLine(im *image.RGBA, line RasterLine, gradient color.Gradient) {
	x0, y0, x1, y1 := line.A.X, line.A.Y, line.B.X, line.B.Y
	dx := abs(x1 - x0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	dy := -abs(y1 - y0)
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	error := dx + dy

	xprogress := float64(0)
	for {
		var c color.Color
		if dx == 0 {
			c = gradient.Interpolate(0.0)
		} else {
			c = gradient.Interpolate(xprogress / float64(dx))
		}

		im.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * error
		if e2 >= dy {
			if x0 == x1 {
				break
			}
			error = error + dy
			x0 += sx
			xprogress += 1
		}
		if e2 <= dx {
			if y0 == y1 {
				break
			}
			error = error + dx
			y0 += sy
		}
	}
}

func toImageDimension(d float64, pixelCount int) *int {
	if d < -1.0 || d > 1.0 {
		return nil
	}
	v := int((d/2 + 0.5) * float64(pixelCount))
	return &v
}

func toImagePixel(p geometry.Pixel, width, height int) *RasterPixel {
	x := toImageDimension(p.X, width)
	y := toImageDimension(p.Y, height)
	if x == nil || y == nil {
		return nil
	}
	return &RasterPixel{
		X: *x,
		Y: *y,
	}
}

func (r Renderer) getWireframeImage(scene scenes.Frame, ip ImagePreset) image.Image {
	img := image.NewRGBA(
		image.Rect(
			0, 0, ip.width, ip.height,
		),
	)
	for _, tri := range scene.Flatten() {
		for _, line := range tri.GetWireframe() {
			sceneA, aDepth := line.A.ToPixel()
			sceneB, bDepth := line.B.ToPixel()
			if sceneA == nil || sceneB == nil {
				fmt.Printf("Skipping line %s since it may be behind the screen", line)
				continue
			}
			pixA := toImagePixel(*sceneA, ip.width, ip.height)
			pixB := toImagePixel(*sceneB, ip.width, ip.height)
			if pixA == nil || pixB == nil {
				fmt.Printf("Skipping line %s since one or both pixels are outside of screen", line)
				continue
			}
			rasterLine := RasterLine{
				*pixA,
				*pixB,
			}

			greenBlack := color.SimpleGradient{
				color.Green,
				color.Black,
			}
			ratio := 8.0
			colorA := greenBlack.Interpolate(2*maths.Sigmoid(aDepth/ratio) - 1)
			colorB := greenBlack.Interpolate(2*maths.Sigmoid(bDepth/ratio) - 1)
			r.renderLine(img, rasterLine, color.SimpleGradient{colorA, colorB})
		}
		bbox := tri.GetBoundingBox()
		pixA := toImagePixel(bbox.TopLeft, ip.width, ip.height)
		pixB := toImagePixel(bbox.BottomRight, ip.width, ip.height)
		if pixA == nil || pixB == nil {
			continue
		}
		r.renderLine(img, RasterLine{*pixA, RasterPixel{pixA.X, pixB.Y}}, color.SimpleGradient{color.Red, color.Red})
		r.renderLine(img, RasterLine{*pixA, RasterPixel{pixB.X, pixA.Y}}, color.SimpleGradient{color.Red, color.Red})
		r.renderLine(img, RasterLine{*pixB, RasterPixel{pixA.X, pixB.Y}}, color.SimpleGradient{color.Red, color.Red})
		r.renderLine(img, RasterLine{*pixB, RasterPixel{pixB.X, pixA.Y}}, color.SimpleGradient{color.Red, color.Red})
	}
	return img
}

func (r Renderer) getTriangleDepthImage(scene scenes.Frame, ip ImagePreset) image.Image {
	img := image.NewRGBA(
		image.Rect(
			0, 0, ip.width, ip.height,
		),
	)
	// set to black bakcground
	pixelColor := color.Black
	for x := 0; x < ip.width; x++ {
		for y := 0; y < ip.height; y++ {

			// insert pixels with flipped y- coord, so y would be -1 at the bottom, +1 at the top of the image
			img.Set(x, ip.height-y-1, pixelColor)
		}
		r.lineChannel <- 1
	}
	windows := subdivideSceneIntoWindows(scene, ip)
	gradient := color.SimpleGradient{color.Black, color.Red}
	for _, window := range windows {
		for x := window.xMin; x < window.xMax; x++ {
			for y := window.yMin; y < window.yMax; y++ {
				nTriangles := len(window.triangles)
				if x == window.xMin || y == window.yMin {
					pixelColor = color.Green
				} else {
					pixelColor = gradient.Interpolate(float64(nTriangles) / 20.0)
				}

				// insert pixels with flipped y- coord, so y would be -1 at the bottom, +1 at the top of the image
				img.Set(x, ip.height-y-1, pixelColor)
			}
		}
	}
	return img
}
