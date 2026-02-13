# DataTracing v2 — Current Implementation Status

## What is implemented

### 1) Canonical tracing model and domain mapping
- Canonical public model exists in `pkg/tracing` (`Trace`, `Span`, `SpanLink`, status, events, attributes).
- Internal domain model aliases the public model to keep application logic decoupled from transport/storage specifics.

### 2) Clean architecture structure and boundaries
- Project is organized into `cmd`, `internal` (application/domain/infrastructure), and public `pkg` modules.
- Dependency guardrail test exists to ensure public `pkg/*` packages do not import `internal/*`.

### 3) Application services
- **Collector service** implemented with:
  - buffered ingestion channel,
  - worker pool,
  - periodic or size-based batching,
  - tail-sampling policy integration,
  - persistence through `TraceStore` abstraction.
- **Query service** implemented with:
  - trace retrieval,
  - DAG reconstruction (`BuildDAG`),
  - trace summary search passthrough.

### 4) Storage adapters
- **Local JSONL file store** (`internal/infrastructure/local`) implemented for save/get/query behavior.
- **In-memory store** (`internal/infrastructure/memory`) implemented for tests and lightweight runtime usage.
- **Store provider** selects backend via environment configuration.
- **ClickHouse adapter scaffold** and schema file exist, but CRUD methods are placeholders.

### 5) SDK and propagation
- SDK tracer exists with start/end lifecycle and exporter integration.
- HTTP exporter exists for collector ingestion.
- Propagation helpers exist for carriers.
- HTTP middleware and instrumentation wrappers exist for Kafka and Cadence.

### 6) Runtime entrypoints
- `cmd/collector` for ingestion service startup.
- `cmd/query-api` for query API startup.

### 7) Tests
- Unit/integration-style tests exist across:
  - application services,
  - infrastructure stores,
  - SDK,
  - propagation,
  - middleware,
  - instrumentation,
  - architecture dependency rules.

## What is missing / incomplete

1. **ClickHouse production readiness**
   - `SaveSpans`, `GetTrace`, and `QueryTraces` are currently stubs in the ClickHouse store.
   - Missing robust integration tests against a real ClickHouse instance.

2. **Collector ingestion formats and resilience**
   - Current ingestion path is minimal; protobuf ingestion path is not implemented.
   - More explicit retry/error handling and observability around dropped spans can be improved.

3. **Sampling sophistication**
   - Tail-sampling exists but policy expressiveness is still basic.
   - Advanced policy conditions (error-rate, latency threshold windows, operation filters with combinations) are not fully implemented.

4. **Operational hardening**
   - Security, authn/authz, quotas, and production SLO instrumentation need expansion.
   - Documentation for deployment topologies and scaling guidance is limited.

5. **End-to-end validation depth**
   - Full E2E tests that exercise SDK -> collector -> persistence -> query path in realistic environments should be expanded.

## Readiness snapshot

- **Architecture/readability**: strong bootstrap foundation.
- **Core workflow**: functional in local + memory modes.
- **Production readiness**: partial; key gap is persistent backend completion (ClickHouse) and operational hardening.
