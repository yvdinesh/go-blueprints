package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(data interface{})
}

type tracer struct {
	out io.Writer
}

func NewTracer(w io.Writer) Tracer {
	return &tracer{out: w}
}

func (t *tracer) Trace(data interface{}) {
	fmt.Fprintln(t.out, data)
}
