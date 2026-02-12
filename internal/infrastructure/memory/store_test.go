package memory

import (
	"context"
	"datatracing/internal/domain"
	"sync"
	"testing"
	"time"
)

func TestTraceStore_Concurrency(t *testing.T) {
	store := NewTraceStore()
	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = store.SaveSpans(context.Background(), []domain.Span{{TraceID: "t1", SpanID: time.Now().String(), Operation: "op", StartTime: time.Now(), EndTime: time.Now()}})
		}(i)
	}
	wg.Wait()
	spans, err := store.GetTrace(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if len(spans) != 50 {
		t.Fatalf("got %d", len(spans))
	}
}
