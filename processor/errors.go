package processor

import "github.com/zalgonoise/x/errs"

const (
	unlockerDomain = errs.Domain("apbrc/processor")

	ErrInvalid = errs.Kind("invalid")

	ErrPath = errs.Entity("path")
)

var (
	ErrInvalidPath = errs.New(unlockerDomain, ErrInvalid, ErrPath)
)
