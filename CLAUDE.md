# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-ted1k is a Go-based energy data synchronization system. It ingests real-time energy measurements (timestamp + watt) from a NATS message broker into TimescaleDB, and can pump/sync data between multiple backends (PostgreSQL, MySQL, JSONL files, IPFS).

## Build and Test Commands

```bash
# Run all tests
go test ./...

# Run tests verbose
go test -v ./...

# Run tests for a single package
go test -v ./postgres/
go test -v ./merge/

# Build main binaries
go build ./cmd/subscribe/
go build ./cmd/pump/

# Docker
docker compose build --pull
docker compose up -d
docker compose logs -f subscribe
```

## Architecture

### Core Abstraction

All data flow is built around two interfaces in `types/types.go`:

- `EntryReader` — produces `[]Entry` via a channel (`Read() <-chan []Entry`)
- `EntryWriter` — consumes `[]Entry` from a channel (`Write(src <-chan []Entry) (int, error)`)

The `Entry` struct is simply `{Stamp time.Time, Watt int}`.

### Backend Packages

Each backend implements the Reader/Writer interfaces:

- **postgres/** — TimescaleDB/PostgreSQL with pgx. Writer uses CopyFrom with automatic fallback to multi-value INSERT (ON CONFLICT DO NOTHING) on duplicate key errors. Batch size: 10,000.
- **mysql/** — MySQL backend using go-sql-driver.
- **jsonl/** — JSON Lines files organized by time grain (year/month/day/hour) under `./data/jsonl/`. Uses EasyJSON for fast marshaling.
- **ipfs/** — IPFS-backed storage with same grain structure.
- **ephemeral/** — In-memory synthetic data generator (31M points) and no-op writer for testing.

### Entry Points (`cmd/`)

- **subscribe** — NATS subscriber → TimescaleDB. Listens on `im.qcic.heartbeat`, filters for `capture.ted1k` host, queues entries for async DB insertion with exponential backoff retry.
- **pump** — Reads from one backend (typically MySQL at `darwin.imetrical.com`), writes to another (typically local TimescaleDB). Flags: `-since` (duration, default 100 days), `-skip-copy-from`.
- **dumpjsonl** — Export from Postgres to JSONL files.

### Supporting Packages

- **merge/** — Compares two sorted Entry streams, classifying differences (Equal, Conflict, MissingInA/B).
- **progress/** — Channel passthrough monitor logging throughput and gap detection.
- **timer/** — Elapsed time and rate formatting (k/s, M/s).
- **logsetup/** — Custom log formatter with ISO 8601 millisecond timestamps.
- **iterator/** — Entry stream iteration interfaces.

## Database

Table `watt` in TimescaleDB with schema `(stamp TIMESTAMPTZ PRIMARY KEY, watt INTEGER NOT NULL DEFAULT 0)`, configured as a TimescaleDB hypertable on `stamp`. Unique constraint on stamp prevents duplicates.

## Environment

- `PGCONN` — PostgreSQL connection string (default: `postgres://postgres:secret@0.0.0.0:5432/ted`)
- `PG.env` / `MYSQL.env` — Docker service credentials
- NATS server: `nats://nats.ts.imetrical.com:4222`

## Key Patterns

- Channel-based concurrency: readers produce on channels, writers consume, with progress monitors as passthrough.
- EasyJSON code generation: `types/types_easyjson.go` is generated — regenerate with `easyjson` tool if `Entry` struct changes.
- The postgres writer's CopyFrom-to-INSERT fallback is intentional: CopyFrom is faster for clean inserts but can't handle duplicates.
