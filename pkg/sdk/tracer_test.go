package sdk

import (
	"context"
	"datatracing/pkg/tracing"
	"sync"
	"testing"
)

type captureExporter struct {
	mu    sync.Mutex
	spans []tracing.Span
}

func (e *captureExporter) Export(_ context.Context, span tracing.Span) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.spans = append(e.spans, span)
	return nil
}

func TestTracer_StartEnd(t *testing.T) {
	exp := &captureExporter{}
	tracer := NewTracer(exp, alwaysSampleSampler{})
	ctx, span := tracer.Start(context.Background(), "op", tracing.SpanKindTransform)
	span.SetAttribute("entity_id", "1")
	span.AddEvent("cache_miss")
	span.SetStatus(tracing.SpanStatusError)
	span.End()
	span.End()
	if _, ok := ctx.Value(struct{}{}).(string); ok { // noop to keep ctx used
		t.Fatal("impossible")
	}
	if len(exp.spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(exp.spans))
	}
	if exp.spans[0].EndTime.IsZero() {
		t.Fatal("end time must be set")
	}
	if exp.spans[0].Status != tracing.SpanStatusError {
		t.Fatalf("unexpected status: %s", exp.spans[0].Status)
	}
}

type alwaysSampleSampler struct{}

func (alwaysSampleSampler) ShouldSample(map[string]string, float64) bool { return true }
