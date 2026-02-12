package httpserver

import (
	"datatracing/internal/application"
	"datatracing/internal/domain"
	"encoding/json"
	"net/http"
)

func CollectorHandler(collector *application.CollectorService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		var span domain.Span
		if err := json.NewDecoder(r.Body).Decode(&span); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if span.TraceID == "" || span.SpanID == "" || span.Operation == "" {
			http.Error(w, "missing required fields", http.StatusBadRequest)
			return
		}
		if ok := collector.Ingest(span); !ok {
			http.Error(w, "collector overloaded", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
}
