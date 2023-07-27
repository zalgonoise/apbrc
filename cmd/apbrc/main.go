package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/monitoring"
	"github.com/zalgonoise/apbrc/processor"
	"github.com/zalgonoise/apbrc/processor/modifiers"
	"github.com/zalgonoise/apbrc/processor/modifiers/engine"
	"github.com/zalgonoise/apbrc/processor/modifiers/input"
)

func main() {
	err, code := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v", err)
	}

	os.Exit(code)
}

func run() (err error, code int) {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		return err, 1
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	}))

	mods := initMods(cfg, logger)

	proc := processor.ProcessorWithLogs(processor.New(cfg, mods...), logger)

	if err = proc.Run(ctx); err != nil {
		return err, 1
	}

	return nil, 0
}

func initMods(cfg *config.Config, logger monitoring.Logger) []processor.Applier {
	mods := make([]processor.Applier, 0, 2)

	if cfg.FrameRate != nil {
		mods = append(mods, modifiers.ModifierWithLogs(
			engine.FrameRate(*cfg.FrameRate),
			logger,
		))
	}

	if cfg.Input != nil {
		mods = append(mods, modifiers.ModifierWithLogs(
			input.Input(*cfg.Input),
			logger,
		))
	}

	return mods
}
