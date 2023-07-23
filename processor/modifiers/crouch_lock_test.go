package modifiers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/texth"
)

func TestNewCrouchLockModifier(t *testing.T) {
	logger := logx.New(texth.New(os.Stdout))

	for _, testcase := range []struct {
		name       string
		holdCrouch bool

		wants KeyValueModifier[string]
	}{
		{
			name:       "Success/HoldCrouch",
			holdCrouch: true,
			wants: KeyValueModifier[string]{
				filePath: crouchLockModifierPath,
				logger:   logger,
				modifiers: []KeyValue[string]{
					{
						Key:    crouchLockKey,
						Value:  holdCrouchValue,
						Format: bindingsModifierFormat,
					},
				},
			},
		},
		{
			name:       "Success/Defaults",
			holdCrouch: false,
			wants: KeyValueModifier[string]{
				filePath: crouchLockModifierPath,
				logger:   logger,
				modifiers: []KeyValue[string]{
					{
						Key:    crouchLockKey,
						Value:  defaultCrouchValue,
						Format: bindingsModifierFormat,
					},
				},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mod := NewCrouchLockModifier(testcase.holdCrouch, logger)

			require.Equal(t, testcase.wants, mod)
		})
	}
}
