package checker

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	patchURL       = "http://apb.patch.gamersfirst.com"
	notesURL       = "http://apb.notes.gamersfirst.com"
	crashReportURL = "http://apb.crashreport.gamersfirst.com"
)

func CheckAll(ctx context.Context, logger *slog.Logger,
	client *http.Client, timeout time.Duration, failFast bool,
) bool {
	var ok = true

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	patchReq, err := http.NewRequestWithContext(ctx, http.MethodGet, patchURL, http.NoBody)
	if err != nil {
		return false
	}

	if err := httpDo(client, patchReq); err != nil {
		logger.ErrorContext(ctx, "failed to ping patch URL", slog.String("error", err.Error()))

		if failFast {
			return false
		}

		ok = false
	}

	notesReq, err := http.NewRequestWithContext(ctx, http.MethodGet, notesURL, http.NoBody)
	if err != nil {
		return false
	}

	if err := httpDo(client, notesReq); err != nil {
		logger.ErrorContext(ctx, "failed to ping notes URL", slog.String("error", err.Error()))

		if failFast {
			return false
		}

		ok = false
	}

	crashReportReq, err := http.NewRequestWithContext(ctx, http.MethodGet, crashReportURL, http.NoBody)
	if err != nil {
		return false
	}

	if err := httpDo(client, crashReportReq); err != nil {
		logger.ErrorContext(ctx, "failed to ping crash report URL", slog.String("error", err.Error()))

		if failFast {
			return false
		}

		ok = false
	}

	return ok
}

func httpDo(client *http.Client, req *http.Request) error {
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > 399 {
		return fmt.Errorf("failing status code: %s", res.Status)
	}

	return nil
}
