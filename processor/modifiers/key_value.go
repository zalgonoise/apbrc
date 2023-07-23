package modifiers

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/zalgonoise/attr"
	"github.com/zalgonoise/logx"
)

// KeyValue is a generic type that holds a key-value pair, used in configuration attributes
type KeyValue[T comparable] struct {
	Key    string
	Value  T
	Format string
}

// KeyValueModifier is a data structure that will apply KeyValue-type modifications to a configuration file
type KeyValueModifier[T comparable] struct {
	filePath string

	modifiers []KeyValue[T]
	logger    logx.Logger
}

// NewKeyValueModifier creates a KeyValueModifier of type T, configured with `filePath` as a base path, a logx.Logger,
// and any number of KeyValue modifiers
func NewKeyValueModifier[T comparable](
	filePath string, logger logx.Logger, modifiers ...KeyValue[T],
) KeyValueModifier[T] {
	return KeyValueModifier[T]{
		filePath: filePath,

		modifiers: modifiers,
		logger:    logger,
	}
}

// Apply modifies the configuration file on `basePath` path (directly under the KeyValueModifier's
// configured base path), returning an error if raised
func (m KeyValueModifier[T]) Apply(basePath string) error {
	if len(m.modifiers) == 0 {
		return nil
	}

	path := basePath + m.filePath

	m.logger.Info("applying modifiers to config file",
		attr.String("path", path),
		attr.Int("num_modifiers", len(m.modifiers)),
	)

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

		for i := range m.modifiers {
			if bytes.Contains(line, []byte(m.modifiers[i].Key)) {
				m.logMatch(line, m.modifiers[i])

				if _, err = fmt.Fprintf(writer,
					m.modifiers[i].Format, m.modifiers[i].Key, m.modifiers[i].Value,
				); err != nil {
					return err
				}

				// remove modifier once it has been consumed
				m.modifiers = append(m.modifiers[:i], m.modifiers[i+1:]...)

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

	m.logger.Info("overwritten configuration file successfully",
		attr.String("path", path),
	)

	return nil
}

func (m KeyValueModifier[T]) logMatch(line []byte, kv KeyValue[T]) {
	m.logger.Info("matched config key",
		attr.String("key", kv.Key),
		attr.New("new_value", kv.Value),
		attr.String("original_value", string(line)),
	)
}
