package main

import (
	"fmt"
	go_color "image/color"
	"time"

	"github.com/libeks/go-scene-renderer/color"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

func computePalette(colors []color.Color) ([]color.Color, error) {
	start := time.Now()
	// set up a random two-dimensional data set (float64 values between 0.0 and 1.0)
	var d clusters.Observations
	for _, c := range colors {
		d = append(d, clusters.Coordinates{
			c.R,
			c.G,
			c.B,
		})
	}

	// Partition the data points into 16 clusters
	// km := kmeans.New()
	// km, err := kmeans.NewWithOptions(0.05, nil) // optimized, runs faster
	km, err := kmeans.NewWithOptions(0.1, nil) // optimized, runs faster
	if err != nil {
		return nil, err
	}
	clusters, err := km.Partition(d, 256)
	if err != nil {
		return nil, err
	}

	palette := make([]color.Color, len(clusters))
	for i, c := range clusters {
		palette[i] = color.Color{
			c.Center[0], c.Center[1], c.Center[2],
		}
		// fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		// fmt.Printf("Matching data points: %+v\n\n", c.Observations)
	}
	fmt.Printf("Palette computation took %s, returned %d clusters\n", time.Since(start), len(palette))
	return palette, nil
}

func convertPalette(colors []color.Color) []go_color.Color {
	out := make([]go_color.Color, len(colors))
	for i, c := range colors {
		out[i] = c
	}
	return out
}
