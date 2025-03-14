package config

import (
	"errors"
	"flag"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrNoFileConfig = errors.New("no file config")

type Config struct {
	Path string `yaml:"path"`

	FrameRate *FrameRateConfig `yaml:"frame_rate"`
	Input     *InputConfig     `json:"input"`
}

type FrameRateConfig struct {
	Cap int `yaml:"cap"`
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type InputConfig struct {
	SprintLock bool `yaml:"sprint_lock"`
	CrouchHold bool `yaml:"crouch_hold"`
	Reset      bool `yaml:"reset"`
}

func NewConfigFromFile() (*Config, error) {
	f, err := os.OpenFile("config.yaml", os.O_RDONLY, 0600)

	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoFileConfig
	}

	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return readConfigFile(buf)
}

func readConfigFile(buf []byte) (*Config, error) {
	config := &Config{}

	if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}

	if config.Path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		config.Path = currentPath
	}

	return config, nil
}

func NewConfig(args []string) (*Config, error) {
	config, err := NewConfigFromFile()

	if err == nil {
		return config, nil
	}

	if !errors.Is(err, ErrNoFileConfig) {
		return nil, err
	}

	fs := flag.NewFlagSet("apply", flag.ExitOnError)

	path := fs.String("dir", "", "path to the game's installation folder")

	// frame rate options
	frameRateCap := fs.Int("cap", 0,
		"frame rate limit value to set when the Smoothed frame rate option is disabled")
	frameRateMin := fs.Int("min", 22,
		"minimum frame rate value to set when the Smoothed frame rate option is enabled")
	frameRateMax := fs.Int("max", 128,
		"maximum frame rate value to set when the Smoothed frame rate option is enabled")

	// input options
	lockSprint := fs.Bool("lock-sprint", false,
		"changes the sprint input configuration to always sprint (lock sprint action)",
	)
	holdCrouch := fs.Bool("hold-crouch", false,
		"changes the crouch input configuration to act as a press-and-hold key (crouch unlock action)",
	)
	resetInput := fs.Bool("reset-input", false,
		"returns any input bindings changes back to their default configuration",
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

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
