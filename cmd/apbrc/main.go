package main

import (
	"fmt"
	"os"

	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/texth"

	"github.com/zalgonoise/apbrc/config"
	"github.com/zalgonoise/apbrc/processor"
)

func main() {
	err, code := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v", err)
	}

	os.Exit(code)
}

func run() (err error, code int) {
	cfg, err := config.NewConfig()
	if err != nil {
		return err, 1
	}

	logger := logx.New(texth.New(os.Stderr))

	proc := processor.New(cfg, logger)

	if err = proc.Run(); err != nil {
		return err, 1
	}

	return nil, 0
}
