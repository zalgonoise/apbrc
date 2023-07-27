package input

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zalgonoise/apbrc/processor/modifiers"
)

func TestNewCrouchLockModifier(t *testing.T) {
	lf := "\n"
	if runtime.GOOS == "windows" {
		lf = "\r\n"
	}

	for _, testcase := range []struct {
		name       string
		holdCrouch bool

		wants modifiers.Attribute
	}{
		{
			name:       "Success/HoldCrouch",
			holdCrouch: true,
			wants: modifiers.KeyValue[string]{
				Key:    crouchLockKey,
				Data:   holdCrouchValue,
				Format: bindingsModifierFormat + lf,
			},
		},
		{
			name:       "Success/Defaults",
			holdCrouch: false,
			wants: modifiers.KeyValue[string]{
				Key:    crouchLockKey,
				Data:   defaultCrouchValue,
				Format: bindingsModifierFormat + lf,
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mod := CrouchLock(testcase.holdCrouch)

			require.Equal(t, testcase.wants, mod)
		})
	}
}
