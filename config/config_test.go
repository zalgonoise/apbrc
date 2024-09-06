package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var conf = []byte(`
path: "/path/to/dir"
frame_rate:
  cap: 0
  min: 60
  max: 144
input:
  sprint_lock: true
  crouch_hold: true
  reset: false
`)

func TestNewConfigFromFile(t *testing.T) {
	for _, testcase := range []struct {
		name  string
		input []byte
		wants *Config
		err   error
	}{
		{
			name:  "simple",
			input: conf,
			wants: &Config{
				Path: "/path/to/dir",
				FrameRate: &FrameRateConfig{
					Cap: 0,
					Min: 60,
					Max: 144,
				},
				Input: &InputConfig{
					SprintLock: true,
					CrouchHold: true,
					Reset:      false,
				},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			c, err := readConfigFile(testcase.input)
			if err != nil {
				require.ErrorIs(t, err, testcase.err)

				return
			}

			require.Equal(t, testcase.wants, c)
		})
	}
}
