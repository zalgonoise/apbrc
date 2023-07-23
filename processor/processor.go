package processor

import (
	"strings"

	"github.com/zalgonoise/logx"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor/modifiers"
)

const (
	topLevelFolder = "APB Reloaded"
)

// Modifier is a type that applies changes to configuration files
type Modifier interface {
	// Apply modifies the configuration file on `basePath` path, returning an error if raised
	Apply(basePath string) error
}

// Processor applies a (set of) Modifier to a (set of) configuration file(s), under a fixed base path
type Processor struct {
	basePath string

	cfg       *config.Config
	modifiers []Modifier
	logger    logx.Logger
}

// New creates a Processor from the input config.Config, and logx.Logger
func New(cfg *config.Config, logger logx.Logger) *Processor {
	if cfg == nil {
		return nil
	}

	mods := make([]Modifier, 0)

	if cfg.FrameRate != nil {
		mods = append(mods,
			modifiers.NewFPSModifier(cfg.FrameRate.MinRate, cfg.FrameRate.MaxRate, cfg.FrameRate.SmoothedRate, logger),
		)
	}

	if cfg.Input != nil {
		mods = append(mods,
			modifiers.NewSprintLockModifier(cfg.Input.SprintLock, logger),
			modifiers.NewCrouchLockModifier(cfg.Input.CrouchHold, logger),
		)
	}

	return &Processor{
		cfg:       cfg,
		modifiers: mods,
		logger:    logger,
	}
}

// Run executes the processor, applying all configured Modifier. It returns an error if raised, on an invalid base path
// or on the first Modifier-returned error
func (p *Processor) Run() error {
	if p.basePath == "" {
		dir, ok := topLevel(p.cfg.Path)
		if !ok {
			return ErrInvalidPath
		}

		p.basePath = dir
	}

	for i := range p.modifiers {
		if err := p.modifiers[i].Apply(p.basePath); err != nil {
			return err
		}
	}

	p.logger.Info("processor executed successfully")

	return nil
}

func topLevel(dir string) (string, bool) {
	sep := "/"
	elems := strings.Split(dir, sep)

	if len(elems) == 1 {
		sep = "\\"
		elems = strings.Split(dir, sep)
	}

	for i := range elems {
		if elems[i] == topLevelFolder {
			return strings.Join(elems[:i+1], sep), true
		}
	}

	return dir, false
}
