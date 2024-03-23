package renderer

type RasterLine struct {
	A RasterPixel
	B RasterPixel
}

type RasterPixel struct {
	X int
	Y int
}
