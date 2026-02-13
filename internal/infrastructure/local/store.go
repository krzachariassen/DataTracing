package local

import (
	"bufio"
	"context"
	"datatracing/internal/domain"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type TraceStore struct {
	mu   sync.Mutex
	path string
}

func NewTraceStore(path string) *TraceStore {
	return &TraceStore{path: path}
}

func (s *TraceStore) SaveSpans(_ context.Context, spans []domain.Span) error {
	if len(spans) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, span := range spans {
		if err := enc.Encode(span); err != nil {
			return err
		}
	}
	return nil
}

func (s *TraceStore) GetTrace(_ context.Context, traceID string) ([]domain.Span, error) {
	spans, err := s.loadSpans()
	if err != nil {
		return nil, err
	}
	out := make([]domain.Span, 0)
	for _, span := range spans {
		if span.TraceID == traceID {
			out = append(out, span)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartTime.Before(out[j].StartTime) })
	return out, nil
}

func (s *TraceStore) QueryTraces(_ context.Context, filter domain.QueryFilter) ([]domain.TraceSummary, error) {
	spans, err := s.loadSpans()
	if err != nil {
		return nil, err
	}
	byTrace := make(map[string][]domain.Span)
	for _, span := range spans {
		byTrace[span.TraceID] = append(byTrace[span.TraceID], span)
	}

	out := make([]domain.TraceSummary, 0)
	for traceID, traceSpans := range byTrace {
		if len(traceSpans) == 0 {
			continue
		}
		matched := false
		start, end := traceSpans[0].StartTime, traceSpans[0].EndTime
		status := string(domain.SpanStatusOK)
		for _, span := range traceSpans {
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
			out = append(out, domain.TraceSummary{TraceID: traceID, StartTime: start, EndTime: end, SpanCount: len(traceSpans), Status: status})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartTime.After(out[j].StartTime) })
	if filter.Limit > 0 && len(out) > filter.Limit {
		out = out[:filter.Limit]
	}
	return out, nil
}

func (s *TraceStore) loadSpans() ([]domain.Span, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	spans := make([]domain.Span, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var span domain.Span
		if err := json.Unmarshal(line, &span); err != nil {
			return nil, err
		}
		spans = append(spans, span)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return spans, nil
}
