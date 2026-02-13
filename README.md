# DataTracing v2

DataTracing v2 is a generic "Jaeger for data pipelines" implemented in Go.

## Included in this bootstrap

- Clean Architecture layout (`cmd`, `internal`, `pkg`)
- Canonical tracing model in `pkg/tracing` (`Trace`, `Span`, `SpanLink`)
- File-backed shared `TraceStore` by default (`./data/traces.jsonl`) with optional in-memory mode via `DATATRACING_STORE=memory`
- Query DAG reconstruction
- Go SDK (`Start`, `End`, attributes, events)
- W3C-like propagation helpers for HTTP/Kafka/workflow carriers
- Collector HTTP ingestion with worker pool and tail sampling
- Query API (`/trace/{trace_id}`, `/search`)
- ClickHouse schema and adapter scaffold
- Kafka/Cadence integration wrappers

## Run

```bash
go test ./...
go run ./cmd/collector
go run ./cmd/query-api
# both binaries use DATATRACING_STORE_PATH (default: ./data/traces.jsonl)
```

## Next steps

- Complete ClickHouse adapter CRUD and integration tests with containerized ClickHouse
- Add protobuf ingestion path in collector
- Add richer tail-sampling policy expressions


## Architecture guardrails

- Public `pkg/*` packages are intentionally decoupled from `internal/*` implementation packages.
- `internal/architecture/dependency_test.go` enforces this rule in CI.
