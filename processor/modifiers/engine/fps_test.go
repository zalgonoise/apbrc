package engine

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
)

func TestNewFPSModifier(t *testing.T) {
	lf := "\n"
	if runtime.GOOS == "windows" {
		lf = "\r\n"
	}

	for _, testcase := range []struct {
		name     string
		min      int
		max      int
		smoothed int

		wants modifiers.Modifier
	}{
		{
			name:     "Success/Simple",
			min:      60,
			max:      300,
			smoothed: 300,
			wants: modifiers.Modifier{
				FilePath: fpsModifierPath,
				Attributes: []modifiers.Attribute{
					modifiers.KeyValue[int]{
						Key:    minRateKey,
						Data:   60,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    maxRateKey,
						Data:   300,
						Format: fpsModifierFormat + lf,
					},
					modifiers.KeyValue[int]{
						Key:    smoothedRateKey,
						Data:   300,
						Format: fpsModifierFormat + lf,
					},
				},
			},
		},
		{
			name: "Success/OnlyMax",
			max:  300,
			wants: modifiers.Modifier{
				FilePath: fpsModifierPath,
				Attributes: []modifiers.Attribute{
					modifiers.KeyValue[int]{
						Key:    maxRateKey,
						Data:   300,
						Format: fpsModifierFormat + lf,
					},
				},
			},
		},
		{
			name: "Success/None",
			wants: modifiers.Modifier{
				FilePath:   fpsModifierPath,
				Attributes: []modifiers.Attribute{},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mod := FrameRate(config.FrameRateConfig{
				MinRate:      testcase.min,
				MaxRate:      testcase.max,
				SmoothedRate: testcase.smoothed,
			})
			require.Equal(t, testcase.wants, mod)
		})
	}
}
