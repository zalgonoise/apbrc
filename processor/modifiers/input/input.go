package input

import (
	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
)

const inputBindingsModifierPath = `/APBGame/Config/DefaultInput.ini`

func Input(cfg config.InputConfig) modifiers.Modifier {
	mod := modifiers.Modifier{
		FilePath:   inputBindingsModifierPath,
		Attributes: make([]modifiers.Attribute, 0, 2),
	}

	if cfg.Reset {
		mod.Attributes = append(mod.Attributes,
			SprintLock(false),
			CrouchLock(false),
		)

		return mod
	}

	if cfg.SprintLock {
		mod.Attributes = append(mod.Attributes, SprintLock(cfg.SprintLock))
	}

	if cfg.CrouchHold {
		mod.Attributes = append(mod.Attributes, CrouchLock(cfg.CrouchHold))
	}

	return mod
}
