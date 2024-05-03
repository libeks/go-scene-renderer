package renderer

import (
	"cmp"
	"fmt"
	"math"
	"slices"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/objects"
	"github.com/libeks/go-scene-renderer/scenes"
)

// Window specifies a renderable area along with the triangles within in
type Window struct {
	// coordinates are in image pixel space
	xMin       int                         // inclusive
	xMax       int                         // non-inclusive
	yMin       int                         // inclusive
	yMax       int                         // non-inclusive
	triangles  []objects.StaticBasicObject // list of triangles whose bounding box intersects the window
	background scenes.Background
}

// GetColor returns the color at the pixel, as well as the number of triangles, and comparisons before a match was made
func (w Window) GetColor(x, y float64) (colors.Color, int, int) {
	minZ := math.MaxFloat64
	checks := 0
	var closestColor *colors.Color
	for _, tri := range w.triangles {
		if tri.GetBoundingBox().MinZDepth > minZ && closestColor != nil {
			// this triangle is behind the one we've found already, break early
			break
		}
		c, zDepth := tri.GetColorDepth(x, y)
		checks += 1
		if c != nil && zDepth < minZ {
			minZ = zDepth
			closestColor = c
		}
	}
	if closestColor != nil {
		return *closestColor, len(w.triangles), checks
	}
	return w.background.GetColor(x, y), len(w.triangles), checks
}

func (w Window) Width() int {
	return w.xMax - w.xMin
}

func (w Window) Height() int {
	return w.yMax - w.yMin
}

func (w Window) String() string {
	return fmt.Sprintf("Window from (%d, %d) to (%d, %d) with objects %v", w.xMin, w.yMin, w.xMax, w.yMax, w.triangles)
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
	tlW := []objects.StaticBasicObject{}
	trW := []objects.StaticBasicObject{}
	blW := []objects.StaticBasicObject{}
	brW := []objects.StaticBasicObject{}
	for _, tri := range w.triangles {
		bbox := tri.GetBoundingBox()
		// fmt.Printf("bbox inside %s\n", bbox)
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
	retWins := []Window{
		{w.xMin, xMid, w.yMin, yMid, tlW, w.background},
		{xMid, w.xMax, w.yMin, yMid, trW, w.background},
		{w.xMin, xMid, yMid, w.yMax, blW, w.background},
		{xMid, w.xMax, yMid, w.yMax, brW, w.background},
	}
	return retWins
}

func initiateWindow(scene scenes.StaticScene, ip ImagePreset) []Window {
	triangles, background := scene.Flatten()
	i := 0
	for _, tri := range triangles {
		bbox := tri.GetBoundingBox()
		if bbox.TopLeft.X > 1 || bbox.TopLeft.Y > 1 {
			continue
		}
		if bbox.BottomRight.X < -1 || bbox.BottomRight.Y < -1 {
			continue
		}
		triangles[i] = tri
		i += 1
	}
	triangles = triangles[:i]
	// put the closest triangles in the front
	slices.SortFunc(triangles, func(a, b objects.StaticBasicObject) int {
		return cmp.Compare(a.GetBoundingBox().MinZDepth, b.GetBoundingBox().MinZDepth)
	})

	return []Window{
		{0, ip.width, 0, ip.height, triangles, background},
	}
}

func subdivideSceneIntoWindows(scene scenes.StaticScene, ip ImagePreset) []Window {
	// start with one window for the whole image. Assume that all objects fall within the image
	windows := initiateWindow(scene, ip)
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
	return finalWindows
}
