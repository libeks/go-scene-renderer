package renderer

import (
	"os"
	"path/filepath"
)

func cleanUpTempFiles(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, f := range files {
		// fmt.Printf("About to remove %s\n", f)
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func createSubdirectories(outFileFormat string) error {
	return os.MkdirAll(filepath.Dir(outFileFormat), os.ModePerm)
}

// convert coordinate from pixel space (0, pixels-1) to image space (-1.0, 1.0)
func getImageSpace(x, pixels int) float64 {
	return 2*float64(x)/float64(pixels) - 1.0
}

// get the width/height of a pixel in image space, in one dimension
func getPixelWiggle(pixels int) float64 {
	return 2.0 / float64(pixels)
}
