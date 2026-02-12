package memory

import (
	"context"
	"datatracing/internal/domain"
	"sort"
	"sync"
)

type TraceStore struct {
	mu    sync.RWMutex
	spans map[string][]domain.Span
}

func NewTraceStore() *TraceStore {
	return &TraceStore{spans: make(map[string][]domain.Span)}
}

func (s *TraceStore) SaveSpans(_ context.Context, spans []domain.Span) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, span := range spans {
		s.spans[span.TraceID] = append(s.spans[span.TraceID], span)
	}
	return nil
}

func (s *TraceStore) GetTrace(_ context.Context, traceID string) ([]domain.Span, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := append([]domain.Span(nil), s.spans[traceID]...)
	sort.Slice(out, func(i, j int) bool { return out[i].StartTime.Before(out[j].StartTime) })
	return out, nil
}

func (s *TraceStore) QueryTraces(_ context.Context, filter domain.QueryFilter) ([]domain.TraceSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.TraceSummary, 0)
	for traceID, spans := range s.spans {
		if len(spans) == 0 {
			continue
		}
		matched := false
		start, end := spans[0].StartTime, spans[0].EndTime
		status := string(domain.SpanStatusOK)
		for _, span := range spans {
			if span.StartTime.Before(start) {
				start = span.StartTime
			}
			if span.EndTime.After(end) {
				end = span.EndTime
			}
			if span.Status == domain.SpanStatusError {
				status = string(domain.SpanStatusError)
			}
			if filter.Operation != "" && span.Operation != filter.Operation {
				continue
			}
			if filter.Status != "" && span.Status != filter.Status {
				continue
			}
			if !filter.From.IsZero() && span.StartTime.Before(filter.From) {
				continue
			}
			if !filter.To.IsZero() && span.StartTime.After(filter.To) {
				continue
			}
			matched = true
		}
		if matched || (filter.Operation == "" && filter.Status == "" && filter.From.IsZero() && filter.To.IsZero()) {
			out = append(out, domain.TraceSummary{TraceID: traceID, StartTime: start, EndTime: end, SpanCount: len(spans), Status: status})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartTime.After(out[j].StartTime) })
	if filter.Limit > 0 && len(out) > filter.Limit {
		out = out[:filter.Limit]
	}
	return out, nil
}
