package input

import (
	"runtime"
	"strings"

	"github.com/zalgonoise/apbrc/processor/modifiers"
)

const (
	sprintLockKey          = `+Bindings=(Name="Sprint"`
	defaultSprintValue     = "InputSprinting | OnRelease InputStopSprinting"
	alwaysSprintValue      = "InputStopSprinting | OnRelease InputSprinting"
	bindingsModifierFormat = "%s,Command=%q)"
)

// SprintLock creates a new input sprint Modifier, which will change the
// sprinting behavior.
//
// Setting `alwaysSprint` to true inverts the behavior of the Sprint key, so you hold it to stop sprinting, and sets to
// always sprint when released. Setting `alwaysSprint` to false applies the default configuration.
func SprintLock(alwaysSprint bool) modifiers.Attribute {
	sb := &strings.Builder{}
	sb.WriteString(bindingsModifierFormat)

	if runtime.GOOS == "windows" {
		sb.WriteByte('\r')
	}
	sb.WriteByte('\n')

	attr := modifiers.KeyValue[string]{
		Key:    sprintLockKey,
		Data:   defaultSprintValue,
		Format: sb.String(),
	}

	if alwaysSprint {
		attr.Data = alwaysSprintValue
	}

	return attr
}
