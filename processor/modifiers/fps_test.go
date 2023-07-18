package modifiers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/texth"
)

func TestNewFPSModifier(t *testing.T) {
	logger := logx.New(texth.New(os.Stdout))

	for _, testcase := range []struct {
		name     string
		min      int
		max      int
		smoothed int

		wants KeyValueModifier[int]
	}{
		{
			name:     "Success/Simple",
			min:      60,
			max:      300,
			smoothed: 300,
			wants: KeyValueModifier[int]{
				filePath: fpsModifierPath,
				logger:   logger,
				modifiers: []KeyValue[int]{
					{
						Key:    minRateKey,
						Value:  60,
						Format: fpsModifierFormat,
					},
					{
						Key:    maxRateKey,
						Value:  300,
						Format: fpsModifierFormat,
					},
					{
						Key:    smoothedRateKey,
						Value:  300,
						Format: fpsModifierFormat,
					},
				},
			},
		},
		{
			name: "Success/OnlyMax",
			max:  300,
			wants: KeyValueModifier[int]{
				filePath: fpsModifierPath,
				logger:   logger,
				modifiers: []KeyValue[int]{
					{
						Key:    maxRateKey,
						Value:  300,
						Format: fpsModifierFormat,
					},
				},
			},
		},
		{
			name: "Success/None",
			wants: KeyValueModifier[int]{
				filePath:  fpsModifierPath,
				logger:    logger,
				modifiers: []KeyValue[int]{},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mod := NewFPSModifier(testcase.min, testcase.max, testcase.smoothed, logger)

			require.Equal(t, testcase.wants, mod)
		})
	}
}
