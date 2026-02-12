package sdk

import (
	"context"
	"datatracing/pkg/tracing"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPExporter_Export(t *testing.T) {
	received := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = true
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exporter := HTTPExporter{URL: srv.URL, Client: srv.Client()}
	if err := exporter.Export(context.Background(), tracing.Span{TraceID: "t", SpanID: "s", Operation: "op"}); err != nil {
		t.Fatal(err)
	}
	if !received {
		t.Fatal("expected request")
	}
}
