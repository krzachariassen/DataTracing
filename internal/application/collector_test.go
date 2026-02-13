package application

import (
	"context"
	"datatracing/internal/domain"
	"sync"
	"testing"
	"time"
)

type spyStore struct {
	mu    sync.Mutex
	spans []domain.Span
}

func (s *spyStore) SaveSpans(_ context.Context, spans []domain.Span) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.spans = append(s.spans, spans...)
	return nil
}
func (s *spyStore) GetTrace(context.Context, string) ([]domain.Span, error) { return nil, nil }
func (s *spyStore) QueryTraces(context.Context, domain.QueryFilter) ([]domain.TraceSummary, error) {
	return nil, nil
}

func TestCollector_IngestAndSample(t *testing.T) {
	store := &spyStore{}
	collector := NewCollectorService(store, TailSamplingPolicy{ErrorAlways: true}, 2, 10*time.Millisecond, 1, 10)
	defer collector.Close()

	ok := collector.Ingest(domain.Span{TraceID: "t1", SpanID: "s1", Operation: "op", Status: domain.SpanStatusOK})
	if !ok {
		t.Fatal("ingest should succeed")
	}
	collector.Ingest(domain.Span{TraceID: "t1", SpanID: "s2", Operation: "op", Status: domain.SpanStatusError})
	time.Sleep(30 * time.Millisecond)
	store.mu.Lock()
	defer store.mu.Unlock()
	if len(store.spans) != 2 {
		t.Fatalf("expected stored spans, got %d", len(store.spans))
	}
}

func TestCollector_DoesNotDropEarlyTraceChunks(t *testing.T) {
	store := &spyStore{}
	collector := NewCollectorService(store, TailSamplingPolicy{ErrorAlways: true}, 1, time.Hour, 2, 10)

	if !collector.Ingest(domain.Span{TraceID: "t1", SpanID: "s1", Operation: "op", Status: domain.SpanStatusOK}) {
		t.Fatal("ingest should succeed")
	}
	if !collector.Ingest(domain.Span{TraceID: "t1", SpanID: "s2", Operation: "op", Status: domain.SpanStatusError}) {
		t.Fatal("ingest should succeed")
	}

	collector.Close()

	store.mu.Lock()
	defer store.mu.Unlock()
	if len(store.spans) != 2 {
		t.Fatalf("expected full trace to be persisted once sampled, got %d spans", len(store.spans))
	}
}
