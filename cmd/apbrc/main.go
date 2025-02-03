package main

import (
	"context"
	"flag"
	"github.com/zalgonoise/apbrc/checker"
	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/log"
	"github.com/zalgonoise/apbrc/processor"
	"github.com/zalgonoise/apbrc/processor/modifiers"
	"github.com/zalgonoise/apbrc/processor/modifiers/engine"
	"github.com/zalgonoise/apbrc/processor/modifiers/input"
	"github.com/zalgonoise/x/cli"
	"log/slog"
	"net/http"
	"time"
)

var modes = []string{"", "apply", "check"}

func main() {
	runner := cli.NewRunner("apbrc",
		cli.WithOneOf(modes...),
		cli.WithExecutors(map[string]cli.Executor{
			"":      cli.Executable(ExecApply),
			"apply": cli.Executable(ExecApply),
			"check": cli.Executable(ExecCheck),
		}),
	)

	cli.Run(runner)
}

func ExecApply(ctx context.Context, logger *slog.Logger, args []string) (int, error) {
	cfg, err := config.NewConfig(args)
	if err != nil {
		return 1, err
	}

	mods := initMods(cfg, logger)

	proc := processor.ProcessorWithLogs(processor.New(cfg, mods...), logger)

	if err = proc.Run(ctx); err != nil {
		return 1, err
	}

	return 0, nil
}

func ExecCheck(ctx context.Context, logger *slog.Logger, args []string) (int, error) {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)

	timeout := fs.Duration("timeout", 30*time.Second, "how long to wait before a request")
	freq := fs.Duration("freq", time.Minute, "how frequently should it ping APB servers until it works")
	count := fs.Int("count", 0, "number of pings to emit")
	keepPinging := fs.Bool("keep-pinging", false, "keep pinging even if servers are online")
	failFast := fs.Bool("fail-fast", false, "fail fast")

	if err := fs.Parse(args); err != nil {
		return 1, err
	}

	if *count > 0 {
		logger.InfoContext(ctx, "pinging", slog.Int("num_pings", *count))

		client := &http.Client{Timeout: 30 * time.Second}

		for range *count {
			checker.CheckAll(ctx, logger, client, *timeout, *failFast)
		}

		return 0, nil
	}

	for {
		logger.InfoContext(ctx, "pinging",
			slog.Bool("fail-fast", *failFast), slog.Duration("frequency", *freq))

		client := &http.Client{Timeout: 30 * time.Second}

		if ok := checker.CheckAll(ctx, logger, client, *timeout, *failFast); ok && !*keepPinging {
			return 0, nil
		}

		time.Sleep(*freq)
	}
}

func initMods(cfg *config.Config, logger log.Logger) []processor.Applier {
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
