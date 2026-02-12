package memory

import (
	"context"
	"datatracing/internal/domain"
	"testing"
	"time"
)

func TestTraceStore_QueryTraces_FilterAndLimit(t *testing.T) {
	store := NewTraceStore()
	now := time.Now().UTC()
	_ = store.SaveSpans(context.Background(), []domain.Span{
		{TraceID: "t1", SpanID: "s1", Operation: "normalize", Status: domain.SpanStatusOK, StartTime: now.Add(-2 * time.Minute), EndTime: now.Add(-time.Minute)},
		{TraceID: "t2", SpanID: "s2", Operation: "normalize", Status: domain.SpanStatusError, StartTime: now.Add(-10 * time.Second), EndTime: now},
	})

	filtered, err := store.QueryTraces(context.Background(), domain.QueryFilter{Operation: "normalize", Status: domain.SpanStatusError, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 {
		t.Fatalf("expected 1, got %d", len(filtered))
	}
	if filtered[0].TraceID != "t2" {
		t.Fatalf("expected t2 first, got %+v", filtered[0])
	}
}
