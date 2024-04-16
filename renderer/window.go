package renderer

import (
	"cmp"
	"math"
	"slices"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/objects"
	"github.com/libeks/go-scene-renderer/scenes"
)

// Window specifies a renderable area along with the triangles within in
type Window struct {
	// coordinates are in image pixel space
	xMin       int                   // inclusive
	xMax       int                   // non-inclusive
	yMin       int                   // inclusive
	yMax       int                   // non-inclusive
	triangles  []objects.BasicObject // list of triangles whose bounding box intersects the window
	background scenes.Background
}

// GetColor returns the color at the pixel, as well as the number of triangles, and comparisons before a match was made
func (w Window) GetColor(x, y float64) (colors.Color, int, int) {
	minZ := math.MaxFloat64
	checks := 0
	var closestColor *colors.Color
	for _, tri := range w.triangles {
		if tri.GetBoundingBox().MinDepth > minZ && closestColor != nil {
			// this triangle is behind the one we've found already, break early
			break
		}
		c, depth := tri.GetColorDepth(x, y)
		checks += 1
		if c != nil && depth < minZ {
			minZ = depth
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
	tlW := []objects.BasicObject{}
	trW := []objects.BasicObject{}
	blW := []objects.BasicObject{}
	brW := []objects.BasicObject{}
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
		{w.xMin, xMid, w.yMin, yMid, tlW, w.background},
		{xMid, w.xMax, w.yMin, yMid, trW, w.background},
		{w.xMin, xMid, yMid, w.yMax, blW, w.background},
		{xMid, w.xMax, yMid, w.yMax, brW, w.background},
	}
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
	slices.SortFunc(triangles, func(a, b objects.BasicObject) int {
		return cmp.Compare(a.GetBoundingBox().MinDepth, b.GetBoundingBox().MinDepth)
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
