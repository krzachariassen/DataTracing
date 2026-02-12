package application

import (
	"context"
	"datatracing/internal/domain"
	"testing"
	"time"
)

type queryStoreStub struct {
	spans    []domain.Span
	summary  []domain.TraceSummary
	traceErr error
	queryErr error
}

func (s queryStoreStub) SaveSpans(context.Context, []domain.Span) error { return nil }
func (s queryStoreStub) GetTrace(context.Context, string) ([]domain.Span, error) {
	return s.spans, s.traceErr
}
func (s queryStoreStub) QueryTraces(context.Context, domain.QueryFilter) ([]domain.TraceSummary, error) {
	return s.summary, s.queryErr
}

func TestQueryService_GetTraceDAGAndSearch(t *testing.T) {
	now := time.Now().UTC()
	store := queryStoreStub{
		spans:   []domain.Span{{TraceID: "t1", SpanID: "s1", StartTime: now}},
		summary: []domain.TraceSummary{{TraceID: "t1", SpanCount: 1}},
	}
	svc := NewQueryService(store)

	dag, err := svc.GetTraceDAG(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if dag.TraceID != "t1" || len(dag.Nodes) != 1 {
		t.Fatalf("unexpected dag: %+v", dag)
	}
	res, err := svc.SearchTraces(context.Background(), domain.QueryFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].TraceID != "t1" {
		t.Fatalf("unexpected result: %+v", res)
	}
}
