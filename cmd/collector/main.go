package main

import (
	"datatracing/internal/application"
	"datatracing/internal/infrastructure/httpserver"
	"datatracing/internal/infrastructure/memory"
	"log"
	"net/http"
	"time"
)

func main() {
	store := memory.NewTraceStore()
	collector := application.NewCollectorService(store, application.TailSamplingPolicy{ErrorAlways: true, LatencyThreshold: 500 * time.Millisecond}, 200, 300*time.Millisecond, 4, 5000)
	defer collector.Close()

	h := httpserver.CollectorHandler(collector)
	log.Println("collector listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
