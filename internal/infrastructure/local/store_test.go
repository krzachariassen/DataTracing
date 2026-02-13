package local

import (
	"context"
	"datatracing/internal/domain"
	"path/filepath"
	"testing"
	"time"
)

func TestTraceStore_PersistsAndQueries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "traces.jsonl")
	store := NewTraceStore(path)
	now := time.Now()
	spans := []domain.Span{
		{TraceID: "t1", SpanID: "s1", Operation: "op", StartTime: now, EndTime: now, Status: domain.SpanStatusOK},
		{TraceID: "t1", SpanID: "s2", Operation: "op", StartTime: now.Add(time.Millisecond), EndTime: now.Add(time.Millisecond), Status: domain.SpanStatusError},
	}
	if err := store.SaveSpans(context.Background(), spans); err != nil {
		t.Fatalf("save spans: %v", err)
	}

	loaded, err := store.GetTrace(context.Background(), "t1")
	if err != nil {
		t.Fatalf("get trace: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(loaded))
	}

	summaries, err := store.QueryTraces(context.Background(), domain.QueryFilter{Status: domain.SpanStatusError})
	if err != nil {
		t.Fatalf("query traces: %v", err)
	}
	if len(summaries) != 1 || summaries[0].TraceID != "t1" {
		t.Fatalf("expected trace t1 in query result, got %#v", summaries)
	}
}
