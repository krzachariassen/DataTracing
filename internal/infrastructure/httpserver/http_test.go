package httpserver

import (
	"bytes"
	"datatracing/internal/application"
	"datatracing/internal/domain"
	"datatracing/internal/infrastructure/memory"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCollectorHandler(t *testing.T) {
	store := memory.NewTraceStore()
	collector := application.NewCollectorService(store, application.TailSamplingPolicy{ErrorAlways: true}, 10, 5*time.Millisecond, 1, 10)
	defer collector.Close()
	h := CollectorHandler(collector)
	span := domain.Span{TraceID: "t", SpanID: "s", Operation: "op", Status: domain.SpanStatusError}
	b, _ := json.Marshal(span)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("got %d", rr.Code)
	}
}

func TestQueryHandler(t *testing.T) {
	store := memory.NewTraceStore()
	_ = store.SaveSpans(nil, []domain.Span{{TraceID: "t", SpanID: "s", Operation: "op", StartTime: time.Now(), EndTime: time.Now()}})
	query := application.NewQueryService(store)
	h := QueryHandler(query)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/trace/t", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d", rr.Code)
	}
}

func TestCollectorHandler_MethodNotAllowed(t *testing.T) {
	store := memory.NewTraceStore()
	collector := application.NewCollectorService(store, application.TailSamplingPolicy{ErrorAlways: true}, 10, 5*time.Millisecond, 1, 10)
	defer collector.Close()
	h := CollectorHandler(collector)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got %d", rr.Code)
	}
}

func TestQueryHandler_Search(t *testing.T) {
	store := memory.NewTraceStore()
	now := time.Now().UTC()
	_ = store.SaveSpans(nil, []domain.Span{{TraceID: "t", SpanID: "s", Operation: "op", Status: domain.SpanStatusError, StartTime: now, EndTime: now}})
	query := application.NewQueryService(store)
	h := QueryHandler(query)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/search?operation=op&status=ERROR", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("got %d", rr.Code)
	}
}
