package modifiers

import (
	"bytes"
	"context"
	"fmt"
)

// KeyValue is a generic type that holds a key-value pair, used in configuration attributes
type KeyValue[T comparable] struct {
	Key    string
	Data   T
	Format string
}

func (m KeyValue[T]) Match(_ context.Context, line []byte) (ok bool, match string) {
	return bytes.Contains(line, []byte(m.Key)), m.Key
}

func (m KeyValue[T]) Value(_ context.Context, _ string) (data []byte, value any) {
	return []byte(fmt.Sprintf(m.Format, m.Key, m.Data)), m.Data
}
