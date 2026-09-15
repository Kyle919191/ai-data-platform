# Project progress

Living checklist for the AI Data Platform. Check a box when that step is verified (run + expected result), not when the file is merely typed.

**You are here:** Milestone 0, step **0.8** — Docker image for the gateway

Repo: https://github.com/Kyle919191/ai-data-platform

How to read this:

- `[x]` done
- `[ ]` not done
- Later milestones stay coarse until we are close; we will split them into typing steps like 0.1, 0.2, … when we get there.
- Do not skip ahead. Kafka starts at Milestone 1.

---

## Milestone 0 — Repository and Go foundations

**Goal:** one Go HTTP service you can run and test locally. **No Kafka.**

**Definition of done:** `make test`, `make run`, `curl localhost:8080/healthz`

- [x] **0.1** Go module (`ai-data-platform`) + `GET /healthz`
- [x] **0.2** Makefile: `make run`
- [x] **0.3** Event struct + `POST /v1/events` (accept JSON, return `event_id`)
- [x] **0.4** Minimal validation (reject bad / incomplete events)
- [x] **0.5** Log accepted events
- [x] **0.6** Unit tests for handlers / validation
- [x] **0.7** Makefile: `make test`
- [ ] **0.8** Docker image for the gateway (design lists this in M0; after tests)

---

## Milestone 1 — AI event ingestion

**Goal:** fake agent traffic hits the gateway; then Kafka is the first real distributed system.

**Definition of done:** 100+ events/sec generated and consumed reliably.

- [ ] Event generator (`cmd/event-generator`) → HTTP gateway (print/log events, still no Kafka)
- [ ] Add Kafka locally (Docker Compose)
- [ ] Gateway publishes to `ai.events.raw` (partition key: `trace_id`)
- [ ] Confirm consume path and basic Kafka ideas (topic, partition, offset)

---

## Milestone 2 — Reliable Kafka ingestion

**Goal:** the gateway still behaves when Kafka is slow, down, or sees duplicates.

**Definition of done:** you can explain exactly what happens during failure.

- [ ] Producer batching
- [ ] Retries + timeouts
- [ ] Idempotent event IDs
- [ ] Dead-letter queue (`ai.events.invalid`)
- [ ] Metrics on the gateway
- [ ] Graceful shutdown
- [ ] Load test; kill Kafka; restart; observe recovery

---

## Milestone 3 — Iceberg lakehouse

**Goal:** events become historical tables, not only a stream.

**Definition of done:** you can draw how an Iceberg table is stored on disk (MinIO + metadata + Parquet).

- [ ] MinIO (local S3)
- [ ] Iceberg catalog + first tables (`raw_events`, later `traces`, `model_calls`)
- [ ] Write events as Parquet via Iceberg
- [ ] Experiments: snapshot, time travel, schema evolution, partition evolution

---

## Milestone 4 — Flink streaming

**Goal:** real-time path `Kafka → Flink → Iceberg`.

**Definition of done:** kill Flink mid-job and show recovery from checkpoints.

- [ ] Trace aggregation job (group by `trace_id`)
- [ ] Windowed metrics job
- [ ] Event time vs processing time, watermarks, late events
- [ ] Checkpoints + crash recovery

---

## Milestone 5 — Spark batch

**Goal:** heavy offline work that should not live in Flink.

**Definition of done:** you can explain why this job is Spark rather than Flink.

- [ ] Historical trace reconstruction from `raw_events`
- [ ] Training dataset builder
- [ ] Historical backfill (e.g. fix bad `cost_usd` for a date range)

---

## Milestone 6 — Airflow

**Goal:** schedule and retry *offline* jobs. Streaming stays on Kafka/Flink.

**Definition of done:** failure/retry is visible; jobs are safe to rerun (idempotent).

- [ ] DAGs: daily aggregation, training dataset, backfill, Iceberg maintenance
- [ ] Task failure blocks downstream publish
- [ ] Observe retry behavior

---

## Milestone 7 — Snowflake + dbt

**Goal:** Iceberg remains the lakehouse source of truth; Snowflake serves modeled analytics.

**Definition of done:** answer business questions from dbt marts.

- [ ] Load curated data into Snowflake
- [ ] dbt staging / intermediate / marts + tests
- [ ] Marts such as model performance, customer cost, agent reliability

---

## Milestone 8 — Go control plane

**Goal:** self-service dataset API (control plane), not more pipelines (data plane).

**Definition of done:** `platform dataset create` provisions the abstraction without hand-wiring every system.

- [ ] Dataset registry + PostgreSQL metadata (not event data)
- [ ] `POST /v1/datasets`, status API
- [ ] CLI talks to the API only
- [ ] Lifecycle + reconciler (desired vs current state)

---

## Milestone 9 — Automated Iceberg maintenance

**Goal:** detect small files and compact without a human always noticing first.

**Definition of done:** generate small files → platform detects → compact → fewer/larger files.

- [ ] Table metrics collector
- [ ] Compaction recommendation + manual compact API
- [ ] Automated policy → Airflow → Spark rewrite

---

## Milestone 10 — Observability

**Goal:** diagnose a broken pipeline from metrics, not only logs.

**Definition of done:** you can find ingestion errors, Kafka lag, Flink checkpoint age, Iceberg health from dashboards.

- [ ] OpenTelemetry in Go services
- [ ] Prometheus + Grafana
- [ ] Dashboards: ingestion, Kafka, Flink, data health, maintenance

---

## Milestone 11 — Exposure labs (one at a time, after the core works)

Each lab: what problem it solves, overlap with our stack, best workload, would AgentHub use it?

- [ ] Trino over Iceberg
- [ ] ClickHouse real-time metrics
- [ ] Druid comparison
- [ ] Airbyte Postgres → Snowflake
- [ ] Databricks vs local Spark

---

## Milestone 12 — Kubernetes (last)

**Goal:** deploy what already works locally. Do not start here.

- [ ] Deploy Go services, Kafka-compatible cluster, Flink, control-plane deps
- [ ] Deployments, Services, ConfigMaps, Secrets, probes, limits, scaling

---

## Cross-cutting (attach to the milestone where they first appear)

Do not treat these as a reason to jump ahead.

- [ ] ADRs under `docs/decisions/`
- [ ] Learning note after each milestone (`docs/learning/`)
- [ ] Failure experiments (duplicate events, late events, killed workers, …)
- [ ] Benchmarks (100 / 1k / 10k events/sec) when the path exists
- [ ] Security later: auth, tenant isolation, PII scrubbing
