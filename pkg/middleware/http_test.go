package middleware

import (
	"context"
	"datatracing/pkg/sdk"
	"datatracing/pkg/tracing"
	"net/http"
	"net/http/httptest"
	"testing"
)

type captureExporter struct{ spans []tracing.Span }

func (e *captureExporter) Export(_ context.Context, s tracing.Span) error {
	e.spans = append(e.spans, s)
	return nil
}

func TestTraceHTTP(t *testing.T) {
	exp := &captureExporter{}
	tracer := sdk.NewTracer(exp, alwaysSampleSampler{})
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h := TraceHTTP(tracer, next)
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("trace-id", "trace-a")
	req.Header.Set("span-id", "parent")
	req.Header.Set("sampled", "1")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("unexpected code %d", rr.Code)
	}
	if len(exp.spans) != 1 || exp.spans[0].ParentID != "parent" {
		t.Fatalf("span not captured correctly: %+v", exp.spans)
	}
}

type alwaysSampleSampler struct{}

func (alwaysSampleSampler) ShouldSample(map[string]string, float64) bool { return true }
