package input

import (
	"runtime"
	"strings"

	"github.com/zalgonoise/apbrc/processor/modifiers"
)

const (
	crouchLockKey      = `+Bindings=(Name="Duck"`
	defaultCrouchValue = "Button m_bDuckButton | InputToggleDuck"
	holdCrouchValue    = "Button m_bDuckButton | InputToggleDuck | OnRelease InputToggleDuck"
)

// CrouchLock creates a new input sprint Modifier, which will change the
// crouching behavior.
//
// Setting `holdCrouch` to true inverts the behavior of the Crouch key, so it no longer acts like an on/off toggle,
// allowing a press-and-hold type of crouch action. Setting `holdCrouch` to false applies the default configuration.
func CrouchLock(holdCrouch bool) modifiers.Attribute {
	sb := &strings.Builder{}
	sb.WriteString(bindingsModifierFormat)

	if runtime.GOOS == "windows" {
		sb.WriteByte('\r')
	}
	sb.WriteByte('\n')

	attr := modifiers.KeyValue[string]{
		Key:    crouchLockKey,
		Data:   defaultCrouchValue,
		Format: sb.String(),
	}

	if holdCrouch {
		attr.Data = holdCrouchValue
	}

	return attr
}
