# DataTracing v2 — Project Purpose and Problem Statement

## Why this project exists
Modern data platforms are distributed across orchestrators, stream processors, batch jobs, APIs, and warehouses. When data quality incidents happen, teams often struggle to answer basic operational questions quickly:

- Where did a specific record fail?
- Which service or task introduced latency?
- Which upstream dependency caused this downstream breakage?
- Was the trace dropped by sampling, ingestion, or storage?

DataTracing v2 is being built to solve this observability gap by providing **trace-level visibility for data pipelines** ("Jaeger for data pipelines").

## Core challenges we are solving

1. **Fragmented execution context**
   - Pipeline steps run across heterogeneous systems (HTTP, Kafka, workflow engines).
   - We need propagation and correlation primitives that survive transport boundaries.

2. **Difficult end-to-end debugging**
   - Logs and metrics alone are insufficient for causal reconstruction.
   - We need a trace graph (DAG) that shows parent-child relationships and cross-links.

3. **Scale and cost constraints**
   - High-throughput pipelines create huge span volumes.
   - We need controlled ingestion and tail sampling so we can keep high-value traces.

4. **Architecture sustainability**
   - Observability systems evolve quickly (new stores, protocols, integrations).
   - We need clean architecture boundaries so infrastructure can change without rewriting core logic.

## Product purpose

DataTracing v2 aims to deliver a dependable tracing foundation for data systems with:

- A canonical tracing model and SDK.
- Context propagation across transport types.
- Collector ingestion with sampling.
- Query APIs to retrieve trace DAGs and search trace summaries.
- Pluggable storage adapters, starting with local file and memory, then ClickHouse.

## Success criteria (high-level)

- Engineers can instrument services quickly with stable SDK interfaces.
- A trace can be reconstructed end-to-end across boundaries.
- Query APIs return consistent, correct DAGs and summaries.
- Storage backends can be swapped with minimal application-layer impact.
- Architecture rules remain enforceable by tests/CI over time.
