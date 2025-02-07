package processor

import (
	"context"
	"fmt"
	"github.com/zalgonoise/x/errs"
	"log/slog"
	"strings"

	"github.com/zalgonoise/apbrc/config"
)

const (
	topLevelFolder = "APB Reloaded"
)

const (
	unlockerDomain = errs.Domain("apbrc/processor")

	ErrInvalid = errs.Kind("invalid")

	ErrPath = errs.Entity("path")
)

var (
	ErrInvalidPath = errs.WithDomain(unlockerDomain, ErrInvalid, ErrPath)
)

// Applier is a type that applies changes to configuration files
type Applier interface {
	// Apply modifies the configuration file on `basePath` path, returning an error if raised
	Apply(ctx context.Context, basePath string) error
}

// Processor applies a (set of) Applier to a (set of) configuration file(s), under a fixed base path
type Processor struct {
	basePath string

	cfg       *config.Config
	logger    *slog.Logger
	modifiers []Applier
}

// New creates a Processor from the input config.Config, slog.Logger and Applier(s)
func New(cfg *config.Config, logger *slog.Logger, mods ...Applier) *Processor {
	if cfg == nil || mods == nil {
		return nil
	}

	return &Processor{
		cfg:       cfg,
		logger:    logger,
		modifiers: mods,
	}
}

// Run executes the processor, applying all configured Applier. It returns an error if raised, on an invalid base path
// or on the first Applier-returned error
func (p *Processor) Run(ctx context.Context) error {
	p.logger.InfoContext(ctx, "executing processor")

	if p.basePath == "" {
		dir, ok := topLevel(p.cfg.Path)
		if !ok {
			return fmt.Errorf("%w: scanning top-level folder from: %s", ErrInvalidPath, p.cfg.Path)
		}

		p.basePath = dir
	}

	for i := range p.modifiers {
		if err := p.modifiers[i].Apply(ctx, p.basePath); err != nil {
			p.logger.ErrorContext(ctx, "processor execution failed",
				slog.String("error", err.Error()),
				// TODO: add modifier info with fmt.Stringer implementation
			)

			return err
		}
	}

	p.logger.InfoContext(ctx, "processor executed successfully")

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
