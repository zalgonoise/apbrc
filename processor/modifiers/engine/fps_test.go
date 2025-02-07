package engine

import (
	"github.com/zalgonoise/apbrc/log"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
)

func TestNewFPSModifier(t *testing.T) {
	logger := log.New()
	lf := "\n"
	if runtime.GOOS == "windows" {
		lf = "\r\n"
	}

	for _, testcase := range []struct {
		name         string
		frameRateCap int
		frameRateMin int
		frameRateMax int

		wants modifiers.Modifier
	}{
		{
			name:         "Success/Simple",
			frameRateCap: 300,
			frameRateMin: 60,
			frameRateMax: 300,
			wants: modifiers.Modifier{
				FilePath: fpsModifierPath,
				Attributes: []modifiers.Attribute{
					modifiers.KeyValue[int]{
						Key:    frameRateCapKey,
						Data:   300,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    frameRateMinKey,
						Data:   60,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    frameRateMaxKey,
						Data:   300,
						Format: fpsModifierFormat + lf,
					},
				},
			},
		},
		{
			name:         "Success/OnlyCap",
			frameRateCap: 300,
			wants: modifiers.Modifier{
				FilePath: fpsModifierPath,
				Attributes: []modifiers.Attribute{
					modifiers.KeyValue[int]{
						Key:    frameRateCapKey,
						Data:   300,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    frameRateMinKey,
						Data:   0,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    frameRateMaxKey,
						Data:   0,
						Format: fpsModifierFormat + lf,
					},
				},
			},
		},
		{
			name: "Success/None",
			wants: modifiers.Modifier{
				FilePath: fpsModifierPath,
				Attributes: []modifiers.Attribute{
					modifiers.KeyValue[int]{
						Key:    frameRateCapKey,
						Data:   0,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    frameRateMinKey,
						Data:   0,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    frameRateMaxKey,
						Data:   0,
						Format: fpsModifierFormat + lf,
					},
				},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mod := FrameRate(config.FrameRateConfig{
				Cap: testcase.frameRateCap,
				Min: testcase.frameRateMin,
				Max: testcase.frameRateMax,
			}, logger)
			require.Equal(t, testcase.wants, mod)
		})
	}
}
