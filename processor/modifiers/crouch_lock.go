package modifiers

import (
	"github.com/zalgonoise/logx"
)

const (
	crouchLockModifierPath = `\APBGame\Config\DefaultInput.ini`
	crouchLockKey          = `+Bindings=(Name="Duck"`
	defaultCrouchValue     = "Button m_bDuckButton | InputToggleDuck"
	holdCrouchValue        = "Button m_bDuckButton | InputToggleDuck | OnRelease InputToggleDuck"
)

// NewCrouchLockModifier creates a new input sprint KeyValueModifier, which will change the
// crouching behavior.
//
// Setting `holdCrouch` to true inverts the behavior of the Crouch key, so it no longer acts like an on/off toggle,
// allowing a press-and-hold type of crouch action. Setting `holdCrouch` to false applies the default configuration.
func NewCrouchLockModifier(holdCrouch bool, logger logx.Logger) KeyValueModifier[string] {
	mod := KeyValue[string]{
		Key:    crouchLockKey,
		Value:  defaultCrouchValue,
		Format: bindingsModifierFormat,
	}

	if holdCrouch {
		mod.Value = holdCrouchValue
	}

	return NewKeyValueModifier[string](crouchLockModifierPath, logger, mod)
}
