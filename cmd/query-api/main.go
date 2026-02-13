package main

import (
	"datatracing/internal/application"
	"datatracing/internal/infrastructure/httpserver"
	"datatracing/internal/infrastructure/store"
	"log"
	"net/http"
)

func main() {
	traceStore := store.NewSharedTraceStore()
	query := application.NewQueryService(traceStore)
	log.Println("query-api listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", httpserver.QueryHandler(query)))
}
