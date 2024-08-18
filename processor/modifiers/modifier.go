package modifiers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

type Attribute interface {
	Match(ctx context.Context, line []byte) (ok bool, match string)
	Value(ctx context.Context, key string) (data []byte, value any)
}

// Modifier is a data structure that will apply KeyValue-type modifications to a configuration file
type Modifier struct {
	FilePath   string
	Attributes []Attribute
}

// NewModifier creates a Modifier of type T, configured with `filePath` as a base path,
// and any number of KeyValue modifiers
func NewModifier(filePath string, attributes ...Attribute) Modifier {
	return Modifier{
		FilePath:   filePath,
		Attributes: attributes,
	}
}

// Apply modifies the configuration file on `basePath` path (directly under the Modifier's
// configured base path), returning an error if raised
func (m Modifier) Apply(ctx context.Context, basePath string) error {
	if len(m.Attributes) == 0 {
		return nil
	}

	path := basePath + m.FilePath

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	writer := bytes.NewBuffer(nil)

scanLoop:
	for {
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
				data, _ := m.Attributes[i].Value(ctx, key)
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

	return nil
}
