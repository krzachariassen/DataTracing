package application

import (
	"context"
	"datatracing/internal/domain"
)

type TraceStore interface {
	SaveSpans(ctx context.Context, spans []domain.Span) error
	GetTrace(ctx context.Context, traceID string) ([]domain.Span, error)
	QueryTraces(ctx context.Context, filter domain.QueryFilter) ([]domain.TraceSummary, error)
}
