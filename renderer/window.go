package renderer

import (
	"math"

	"github.com/libeks/go-scene-renderer/colors"
	"github.com/libeks/go-scene-renderer/objects"
	"github.com/libeks/go-scene-renderer/scenes"
)

// Window specifies a renderable area along with the triangles within in
type Window struct {
	// coordinates are in image pixel space
	xMin       int                       // inclusive
	xMax       int                       // non-inclusive
	yMin       int                       // inclusive
	yMax       int                       // non-inclusive
	triangles  []*objects.StaticTriangle // list of triangles whose bounding box intersects the window
	background scenes.Background
}

func (w Window) GetColor(x, y float64) colors.Color {
	minZ := math.MaxFloat64
	var closestColor *colors.Color
	for _, tri := range w.triangles {
		c, depth := tri.GetColorDepth(x, y)
		if c != nil && depth < minZ {
			minZ = depth
			closestColor = c
		}
	}
	if closestColor != nil {
		// todo, do something with minZ
		return *closestColor
	}
	return w.background.GetColor(x, y)
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
	tlW := []*objects.StaticTriangle{}
	trW := []*objects.StaticTriangle{}
	blW := []*objects.StaticTriangle{}
	brW := []*objects.StaticTriangle{}
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
	return []Window{
		{0, ip.width, 0, ip.height, triangles[:i], background},
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
	// fmt.Printf(
	// 	"Image has %d windows, max # of triangles is %d, total work (pixels * triangles) is %d\n",
	// 	len(finalWindows), maxTriangles, totalWork,
	// )
	return finalWindows
}
