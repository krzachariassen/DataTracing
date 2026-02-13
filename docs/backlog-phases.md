# DataTracing v2 — Incremental Backlog (Phased)

This backlog is intentionally split into small, sequential phases so work can move safely with TDD and frequent integration.

## Phase 0 — Baseline and quality gates

- [ ] Add/update contributor docs for local run/test commands.
- [ ] Ensure CI runs `go test ./...` and dependency guardrails.
- [ ] Add coverage reporting and minimum threshold (informational first).
- [ ] Add lint/static analysis step (`go vet`, optional golangci-lint).

## Phase 1 — Complete ClickHouse TraceStore MVP

- [ ] Implement `SaveSpans` with batch inserts.
- [ ] Implement `GetTrace` by `trace_id`, sorted by start time.
- [ ] Implement `QueryTraces` with operation/status/time filters and limit.
- [ ] Add integration tests using real ClickHouse container.
- [ ] Validate schema indexing/partition assumptions with realistic load.

## Phase 2 — Collector protocol and reliability enhancements

- [ ] Add protobuf ingestion endpoint alongside JSON path.
- [ ] Validate payload size limits and malformed input behavior.
- [ ] Add clear metrics/logging for accepted, dropped, and sampled spans.
- [ ] Introduce configurable backpressure behavior for saturated queues.
- [ ] Add tests for graceful shutdown and flush guarantees.

## Phase 3 — Sampling policy evolution

- [ ] Add composable policy DSL/config (AND/OR conditions).
- [ ] Implement policies for latency threshold, error presence, and operation match.
- [ ] Add deterministic policy tests with trace fixtures.
- [ ] Expose policy config through collector runtime settings.

## Phase 4 — Query API and user-facing retrieval improvements

- [ ] Expand search filters (service name, duration range, attribute matching).
- [ ] Add pagination/cursor semantics to search endpoint.
- [ ] Add endpoint contract tests for query handlers.
- [ ] Add validation and consistent error envelope for API responses.

## Phase 5 — SDK and instrumentation usability

- [ ] Add richer examples for HTTP/Kafka/Cadence instrumentation.
- [ ] Add resilience features in exporter (timeouts, retries, bounded queue).
- [ ] Improve propagation interoperability tests across carrier types.
- [ ] Add semantic conventions for common data pipeline operations.

## Phase 6 — End-to-end and performance confidence

- [ ] Add E2E scenario tests: SDK -> collector -> store -> query.
- [ ] Add benchmark suite for ingestion and query hotspots.
- [ ] Define target SLOs (ingest latency, query p95, drop rate).
- [ ] Tune batching and worker defaults from benchmark data.

## Phase 7 — Production hardening and operations

- [ ] Add authentication/authorization strategy for ingest/query APIs.
- [ ] Add multi-environment configuration strategy and secrets handling.
- [ ] Add deployment templates (container/k8s) and operational runbook.
- [ ] Add disaster-recovery and data-retention guidance.

## Phase 8 — Optional ecosystem extensions

- [ ] Add additional storage adapters if required by adoption.
- [ ] Add OTEL bridge/import compatibility where useful.
- [ ] Add UI or existing visualization integration for DAG exploration.

---

## Suggested immediate next 3 sprints

1. **Sprint A**: Phase 1 (ClickHouse MVP) + Phase 0 quality gates.
2. **Sprint B**: Phase 2 reliability + Phase 3 basic policy DSL.
3. **Sprint C**: Phase 4 query enhancements + Phase 6 first E2E path.
