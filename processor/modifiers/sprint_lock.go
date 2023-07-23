package modifiers

import (
	"github.com/zalgonoise/logx"
)

const (
	sprintLockModifierPath = `\APBGame\Config\DefaultInput.ini`
	sprintLockKey          = `+Bindings=(Name="Sprint"`
	defaultSprintValue     = "InputSprinting | OnRelease InputStopSprinting"
	alwaysSprintValue      = "InputStopSprinting | OnRelease InputSprinting"
	bindingsModifierFormat = "%s,Command=\"%s\")\r\n"
)

// NewSprintLockModifier creates a new input sprint KeyValueModifier, which will change the
// sprinting behavior.
//
// Setting `alwaysSprint` to true inverts the behavior of the Sprint key, so you hold it to stop sprinting, and sets to
// always sprint when released. Setting `alwaysSprint` to false applies the default configuration.
func NewSprintLockModifier(alwaysSprint bool, logger logx.Logger) KeyValueModifier[string] {
	mod := KeyValue[string]{
		Key:    sprintLockKey,
		Value:  defaultSprintValue,
		Format: bindingsModifierFormat,
	}

	if alwaysSprint {
		mod.Value = alwaysSprintValue
	}

	return NewKeyValueModifier[string](sprintLockModifierPath, logger, mod)
}
