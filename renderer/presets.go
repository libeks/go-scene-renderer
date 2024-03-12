package renderer

type ImagePreset struct {
	width        int
	height       int
	interpolateN int
}

type VideoPreset struct {
	ImagePreset
	nFrameCount int
	frameRate   int
}

type Pixel struct {
	X int
	Y int
}

var (
	ImagePresetTest = ImagePreset{
		width:        200,
		height:       200,
		interpolateN: 1,
	}
	ImagePresetHiDef = ImagePreset{
		width:        1000,
		height:       1000,
		interpolateN: 16,
	}
	VideoPresetTest = VideoPreset{
		ImagePreset: ImagePresetTest,
		nFrameCount: 30,
		frameRate:   15,
	}
	VideoPresetHiDef = VideoPreset{
		ImagePreset: ImagePresetHiDef,
		nFrameCount: 400,
		frameRate:   30,
	}
)
