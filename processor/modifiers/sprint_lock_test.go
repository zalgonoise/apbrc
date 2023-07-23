package modifiers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/texth"
)

func TestNewSprintLockModifier(t *testing.T) {
	logger := logx.New(texth.New(os.Stdout))

	for _, testcase := range []struct {
		name         string
		alwaysSprint bool

		wants KeyValueModifier[string]
	}{
		{
			name:         "Success/LockSprint",
			alwaysSprint: true,
			wants: KeyValueModifier[string]{
				filePath: sprintLockModifierPath,
				logger:   logger,
				modifiers: []KeyValue[string]{
					{
						Key:    sprintLockKey,
						Value:  alwaysSprintValue,
						Format: bindingsModifierFormat,
					},
				},
			},
		},
		{
			name:         "Success/Defaults",
			alwaysSprint: false,
			wants: KeyValueModifier[string]{
				filePath: sprintLockModifierPath,
				logger:   logger,
				modifiers: []KeyValue[string]{
					{
						Key:    sprintLockKey,
						Value:  defaultSprintValue,
						Format: bindingsModifierFormat,
					},
				},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mod := NewSprintLockModifier(testcase.alwaysSprint, logger)

			require.Equal(t, testcase.wants, mod)
		})
	}
}
