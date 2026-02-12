package clickhouse

import (
	"context"
	"database/sql"
	"datatracing/internal/domain"
)

type TraceStore struct {
	db *sql.DB
}

func NewTraceStore(db *sql.DB) *TraceStore { return &TraceStore{db: db} }

func (s *TraceStore) SaveSpans(ctx context.Context, spans []domain.Span) error {
	_ = ctx
	_ = spans
	// Intentionally minimal for v2 bootstrapping; use schema in schema.sql.
	return nil
}

func (s *TraceStore) GetTrace(ctx context.Context, traceID string) ([]domain.Span, error) {
	_ = ctx
	_ = traceID
	return nil, nil
}

func (s *TraceStore) QueryTraces(ctx context.Context, filter domain.QueryFilter) ([]domain.TraceSummary, error) {
	_ = ctx
	_ = filter
	return nil, nil
}
