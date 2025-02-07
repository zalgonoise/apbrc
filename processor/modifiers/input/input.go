package input

import (
	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
	"log/slog"
)

const inputBindingsModifierPath = `/APBGame/Config/DefaultInput.ini`

func Input(cfg config.InputConfig, logger *slog.Logger) modifiers.Modifier {
	if cfg.Reset {
		return modifiers.New(inputBindingsModifierPath, logger,
			SprintLock(false),
			CrouchLock(false),
		)
	}

	attrs := make([]modifiers.Attribute, 0, 2)

	if cfg.SprintLock {
		attrs = append(attrs, SprintLock(cfg.SprintLock))
	}

	if cfg.CrouchHold {
		attrs = append(attrs, CrouchLock(cfg.CrouchHold))
	}

	return modifiers.New(inputBindingsModifierPath, logger, attrs...)
}
