package engine

import (
	"runtime"
	"strings"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
)

const (
	minRateKey      = "MinSmoothedFrameRate"
	smoothedRateKey = "MaxSmoothedFrameRate"
	maxRateKey      = "MaxClientFrameRate"

	defaultMinRate      = 22
	defaultMaxRate      = 128
	defaultSmoothedRate = 100

	fpsModifierPath   = `/Engine/Config/BaseEngine.ini`
	fpsModifierFormat = "%s=%d"
)

// FrameRate creates a new frame-rate Modifier, which will change the min, max and smoothed frame rate
// configuration values
func FrameRate(cfg config.FrameRateConfig) modifiers.Modifier {
	sb := &strings.Builder{}
	sb.WriteString(fpsModifierFormat)

	if runtime.GOOS == "windows" {
		sb.WriteByte('\r')
	}
	sb.WriteByte('\n')

	mods := make([]modifiers.Attribute, 0, 3)

	if cfg.MinRate > 0 && cfg.MinRate != defaultMinRate {
		mods = append(mods, modifiers.KeyValue[int]{
			Key:    minRateKey,
			Data:   cfg.MinRate,
			Format: sb.String(),
		})
	}

	if cfg.MaxRate > 0 && cfg.MaxRate != defaultMaxRate {
		mods = append(mods, modifiers.KeyValue[int]{
			Key:    maxRateKey,
			Data:   cfg.MaxRate,
			Format: sb.String(),
		})
	}

	if cfg.SmoothedRate > 0 && cfg.SmoothedRate != defaultSmoothedRate {
		mods = append(mods, modifiers.KeyValue[int]{
			Key:    smoothedRateKey,
			Data:   cfg.SmoothedRate,
			Format: sb.String(),
		})
	}

	return modifiers.NewModifier(fpsModifierPath, mods...)
}
