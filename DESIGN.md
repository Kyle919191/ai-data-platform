# AI Data Platform

## A Self-Service Data Platform for LLM / Agent Telemetry, Analytics, and Training Data

**Primary implementation language:** Go  
**Project type:** Data Platform / Data Infrastructure Engineering  
**Primary use case:** AI agent and LLM event data  
**Target environment:** Local development first, then cloud/container deployment  
**Learning philosophy:** Build one coherent platform incrementally. Every new technology must solve a concrete problem introduced by the previous milestone.

---

# 1. Project Mission

Build a production-inspired data platform for applications that run LLMs and AI agents.

The platform should accept high-volume AI telemetry such as:

- user prompts;
- model responses;
- model calls;
- tool calls;
- retrieval operations;
- token usage;
- model latency;
- errors;
- user feedback;
- evaluation scores;
- agent traces;
- cost information.

It must support both:

1. **real-time operational analytics**, such as model latency, error rate, throughput, and cost;
2. **offline analytical and ML datasets**, such as training pairs, evaluation datasets, customer usage reports, and historical model-performance comparisons.

The project is not primarily an analytics project.

The goal is to learn how a **data platform engineer builds infrastructure that allows other engineers, analysts, and ML teams to reliably produce and consume data.**

---

# 2. Core Learning Goals

By the end of the project, the engineer should understand:

### Distributed ingestion

- Kafka partitions
- offsets
- consumer groups
- ordering guarantees
- producer acknowledgments
- batching
- retries
- idempotency
- backpressure
- dead-letter queues

### Stream processing

- event time
- processing time
- watermarks
- late-arriving events
- windows
- stateful processing
- checkpointing
- recovery
- exactly-once vs at-least-once semantics

### Batch processing

- distributed transformations
- large joins
- partition pruning
- shuffle
- backfills
- incremental processing
- job retries

### Data lakehouse

- Parquet
- Apache Iceberg
- snapshots
- manifests
- metadata
- schema evolution
- partition evolution
- optimistic concurrency
- compaction
- small-file problems
- snapshot expiration
- multi-engine access

### Warehousing

- Snowflake
- warehouse modeling
- data marts
- SQL transformations
- analytical consumption patterns

### Transformation

- dbt
- staging/intermediate/mart layers
- incremental models
- testing
- lineage
- documentation

### Orchestration

- Airflow DAGs
- scheduling
- retries
- sensors
- backfills
- dependency management
- idempotent jobs

### Data platform engineering

- control plane vs data plane
- self-service APIs
- dataset registration
- pipeline lifecycle management
- health checks
- metadata
- platform observability
- automated maintenance
- failure recovery
- infrastructure abstraction

---

# 3. Required Technologies

## Core

The final system must meaningfully use:

- Go
- Apache Kafka
- Apache Flink
- Apache Spark
- Apache Iceberg
- Apache Airflow
- dbt
- Snowflake

## Exposure technologies

These do not need to be part of the critical production path.

They should instead be explored through focused labs built around the same dataset:

- Databricks
- ClickHouse
- Apache Druid
- Trino
- Airbyte

The goal is understanding **why and when these technologies are appropriate**, not merely adding their logos to the architecture.

---

# 4. High-Level Use Case

Assume a fictional company called **AgentHub**.

AgentHub operates AI agents for customers.

A single agent request may execute:

```text
User prompt
    ↓
Retrieval
    ↓
LLM call
    ↓
Tool call
    ↓
LLM call
    ↓
Tool call
    ↓
LLM call
    ↓
Final answer
    ↓
User feedback
```

Each step emits an event.

The platform must ingest those events and support questions such as:

```text
What is GPT-X's p95 latency?

Which models are most expensive per successful task?

Which tools fail most often?

Which customers consume the most tokens?

What percentage of agent runs successfully complete?

Did a new model deployment increase error rates?

Can we create a high-quality training dataset from successful sessions?

Can we reconstruct every event belonging to a particular agent trace?

How far behind is the streaming pipeline?

Which Iceberg tables have accumulated too many small files?
```

---

# 5. Overall Architecture

```text
                         ┌─────────────────────┐
                         │   AI Applications   │
                         │                     │
                         │ agents / LLM apps   │
                         └──────────┬──────────┘
                                    │
                                    │ HTTP / gRPC
                                    ▼
                         ┌─────────────────────┐
                         │ Ingestion Gateway   │
                         │        Go           │
                         └──────────┬──────────┘
                                    │
                                    ▼
                              ┌──────────┐
                              │  Kafka   │
                              └────┬─────┘
                                   │
                 ┌─────────────────┼────────────────┐
                 │                 │                │
                 ▼                 ▼                ▼
             Flink            Raw Sink        Dead Letter
          streaming jobs         │              Queue
                 │               │
          ┌──────┴──────┐        ▼
          │             │      Iceberg
          ▼             ▼         │
     ClickHouse       Iceberg     │
     real-time         hot data   │
      metrics                     │
                                  ▼
                                Spark
                           batch processing
                                  │
                                  ▼
                               Iceberg
                           curated datasets
                                  │
                     ┌────────────┼─────────────┐
                     │            │             │
                     ▼            ▼             ▼
                   Trino       Snowflake    Databricks
                                  │
                                  ▼
                                 dbt
                                  │
                                  ▼
                           Analytical marts


 Operational DB / SaaS
          │
       Airbyte
          │
          ▼
      Snowflake


              ┌───────────────────────────────┐
              │         Go Control Plane      │
              │                               │
              │ datasets                      │
              │ pipelines                     │
              │ health                        │
              │ maintenance                   │
              │ metadata                      │
              └──────────────┬────────────────┘
                             │
                             ▼
                          Airflow
                     orchestration layer
```

---

# 6. Architectural Principle: Data Plane vs Control Plane

This distinction is central to the project.

## Data plane

The data plane performs actual data movement and processing.

It contains:

```text
Kafka
Flink
Spark
Iceberg
Snowflake
ClickHouse
Trino
```

Example:

```text
Kafka → Flink → Iceberg
```

That is a data-plane operation.

## Control plane

The control plane manages the lifecycle of those systems.

It is primarily written by us in Go.

Example:

```text
POST /v1/datasets
```

might cause the control plane to:

```text
create dataset metadata
↓
create Kafka topic configuration
↓
create Iceberg table
↓
register maintenance policy
↓
register monitoring
↓
return dataset status
```

This distinction should remain explicit throughout the project.

---

# 7. Repository Structure

Use a monorepo.

```text
ai-data-platform/
│
├── README.md
├── Makefile
├── docker-compose.yml
├── go.mod
│
├── cmd/
│   ├── ingestion-gateway/
│   │   └── main.go
│   │
│   ├── control-plane/
│   │   └── main.go
│   │
│   ├── event-generator/
│   │   └── main.go
│   │
│   └── platform-cli/
│       └── main.go
│
├── internal/
│   ├── ingestion/
│   │   ├── api/
│   │   ├── kafka/
│   │   ├── validation/
│   │   └── batching/
│   │
│   ├── controlplane/
│   │   ├── dataset/
│   │   ├── pipeline/
│   │   ├── health/
│   │   ├── metadata/
│   │   ├── maintenance/
│   │   └── scheduler/
│   │
│   ├── events/
│   ├── config/
│   ├── telemetry/
│   └── storage/
│
├── schemas/
│   ├── event.schema.json
│   ├── trace.schema.json
│   ├── model_call.schema.json
│   ├── tool_call.schema.json
│   └── feedback.schema.json
│
├── flink/
│   ├── trace-aggregator/
│   ├── realtime-metrics/
│   └── dead-letter-handler/
│
├── spark/
│   ├── session-builder/
│   ├── training-dataset/
│   ├── backfill/
│   └── compaction/
│
├── airflow/
│   └── dags/
│       ├── training_dataset.py
│       ├── daily_aggregation.py
│       ├── iceberg_maintenance.py
│       └── historical_backfill.py
│
├── dbt/
│   ├── dbt_project.yml
│   ├── models/
│   │   ├── staging/
│   │   ├── intermediate/
│   │   └── marts/
│   └── tests/
│
├── infra/
│   ├── docker/
│   ├── terraform/
│   └── kubernetes/
│
├── observability/
│   ├── prometheus/
│   ├── grafana/
│   └── otel/
│
├── labs/
│   ├── clickhouse/
│   ├── druid/
│   ├── trino/
│   ├── databricks/
│   └── airbyte/
│
├── benchmarks/
│
├── docs/
│   ├── architecture.md
│   ├── event-model.md
│   ├── reliability.md
│   ├── experiments.md
│   └── decisions/
│
└── tests/
    ├── integration/
    ├── load/
    └── failure/
```

---

# 8. Canonical Event Model

Every event should have a common envelope.

```go
type Event struct {
    EventID       string          `json:"event_id"`
    TraceID       string          `json:"trace_id"`
    SessionID     string          `json:"session_id"`
    OrganizationID string         `json:"organization_id"`
    UserID        string          `json:"user_id"`

    EventType     string          `json:"event_type"`
    EventTime     time.Time       `json:"event_time"`

    SchemaVersion int             `json:"schema_version"`

    Payload       json.RawMessage `json:"payload"`
}
```

Example:

```json
{
  "event_id": "evt_82a13",
  "trace_id": "trace_a91",
  "session_id": "sess_887",
  "organization_id": "org_42",
  "user_id": "user_812",
  "event_type": "model_call_completed",
  "event_time": "2026-09-14T18:30:22.123Z",
  "schema_version": 1,
  "payload": {
    "model": "model-x",
    "provider": "provider-a",
    "prompt_tokens": 820,
    "completion_tokens": 314,
    "latency_ms": 1281,
    "cost_usd": 0.0143,
    "success": true
  }
}
```

---

# 9. Event Types

Initial event types:

```text
trace_started
trace_completed

model_call_started
model_call_completed
model_call_failed

tool_call_started
tool_call_completed
tool_call_failed

retrieval_started
retrieval_completed

feedback_received

evaluation_completed
```

Do not introduce all of them immediately.

Start only with:

```text
trace_started
model_call_completed
trace_completed
```

Add complexity later.

---

# 10. Kafka Design

Initial topics:

```text
ai.events.raw
ai.events.validated
ai.events.invalid
ai.events.metrics
```

Later:

```text
ai.model.calls
ai.tool.calls
ai.feedback
ai.evaluations
```

## Partitioning

Primary partition key:

```text
trace_id
```

Reason:

Events belonging to one trace should generally retain Kafka partition ordering.

This enables Flink to reconstruct traces more naturally.

## Questions the project must eventually answer

The engineer should be able to explain:

```text
Why use multiple partitions?

How does partition count affect parallelism?

What happens when a consumer dies?

What does consumer-group rebalancing mean?

What happens if a producer retries?

What guarantees does Kafka provide about ordering?

What does an offset represent?

When should an offset be committed?

What happens when consumers are slower than producers?
```

---

# 11. Go Ingestion Gateway

The first major service.

## Responsibilities

```text
receive event
↓
validate request
↓
assign IDs if necessary
↓
validate schema
↓
attach ingestion metadata
↓
serialize
↓
publish to Kafka
↓
return acknowledgment
```

## API

### POST /v1/events

Request:

```json
{
  "trace_id": "trace_123",
  "session_id": "sess_456",
  "event_type": "model_call_completed",
  "event_time": "2026-09-14T18:30:22Z",
  "payload": {
    "model": "model-x",
    "latency_ms": 822
  }
}
```

Response:

```json
{
  "event_id": "evt_8391",
  "accepted": true
}
```

## Later additions

```text
POST /v1/events/batch

GET /healthz

GET /readyz

GET /metrics
```

## Technical concepts to implement

Start simple.

Then progressively introduce:

```text
goroutines
channels
bounded queues
producer batching
timeouts
context cancellation
retry policies
structured logging
graceful shutdown
rate limiting
backpressure
```

---

# 12. Event Generator

Create a Go program that simulates AI-agent traffic.

Example invocation:

```bash
go run ./cmd/event-generator \
  --rate 100 \
  --agents 20 \
  --duration 10m
```

The generator should simulate:

```text
normal traces

slow model responses

tool failures

model failures

delayed events

duplicate events

out-of-order events
```

This service becomes critical for testing distributed-system behavior.

---

# 13. Flink Responsibilities

Flink handles the real-time path.

## Job 1: Trace reconstruction

Input:

```text
Kafka ai.events.validated
```

Group by:

```text
trace_id
```

Maintain state:

```text
trace start
model calls
tool calls
token usage
cost
errors
completion
```

Output:

```text
completed trace
```

Example:

```json
{
  "trace_id": "trace_123",
  "duration_ms": 4218,
  "model_calls": 3,
  "tool_calls": 2,
  "prompt_tokens": 1842,
  "completion_tokens": 719,
  "cost_usd": 0.038,
  "success": true
}
```

---

# 14. Flink Event-Time Exercise

This is mandatory.

Create events where:

```text
event A timestamp = 10:00:01
event B timestamp = 10:00:03
event C timestamp = 10:00:02
```

but they arrive:

```text
A
B
C
```

Learn why:

```text
event time != arrival time
```

Then configure watermarks.

Test:

```text
late events
```

and determine:

```text
How late is acceptable?

What happens to events beyond the watermark?

Should late events modify previous aggregates?

Should they go to a side output?
```

---

# 15. Flink Real-Time Metrics Job

Compute rolling metrics.

For each model:

```text
request_count
success_count
failure_count

p50 latency
p95 latency
p99 latency

prompt_tokens
completion_tokens

estimated_cost

requests_per_second
```

Windows:

```text
1 minute
5 minutes
15 minutes
```

Output initially to logs.

Later output to ClickHouse.

---

# 16. Iceberg Data Lake

Iceberg is the canonical historical data layer.

Initial tables:

```text
raw_events
traces
model_calls
```

Later:

```text
tool_calls
feedback
evaluations
training_examples
```

---

# 17. Object Storage

For local development:

```text
MinIO
```

Treat MinIO as local S3-compatible object storage.

Example:

```text
s3://agenthub-lake/raw_events/
s3://agenthub-lake/traces/
s3://agenthub-lake/model_calls/
```

Do not treat object-store directories as database tables.

Iceberg metadata defines the table.

---

# 18. Iceberg Concepts to Study Directly

For every concept below, create a small experiment.

## Snapshot

Insert data.

Inspect the snapshot.

Insert more data.

Inspect the new snapshot.

Query the old snapshot.

## Schema evolution

Start:

```text
model
latency_ms
tokens
```

Add:

```text
provider
region
```

Verify old data still works.

## Partition evolution

Start:

```text
day(event_time)
```

Later change partition strategy.

Understand why Iceberg can evolve partitioning without rewriting every existing query.

## Optimistic concurrency

Attempt concurrent writes.

Observe commit behavior.

Understand why metadata commits can conflict.

## Small files

Generate thousands of tiny writes.

Measure:

```text
number of files
average file size
query behavior
metadata overhead
```

Then compact.

---

# 19. Spark Responsibilities

Spark owns compute-heavy offline work.

Do not use Spark for everything.

## Spark job 1: historical trace rebuild

Read:

```text
raw_events
```

Reconstruct traces historically.

Purpose:

Compare batch reconstruction with Flink real-time reconstruction.

This teaches:

```text
batch vs streaming
```

---

# 20. Spark Job 2: Training Dataset Builder

Input:

```text
traces
model_calls
feedback
evaluations
```

Output:

```text
training_examples
```

Example:

```json
{
  "prompt": "...",
  "response": "...",
  "tool_trace": [],
  "model": "model-x",
  "user_feedback": 1,
  "evaluation_score": 0.94,
  "quality_bucket": "high"
}
```

Pipeline operations might include:

```text
join trace events

remove failed runs

deduplicate

filter PII

filter low-quality feedback

calculate quality score

sample dataset

write curated Iceberg table
```

---

# 21. Spark Job 3: Historical Backfill

Suppose a bug causes:

```text
cost_usd
```

to be incorrectly calculated for one week.

Write a Spark backfill that:

```text
reads affected period
↓
recalculates cost
↓
rewrites affected rows
↓
publishes corrected dataset
```

This is an extremely important data-engineering scenario.

---

# 22. Snowflake

Snowflake becomes the managed analytical warehouse.

Snowflake does not replace Iceberg.

Conceptual division:

```text
Iceberg
canonical lakehouse datasets

Snowflake
curated analytical warehouse
```

---

# 23. Snowflake Data Model

Create tables such as:

```text
raw_model_calls

raw_agent_traces

raw_feedback
```

Then dbt builds:

```text
stg_model_calls
stg_agent_traces
stg_feedback

int_agent_sessions
int_model_usage

fact_agent_runs
fact_model_calls
fact_tool_calls

dim_model
dim_organization

mart_model_performance
mart_customer_cost
mart_agent_reliability
```

---

# 24. dbt Project

Example dependency:

```text
raw_model_calls
        │
        ▼
stg_model_calls
        │
        ▼
int_model_usage
        │
        ▼
fact_model_calls
        │
        ▼
mart_model_performance
```

Tests should include:

```text
not_null
unique
relationships
accepted_values
```

Create custom tests for:

```text
latency >= 0

tokens >= 0

cost >= 0

successful trace must have completion timestamp
```

---

# 25. Airflow

Airflow coordinates offline workflows.

It should not control the per-event streaming path.

Kafka and Flink run continuously.

Airflow handles jobs such as:

```text
daily aggregation

training dataset generation

historical backfill

Iceberg maintenance

dbt transformation

quality checks
```

---

# 26. Example Airflow DAG

```text
wait_for_daily_partition
          │
          ▼
spark_build_sessions
          │
          ▼
validate_sessions
          │
          ▼
spark_build_training_data
          │
          ▼
publish_to_snowflake
          │
          ▼
dbt_run
          │
          ▼
dbt_test
          │
          ▼
publish_success
```

Failure of a task should prevent downstream publication.

---

# 27. Go Control Plane

This is the flagship component.

It distinguishes the project from a normal data-engineering pipeline.

The control plane manages datasets.

---

# 28. Dataset Resource

Internal representation:

```go
type Dataset struct {
    ID              string
    Name            string
    Description     string

    KafkaTopic      string
    IcebergTable    string

    PartitionKey    string
    RetentionDays   int

    SchemaVersion   int

    Status          DatasetStatus

    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

---

# 29. Control Plane API

## Create dataset

```http
POST /v1/datasets
```

Request:

```json
{
  "name": "model-events",
  "description": "All model-call telemetry",
  "partition_key": "trace_id",
  "retention_days": 30
}
```

Response:

```json
{
  "id": "ds_123",
  "name": "model-events",
  "status": "PROVISIONING"
}
```

Eventually the platform provisions:

```text
Kafka configuration
Iceberg table
monitoring configuration
maintenance policy
metadata record
```

---

# 30. Dataset Status API

```http
GET /v1/datasets/{id}/status
```

Response:

```json
{
  "dataset": "model-events",

  "kafka": {
    "partitions": 8,
    "consumer_lag": 241
  },

  "flink": {
    "status": "RUNNING",
    "checkpoint_age_seconds": 11
  },

  "iceberg": {
    "snapshot_count": 328,
    "file_count": 12814,
    "small_file_count": 3128,
    "average_file_size_mb": 13.2
  },

  "freshness": {
    "latest_event_age_seconds": 2.3
  },

  "maintenance": {
    "compaction_recommended": true
  }
}
```

---

# 31. CLI

Build a CLI client.

Examples:

```bash
platform dataset create \
  --name model-events \
  --partition-key trace_id
```

```bash
platform dataset list
```

```bash
platform dataset status model-events
```

```bash
platform dataset compact model-events
```

```bash
platform dataset backfill model-events \
  --from 2026-09-01 \
  --to 2026-09-07
```

The CLI simply talks to the control-plane API.

Do not duplicate business logic in the CLI.

---

# 32. Metadata Database

The control plane needs its own transactional state.

Use:

```text
PostgreSQL
```

Store:

```text
datasets
pipelines
maintenance_jobs
schema_versions
job_runs
```

Do not store the actual analytical event data here.

---

# 33. Dataset Lifecycle

State machine:

```text
CREATING
   ↓
PROVISIONING
   ↓
ACTIVE
   ↓
DEGRADED
   ↓
ACTIVE

ACTIVE
   ↓
DELETING
   ↓
DELETED
```

Failed provisioning:

```text
PROVISIONING
   ↓
FAILED
```

This becomes an opportunity to learn control-plane reconciliation.

---

# 34. Reconciliation

Do not assume an API request completes every distributed operation immediately.

Example:

```text
POST dataset

database says:
desired state = ACTIVE

current state = PROVISIONING
```

Background reconciler checks:

```text
Does Kafka topic exist?

Does Iceberg table exist?

Are expected jobs configured?

Are health checks passing?
```

Then:

```text
current state = ACTIVE
```

This introduces a Kubernetes-style control-loop concept without requiring Kubernetes initially.

---

# 35. Automatic Iceberg Maintenance

This should become the advanced feature of the project.

The control plane periodically evaluates tables.

Metrics:

```text
small_file_count
average_file_size
snapshot_count
time_since_last_compaction
write_rate
table_size
```

---

# 36. Initial Compaction Policy

Start simple:

```text
if small_file_count > threshold:
    compact
```

Later:

```text
if small_file_ratio > 0.40:
    compact
```

Then:

```text
if average_file_size < target_file_size
AND write_rate < low_activity_threshold:
    compact
```

This creates a basic policy engine.

---

# 37. Maintenance Flow

```text
Go Control Plane
      │
      ▼
Maintenance Decision
      │
      ▼
Airflow API
      │
      ▼
Spark Compaction Job
      │
      ▼
Iceberg Rewrite
      │
      ▼
New Snapshot
      │
      ▼
Metrics Updated
```

---

# 38. Idempotency

Every major write path must eventually address idempotency.

Example problem:

```text
producer sends event
↓
network timeout
↓
producer doesn't know whether Kafka received it
↓
producer retries
```

Possible duplicate.

Every event has:

```text
event_id
```

Downstream processing should support deduplication.

You must be able to explain why:

```text
"exactly once"
```

is much more complicated than merely setting one configuration flag.

---

# 39. Dead-Letter Queue

Invalid events must not silently disappear.

Flow:

```text
event
  │
  ▼
validation
  │
  ├── valid ───► ai.events.validated
  │
  └── invalid ─► ai.events.invalid
```

DLQ records:

```text
original event
error type
error message
ingestion timestamp
schema version
```

Build a CLI:

```bash
platform dlq inspect
```

Later:

```bash
platform dlq replay
```

---

# 40. Schema Evolution

Version every event.

Example:

```text
schema_version = 1
```

Later add:

```text
schema_version = 2
```

Do not break existing consumers immediately.

Exercise:

Version 1:

```json
{
  "model": "x",
  "latency_ms": 500
}
```

Version 2:

```json
{
  "model": "x",
  "latency_ms": 500,
  "region": "us-central"
}
```

Understand backward compatibility.

---

# 41. Observability

Instrument your Go services.

Use:

```text
OpenTelemetry
Prometheus
Grafana
```

Metrics:

```text
ingestion_requests_total

ingestion_errors_total

kafka_publish_latency

events_published_total

events_rejected_total

controlplane_requests_total

dataset_provisioning_duration

maintenance_jobs_total

maintenance_failures_total
```

---

# 42. Platform Dashboard

Create Grafana dashboards covering:

## Ingestion

```text
events/sec

publish latency

error rate
```

## Kafka

```text
consumer lag

partition throughput
```

## Flink

```text
checkpoint duration

checkpoint failures

records processed

backpressure
```

## Data health

```text
freshness

small files

latest snapshot age
```

## Maintenance

```text
compactions

duration

bytes rewritten

failures
```

---

# 43. Failure Testing

Distributed systems knowledge comes from breaking things.

Create deliberate failure experiments.

## Kafka unavailable

Expected:

```text
gateway retry
bounded queue
eventual failure response
metrics increase
```

## Flink worker crashes

Observe:

```text
checkpoint recovery
state restoration
Kafka replay
```

## Duplicate event

Verify:

```text
deduplication behavior
```

## Late event

Verify:

```text
watermark behavior
```

## Spark job fails halfway

Verify:

```text
Airflow retry
idempotent rerun
```

## Iceberg commit conflict

Observe:

```text
optimistic concurrency behavior
```

## Snowflake unavailable

Ensure upstream Iceberg datasets remain valid.

This reinforces decoupling.

---

# 44. Data Quality

Create checks for:

```text
event_id not null

trace_id not null

latency >= 0

token count >= 0

cost >= 0

known event_type

known schema_version
```

Also create semantic checks:

```text
trace_completed must reference existing trace_started

successful trace should contain >= 1 model call

trace completion >= trace start
```

---

# 45. ClickHouse Lab

Purpose:

Explore low-latency real-time analytical serving.

Flow:

```text
Kafka
  ↓
Flink
  ↓
ClickHouse
```

Load aggregated metrics.

Benchmark queries such as:

```sql
SELECT
    model,
    count(*),
    quantile(0.95)(latency_ms)
FROM model_metrics
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY model;
```

Compare:

```text
ClickHouse
vs
Snowflake
vs
Trino over Iceberg
```

Write conclusions.

---

# 46. Druid Lab

Load the same real-time telemetry dataset into Druid.

Benchmark:

```text
time filtering

GROUP BY model

high-cardinality organization filters

time-series aggregation

approximate percentiles
```

Compare Druid with ClickHouse.

The objective is not declaring one universally better.

Explain:

```text
which workload favors which system?
```

---

# 47. Trino Lab

Connect Trino to Iceberg.

Use SQL to query:

```text
raw_events

traces

training_examples
```

Compare:

```text
Trino interactive SQL

vs

Spark batch processing
```

Learn why one engine may be more appropriate for interactive queries while another is better for large transformation jobs.

---

# 48. Databricks Lab

Take one existing Spark job.

Example:

```text
training dataset builder
```

Run locally first.

Then port it to Databricks.

Compare:

```text
local Spark

vs

managed Spark
```

Explore:

```text
jobs

clusters

Spark UI

catalog/governance

Iceberg interoperability

operational experience
```

Document what Databricks provides that your local setup does not.

---

# 49. Airbyte Lab

Create a transactional PostgreSQL database representing AgentHub's application database.

Tables:

```text
organizations

users

subscriptions

plans
```

Replicate:

```text
Postgres
   ↓
Airbyte
   ↓
Snowflake
```

Then dbt joins business data with AI usage data.

Example mart:

```text
organization_id

subscription_plan

monthly_ai_requests

monthly_tokens

monthly_cost

success_rate
```

This teaches the distinction between:

```text
event streaming

vs

database replication / ELT
```

---

# 50. Local Development Strategy

Do not start the entire architecture at once.

Use Docker Compose progressively.

Initial environment:

```text
Kafka

Postgres
```

Then add:

```text
MinIO
Iceberg catalog
```

Then:

```text
Flink
```

Then:

```text
Spark
```

Later:

```text
Airflow
Prometheus
Grafana
```

Snowflake remains external.

Exposure technologies are introduced only after the core platform works.

---

# 51. Milestone 0 — Repository and Go Foundations

Goal:

```text
create clean repository
run Go service
run tests
build Docker image
```

Build:

```text
GET /healthz
```

Learn:

```text
Go modules
package structure
HTTP server
configuration
logging
tests
Docker
```

Do not add Kafka yet.

Definition of done:

```text
make test

make run

curl localhost:<port>/healthz
```

works.

---

# 52. Milestone 1 — AI Event Ingestion

Build:

```text
event generator
     ↓
Go ingestion gateway
```

No Kafka initially.

Print validated events.

Then add Kafka:

```text
event generator
     ↓
Go ingestion gateway
     ↓
Kafka
```

Definition of done:

```text
100+ events/sec can be generated and consumed reliably.
```

Learn:

```text
Kafka fundamentals
```

---

# 53. Milestone 2 — Reliable Kafka Ingestion

Add:

```text
batching

retry

timeouts

idempotent IDs

DLQ

metrics

graceful shutdown
```

Run load test.

Kill Kafka.

Restart Kafka.

Observe recovery.

Definition of done:

You can explain exactly what happens during failure.

---

# 54. Milestone 3 — Iceberg Lakehouse

Add:

```text
MinIO

Iceberg

Parquet
```

Write events into Iceberg.

Perform experiments:

```text
snapshot

time travel

schema evolution

partition evolution
```

Definition of done:

You can draw how an Iceberg table is physically represented.

---

# 55. Milestone 4 — Flink Streaming

Pipeline:

```text
Kafka
  ↓
Flink
  ↓
Iceberg
```

Implement:

```text
trace aggregation
```

Then:

```text
windowed metrics
```

Study:

```text
watermarks

late events

checkpoints

state
```

Definition of done:

Kill Flink during processing and demonstrate recovery.

---

# 56. Milestone 5 — Spark Batch Processing

Create:

```text
historical trace reconstruction

training dataset builder
```

Read/write Iceberg.

Run historical backfill.

Definition of done:

Explain why this job is Spark rather than Flink.

---

# 57. Milestone 6 — Airflow

Create DAGs for:

```text
daily aggregation

training dataset

backfill

Iceberg maintenance
```

Definition of done:

Failure/retry behavior is observable and jobs are idempotent.

---

# 58. Milestone 7 — Snowflake + dbt

Load curated analytical data.

Create dbt models.

Build analytics marts.

Definition of done:

Answer business questions through modeled Snowflake tables.

---

# 59. Milestone 8 — Go Control Plane

Create:

```text
dataset registry

dataset API

dataset status API

CLI

health aggregation
```

Definition of done:

```bash
platform dataset create
```

creates a platform dataset abstraction without users manually configuring every underlying technology.

---

# 60. Milestone 9 — Automated Maintenance

Implement:

```text
Iceberg table metrics collector

compaction recommendation

manual compact API

automated compaction policy
```

Definition of done:

Generate small files, let the platform detect them, run compaction, and demonstrate improvement.

---

# 61. Milestone 10 — Platform Observability

Add:

```text
Prometheus

Grafana

OpenTelemetry
```

Definition of done:

You can diagnose a broken pipeline from platform metrics.

---

# 62. Milestone 11 — Technology Exposure Labs

One at a time:

```text
Trino

ClickHouse

Druid

Airbyte

Databricks
```

Each lab must answer:

```text
What problem does this technology solve?

Where does it overlap with our current platform?

What workload is it particularly good at?

What would it replace, if anything?

What operational complexity does it add?

Would we use it in production for AgentHub?
```

---

# 63. Milestone 12 — Kubernetes

Only after everything works locally.

Deploy:

```text
Go services

Kafka or compatible Kafka implementation

Flink

control-plane dependencies
```

Do not necessarily self-host Snowflake or other managed systems.

Learn:

```text
Deployments

Services

ConfigMaps

Secrets

health probes

resource limits

horizontal scaling
```

Kubernetes is deliberately late in the project.

---

# 64. Testing Strategy

Use four categories.

## Unit tests

Examples:

```text
schema validation

maintenance policy

dataset state transitions

cost calculations
```

## Integration tests

Examples:

```text
Go → Kafka

Kafka → Flink

Spark → Iceberg

Control Plane → metadata DB
```

## End-to-end tests

Example:

```text
generate agent trace
↓
ingest
↓
Kafka
↓
Flink
↓
Iceberg
↓
Spark
↓
Snowflake
↓
dbt
```

Validate output.

## Failure tests

Explicitly crash services.

---

# 65. Benchmarking

Create reproducible load profiles.

Example:

```text
small
100 events/sec

medium
1,000 events/sec

large
10,000 events/sec
```

Record:

```text
throughput

p50 ingestion latency

p95 ingestion latency

Kafka lag

Flink processing latency

CPU

memory
```

The absolute numbers are not the main goal.

Understanding bottlenecks is.

---

# 66. Performance Experiments

Conduct experiments such as:

## Kafka partition count

Compare:

```text
1

4

8

16
```

partitions.

## Go batching

Compare:

```text
1 event/batch

10

100

1000
```

## Iceberg file size

Compare many tiny files versus fewer large files.

## Spark partition count

Observe:

```text
parallelism

shuffle

task overhead
```

## Query engines

Compare:

```text
Trino

Spark SQL

Snowflake

ClickHouse
```

for selected workloads.

---

# 67. Security — Later Phase

Do not overcomplicate early development.

Eventually introduce:

```text
API authentication

organization isolation

secrets management

TLS

PII masking

audit logging
```

AI data can contain sensitive prompts.

Add a simple PII-scrubbing stage to Spark later.

---

# 68. Cost Model

Store per-model pricing.

Example table:

```text
model_prices

model
prompt_price_per_million_tokens
completion_price_per_million_tokens
effective_from
```

Use this to calculate:

```text
cost per request

cost per trace

cost per organization

cost per successful task
```

This makes the AI use case much more realistic.

---

# 69. Dataset for ML / Training

The final curated dataset should demonstrate why the platform exists.

Example:

```text
training_examples
```

Fields:

```text
trace_id

prompt

response

tool_history

model

prompt_tokens

completion_tokens

cost

latency

user_feedback

evaluation_score

quality_score

created_at
```

This makes it possible to say:

```text
The same platform supports both operational analytics and AI training/evaluation datasets.
```

---

# 70. Important Design Decisions

Create ADR files under:

```text
docs/decisions/
```

Example:

```text
0001-use-kafka-for-event-backbone.md

0002-use-iceberg-as-canonical-table-format.md

0003-use-go-for-control-plane.md

0004-separate-flink-and-spark-workloads.md

0005-snowflake-is-analytical-serving-not-source-of-truth.md

0006-use-airflow-for-batch-orchestration.md
```

Each ADR:

```text
Context

Decision

Alternatives considered

Consequences
```

This forces architectural thinking.

---

# 71. Things the Agent Must NOT Do

The coding agent must not:

```text
generate the entire project immediately

introduce unnecessary abstraction before the concept is understood

install every technology on day one

hide distributed-systems behavior behind libraries without explanation

skip testing

skip failure experiments

move to the next milestone before current acceptance criteria work

rewrite working components solely for style

use Kubernetes during early milestones

silently change architecture from this design document
```

---

# 72. Required Agent Teaching Behavior

The IDE coding agent should behave like a senior data-platform engineer mentoring a junior engineer.

For every implementation step:

### First explain

Before writing code, explain:

```text
what we are building

why it is needed

where it sits in the architecture

what concept I should learn

what alternatives exist
```

### Then propose a small step

Never implement an entire milestone at once.

Example:

Bad:

```text
Let's build Kafka ingestion, Flink, Iceberg, and monitoring.
```

Good:

```text
First we will create one Kafka topic and write one Go producer.

After that works, we will write one consumer.

Only after we understand offsets will we integrate the gateway.
```

### Make me participate

When reasonable, ask me to predict behavior before running an experiment.

Examples:

```text
What do you think happens if the consumer crashes before committing the offset?

What ordering guarantee do you expect across two Kafka partitions?

What happens to a Flink window if an event arrives after the watermark?
```

Then test the prediction.

### Explain code deeply

For important code, explain:

```text
what each abstraction is doing

why concurrency is needed

where failure can happen

what state exists

what guarantees exist
```

Do not merely generate files.

---

# 73. Agent Workflow for Every Milestone

The agent should follow:

```text
1. Explain architecture concept.

2. Show the exact mini-goal.

3. Explain files we will create.

4. Let me implement or inspect critical sections when useful.

5. Run the component.

6. Verify normal behavior.

7. Break the component intentionally.

8. Explain what happened.

9. Add tests.

10. Write a short learning note.

11. Only then proceed.
```

---

# 74. Learning Notes

After each milestone create:

```text
docs/learning/
```

Example:

```text
01-kafka.md

02-iceberg.md

03-flink.md

04-spark.md
```

Each note should contain:

```text
What problem does this technology solve?

How did we use it?

What did I misunderstand initially?

Important concepts

Failure behavior

What would happen at larger scale?

Alternatives

Interview questions I should now be able to answer
```

---

# 75. Example Kafka Learning Questions

After Kafka milestone I should be able to answer:

```text
What is a Kafka partition?

Why does partition count affect scalability?

What is an offset?

What is a consumer group?

What happens during consumer rebalance?

What ordering guarantees exist?

What is consumer lag?

Why can retries produce duplicates?

What does an idempotent producer solve?

What does it not solve?

Why isn't exactly-once trivial?
```

---

# 76. Example Flink Learning Questions

I should be able to answer:

```text
event time vs processing time?

what is a watermark?

what is stateful stream processing?

what is checkpointing?

how does Flink recover?

what is backpressure?

how does Kafka replay interact with Flink recovery?

what makes end-to-end exactly-once difficult?
```

---

# 77. Example Iceberg Learning Questions

I should be able to answer:

```text
What problem does Iceberg solve?

Why isn't Parquet alone sufficient?

What is an Iceberg snapshot?

What are manifests?

How does time travel work?

How can schema evolution avoid rewriting data?

What is hidden partitioning?

What is the small-file problem?

How does compaction work?

What is optimistic concurrency?
```

---

# 78. Example Spark Learning Questions

I should understand:

```text
driver

executor

task

stage

partition

shuffle

wide vs narrow transformation

data skew

predicate pushdown

partition pruning

why too many tiny partitions are expensive
```

---

# 79. Example Platform Engineering Questions

By the end I should answer:

```text
What is a data platform?

What is the difference between control plane and data plane?

Why have both Flink and Spark?

Why have both Iceberg and Snowflake?

Why use Kafka instead of sending directly to Snowflake?

Why use Airflow if Kafka already moves data?

Why use dbt if Spark can transform data?

How should a backfill work?

How should schema changes be rolled out safely?

How should failed records be handled?

How would you make a pipeline self-service?

How would you detect stale datasets?

How would you handle small files?

How would you debug consumer lag?

How would you scale ingestion?
```

---

# 80. Project Completion Demo

The final demonstration should start from zero.

Create dataset:

```bash
platform dataset create \
  --name agent-events \
  --partition-key trace_id
```

Start workload:

```bash
agent-generator \
  --rate 1000 \
  --duration 10m
```

Show:

```text
Kafka receiving events

Flink processing streams

Iceberg snapshots appearing

real-time metrics

Spark building historical dataset

Airflow orchestration

dbt models

Snowflake analytics
```

Then deliberately create many small files.

Show:

```bash
platform dataset status agent-events
```

Example:

```text
Dataset: agent-events

Kafka
  partitions             8
  lag                    19

Flink
  status                 healthy
  checkpoint age         9 sec

Iceberg
  files                  9,842
  small files            4,992
  average size           8.1 MB

Freshness
  latest event           1.8 sec

Maintenance
  compaction             RECOMMENDED
```

Then:

```bash
platform dataset compact agent-events
```

Show:

```text
Airflow job

Spark rewrite

Iceberg snapshot commit
```

Finally:

```text
files                  1,102

average file size      72 MB
```

This is the flagship demo.

---

# 81. What Makes the Project Strong

The project should not be presented as:

```text
"I used Kafka, Spark, Flink, Snowflake, dbt and Airflow."
```

It should be presented as:

```text
"I built a self-service AI data platform around a Go control plane.

The platform ingests distributed LLM and agent telemetry through Kafka, performs stateful real-time processing in Flink, persists canonical historical datasets in Iceberg, performs large-scale offline transformations and backfills with Spark, orchestrates offline workflows through Airflow, and serves curated analytical datasets through Snowflake and dbt.

I also built table-health monitoring and an automated Iceberg maintenance subsystem that detects small-file degradation and schedules Spark compaction jobs.

I benchmarked alternative analytical architectures using Trino, ClickHouse, Druid and Databricks, and used Airbyte to integrate transactional application data."
```

That is the narrative.

---

# 82. First Step

Do **not** begin with Kafka.

The first implementation session should only do this:

```text
repository
↓
Go module
↓
HTTP server
↓
POST /v1/events
↓
validate one event
↓
log accepted event
↓
unit test
```

No distributed systems yet.

Once that works:

```text
Go Event Generator
      ↓
Go Ingestion Gateway
```

Then introduce Kafka.

---

# 83. Initial Prompt for the Coding Agent

Use the following instruction after providing this design document:

> You are my senior data-platform engineer mentor. We are building the AI Data Platform described in this design document.
>
> I am building this primarily to learn data platform engineering, not merely to finish the project quickly.
>
> Therefore, do not generate the entire project for me. Work milestone by milestone and in very small increments.
>
> Before every meaningful implementation step, explain:
>
> 1. what problem we are solving;
> 2. where this component sits in the architecture;
> 3. the underlying data/distributed-systems concept;
> 4. what we are going to implement now;
> 5. what we are deliberately postponing.
>
> For important systems concepts, ask me to predict behavior before we run experiments.
>
> Every milestone should include normal-path testing and at least one intentional failure experiment where appropriate.
>
> Prefer simple implementations before abstractions.
>
> Do not introduce Kubernetes until the core platform works locally.
>
> Do not introduce exposure technologies until the core platform works.
>
> Follow the architecture in this design document unless there is a strong engineering reason to change it; if so, explain the tradeoff before making the change.
>
> Start with Milestone 0 only.
>
> For our first session, help me:
>
> - initialize the repository;
> - choose a clean Go project structure;
> - create a minimal HTTP service;
> - implement GET /healthz;
> - implement a minimal POST /v1/events endpoint;
> - define the first Event struct;
> - validate a minimal event;
> - write unit tests;
> - run everything locally.
>
> Do not add Kafka yet.
>
> Teach me each step rather than simply writing all the code.

---

# 84. Final Architecture Philosophy

The architecture is deliberately large.

The implementation path is deliberately small.

At any moment, only introduce one major new concept.

The progression should feel like:

```text
HTTP service

→ event ingestion

→ Kafka

→ reliable streaming

→ Iceberg storage

→ Flink stateful processing

→ Spark batch processing

→ Airflow orchestration

→ Snowflake/dbt analytics

→ Go control plane

→ platform reliability

→ automatic maintenance

→ observability

→ alternative-engine experiments

→ Kubernetes
```

The project succeeds if, at the end, the engineer can reason about why the architecture works, how it fails, and what tradeoffs each component introduces.

The project fails if it merely becomes a Docker Compose file containing fifteen technologies.