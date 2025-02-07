package engine

import (
	"log/slog"
	"runtime"
	"strings"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
)

const (
	frameRateCapKey = "MaxClientFrameRate"
	frameRateMinKey = "MinSmoothedFrameRate"
	frameRateMaxKey = "MaxSmoothedFrameRate"

	defaultFrameRateCap = 100
	defaultFrameRateMin = 22
	defaultFrameRateMax = 128

	fpsModifierPath   = `/Engine/Config/BaseEngine.ini`
	fpsModifierFormat = "%s=%d"
)

// FrameRate creates a new frame-rate Modifier, which affects the engine's configuration for frame rate;
// updating its smoothed minimum and maximum frame rates, as well as the frame rate cap.
func FrameRate(cfg config.FrameRateConfig, logger *slog.Logger) modifiers.Modifier {
	sb := &strings.Builder{}
	sb.WriteString(fpsModifierFormat)

	if runtime.GOOS == "windows" {
		sb.WriteByte('\r')
	}
	sb.WriteByte('\n')

	mods := make([]modifiers.Attribute, 0, 3)

	if cfg.Cap >= 0 && cfg.Cap != defaultFrameRateCap {
		mods = append(mods, modifiers.KeyValue[int]{
			Key:    frameRateCapKey,
			Data:   cfg.Cap,
			Format: sb.String(),
		})
	}

	if cfg.Min >= 0 && cfg.Min != defaultFrameRateMin {
		mods = append(mods, modifiers.KeyValue[int]{
			Key:    frameRateMinKey,
			Data:   cfg.Min,
			Format: sb.String(),
		})
	}

	if cfg.Max >= 0 && cfg.Max != defaultFrameRateMax {
		mods = append(mods, modifiers.KeyValue[int]{
			Key:    frameRateMaxKey,
			Data:   cfg.Max,
			Format: sb.String(),
		})
	}

	return modifiers.New(fpsModifierPath, logger, mods...)
}
