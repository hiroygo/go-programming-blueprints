package trace

import (
	"fmt"
	"io"
)

// Tracer はコード内での出来事を記録する
type Tracer interface {
	Trace(...interface{}) error
}

func New(w io.Writer) Tracer {
	// w が nil なら何もしない Tracer を返す
	if w == nil {
		return &nilTracer{}
	}
	return &tracer{w: w}
}

type nilTracer struct{}

func (t *nilTracer) Trace(...interface{}) error {
	return nil
}

type tracer struct {
	w io.Writer
}

func (t *tracer) Trace(a ...interface{}) error {
	if _, err := io.WriteString(t.w, fmt.Sprint(a...)); err != nil {
		return fmt.Errorf("WriteString error: %w", err)
	}
	if _, err := io.WriteString(t.w, "\n"); err != nil {
		return fmt.Errorf("WriteString error: %w", err)
	}
	return nil
}
