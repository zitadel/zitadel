# Identity Signals — Internal Design Document

> **Status: Preview** — This subsystem is under active development. APIs, storage format, and configuration may change without notice.

## Overview

Identity Signals provides identity-aware observability for ZITADEL administrators. It captures behavioral data from authentication flows and API operations, stores it in a columnar analytics engine (DuckLake), and exposes it through a query API and Console UI.

**This is NOT a general-purpose telemetry system.** It captures identity-scoped signals — who did what, when, from where — to support admin investigations, compliance, and future threat detection.

## Architecture

```
┌──────────────────────────────────────────────────────┐
│  Signal Sources                                       │
│  ┌─────────────────┐  ┌────────────────────────────┐ │
│  │ HTTP Interceptor │  │ Event Hook                 │ │
│  │ (requests)       │  │ (eventstore post-push)     │ │
│  └────────┬────────┘  └──────────┬─────────────────┘ │
└───────────┼──────────────────────┼───────────────────┘
            │ Signal               │ Signal
            ▼                      ▼
┌──────────────────────────────────────────────────────┐
│  Emitter                                              │
│  - Buffered channel (default 4096, fire-and-forget)   │
│  - Debouncer batches signals by time/size             │
│  - Drops signals when channel full (counter metric)   │
└────────────────────────┬─────────────────────────────┘
                         │ []Signal (batch)
                         ▼
┌──────────────────────────────────────────────────────┐
│  DuckLakeStore (SignalSink + SignalReader)             │
│  ┌─────────────┐  ┌──────────┐  ┌─────────────────┐ │
│  │ DuckDB      │  │ Postgres │  │ Parquet Files    │ │
│  │ (in-process)│  │ (catalog)│  │ (data, FS or S3)│ │
│  └─────────────┘  └──────────┘  └─────────────────┘ │
└────────────────────────┬─────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│  Signal Service (connectRPC)                          │
│  - ListSignals (search + pagination)                  │
│  - AggregateSignals (group-by, time-bucket, metrics)  │
│  - Permission: iam.read                               │
└────────────────────────┬─────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────┐
│  Console UI (Angular)                                 │
│  - Overview: dashboard with stat cards + charts       │
│  - Explore: ad-hoc aggregation queries                │
│  - Logs: filterable signal table with detail expand   │
│  - Activity: per-entity timeline with trace grouping  │
└──────────────────────────────────────────────────────┘
```

## Signal Streams

Two streams capture different aspects of identity activity:

| Stream     | Source              | Contains                                  | Default Retention |
|------------|---------------------|-------------------------------------------|--------------------|
| `requests` | HTTP interceptor    | API calls, IP, user agent, duration_ms    | 30 days            |
| `events`   | Eventstore hook     | Domain events, payload, aggregate context | 90 days            |

Each stream can be independently enabled/disabled and has its own retention policy.

## Storage Design

### DuckLake Architecture

DuckLake combines three components:

1. **DuckDB** (in-process) — Columnar analytics engine. Runs queries, handles Parquet I/O. Requires CGO (`//go:build cgo`; stubs in `ducklake_store_nocgo.go` for `CGO_ENABLED=0` builds).

2. **PostgreSQL** (catalog) — DuckLake extension stores table metadata in a dedicated `signals` schema. Created by `cmd/initialise/sql/11_signals.sql` during `zitadel init`.

3. **Parquet files** (data) — Signal data written as columnar Parquet files. Backends: local filesystem (`fs`) or S3-compatible storage (`s3`).

### Write Path

1. Signal arrives at emitter channel (non-blocking, fire-and-forget)
2. Debouncer accumulates batch (by time or size threshold)
3. Batch flushed to DuckLake via `INSERT INTO signals VALUES (?, ?, ...)`
4. DuckDB writes Parquet file to configured data path
5. Background compaction merges small files (configurable threshold, default 10 files)

### Read Path

1. Signal Service receives ListSignals/AggregateSignals RPC
2. Instance ID injected from auth context (tenant isolation)
3. DuckLakeStore builds SQL query with filters
4. DuckDB scans Parquet files via DuckLake catalog
5. Results returned as proto messages

### Instance Isolation

All queries are scoped by `instance_id` (injected server-side from the auth context). There is no way to query across instances.

## ID Extraction Strategy

Events use inconsistent JSON field names across aggregate types. The extraction priority in `event_hook.go`:

1. **Aggregate type prefix** — `user*` → aggregate ID is the user ID; `session*` → aggregate ID is the session ID
2. **JSON payload parsing** — Generic `map[string]json.RawMessage` with `firstStringField()` checking all known variants:
   - User ID: `userID`, `user_id`, `userId`, `hint_user_id`
   - Session ID: `sessionID`, `session_id`
   - Client ID: `clientID`, `clientId`, `client_id`
3. **Creator fallback** — `event.Creator()` returns the acting user from auth context. Used when payload is empty (e.g., `auth_request.code.exchanged`). Filtered to exclude system identifiers (`SYSTEM`).
4. **Organization** — Always from `Aggregate().ResourceOwner` (never empty).

## Security Model

- **API permission**: Both RPCs require `iam.read` (instance admin)
- **Route guard**: Console UI uses `authGuard` + `roleGuard` with `iam.read`
- **Self-exclusion**: The HTTP interceptor skips signal API calls (`/zitadel.signal.*`) to prevent self-recording loops
- **Instance isolation**: All queries filtered by `instance_id` from auth context

## Configuration Reference

Top-level key: `IdentitySignals` in `cmd/defaults.yaml`.

| Key | Default | Description |
|-----|---------|-------------|
| `Enabled` | `false` | Master switch for signal collection |
| `GeoCountryHeader` | `""` | HTTP header for country code (e.g., `CF-IPCountry`) |
| `Store.ChannelSize` | `4096` | Emitter channel buffer size |
| `Store.Debounce.MinFrequency` | `1s` | Max time between batch flushes |
| `Store.Debounce.MaxBulkSize` | `100` | Batch size threshold |
| `Store.DuckLake.Enabled` | `false` | Enable DuckLake storage backend |
| `Store.DuckLake.DataPath` | `/var/lib/zitadel/signals` | Parquet file root |
| `Store.DuckLake.MetadataSchema` | `signals` | PostgreSQL schema for catalog |
| `Store.DuckLake.Backend` | `fs` | `fs` or `s3` |
| `Store.DuckLake.FlushInterval` | `30s` | How often emitter flushes to Parquet |
| `Store.DuckLake.CompactionInterval` | `1h` | How often compaction runs |
| `Store.DuckLake.CompactionThreshold` | `10` | Min files to trigger compaction |
| `Streams.Requests.Retention` | `720h` | Request signal retention (30 days) |
| `Streams.Events.Retention` | `2160h` | Event signal retention (90 days) |
| `Retention.PruneInterval` | `6h` | How often pruning worker runs |

## File Map

| File | Purpose |
|------|---------|
| `config.go` | Configuration structs |
| `emitter.go` | Buffered signal emission + debouncing |
| `signal.go` | Signal struct and stream/outcome constants |
| `signal_interceptor.go` | HTTP/connectRPC middleware (request stream) |
| `event_hook.go` | Eventstore post-push hook (event stream) |
| `ducklake_store.go` | DuckLake storage (CGO build) |
| `ducklake_store_nocgo.go` | Stub for CGO_ENABLED=0 builds |
| `event_hook_test.go` | Tests for ID extraction and outcome classification |

## Known Limitations (Preview)

- **Single-writer DuckDB**: Only one process can write at a time. Not suitable for horizontally-scaled write paths without external coordination.
- **No real-time streaming**: Signals are batched; minimum latency is the flush interval (default 30s).
- **Empty user IDs**: Some events (`auth_request.code.exchanged`, `auth_request.succeeded`) have empty payloads. The creator fallback helps but system-initiated events still show no user.
- **No retention enforcement on S3**: Pruning currently only works for the `fs` backend.
- **Schema evolution**: Adding new Signal fields requires DuckLake schema migration. No automated migration tooling yet.
- **CGO dependency**: DuckDB requires CGO. Production builds use a separate Dockerfile (`Dockerfile.signals`) with `CGO_ENABLED=1`.

## Future Work

- Feature flag gating (per-instance enablement via Feature API)
- Threat detection engine (using signals as input to rule evaluation)
- Webhook/alert integration for anomaly detection
- Multi-writer support (DuckDB WAL sharing or external queue)
- Dashboard presets and saved queries
- Signal export (CSV, OTEL)
