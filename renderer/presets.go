package renderer

import (
	"fmt"
	"strconv"
	"strings"
)

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
	ImagePresetIntermediate = ImagePreset{
		width:        500,
		height:       500,
		interpolateN: 1,
	}
	ImagePresetHiDef = ImagePreset{
		width:        1000,
		height:       1000,
		interpolateN: 4,
	}
	VideoPresetTest = VideoPreset{
		ImagePreset: ImagePresetTest,
		nFrameCount: 30,
		frameRate:   15,
	}
	VideoPresetIntermediate = VideoPreset{
		ImagePreset: ImagePresetIntermediate,
		nFrameCount: 100,
		frameRate:   30,
	}
	VideoPresetHiDef = VideoPreset{
		ImagePreset: ImagePresetHiDef,
		nFrameCount: 400,
		frameRate:   30,
	}

	defaultImagePreset = ImagePresetTest
	defaultVideoPreset = VideoPresetTest
)

func ParseImagePreset(flagVal string) (ImagePreset, error) {
	if strings.Contains(flagVal, ",") {
		chunks := strings.Split(flagVal, ",")
		if len(chunks) != 3 {
			return ImagePreset{}, fmt.Errorf("expect two commas in image flag, got '%s'", flagVal)
		}
		intChunks := make([]int, len(chunks))
		for i, chunk := range chunks {
			val, err := strconv.Atoi(chunk)
			if err != nil {
				return ImagePreset{}, err
			}
			intChunks[i] = val
		}
		width, height, interpolate := intChunks[0], intChunks[1], intChunks[2]
		return ImagePreset{
			width:        width,
			height:       height,
			interpolateN: interpolate,
		}, nil
	}
	switch flagVal {
	case "default":
		return defaultImagePreset, nil
	case "test":
		return ImagePresetTest, nil
	case "hidef":
		return ImagePresetHiDef, nil
	default:
		return ImagePreset{}, fmt.Errorf("could not parse image format '%s'", flagVal)
	}
}

func ParseVideoPreset(flagVal string) (VideoPreset, error) {
	if strings.Contains(flagVal, ",") {
		chunks := strings.Split(flagVal, ",")
		if len(chunks) != 5 {
			return VideoPreset{}, fmt.Errorf("expect two commas in video flag, got '%s'", flagVal)
		}
		intChunks := make([]int, len(chunks))
		for i, chunk := range chunks {
			val, err := strconv.Atoi(chunk)
			if err != nil {
				return VideoPreset{}, err
			}
			intChunks[i] = val
		}
		width, height, interpolate, frames, frameRate := intChunks[0], intChunks[1], intChunks[2], intChunks[3], intChunks[4]
		return VideoPreset{
			ImagePreset: ImagePreset{
				width,
				height,
				interpolate,
			},
			nFrameCount: frames,
			frameRate:   frameRate,
		}, nil
	}
	switch flagVal {
	case "default":
		return defaultVideoPreset, nil
	case "test":
		return VideoPresetTest, nil
	case "intermediate":
		return VideoPresetIntermediate, nil
	case "hidef":
		return VideoPresetHiDef, nil
	default:
		return VideoPreset{}, fmt.Errorf("could not parse video format '%s'", flagVal)
	}
}
