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

	patchReq, err := http.NewRequestWithContext(ctx, http.MethodHead, patchURL, http.NoBody)
	if err != nil {
		return false
	}

	code, err := httpDo(client, patchReq)
	if err != nil {
		logger.ErrorContext(ctx, "failed to ping patch URL",
			slog.String("error", err.Error()), slog.Int("status_code", code))

		if failFast {
			return false
		}

		ok = false
	} else {
		logger.InfoContext(ctx, "pinged patch URL successfully", slog.Int("status_code", code))
	}

	notesReq, err := http.NewRequestWithContext(ctx, http.MethodHead, notesURL, http.NoBody)
	if err != nil {
		return false
	}

	code, err = httpDo(client, notesReq)
	if err != nil {
		logger.ErrorContext(ctx, "failed to ping notes URL",
			slog.String("error", err.Error()), slog.Int("status_code", code))

		if failFast {
			return false
		}

		ok = false
	} else {
		logger.InfoContext(ctx, "pinged notes URL successfully", slog.Int("status_code", code))
	}

	crashReportReq, err := http.NewRequestWithContext(ctx, http.MethodHead, crashReportURL, http.NoBody)
	if err != nil {
		return false
	}

	code, err = httpDo(client, crashReportReq)
	if err != nil {
		logger.ErrorContext(ctx, "failed to ping crash report URL",
			slog.String("error", err.Error()), slog.Int("status_code", code))

		if failFast {
			return false
		}

		ok = false
	} else {
		logger.InfoContext(ctx, "pinged crash report URL successfully", slog.Int("status_code", code))
	}

	return ok
}

func httpDo(client *http.Client, req *http.Request) (int, error) {
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.StatusCode > 499 {
		return res.StatusCode, fmt.Errorf("failing status code: %s", res.Status)
	}

	return res.StatusCode, nil
}
