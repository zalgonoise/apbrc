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
	Cap int
	Min int
	Max int
}

type InputConfig struct {
	SprintLock bool
	CrouchHold bool
	Reset      bool
}

func NewConfig() (*Config, error) {
	path := flag.String("dir", "", "path to the game's installation folder")

	// frame rate options
	frameRateCap := flag.Int("frameRateCap", 0,
		"frame rate limit value to set when the Smoothed frame rate option is disabled")
	frameRateMin := flag.Int("frameRateMin", 22,
		"minimum frame rate value to set when the Smoothed frame rate option is enabled")
	frameRateMax := flag.Int("frameRateCap", 128,
		"maximum frame rate value to set when the Smoothed frame rate option is enabled")

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
			Min: *frameRateMin,
			Cap: *frameRateCap,
			Max: *frameRateMax,
		},
		Input: newInputConfig(*lockSprint, *holdCrouch, *resetInput),
	}, nil
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
