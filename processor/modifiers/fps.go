package modifiers

import "github.com/zalgonoise/logx"

const (
	minRateKey      = "MinSmoothedFrameRate"
	smoothedRateKey = "MaxSmoothedFrameRate"
	maxRateKey      = "MaxClientFrameRate"

	defaultMinRate      = 22
	defaultMaxRate      = 128
	defaultSmoothedRate = 100

	fpsModifierPath   = `\Engine\Config\BaseEngine.ini`
	fpsModifierFormat = "%s=%d\r\n"
)

// NewFPSModifier creates a new frame-rate KeyValueModifier, which will change the min, max and smoothed frame rate
// configuration values
func NewFPSModifier(min, max, smoothed int, logger logx.Logger) KeyValueModifier[int] {
	modifiers := make([]KeyValue[int], 0, 3)

	if min > 0 && min != defaultMinRate {
		modifiers = append(modifiers, KeyValue[int]{
			Key:    minRateKey,
			Value:  min,
			Format: fpsModifierFormat,
		})
	}

	if max > 0 && max != defaultMaxRate {
		modifiers = append(modifiers, KeyValue[int]{
			Key:    maxRateKey,
			Value:  max,
			Format: fpsModifierFormat,
		})
	}

	if smoothed > 0 && smoothed != defaultSmoothedRate {
		modifiers = append(modifiers, KeyValue[int]{
			Key:    smoothedRateKey,
			Value:  smoothed,
			Format: fpsModifierFormat,
		})
	}

	return NewKeyValueModifier[int](fpsModifierPath, logger, modifiers...)
}
