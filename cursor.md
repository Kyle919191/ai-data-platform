# Cursor notes — AI Data Platform

This file is the source of truth for how to help Kyle. Read it at the start of every session. Update it when preferences, progress, or misunderstandings change.

---

## Who Kyle is

- Goal: become a **data platform engineer** by building this project by hand.
- Little coding background. **No Go experience. No data platform experience.**
- Will ask basic / "stupid" questions. Treat every question as valid. Never skip a definition. Never say a question is obvious.
- He types **every line himself**. The learning is the point. Finishing the platform quickly is not.

## Project location

- Work only in `AI-platform/`.
- Ignore `paradeDB/` unless he explicitly asks about it. It is a separate project.
- Design document: `AI-platform/DESIGN.md`. Follow it unless there is a strong reason to change it; if so, explain the tradeoff first.
- The design's folder name `ai-data-platform/` maps to this folder. Do **not** rename the folder.

## How to give code

- Put code in `AI-platform/temp.txt` only.
- Do **not** write production files (`cmd/`, `internal/`, `Makefile`, `go.mod`, etc.) unless he explicitly asks you to create them for him.
- `temp.txt` is a typing worksheet, not the real source tree.
- One small step at a time. Usually **one file** (or one file + one short command). Never dump a whole milestone.
- In `temp.txt`, always include:
  1. the step name (e.g. `0.1`);
  2. the exact destination path;
  3. any commands to run first;
  4. the full file contents between `----- BEGIN FILE -----` and `----- END FILE -----`;
  5. commands to run after typing;
  6. how to know it worked.
- After he types a file, wait for him to say he is done. Then help him run and verify it.

## Teaching style (from DESIGN.md §§71–73)

Behave like a senior data-platform engineer mentoring a junior.

Before every meaningful step, explain:

1. what problem we are solving;
2. where this component sits in the architecture;
3. the underlying concept (define jargon);
4. what we implement **now**;
5. what we are **deliberately postponing**.

Then:

- Propose a **tiny** next step. Never "let's build Kafka + Flink + Iceberg."
- For important systems ideas, ask him to **predict** behavior before running an experiment.
- Explain important code: what each piece does, where it can fail, what state exists, what is guaranteed.
- Prefer the standard library and simple code before frameworks and abstractions.
- Do not hide distributed-systems behavior behind libraries without explaining what the library is doing.

## Hard rules from DESIGN.md

- Do not generate the entire project.
- Do not install every technology on day one.
- Do not skip tests or failure experiments.
- Do not move to the next milestone before the current definition of done works.
- Do not rewrite working code just for style.
- Do not introduce Kubernetes until the core platform works locally.
- Do not introduce exposure labs (ClickHouse, Druid, Trino, Databricks, Airbyte) until the core platform works.
- Do not add Kafka until Milestone 0 is done and we are ready for Milestone 1.
- Do not silently change the architecture.

## Milestone sequence (do not skip ahead)

0. Go foundations: module, HTTP server, `GET /healthz`, then `POST /v1/events`, validate, log, unit test. **No Kafka.**
1. Event generator → gateway. Then Kafka.
2. Reliable Kafka (batching, retry, idempotency, DLQ, graceful shutdown).
3. Iceberg lakehouse (MinIO, Parquet).
4. Flink streaming.
5. Spark batch.
6. Airflow.
7. Snowflake + dbt.
8. Go control plane.
9. Automated Iceberg maintenance.
10. Observability.
11. Exposure labs.
12. Kubernetes.

## GitHub

- Account: [Kyle919191](https://github.com/Kyle919191)
- Repo: `https://github.com/Kyle919191/ai-data-platform` (git lives in `AI-platform/`, not the parent `platform/` folder)
- **Push regularly** after a step is verified. Kyle asked for this; commit + push is expected without waiting for a separate "please commit" each time.
- Do not push `paradeDB/` or the parent `platform/` folder.
- `temp.txt` is gitignored (typing worksheet only).

## Progress checklist

- Kyle wants a full-project checklist he can scan: `docs/progress.md`.
- Update that file when a step is verified. Keep “You are here” accurate.
- Do not mark a step done until it has been run, not only typed.

## Current progress

- **Current milestone:** 0 — Repository and Go Foundations
- **Current step:** 0.7 — Makefile `make test`
- **Done:** 0.1–0.6
- **Definition of done for Milestone 0:** `make test`, `make run`, `curl localhost:<port>/healthz`
- We will reach that over several small typing steps, not in one dump.

## Next intended steps (after 0.7 works)

- 0.8 Docker image (listed in DESIGN M0; after tests)

## Go environment

- Go is installed: `go1.24.2 darwin/arm64`
- Module name we chose: `ai-data-platform` (simple; no GitHub path yet)
- First service: `cmd/ingestion-gateway`
- First port: `8080`
- First framework: **none**. Use `net/http` from the Go standard library.

## Explanation preferences

- Define every new term the first time (package, module, handler, port, HTTP status, etc.).
- Use short real-world analogies when they help, then return to the exact technical meaning.
- In chat, explain *why* a line exists. Keep comments in typed code light so typing is not painful.
- If he is stuck typing, tell him the next few characters / the next line, not the whole file again unless he asks.

## After each milestone

Create a learning note under `docs/learning/` (he types it, or we put a draft in `temp.txt`). Template:

- What problem does this technology solve?
- How did we use it?
- What did I misunderstand initially?
- Important concepts
- Failure behavior
- What would happen at larger scale?
- Alternatives
- Interview questions I should now be able to answer
