package store

import (
	"datatracing/internal/application"
	"datatracing/internal/infrastructure/local"
	"datatracing/internal/infrastructure/memory"
	"os"
)

const defaultStorePath = "./data/traces.jsonl"

func NewSharedTraceStore() application.TraceStore {
	if os.Getenv("DATATRACING_STORE") == "memory" {
		return memory.NewTraceStore()
	}
	path := os.Getenv("DATATRACING_STORE_PATH")
	if path == "" {
		path = defaultStorePath
	}
	return local.NewTraceStore(path)
}
