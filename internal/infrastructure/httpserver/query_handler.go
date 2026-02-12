package httpserver

import (
	"datatracing/internal/application"
	"datatracing/internal/domain"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func QueryHandler(query *application.QueryService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/trace/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/trace/")
		if id == "" {
			http.Error(w, "trace id required", http.StatusBadRequest)
			return
		}
		dag, err := query.GetTraceDAG(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(dag)
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		filter := domain.QueryFilter{Operation: q.Get("operation"), Status: domain.SpanStatus(q.Get("status")), Limit: 100}
		if from := q.Get("from"); from != "" {
			if v, err := time.Parse(time.RFC3339, from); err == nil {
				filter.From = v
			}
		}
		if to := q.Get("to"); to != "" {
			if v, err := time.Parse(time.RFC3339, to); err == nil {
				filter.To = v
			}
		}
		result, err := query.SearchTraces(r.Context(), filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(result)
	})
	return mux
}
