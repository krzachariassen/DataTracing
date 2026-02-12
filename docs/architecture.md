# Architecture

## Layers

- **Public API / SDK (`pkg/*`)**: stable packages consumed by applications (`sdk`, `propagation`, `middleware`, instrumentation, and shared tracing model).
- **Domain (`internal/domain`)**: aliases the canonical model and keeps internal code decoupled from transport/storage.
- **Application (`internal/application`)**: use-cases (query, collector, sampling) and `TraceStore` port.
- **Infrastructure (`internal/infrastructure`)**: adapters (memory, clickhouse, http handlers).
- **Entrypoints (`cmd/*`)**: wiring and process startup only.

## Dependency Rules (Clean Architecture)

1. `cmd` can depend on `internal/application` and `internal/infrastructure`.
2. `internal/infrastructure` can depend on `internal/application` and `internal/domain`.
3. `internal/application` can depend only on `internal/domain`.
4. `pkg/*` must **not** import `internal/*` packages.
5. Domain/business types exposed publicly live in `pkg/tracing`; `internal/domain` re-exports aliases for internal use.

These rules make refactors safer: infra can be replaced without touching business logic, and SDK packages remain usable without internal coupling.

## Data Flow

1. SDK creates spans and applies head sampling.
2. SDK exports spans to collector over HTTP.
3. Collector validates and ingests spans into a buffered channel.
4. Worker pool batches spans and applies tail sampling.
5. Sampled traces are persisted through the `TraceStore` port.
6. Query API loads spans by trace and reconstructs DAG in application layer.

## Enforcement

`internal/architecture/dependency_test.go` validates that `pkg/*` does not import `internal/*`, preventing accidental boundary regressions.
