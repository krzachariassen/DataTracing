package clickhouse

import (
	"context"
	"datatracing/internal/domain"
	"testing"
)

func TestTraceStore_ScaffoldMethods(t *testing.T) {
	store := NewTraceStore(nil)
	if err := store.SaveSpans(context.Background(), []domain.Span{{TraceID: "t", SpanID: "s"}}); err != nil {
		t.Fatal(err)
	}
	spans, err := store.GetTrace(context.Background(), "t")
	if err != nil {
		t.Fatal(err)
	}
	if spans != nil {
		t.Fatalf("expected nil spans in scaffold, got %+v", spans)
	}
	res, err := store.QueryTraces(context.Background(), domain.QueryFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if res != nil {
		t.Fatalf("expected nil summaries in scaffold, got %+v", res)
	}
}
