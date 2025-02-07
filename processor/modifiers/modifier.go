package modifiers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Attribute interface {
	Match(ctx context.Context, line []byte) (ok bool, match string)
	Value(ctx context.Context, key string) (data []byte, value any)
}

// Modifier is a data structure that will apply KeyValue-type modifications to a configuration file
type Modifier struct {
	logger *slog.Logger

	FilePath   string
	Attributes []Attribute
}

// New creates a Modifier, configured with `filePath` as a base path,
// and any number of Attribute modifiers
func New(filePath string, logger *slog.Logger, attributes ...Attribute) Modifier {
	return Modifier{
		logger:     logger,
		FilePath:   filePath,
		Attributes: attributes,
	}
}

// Apply modifies the configuration file on `basePath` path (directly under the Modifier's
// configured base path), returning an error if raised
func (m Modifier) Apply(ctx context.Context, basePath string) error {
	if len(m.Attributes) == 0 {
		m.logger.InfoContext(ctx,
			"no modifiers to apply",
			slog.Int("num_attributes", len(m.Attributes)),
		)

		return nil
	}

	path := basePath + m.FilePath

	m.logger.InfoContext(ctx,
		"reading config file to apply modifiers",
		slog.String("path", path),
		slog.Int("num_attributes", len(m.Attributes)),
	)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	writer := bytes.NewBuffer(nil)
	counter := -1

scanLoop:
	for {
		counter++

		var line []byte
		line, err = reader.ReadBytes('\n')

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		for i := range m.Attributes {
			if ok, key := m.Attributes[i].Match(ctx, line); ok {
				m.logger.DebugContext(ctx, "attribute matched configuration line",
					slog.Int("line_number", counter),
					slog.String("data", string(line)),
				)

				data, _ := m.Attributes[i].Value(ctx, key)

				m.logger.DebugContext(ctx, "applying value",
					slog.String("data", string(line)),
					slog.String("final_value", string(data)),
				)

				if _, err = writer.Write(data); err != nil {
					return err
				}

				// remove modifier once it has been consumed
				m.Attributes = append(m.Attributes[:i], m.Attributes[i+1:]...)

				continue scanLoop
			}
		}

		if _, err = writer.Write(line); err != nil {
			return err
		}
	}

	m.logger.InfoContext(ctx,
		"applying modifiers to config file",
		slog.String("path", path),
	)

	if err = f.Close(); err != nil {
		return err
	}

	f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("%w: %s", err, path)
	}

	if _, err = f.Write(writer.Bytes()); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	m.logger.InfoContext(ctx,
		"overwritten configuration file successfully",
		slog.String("path", path),
	)

	return nil
}
