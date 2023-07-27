package processor

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/zalgonoise/apbrc/monitoring"
)

type Runner interface {
	Run(ctx context.Context) error
}

type processorWithLogs struct {
	r      Runner
	logger monitoring.Logger
}

func (p processorWithLogs) Run(ctx context.Context) error {
	p.logger.InfoContext(ctx, "executing processor")

	if err := p.r.Run(ctx); err != nil {
		p.logger.ErrorContext(ctx, "processor execution failed",
			slog.String("error", err.Error()),
		)

		return err
	}

	p.logger.InfoContext(ctx, "processor executed successfully")

	return nil
}

func ProcessorWithLogs(r Runner, logger monitoring.Logger) Runner {
	return processorWithLogs{r, logger}
}
