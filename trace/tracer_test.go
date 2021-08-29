package trace

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tr := New(&bytes.Buffer{})
	rawTracer, ok := tr.(*tracer)
	if !ok {
		t.Fatalf("%T.(*tracer) error", rawTracer)
	}

	tr = New(nil)
	rawNilTracer, ok := tr.(*nilTracer)
	if !ok {
		t.Fatalf("%T.(*nilTracer) error", rawNilTracer)
	}
}

func TestTrace(t *testing.T) {
	b := &bytes.Buffer{}
	tracer := New(b)

	input1 := "こんにちは trace パッケージ"
	input2 := 123
	expected := fmt.Sprintf("%v%v\n%v\n", input1, input2, input2)

	if err := tracer.Trace(input1, input2); err != nil {
		t.Fatalf("want Trace(%v, %v) = nil, got %v", input1, input2, err)
	}
	if err := tracer.Trace(input2); err != nil {
		t.Fatalf("want Trace(%v) = nil, got %v", input2, err)
	}
	if s := b.String(); s != expected {
		t.Fatalf("Trace wrote wrong string %q, expected %q", s, expected)
	}
}
