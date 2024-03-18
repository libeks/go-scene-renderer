package geometry

// Pixel is expressed in scene coordinates, so only pixels from x in (-1, 1) and y in (-1, 1) will be on screen
type Pixel struct {
	// expressed in
	X float64
	Y float64
}
