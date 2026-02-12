package main

import (
	"datatracing/internal/application"
	"datatracing/internal/infrastructure/httpserver"
	"datatracing/internal/infrastructure/memory"
	"log"
	"net/http"
)

func main() {
	store := memory.NewTraceStore()
	query := application.NewQueryService(store)
	log.Println("query-api listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", httpserver.QueryHandler(query)))
}
