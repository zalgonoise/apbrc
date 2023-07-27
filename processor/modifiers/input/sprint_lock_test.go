package input

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zalgonoise/apbrc/processor/modifiers"
)

func TestNewSprintLockModifier(t *testing.T) {
	lf := "\n"
	if runtime.GOOS == "windows" {
		lf = "\r\n"
	}

	for _, testcase := range []struct {
		name         string
		alwaysSprint bool

		wants modifiers.Attribute
	}{
		{
			name:         "Success/LockSprint",
			alwaysSprint: true,
			wants: modifiers.KeyValue[string]{
				Key:    sprintLockKey,
				Data:   alwaysSprintValue,
				Format: bindingsModifierFormat + lf,
			},
		},
		{
			name:         "Success/Defaults",
			alwaysSprint: false,
			wants: modifiers.KeyValue[string]{
				Key:    sprintLockKey,
				Data:   defaultSprintValue,
				Format: bindingsModifierFormat + lf,
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mod := SprintLock(testcase.alwaysSprint)

			require.Equal(t, testcase.wants, mod)
		})
	}
}
