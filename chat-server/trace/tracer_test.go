package trace_test

import (
	"bytes"
	"github.com/yvdinesh/go-blueprints/chat-server/trace"
	"testing"
)

func TestNewTracer(t *testing.T) {
	buf := bytes.Buffer{}
	tracer := trace.NewTracer(&buf)
	tracer.Trace("test")
	if buf.String() != "test\n" {
		t.Errorf("Tracer should write test, but wrote %v\n", buf.String())
	}
}
