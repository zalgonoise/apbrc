package config

import (
	"flag"
	"os"
)

type Config struct {
	Path string

	FrameRate *FrameRateConfig
	Input     *InputConfig
}

type FrameRateConfig struct {
	MaxRate      int
	MinRate      int
	SmoothedRate int
}

type InputConfig struct {
	SprintLock bool
	CrouchHold bool
	Reset      bool
}

func newInputConfig(lockSprint, holdCrouch, reset bool) *InputConfig {
	if reset {
		return &InputConfig{}
	}

	if !lockSprint && !holdCrouch {
		return nil
	}

	return &InputConfig{
		SprintLock: lockSprint,
		CrouchHold: holdCrouch,
	}
}

func NewConfig() (*Config, error) {
	path := flag.String("dir", "", "path to the game's installation folder")

	// frame rate options
	min := flag.Int("min", 60, "min frame rate value to set")
	max := flag.Int("max", 300, "max frame rate value to set")
	smoothed := flag.Int("smoothed", 300, "smoothed frame rate value to set")

	// input options
	lockSprint := flag.Bool("lock-sprint", false,
		"changes the sprint input configuration to always sprint (lock sprint action)",
	)
	holdCrouch := flag.Bool("hold-crouch", false,
		"changes the crouch input configuration to act as a press-and-hold key (crouch unlock action)",
	)
	resetInput := flag.Bool("reset-input", false,
		"returns any input bindings changes back to their default configuration",
	)

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
		Input: newInputConfig(*lockSprint, *holdCrouch, *resetInput),
	}, nil
}
