package config

import (
	"flag"
	"os"
)

type Config struct {
	Path string

	FrameRate *FrameRateConfig
}

type FrameRateConfig struct {
	MaxRate      int
	MinRate      int
	SmoothedRate int
}

func NewConfig() (*Config, error) {
	path := flag.String("dir", "", "path to the game's installation folder")

	// frame rate options
	min := flag.Int("min", 60, "min frame rate value to set")
	max := flag.Int("max", 300, "max frame rate value to set")
	smoothed := flag.Int("smoothed", 300, "smoothed frame rate value to set")

	flag.Parse()

	// use working directory if path is unset
	if *path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		*path = currentPath
	}

	return &Config{
		Path: *path,
		FrameRate: &FrameRateConfig{
			MinRate:      *min,
			MaxRate:      *max,
			SmoothedRate: *smoothed,
		},
	}, nil
}
