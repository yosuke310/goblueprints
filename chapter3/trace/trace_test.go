package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("New returns nil.")
	} else {
		tracer.Trace("Hello, this is tracer.")
		if buf.String() != "Hello, this is tracer.\n" {
			t.Errorf("'%s' was incorrect outputs.", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("foo")
}
