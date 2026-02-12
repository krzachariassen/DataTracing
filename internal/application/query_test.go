package application

import (
	"datatracing/internal/domain"
	"testing"
	"time"
)

func TestBuildDAG(t *testing.T) {
	now := time.Now().UTC()
	spans := []domain.Span{
		{TraceID: "t1", SpanID: "root", Operation: "root", StartTime: now},
		{TraceID: "t1", SpanID: "child", ParentID: "root", Operation: "child", StartTime: now.Add(time.Millisecond)},
		{TraceID: "t1", SpanID: "orphan", ParentID: "missing", Operation: "orphan", StartTime: now.Add(2 * time.Millisecond)},
	}
	dag := BuildDAG(spans)
	if len(dag.Roots) != 1 || dag.Roots[0].Span.SpanID != "root" {
		t.Fatalf("unexpected roots: %+v", dag.Roots)
	}
	if len(dag.Roots[0].Children) != 1 {
		t.Fatalf("expected child attached")
	}
	if len(dag.Orphans) != 1 || dag.Orphans[0].Span.SpanID != "orphan" {
		t.Fatalf("expected orphan")
	}
}
