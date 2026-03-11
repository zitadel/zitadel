# Design Doc: Tiered Signal Store for Risk Evaluation

**Status:** Draft (POC)  
**Authors:** ZITADEL Core Team  
**Date:** 2026-03-08

---

> **Package naming note**
>
> - `internal/signals` is the behavioral data plane: capture, persist, archive, and replay signals.
> - `internal/detection` is the decision plane: it builds context from those signals, runs rules/classifiers, and returns findings plus allow/block/challenge decisions.
> - Rules live in `internal/detection` because they define how ZITADEL detects suspicious behavior from signal history. They are decision logic, not signal transport or storage.
> - `internal/ratelimit`, `internal/captcha`, and `internal/llm` stay separate because they are reusable supporting capabilities consumed by detection.
> - Instance administrators manage detection policy through `zitadel.settings.v2.SettingsService` and the Console's **Instance Settings > Detection / Detection Rules** pages. Deployment concerns like the storage backend, Redis, captcha provider wiring, and LLM credentials remain startup configuration.

## 1. Problem Statement

ZITADEL's risk engine needs a complete picture of user and session activity to make
informed decisions — rate limiting, SLM/LLM classification, captcha challenges, and
anomaly detection. Today, the signal store (`internal/signals/store.go`) is **in-memory
only** (`MemoryStore`), which means:

- **Signals are lost on restart.** No persistence, no historical analysis.
- **No cross-instance correlation.** In multi-node deployments, each node sees only
  its own signals.
- **No archival.** There's no way to retain signals for audit or long-term pattern
  analysis.
- **Limited enrichment.** Only auth-flow signals are captured; API reads,
  notifications, and other operations are not part of the risk context.

Meanwhile, high-volume signals (HTTP access logs, API reads, notification sends) are
either discarded or sent to external systems via OTel — but they're **not queryable
by the risk engine**.

## 2. Goals

1. **Persist signals durably** without adding mandatory dependencies beyond
   PostgreSQL.
2. **Protect the transaction database** — signal writes must never degrade API
   latency.
3. **Enable a complete risk context** — auth events, API reads, notifications,
   session lifecycle — all in one queryable stream per user/session.
4. **Support tiered storage** — hot (in-memory / Redis) → warm (PostgreSQL) → cold
   (Parquet on FS/S3).
5. **Archive old signals efficiently** via River jobs, freeing PG storage.
6. **Follow existing patterns** — connector abstraction, debouncer, River workers.

## 3. Non-Goals

- Real-time streaming to external consumers (use OTel/instrumentation for that).
- Replacing the existing OTel instrumentation pipeline.
- Building a general-purpose analytics engine.
- Supporting non-PostgreSQL primary databases.
- Replacing `eventstore.events2` — domain events remain the audit log for state
  changes.

---

## 4. Current State

### 4.1 Signal Store (in-memory only)

```go
// internal/signals/store.go
type Store interface {
    Snapshot(ctx context.Context, signal Signal) (Snapshot, error)
    Save(ctx context.Context, signal Signal, findings []Finding) error
}
```

`MemoryStore` holds signals in `map[string][]RecordedSignal` keyed by `userID` and
`sessionID`. Signals are pruned by time window (`HistoryWindow` +
`ContextChangeWindow`) and per-user/session caps (`MaxSignalsPerUser`,
`MaxSignalsPerSession`).

### 4.2 Risk Context

The `RiskContext` struct (`internal/detection/context.go`) is built from a `Snapshot` and
provides counters (failure/success), delta flags (IP/UA/fingerprint changes),
cardinality (distinct IPs/countries), and behavioral signals (login velocity, hour of
day).

### 4.3 Overlap with `eventstore.events2`

The v3 eventstore already captures **domain state changes** as immutable events:

| Already in events2 | NOT in events2 (gaps) |
|--------------------|-----------------------|
| `session.user.checked` | HTTP access logs (path, status, timing) |
| `session.password.checked` | Read operations (who viewed what) |
| `session.totp.checked` | IP addresses / geolocation context |
| `user.password.changed` | Request velocity / behavioral patterns |
| `user.locked` / `user.unlocked` | Rate limit violations |
| `usergrant.created` / `removed` | Notification delivery status |
| `session.terminated` | Cross-operation correlation |

**Key distinction:** `events2` records **what changed** (domain mutations).
Signals record **what happened** (operational behavior). The risk engine needs both
dimensions, but they serve different purposes and have different volume/retention
characteristics.

Rather than making the risk engine query both `events2` and a signal table, **the
signal emitter fires a lightweight signal when relevant domain events occur**. This
keeps the risk engine reading from a single source (the signal table) without needing
to understand the eventstore query model.

### 4.4 Related Patterns Already in the Codebase

| Pattern | Location | Relevance |
|---------|----------|-----------|
| **Connector abstraction** | `internal/cache/connector/` | PG default, Redis optional, Memory fallback, Noop disabled |
| **Debouncer** | `internal/logstore/debouncer.go` | Generic `debouncer[T]` with time + size flush triggers |
| **River queue** | `internal/queue/` | PG-native async job queue with worker registration |
| **Instrumentation** | `backend/v3/instrumentation/` | OTel logs/traces/metrics with `StreamRisk` |
| **Unlogged PG tables** | `internal/cache/pg/` | Used for cache storage — no WAL overhead |

### 4.5 Existing Dependencies

| Dependency | Status |
|------------|--------|
| PostgreSQL | Required (always available) |
| Redis | Optional (connector + circuit breaker exist) |
| River | Available (`github.com/riverqueue/river`) |
| Minio S3 client | Available (`github.com/minio/minio-go/v7`) |
| DuckDB | **Not in go.mod** — new dependency for cold queries |
| Parquet (Go) | **Not in go.mod** — new dependency for archival |

---

## 5. Proposed Architecture

### 5.1 Overview

```
Signal Sources
(HTTP middleware, auth flow, API handlers, notification service, domain events)
         │
         ▼
┌─────────────────────────────────┐
│       Signal Emitter            │
│  (fire-and-forget, bounded)     │
│                                 │
│  signal → buffered channel ──┐  │
│     if full → drop + metric  │  │
└──────────────────────────────┼──┘
                               │
              ┌────────────────┴────────────────┐
              │     Background Goroutine         │
              │     (drains channel)             │
              └──────┬──────────┬───────────────┘
                     │          │
          ┌──────────┴──┐  ┌───┴───────────┐
          │  With Redis │  │ Without Redis  │
          │  (optional) │  │   (default)    │
          └──────┬──────┘  └───┬───────────┘
                 │             │
                 ▼             ▼
          ┌───────────┐  ┌──────────────────┐
          │  Redis    │  │  In-memory       │
          │  Stream   │  │  Ring Buffer     │
          │  (XADD)   │  │  + Debouncer     │
          └─────┬─────┘  └────────┬─────────┘
                │                 │
                │  River job      │  Batch INSERT
                │  (consumer)     │  (debounced)
                ▼                 ▼
          ┌───────────────────────────┐
          │  PostgreSQL Signal Table  │
          │  (unlogged, partitioned)  │
          │  ← Risk engine reads     │
          └─────────────┬─────────────┘
                        │
                        │  River periodic job
                        │  (archival)
                        ▼
          ┌───────────────────────────┐
          │  Cold Storage             │
          │  PG → Parquet             │
          │  FS or S3 (Minio)         │
          │  DuckDB for cold queries  │
          └───────────────────────────┘
```

### 5.2 Tier Responsibilities

| Tier | Storage | Retention | Purpose | Query Pattern |
|------|---------|-----------|---------|---------------|
| **Hot** | In-memory ring buffer or Redis Stream | Seconds–minutes | Decouple write path from PG | Not queried directly by risk engine |
| **Warm** | PostgreSQL (unlogged, partitioned) | Hours–days (configurable) | Risk engine reads, real-time evaluation | Indexed by `(instance_id, caller_id, timestamp)` and `(instance_id, session_id, timestamp)` |
| **Cold** | Parquet files on FS/S3 | Months–years | Historical analysis, audit, forensics | DuckDB with partition pruning |

---

## 6. One Table vs. One Table Per Stream

### 6.1 Decision: Single Table with `stream` Column

The instrumentation system defines 7 streams (`runtime`, `ready`, `request`,
`event_handler`, `queue`, `risk`, `event_pusher`). For signal storage, only a subset
is relevant (primarily `request` and `risk`). Two options were considered:

**Option A — Table per stream:** `signals_request`, `signals_risk`, `signals_audit`

- Pro: Schema tailored per stream, independent retention, easier to reason about volume.
- Con: Risk engine must JOIN/UNION across tables. Multiple archival jobs. Schema drift.

**Option B — Single table with `stream` column** ✅

- Pro: One index set, one archival job, one emitter. Risk engine reads one table.
  Different retention per stream handled by the archival job config.
- Con: Mixed volumes in one table (high-volume access logs alongside lower-volume risk signals).

**Why Option B wins:** The risk engine's primary query is "give me all signals for
caller X in time window Y" — this spans stream types. A user's access log entries,
auth events, and notification signals together form the behavioral picture. Splitting
by stream forces the risk engine to re-assemble what was naturally unified.

PG partitioning handles volume (partition by **time**, not by stream). The `stream`
column enables filtered queries when needed. Archival can apply different retention
per stream within the same River job.

### 6.2 Stream Types for Signals

| Stream | Description | Volume | Source |
|--------|-------------|--------|--------|
| `request` | HTTP/gRPC access logs | High | HTTP middleware |
| `auth` | Authentication flow events | Medium | Auth handlers, session commands |
| `account` | Account changes (from domain events) | Low | Event hooks on user/grant commands |
| `notification` | Notification lifecycle | Low | Notification service |

---

## 7. Signal Schema

### 7.1 Signal Struct (extended)

Every request in ZITADEL has an authenticated caller — even login/register flows use
the login UI's service account. There is no anonymous phase requiring back-fill.

```go
type Signal struct {
    // Identity (always present)
    InstanceID    string
    CallerID      string        // user ID or service account ID — always known
    SessionID     string        // set during auth flows
    FingerprintID string

    // Classification
    Stream        string        // "request", "auth", "account", "notification"
    Operation     string        // e.g., "login.started", "api.read", "notification.sent"
    Resource      string        // e.g., "users.list", "session.create"
    Outcome       Outcome       // success | failure | blocked | challenged

    // Context
    Timestamp     time.Time
    IP            string
    UserAgent     string

    // Tier 1 enrichment (from HTTP headers)
    AcceptLanguage string
    Country        string        // ISO 3166-1 alpha-2 (from GeoCountryHeader)
    ForwardedChain []string
    Referer        string
    SecFetchSite   string
    IsHTTPS        bool
}
```

### 7.2 PostgreSQL Table

```sql
CREATE UNLOGGED TABLE IF NOT EXISTS signals.signals (
    instance_id     TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp(),

    -- Identity
    caller_id       TEXT        NOT NULL,
    session_id      TEXT,
    fingerprint_id  TEXT,

    -- Classification
    stream          TEXT        NOT NULL,  -- 'request', 'auth', 'account', 'notification'
    operation       TEXT        NOT NULL,
    resource        TEXT,
    outcome         TEXT        NOT NULL,

    -- Context
    ip              INET,
    user_agent      TEXT,
    country         TEXT,
    metadata        JSONB       -- extensible: accept_language, referer, forwarded_chain, etc.
) PARTITION BY RANGE (created_at);

-- Risk engine query indices
CREATE INDEX idx_signals_caller
    ON signals.signals (instance_id, caller_id, created_at DESC);
CREATE INDEX idx_signals_session
    ON signals.signals (instance_id, session_id, created_at DESC)
    WHERE session_id IS NOT NULL;
-- Stream-filtered queries (e.g., "only auth signals for this user")
CREATE INDEX idx_signals_stream
    ON signals.signals (instance_id, caller_id, stream, created_at DESC);
```

**Why UNLOGGED:** Signal data is transient (hot/warm tier). It will be archived to
Parquet before being dropped. WAL overhead is unnecessary — if PG crashes, we lose
the current partition's signals (acceptable for risk heuristics). The cold tier
(Parquet) is the durable archive.

**Why a dedicated `signals` schema:** Keeps signal tables isolated from the
`eventstore`, `projections`, and `zitadel` schemas. Clean separation of concerns.

### 7.3 Partition Management

Rolling time partitions are created automatically:

```sql
-- Example: hourly partitions
CREATE TABLE signals.signals_2026030806 PARTITION OF signals.signals
    FOR VALUES FROM ('2026-03-08 06:00:00') TO ('2026-03-08 07:00:00');
```

A River periodic job creates future partitions and drops archived ones. This avoids
VACUUM entirely — `DROP TABLE` on a partition is instantaneous.

### 7.4 Redis Stream Keys

```
signals:{instance_id}    -- one stream per instance
```

Each entry contains signal fields as a flat hash. Capped with `MAXLEN ~ N`
(configurable, default 100000).

---

## 8. Write Path

### 8.1 Signal Emission (Request Path)

The request goroutine NEVER blocks on storage. All signal writes go through a bounded
channel:

```go
type Emitter struct {
    ch      chan Signal
    dropped atomic.Int64  // exposed as OTel metric: signal_store_dropped_total
}

func (e *Emitter) Emit(signal Signal) {
    select {
    case e.ch <- signal:
    default:
        e.dropped.Add(1)
    }
}
```

The channel size is configurable (default: 4096). If full, the signal is dropped and
a counter incremented — visible as a metric for capacity alerting.

### 8.2 Background Drain

A single goroutine drains the channel and writes to the configured hot tier using the
existing debouncer pattern (`internal/logstore/debouncer.go`):

- **Time trigger:** flush every `MinFrequency` (e.g., 1s)
- **Size trigger:** flush when `MaxBulkSize` reached (e.g., 100 signals)

### 8.3 Sink Implementations

```go
type SignalSink interface {
    WriteBatch(ctx context.Context, signals []Signal) error
}
```

| Sink | Behavior |
|------|----------|
| **PG (default)** | `COPY FROM` batch insert into the partitioned signal table. Uses a dedicated connection pool (2-3 conns, separate from transaction pool). |
| **Redis Stream** | `XADD` per signal with `MAXLEN ~ N`. Circuit breaker fallback to PG on failure. |

### 8.4 Redis → PG Bridge (River Consumer)

When Redis is configured as the hot tier, a River periodic job drains the stream into
PG:

```go
type SignalDrainArgs struct{}

func (SignalDrainArgs) Kind() string { return "signal.drain_redis" }

type SignalDrainWorker struct {
    river.WorkerDefaults[SignalDrainArgs]
    // ...
}

func (w *SignalDrainWorker) Work(ctx context.Context, job *river.Job[SignalDrainArgs]) error {
    // XREADGROUP from Redis Stream
    // Batch INSERT into PG signal table (COPY FROM)
    // XACK processed entries
}
```

Scheduled as a River periodic job (e.g., every 5s). `MaxWorkers: 1` to avoid
contention.

---

## 9. Read Path

### 9.1 Risk Engine Queries (Warm Tier — PG)

The `Store` interface is extended for query-oriented reads:

```go
type Store interface {
    // Existing
    Snapshot(ctx context.Context, signal Signal) (Snapshot, error)
    Save(ctx context.Context, signal Signal, findings []Finding) error

    // New: query-oriented reads
    CallerSignals(ctx context.Context, instanceID, callerID string, since time.Time, limit int) ([]RecordedSignal, error)
    SessionSignals(ctx context.Context, instanceID, sessionID string, since time.Time, limit int) ([]RecordedSignal, error)
}
```

The PG implementation uses index scans:

```sql
SELECT * FROM signals.signals
WHERE instance_id = $1 AND caller_id = $2 AND created_at > $3
ORDER BY created_at DESC
LIMIT $4;
```

Sub-millisecond for typical window sizes (50-100 signals over 30 minutes).

### 9.2 Enriched Risk Context

With the full signal stream, the `RiskContext` gains new dimensions:

```go
type RiskContext struct {
    // ... existing fields (counters, deltas, cardinality, behavioral) ...

    // New: cross-operation enrichment
    RecentAPIReads           int       // API read count in window
    RecentNotifications      int       // notifications sent in window
    PasswordChangeInWindow   bool      // password was changed recently
    MFAEnrolledInWindow      bool      // MFA was added recently
    DataAccessVelocity       float64   // API reads per minute
    DistinctResources        int       // unique resources accessed
}
```

This enables rules like:

```yaml
# Detect data exfiltration pattern
- id: data-exfil
  expr: 'DataAccessVelocity > 60 && DistinctResources > 10'
  engine: llm

# Captcha after suspicious pattern
- id: suspicious-pattern
  expr: 'FailureCount >= 3 && IPChanged && !MFAEnrolledInWindow'
  engine: captcha

# Notification flood
- id: notification-flood
  expr: 'RecentNotifications > 20'
  engine: rate_limit
  window: 1h
  limit: 20
  key: 'caller:{{.Current.CallerID}}'
```

### 9.3 Historical Queries (Cold Tier — DuckDB)

For historical analysis (not real-time risk), DuckDB queries Parquet directly:

```go
type ColdStore struct {
    db *sql.DB  // DuckDB connection (embedded)
}

func (s *ColdStore) QueryHistory(ctx context.Context, q HistoryQuery) ([]RecordedSignal, error) {
    // DuckDB reads Parquet from FS or S3:
    // SELECT * FROM read_parquet('s3://signals/instance=.../year=2026/month=03/*.parquet')
    // WHERE caller_id = ? AND created_at BETWEEN ? AND ?
}
```

DuckDB is embedded (Go driver: `github.com/marcboeker/go-duckdb`). It reads Parquet
from local FS and S3 natively.

---

## 10. Archival Path (River Periodic Job)

### 10.1 PG → Parquet Offload

A River periodic job archives old signal partitions:

```go
type SignalArchiveArgs struct{}

func (SignalArchiveArgs) Kind() string { return "signal.archive" }

func (w *SignalArchiveWorker) Work(ctx context.Context, job *river.Job[SignalArchiveArgs]) error {
    // 1. Identify partitions older than retention window
    // 2. SELECT * FROM signals_<partition> → write Parquet file
    // 3. Upload Parquet to FS/S3
    // 4. DROP TABLE signals_<partition> (instant, no VACUUM)
}
```

### 10.2 Per-Stream Retention

The archival job can apply different retention per stream. For example, keep `auth`
signals in PG for 7 days but `request` signals for only 24 hours:

```yaml
Archive:
  Retention:
    request: 24h
    auth: 168h      # 7 days
    account: 720h   # 30 days
    notification: 168h
```

Streams with shorter retention are archived (and their rows removed) earlier. Since
the table is partitioned by time (not stream), the archival job uses
`DELETE FROM ... WHERE stream = ? AND created_at < ?` for per-stream cleanup within a
partition, and `DROP TABLE` for fully-expired partitions.

### 10.3 Parquet Partitioning on Disk

```
signals/
├── instance=ins_abc123/
│   ├── year=2026/
│   │   └── month=03/
│   │       ├── day=07/
│   │       │   ├── hour=14.parquet
│   │       │   └── hour=15.parquet
│   │       └── day=08/
│   │           └── ...
```

### 10.4 Archive Storage Interface

```go
type ArchiveStorage interface {
    Write(ctx context.Context, path string, data io.Reader) error
}
```

Implementations:
- **FSStorage:** writes to local filesystem (configurable base path)
- **S3Storage:** writes to S3-compatible storage (uses existing `minio-go`)

---

## 11. Protection Mechanisms

### 11.1 Request Path Isolation

The signal write path is fully decoupled from the request transaction:

```
Request goroutine                    Background goroutine
─────────────────                    ────────────────────
handle request                       drain channel loop
  │                                    │
  ├─ business logic (PG txn)           ├─ debounce batch
  │                                    │
  ├─ emit signal (channel send)        ├─ COPY INTO signals table
  │   └─ non-blocking select           │   └─ separate conn pool
  │   └─ drop if full                  │
  │                                    │
  └─ return response                   └─ continue draining
```

The request goroutine never touches the signal table, never waits on Redis, never
blocks.

### 11.2 Per-Tier Safeguards

| Tier | Threat | Safeguard |
|------|--------|-----------|
| **Channel** | Backpressure | Fixed buffer (4096). Drop + metric on full. |
| **In-memory** | OOM | Ring buffer with per-user/session caps. Evicts oldest. |
| **Redis** | Memory exhaustion | `XADD MAXLEN ~ 100000` (capped stream). Circuit breaker (existing `redis.CBConfig`). Fallback to PG. |
| **PG** | Transaction DB impact | **UNLOGGED table** (no WAL). **Separate conn pool** (2-3 conns max). **Batch inserts** (COPY FROM). **Time partitions** (DROP, never VACUUM). |
| **Parquet/S3** | Disk/bandwidth | Periodic archival (not continuous). Configurable retention. ZSTD compression. |

### 11.3 Observability

All tiers emit metrics through the instrumentation system (`StreamRisk`):

| Metric | Type | Description |
|--------|------|-------------|
| `signal_store_emitted_total` | Counter | Signals successfully enqueued |
| `signal_store_dropped_total` | Counter | Signals dropped (channel full) |
| `signal_store_batch_size` | Histogram | Batch sizes flushed to PG |
| `signal_store_batch_latency_ms` | Histogram | Batch write duration |
| `signal_store_pg_partitions` | Gauge | Active PG partitions |
| `signal_store_archive_duration_ms` | Histogram | Archival job duration |
| `signal_store_redis_circuit_open` | Gauge | Redis circuit breaker state |

---

## 12. Configuration

Follows existing `cmd/defaults.yaml` patterns:

```yaml
Risk:
  Enabled: false
  # ... existing risk config ...

  SignalStore:
    # Channel buffer for fire-and-forget emission
    ChannelSize: 4096

    # Hot tier mode
    # "direct" = debouncer → PG batch insert (default, no Redis needed)
    # "redis"  = Redis Stream → River consumer → PG batch insert
    Mode: "direct"

    # Debouncer settings (for "direct" mode and PG batch writes)
    Debounce:
      MinFrequency: 1s
      MaxBulkSize: 100

    # Redis Stream settings (only when Mode: "redis")
    Redis:
      StreamMaxLen: 100000
      ConsumerGroup: "signal-drain"
      DrainInterval: 5s

    # Warm tier: PG signal table
    Postgres:
      MaxConns: 3
      PartitionInterval: 1h

    # Cold tier: Parquet archival
    Archive:
      Enabled: false
      Backend: "fs"          # "fs" or "s3"
      FSPath: "/var/lib/zitadel/signals"
      S3:
        Endpoint: ""
        Bucket: "zitadel-signals"
        AccessKey: ""
        SecretKey: ""
        UseSSL: true
      Compression: "zstd"    # "snappy", "zstd", "gzip", "none"
      Interval: 1h
      Retention:
        request: 24h
        auth: 168h
        account: 720h
        notification: 168h
```

---

## 13. Risk Engine Integration

### 13.1 Extended Evaluation Flow

```
Signal arrives
     │
     ▼
┌────────────────┐
│ Emit to store  │  (fire-and-forget → channel)
└────────┬───────┘
         │
         ▼
┌────────────────┐
│ Store.Snapshot  │  (read from PG: caller + session signals)
└────────┬───────┘
         │
         ▼
┌────────────────────┐
│ buildRiskContext() │  (enriched with API reads, notifications, etc.)
└────────┬───────────┘
         │
         ▼
┌────────────────────────────────────────────────────────┐
│                  Rule Chain                             │
│                                                        │
│  ┌──────────┐  ┌─────────────┐  ┌───────────┐        │
│  │  Block   │  │ Rate Limit  │  │  Captcha  │        │
│  │ (expr)   │  │ (fixed)     │  │ (new)     │        │
│  └────┬─────┘  └──────┬──────┘  └─────┬─────┘        │
│       └───────┬───────┘               │               │
│               ▼                       │               │
│         ┌───────────┐                 │               │
│         │  SLM/LLM  │                 │               │
│         │ (Ollama)   │                 │               │
│         └─────┬─────┘                 │               │
│               └───────┬───────────────┘               │
│                       ▼                               │
│                 ┌──────────┐                          │
│                 │ Decision │                          │
│                 └──────────┘                          │
└────────────────────────────────────────────────────────┘
         │
         ▼
  Allow / Block / Challenge(captcha) / RateLimit(429)
```

### 13.2 New Engine Type: Captcha

```go
const (
    EngineBlock     EngineType = "block"
    EngineRateLimit EngineType = "rate_limit"
    EngineLLM       EngineType = "llm"
    EngineLog       EngineType = "log"
    EngineCaptcha   EngineType = "captcha"     // NEW
)
```

When a captcha rule fires, the `Decision` includes a challenge requirement that the
auth flow must satisfy before proceeding.

---

## 14. Implementation Wiring

### Signal Emission Hook Points

| Hook Point | File | How |
|---|---|---|
| **V2 API requests** | `internal/api/api.go:234` | `signals.SignalConnectUnaryInterceptor(emitter)` in the Connect middleware chain. Fires after authorization, captures operation, caller, resource. Stream: `request`. |
| **Auth flow (session)** | `internal/command/session.go:463-478` | `recordSessionRisk()` emits signals on session create/set outcomes. Stream: `auth`. Called after `enforceSessionRisk()` for both allowed and blocked decisions. |
| **Signal interceptor** | `internal/signals/signal_interceptor.go` | Extracts HTTP headers (IP, UA, Accept-Language, Country, Sec-Fetch-Site, X-Forwarded-For) and emits fire-and-forget to the emitter channel. |

### Risk Enforcement

| Component | File | Behavior |
|---|---|---|
| **Risk evaluation** | `internal/command/session.go:430-461` | `enforceSessionRisk()` calls `Evaluate()` before each session mutation. |
| **Block decision** | `session.go:469` | Returns `PermissionDenied` (COMMAND-RISK0) with `OutcomeBlocked` signal. |
| **Challenge decision** | `session.go:458-467` | Returns `PreconditionFailed` (COMMAND-RISK1) with `OutcomeChallenged` signal. Client must present captcha. |
| **Fail-open** | `session.go:439-446` | When `FailOpen=true` and evaluation errors, logs warning and allows the request. |

### Service Initialization

```
cmd/start/start.go
  └── internal/command/command.go:StartCommands()
        └── risk.New(cfg, store, llm, pgDSN, redisClient)
              ├── DuckLakeStore (Parquet + PG catalog via DuckDB)
              ├── Emitter → DuckLakeStore
              ├── CompactionWorker (Parquet file merging)
              └── CaptchaVerifier (Turnstile/hCaptcha/reCAPTCHA)

cmd/start/start.go (worker registration)
  ├── risk.RegisterCompactionWorker(ctx, q, svc)
  │   (after q.Start())
  └── risk.StartCompactionSchedule(ctx, q, svc)
```

### Data Flow

```
HTTP Request
  │
  ├─[V2 API]─→ SignalConnectUnaryInterceptor ─→ Emitter.Emit(signal)
  │                                                    │
  ├─[Session]─→ enforceSessionRisk() ─→ Evaluate()    │
  │             recordSessionRisk() ─→ Emitter.Emit() │
  │                                                    ▼
  │                                          Bounded Channel (4096)
  │                                                    │
  │                                          ┌─────────┴─────────┐
  │                                     [Mode=pg]           [Mode=redis]
  │                                          │                    │
  │                                     PGStore.WriteBatch   GuardedSink
  │                                          │              (circuit breaker)
  │                                          │                    │
  │                                          ▼              RedisStreamSink
  │                                    signals.signals       (XADD MAXLEN ~)
  │                                    (UNLOGGED, partitioned)    │
  │                                          │              DrainWorker
  │                                          │              (XREADGROUP→PG)
  │                                          │                    │
  │                                          ▼                    ▼
  │                                    signals.signals ◄──────────┘
  │                                          │
  │                                   ArchiveWorker (periodic)
  │                                          │
  │                                   ┌──────┴──────┐
  │                              [Backend=fs]   [Backend=s3]
  │                                   │              │
  │                              Parquet files   S3/MinIO
  │                              (ZSTD compressed)
  └─────────────────────────────────────────────────────────
```

---

## 15. POC Phases

### Phase 1: PG Signal Table + Risk Engine Integration ✅ Implemented

1. Define the `Signal` struct extension — add `Stream`, `Resource`, `CallerID`.
2. Create the `signals` schema and partitioned table (migration).
3. Implement `PGStore` — satisfies the existing `Store` interface.
4. Implement the emitter — buffered channel + debouncer + batch COPY.
5. Wire signal emission from HTTP/gRPC middleware and auth flow handlers.
6. Emit lightweight signals on relevant domain events (password change, MFA enroll).
7. Extend `RiskContext` with cross-operation enrichment fields.
8. Add `Risk.SignalStore` config section in `defaults.yaml`.
9. Add partition management (create future, drop expired).

### Phase 2: Redis Hot Tier ✅ Implemented

10. Implement Redis Stream sink — `XADD` with `MAXLEN`, circuit breaker fallback.
11. Implement River drain worker — `XREADGROUP` → PG batch insert → `XACK`.
12. Add `Mode: "redis"` configuration toggle.

### Phase 3: Parquet Archival ✅ Implemented

13. Implement Parquet writer (using `parquet-go` or DuckDB `COPY TO`).
14. Implement archive storage — FS and S3 (Minio) backends.
15. Implement River archival worker — reads old PG partitions, writes Parquet, drops.
16. Implement DuckDB cold reader for historical queries.
17. Add per-stream retention configuration.

### Phase 4: Captcha Engine ✅ Implemented

18. Add `EngineCaptcha` rule type to the risk engine.
19. Integrate captcha challenge into the auth flow decision path.

### Current Limitations (POC)
- **S3 archive backend**: Config is parsed but falls back to FS storage. Minio client injection from static storage not yet wired.
- **DuckDB cold queries**: Not integrated. Cold data in Parquet is queryable via external tools (DuckDB CLI, Spark, pandas).
- **Captcha client-side**: `EngineCaptcha` produces challenge findings and the server returns `PreconditionFailed`. Client-side widget integration (Turnstile/hCaptcha/reCAPTCHA JavaScript) is not yet implemented in the login UI.
- **Redis signal store**: Requires the `cache` profile with Redis enabled. GuardedSink drops signals (with counter) when Redis is unavailable.

---

## 15a. Rate Limiter Architecture

### Multi-Backend Design

The `rate_limit` rule engine uses a `RateLimiterStore` interface with three backends,
following the same tiering pattern as the signal store. All three backends implement
the same **fixed-window** semantics.

```
┌──────────────────────────────────────────────────────────┐
│  Rule Engine  │  rate_limit rule matched                 │
│               │  suffix = "ip:{{.Current.IP}}"           │
│               ▼                                          │
│  ┌─────────────────────────────────────────────────────┐ │
│  │ RateLimiterStore.Check(ctx, storageKey, window, max)│ │
│  └──────────┬──────────────┬──────────────┬────────────┘ │
│             ▼              ▼              ▼              │
│    ┌─────────────┐  ┌────────────┐  ┌───────────────┐   │
│    │ Memory      │  │ Redis      │  │ PG (UNLOGGED) │   │
│    │ (default)   │  │ INCR+TTL   │  │ INSERT ON     │   │
│    │ sharded map │  │ Lua script │  │ CONFLICT      │   │
│    └─────────────┘  └────────────┘  └───────────────┘   │
└──────────────────────────────────────────────────────────┘
```

### RateLimiterStore Interface

```go
type RateLimiterStore interface {
    Check(ctx context.Context, key string, window time.Duration, max int) (count int, allowed bool)
    Prune(ctx context.Context)
}
```

### Backend Details

| Backend | Shared? | Latency | Dependencies | Best For |
|---------|---------|---------|--------------|----------|
| `memory` (default) | No — per-instance | ~µs | None | Single-node, dev, low-scale |
| `redis` | Yes — all instances | ~ms | Redis connection | Multi-node with Redis |
| `pg` | Yes — all instances | ~ms | PG (already required) | Multi-node without Redis |

**Memory** (`MemoryRateLimiter`):
- 64 FNV-sharded mutexes for minimal lock contention.
- Pruned by `maintenanceLoop` every 5 minutes.
- Counters lost on restart — acceptable for rate limiting.

**Redis** (`RedisRateLimiter`):
- Atomic Lua script: `INCR` + `EXPIRE` on first access.
- Keys auto-expire via TTL = window duration; Prune is a no-op.
- Redis stores canonical storage keys under `zitadel:ratelimit:<storageKey>`.
- If Redis is configured but unavailable at startup, ZITADEL downgrades to
  `memory` and logs the configured/effective backend modes.
- Runtime Redis call failures fail open (log warning, allow request).

**PG** (`PGRateLimiter`):
- `UNLOGGED` table `signals.rate_limit_counters` — no WAL overhead.
- Atomic `INSERT ... ON CONFLICT DO UPDATE` resets expired windows inline.
- Pruned by `maintenanceLoop` (`DELETE WHERE window expired`).
- If PG is configured but unavailable at startup, ZITADEL downgrades to
  `memory` and logs the configured/effective backend modes.
- Runtime PG query failures fail open.

### Key Templates and expr Integration

Rate limit rules define **what to limit by** using Go `text/template` on `RiskContext`.
The rendered template is treated as the **operator-controlled suffix**. ZITADEL
always prepends:

- `instance_id`
- `rule.ID`
- `window`

before persisting the counter key to memory / Redis / PG. This prevents
cross-tenant and cross-rule collisions.

Example suffix templates:

```yaml
rules:
  - id: ip-flood
    expr: 'DistinctIPs > 3'
    engine: rate_limit
    rate_limit:
      key: "ip:{{.Current.IP}}"           # per IP address
      window: 5m
      max: 100

  - id: user-auth-flood
    expr: 'FailureCount >= 3'
    engine: rate_limit
    rate_limit:
      key: "user:{{.Current.UserID}}"      # per user
      window: 10m
      max: 20

  - id: session-burst
    expr: 'SessionSignalCount > 50'
    engine: rate_limit
    rate_limit:
      key: "session:{{.Current.SessionID}}" # per session
      window: 1m
      max: 60

  - id: geo-anomaly
    expr: 'CountryChanged && DistinctCountries > 2'
    engine: rate_limit
    rate_limit:
      key: "geo:{{.Current.UserID}}:{{.Current.Country}}" # per user+country
      window: 1h
      max: 5
```

Available template fields match `RiskContext`: `Current.IP`, `Current.UserID`,
`Current.SessionID`, `Current.Country`, `Current.Operation`, `Current.UserAgent`,
`Current.FingerprintID`, and all computed fields.

### Configuration

```yaml
SystemDefaults:
  Risk:
    RateLimit:
      Mode: memory  # memory | redis | pg
```

Environment variable: `ZITADEL_SYSTEMDEFAULTS_RISK_RATELIMIT_MODE`

When `Mode=redis|pg` but the selected backend is unavailable at startup, ZITADEL
degrades to `memory` and logs both the configured and effective backend.

---

## 16. New Dependencies

| Dependency | Purpose | Phase |
|------------|---------|-------|
| `github.com/parquet-go/parquet-go` | Pure Go Parquet writer | Phase 3 |
| `github.com/marcboeker/go-duckdb` | Embedded DuckDB for cold queries | Phase 3 |

No new dependencies required for Phase 1 (PG) or Phase 2 (Redis).

---

## 17. Open Questions

1. **Partition granularity** — 1-hour vs. daily partitions? Hourly is cleaner for
   archival but creates more PG objects. With unlogged tables this should be fine.
2. **Cross-node signal visibility** — In multi-node deployments, PG is the shared
   store. The in-memory ring buffer is per-node. Should the risk engine always go to
   PG, or use a local in-memory cache with TTL?
3. **DuckDB CGO dependency** — `go-duckdb` uses CGO. Is this acceptable for the
   ZITADEL binary? Alternatives: shell out to `duckdb` CLI, or defer cold queries to
   a sidecar.
4. **Parquet schema evolution** — When new fields are added to `Signal`, Parquet
   supports adding columns natively, but we need a versioning/migration strategy.
5. **Per-instance retention** — Should signal retention be configurable per instance
   (tenant), or global?
6. **Captcha provider** — Which service(s) to integrate? hCaptcha, Turnstile,
   reCAPTCHA? Or a pluggable interface?

---

## 18. Signal Operation Taxonomy

| Category | Operation | Description |
|----------|-----------|-------------|
| **Auth** | `login.started` | Login flow initiated |
| | `password.verified` | Password check (success/failure) |
| | `mfa.prompted` | MFA challenge sent |
| | `mfa.verified` | MFA verification (success/failure) |
| | `passkey.verified` | Passkey/WebAuthn verification |
| | `session.created` | Session established |
| | `session.terminated` | Session ended (logout/expiry) |
| | `token.issued` | Access/refresh token issued |
| | `token.refreshed` | Token refresh |
| **Account** | `password.changed` | Password change |
| | `mfa.enrolled` | MFA method added |
| | `mfa.removed` | MFA method removed |
| | `email.changed` | Email address changed |
| | `phone.changed` | Phone number changed |
| **API** | `api.read` | Read API call |
| | `api.write` | Write API call |
| | `api.delete` | Delete API call |
| **Notification** | `notification.sent` | Notification dispatched |
| | `notification.clicked` | Notification link clicked |
| **Grant** | `grant.created` | User grant created |
| | `grant.removed` | User grant removed |

---

## 18a. DuckLake Signal Store Architecture (Active)

> **Note:** This is the active signal store architecture. Sections 1–18 below describe
> the original PG/Redis/Parquet tiered design which has been removed. The DuckLake
> architecture replaces all of it with a simpler pipeline: Go buffer → DuckLake
> (PG catalog + Parquet data files) → CompactionWorker.

### Motivation

The original 3-tier signal store (Redis hot → PG warm → Parquet cold) was complex:
3 workers (partition, drain, archive), PG UNLOGGED tables, Redis streams, and
filesystem/S3 Parquet archival. This put unnecessary write pressure on the OLTP
database and required Redis for optimal throughput.

### DuckLake Architecture

**DuckLake** (DuckDB extension, released May 2025) provides a lakehouse format
that stores catalog metadata in PostgreSQL and data as Parquet files. Since ZITADEL
already has PG, this requires **no new infrastructure**.

```
ZITADEL Pods
  ┌─────────────────────────────────────────┐
  │  Signal Emitter (fire-and-forget channel)│
  │  → Debouncer (batch + time-based flush) │
  │  → DuckLakeStore.WriteBatch()           │
  │    → go-duckdb Appender                 │
  │    → DuckLake INSERT → Parquet files    │
  └─────────────────────────────────────────┘
         │ data files              │ catalog metadata
         ▼                         ▼
  ┌──────────────┐         ┌──────────────┐
  │ Filesystem   │         │ PostgreSQL   │
  │ or S3        │         │ (MVCC)       │
  │ (Parquet)    │         │              │
  └──────────────┘         └──────────────┘
         │                         │
         └────────────┬────────────┘
                      ▼
              ┌──────────────┐
              │ DuckDB SQL   │
              │ (reads)      │
              │ Signals API  │
              └──────────────┘
```

### Key Design Decisions

1. **DuckLake over Apache Iceberg**: Iceberg requires a separate catalog service
   (e.g., Lakekeeper, a Rust sidecar). DuckLake uses PG natively — no new service.
   Data files are still standard Parquet.

2. **Multi-pod writer safety**: DuckLake coordinates concurrent writers via PG MVCC
   (optimistic concurrency). No leader election or distributed locks needed.

3. **Single worker**: One `CompactionWorker` (River job, hourly) replaces three
   workers (partition, drain, archive). It merges small Parquet files into larger
   time-aligned files.

4. **No Redis tier**: The DuckLake path doesn't use Redis. Buffer → periodic flush
   is sufficient since Parquet appends are cheap.

5. **Backward compatibility**: When `DuckLake.Enabled = false`, the existing PG
   signal store continues to work unchanged.

### Configuration

```yaml
SignalStore:
  Enabled: true
  DuckLake:
    Enabled: true
    DataPath: /var/lib/zitadel/signals  # or s3://bucket/signals/
    Backend: fs  # fs or s3
    FlushInterval: 30s
    CompactionInterval: 1h
    S3:
      Endpoint: ""
      Bucket: zitadel-signals
      AccessKey: ""
      SecretKey: ""
      UseSSL: true
```

### Query Capabilities

DuckLakeStore exposes two query methods for the Signals API:

- **SearchSignals**: Filtered, paginated signal retrieval with DuckDB SQL
- **AggregateSignals**: Group-by field or time_bucket with count/distinct_count

Both accept `SignalFilters` (instance_id, user_id, session_id, ip, operation,
stream, outcome, country, time range) and return structured results.

### go-duckdb Dependency

- Package: `github.com/duckdb/duckdb-go/v2` (v2.5.5)
- Requires CGO (adds ~50MB binary size)
- Uses `database/sql` interface for queries and DuckDB Appender API for bulk writes

---

## 19. Future: Multi-Engine Pipeline per Rule

### Problem

Today each rule has exactly one engine (`block`, `rate_limit`, `llm`, `log`, `captcha`).
Real-world scenarios need chained evaluation — e.g., "rate-limit first, then classify
with the SLM, then block or challenge based on the classification."

### Proposed Model

```yaml
- id: suspicious-login
  expr: 'FailureCount >= 3 && IPChanged'
  pipeline:
    - engine: rate_limit
      rate_limit:
        key: "user:{{.Current.UserID}}"
        window: 5m
        max: 10
    - engine: llm
      context_template: "User {{.Current.UserID}} from {{.Current.IP}} ..."
    - engine: captcha
      finding:
        name: suspicious_login_challenge
        message: "Please complete a captcha to continue"
```

**Semantics:**
- Stages execute in order. Each stage produces an optional `Finding`.
- If a stage produces a **blocking** finding (`Block: true`), the pipeline
  short-circuits — subsequent stages are skipped.
- Non-blocking findings (rate limit within budget, LLM classified "low risk") allow
  the pipeline to continue.
- All findings from executed stages are collected into the `Decision`.
- A single-engine rule (`engine: block`) is syntactic sugar for a one-stage pipeline.

**Open questions:**
- Should the LLM classification result be available to subsequent stages as a
  template variable (e.g., `{{.LLMClassification}}`)?
- Should pipeline stages support conditional execution (`when: ...`)?
- How does `engine: log` interact — always runs, never short-circuits?

### Backward Compatibility

The existing single `engine` + inline config is the one-stage form. Parsing checks
for `pipeline` first; if absent, falls back to the single-engine model. No migration
needed.

---

## 20. Future: Composite Rate Limit Keys

Rate limit keys already support Go `text/template` on `RiskContext`. Composite keys
that combine multiple dimensions are supported today:

```yaml
# Per IP + user combination
key: "ip:{{.Current.IP}}:user:{{.Current.UserID}}"

# Per user + country (detect geo-hopping)
key: "geo:{{.Current.UserID}}:{{.Current.Country}}"

# Per session (detect session abuse)
key: "session:{{.Current.SessionID}}"

# Per IP + operation (detect endpoint abuse)
key: "ip:{{.Current.IP}}:op:{{.Current.Operation}}"
```

The canonical storage key always prepends `instance_id`, `rule.ID`, and `window` to
prevent cross-tenant and cross-rule collisions. Available template fields match all
`RiskContext` fields: `Current.IP`, `Current.UserID`, `Current.SessionID`,
`Current.Country`, `Current.Operation`, `Current.UserAgent`, `Current.FingerprintID`,
`Current.Stream`, `Current.Resource`, and all computed fields.

---

## 21. Future: Rule Priority and Conflict Resolution

### Current Model

All rules evaluate independently. If any finding has `Block: true`, the request is
blocked ("most restrictive wins"). There is no priority ordering — rules are evaluated
in slice order, but all rules run regardless of prior matches.

### Proposed Enhancements

1. **Explicit priority field**: Rules with lower priority numbers evaluate first.
   When a blocking finding is produced, remaining rules are skipped (short-circuit).

2. **Weight-based scoring**: Instead of binary block/allow, each finding contributes
   a risk score. The final score is compared against a configurable threshold.
   ```yaml
   - id: ip-change
     expr: 'IPChanged'
     weight: 30
   - id: country-change
     expr: 'CountryChanged'
     weight: 40
   - id: failure-burst
     expr: 'FailureCount >= 5'
     weight: 50
   # Total >= 70 → challenge, >= 90 → block
   ```

3. **Override rules**: A rule can explicitly override another rule's finding:
   ```yaml
   - id: trusted-ip-allowlist
     expr: 'Current.IP in ["10.0.0.0/8"]'
     overrides: ["failure-burst", "ip-change"]
   ```

### Decision

Defer to a later phase. The current "most restrictive wins" model is correct for
security defaults. Priority/scoring adds flexibility but also complexity and potential
misconfiguration risk.

---

## 22. Future: Scope-Aware Rate Limits

Rate limit rules can target specific scopes using the `expr` field:

```yaml
# HTTP path-specific rate limit (API abuse)
- id: api-user-list-flood
  expr: 'Current.Stream == "request" && Current.Operation == "GET /v2/users"'
  engine: rate_limit
  rate_limit:
    key: "caller:{{.Current.CallerID}}"
    window: 1m
    max: 60

# Auth-flow rate limit (credential stuffing)
- id: auth-credential-stuffing
  expr: 'Current.Stream == "auth" && FailureCount >= 3'
  engine: rate_limit
  rate_limit:
    key: "ip:{{.Current.IP}}"
    window: 10m
    max: 20

# Notification rate limit (spam prevention)
- id: notification-flood
  expr: 'Current.Stream == "notification"'
  engine: rate_limit
  rate_limit:
    key: "caller:{{.Current.CallerID}}"
    window: 1h
    max: 100

# Account takeover pattern (MFA + password change)
- id: account-takeover
  expr: 'PasswordChangeInWindow && MFAEnrolledInWindow && CountryChanged'
  engine: block
  finding:
    name: account_takeover_pattern
    message: "Suspicious account changes from new location"
    block: true
```

The `Current.Stream` and `Current.Operation` fields are set by the signal emitter at
each hook point. Available streams: `request`, `auth`, `account`, `notification`. See
§18 for the operation taxonomy.

---

## 23. Future: Signals API (v2)

### Motivation

Instance administrators need visibility into the behavioral signal stream for:
- **Incident investigation**: "Show me all failed auth events for user X in the last
  24 hours."
- **Pattern discovery**: "Which users have the most failed logins this week?"
- **Compliance auditing**: "List all admin API calls from non-corporate IPs."
- **Proactive threat hunting**: "Find users with logins from > 3 countries in 1 hour."

Today, signals are only consumed by the detection engine internally. There is no
admin-facing API to query, filter, or aggregate them.

### Proposed API

A new `SignalService` in the v2 API, following existing ZITADEL API patterns:

```protobuf
service SignalService {
  // List signals with filtering and pagination.
  rpc ListSignals(ListSignalsRequest) returns (ListSignalsResponse);

  // Aggregate signals (e.g., count by user, by outcome, by country).
  rpc AggregateSignals(AggregateSignalsRequest) returns (AggregateSignalsResponse);

  // Get a single signal by ID (if signals have unique IDs — currently they don't,
  // so this may use a composite key: instance_id + created_at + caller_id).
  rpc GetSignal(GetSignalRequest) returns (GetSignalResponse);

  // Analyze signals with LLM forensics (see §24).
  rpc AnalyzeSignals(AnalyzeSignalsRequest) returns (AnalyzeSignalsResponse);
}
```

#### ListSignals

```protobuf
message ListSignalsRequest {
  // Pagination
  zitadel.object.v2.ListQuery query = 1;

  // Filters (all optional, AND-combined)
  optional string caller_id = 2;
  optional string session_id = 3;
  optional string stream = 4;           // "request", "auth", "account", "notification"
  optional string operation = 5;
  optional string outcome = 6;          // "success", "failure", "blocked", "challenged"
  optional string ip = 7;
  optional string country = 8;
  optional google.protobuf.Timestamp since = 9;
  optional google.protobuf.Timestamp until = 10;
}

message ListSignalsResponse {
  zitadel.object.v2.ListDetails details = 1;
  repeated Signal signals = 2;
}

message Signal {
  string caller_id = 1;
  string session_id = 2;
  string stream = 3;
  string operation = 4;
  string outcome = 5;
  string ip = 6;
  string user_agent = 7;
  string country = 8;
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Struct metadata = 10;  // extensible fields
}
```

#### AggregateSignals

Pre-defined aggregation queries for common admin workflows:

```protobuf
message AggregateSignalsRequest {
  // Time window
  google.protobuf.Timestamp since = 1;
  google.protobuf.Timestamp until = 2;

  // Aggregation type
  SignalAggregation aggregation = 3;

  // Optional filters (same as ListSignals)
  optional string stream = 4;
  optional string outcome = 5;
}

enum SignalAggregation {
  SIGNAL_AGGREGATION_UNSPECIFIED = 0;
  // Count signals grouped by caller_id (top users by activity)
  SIGNAL_AGGREGATION_BY_CALLER = 1;
  // Count signals grouped by outcome (success/failure distribution)
  SIGNAL_AGGREGATION_BY_OUTCOME = 2;
  // Count signals grouped by country
  SIGNAL_AGGREGATION_BY_COUNTRY = 3;
  // Count signals grouped by IP
  SIGNAL_AGGREGATION_BY_IP = 4;
  // Count signals grouped by operation
  SIGNAL_AGGREGATION_BY_OPERATION = 5;
  // Failed auth events grouped by caller (credential stuffing detection)
  SIGNAL_AGGREGATION_FAILED_BY_CALLER = 6;
}

message AggregateSignalsResponse {
  repeated SignalBucket buckets = 1;
}

message SignalBucket {
  string key = 1;      // the grouped value (caller_id, country, IP, etc.)
  int64 count = 2;
  // Optional: additional context per bucket
  optional google.protobuf.Timestamp first_seen = 3;
  optional google.protobuf.Timestamp last_seen = 4;
}
```

### Implementation Notes

- Reads from the `signals.signals` PG table (warm tier).
- Uses existing PG indices (`idx_signals_caller`, `idx_signals_session`,
  `idx_signals_stream`).
- Aggregation queries use `GROUP BY` with index scans — sub-second for typical
  windows (24h, 7d) on partitioned tables.
- Cold-tier queries (Parquet via DuckDB) are a future extension — the API returns
  warm-tier data initially.
- Authorization: requires IAM_OWNER role (instance admin). Scoped by instance_id
  automatically via `authz.GetInstance(ctx)`.

---

## 24. Future: LLM Forensics

### Motivation

Beyond real-time risk classification (§13.2), the LLM/SLM can be used for
**post-hoc forensic analysis** of signal history. An admin investigating an incident
should be able to ask: "Summarize the risk profile of user X over the last 7 days" or
"What patterns do you see in these failed login attempts?"

### Proposed Workflow

```
Admin (Console or API)
     │
     ├─ "Analyze signals for user X since 2026-03-01"
     │
     ▼
┌────────────────────┐
│  SignalService      │
│  .AnalyzeSignals() │
└────────┬───────────┘
         │
         ├─ Query signals from PG (warm) or DuckDB (cold)
         │
         ├─ Build forensic prompt:
         │    system: "You are a security analyst reviewing user activity..."
         │    user:   [serialized signal history, aggregations, context]
         │
         ├─ Call LLM (same Ollama endpoint, but with larger context window
         │   and higher token limit than real-time classification)
         │
         └─ Return structured analysis:
              - Risk assessment (low/medium/high)
              - Key observations (unusual patterns)
              - Timeline of suspicious events
              - Recommended actions
```

### API

```protobuf
message AnalyzeSignalsRequest {
  string caller_id = 1;
  google.protobuf.Timestamp since = 2;
  google.protobuf.Timestamp until = 3;
  optional string question = 4;  // freeform question, e.g. "Why were there so many failures?"
}

message AnalyzeSignalsResponse {
  string risk_level = 1;          // "low", "medium", "high"
  float confidence = 2;
  string summary = 3;             // natural language summary
  repeated string observations = 4;
  repeated string recommendations = 5;
}
```

### Implementation Notes

- Reuses the existing `llm.LLMClient` interface and Ollama backend.
- Forensic prompts are larger than real-time (100s of signals vs. 8). Use a larger
  model or longer context window when configured.
- Rate-limit forensic queries (expensive). One concurrent analysis per instance.
- Results are NOT cached — each analysis is fresh.
- Can be triggered from Console UI ("Analyze" button on signal list view).

---

## 25. Future: Scheduled Background Pattern Detection

### Motivation

Some threat patterns only emerge over time and can't be detected per-request:
- Slow credential stuffing (1 attempt per minute across hours)
- Account enumeration across many IPs
- Privilege escalation sequences (grant changes → data access spikes)
- Dormant account reactivation in bulk

### Proposed Architecture

River periodic jobs that run aggregation queries and produce **alerts** (§26):

```go
type PatternDetectionArgs struct {
    PatternID string  // e.g., "slow-credential-stuffing"
}

func (PatternDetectionArgs) Kind() string { return "detection.pattern_scan" }

type PatternDetectionWorker struct {
    river.WorkerDefaults[PatternDetectionArgs]
    db       *sql.DB
    llm      llm.LLMClient  // optional, for classification
    alerter  AlertSink       // writes to alert storage
    patterns []PatternDef
}
```

### Pattern Definitions

Patterns are defined as SQL aggregation queries + threshold checks:

```yaml
PatternDetection:
  Enabled: false
  Patterns:
    - id: slow-credential-stuffing
      description: "Detect slow-rate credential stuffing across long windows"
      query: |
        SELECT caller_id, COUNT(*) as failures, COUNT(DISTINCT ip) as ips
        FROM signals.signals
        WHERE stream = 'auth' AND outcome = 'failure'
          AND created_at > NOW() - INTERVAL '6 hours'
          AND instance_id = $1
        GROUP BY caller_id
        HAVING COUNT(*) >= 10 AND COUNT(DISTINCT ip) >= 3
      schedule: "@every 15m"
      alert:
        severity: high
        name: slow_credential_stuffing
        message: "User {{.caller_id}}: {{.failures}} failures from {{.ips}} IPs in 6h"

    - id: bulk-account-enumeration
      description: "Detect enumeration attempts across many usernames"
      query: |
        SELECT ip, COUNT(DISTINCT caller_id) as users, COUNT(*) as attempts
        FROM signals.signals
        WHERE stream = 'auth' AND outcome = 'failure'
          AND created_at > NOW() - INTERVAL '1 hour'
          AND instance_id = $1
        GROUP BY ip
        HAVING COUNT(DISTINCT caller_id) >= 20
      schedule: "@every 5m"
      alert:
        severity: critical
        name: bulk_enumeration
        message: "IP {{.ip}}: {{.attempts}} attempts against {{.users}} users in 1h"

    - id: privilege-escalation-sequence
      description: "Grant changes followed by unusual data access"
      query: |
        SELECT s.caller_id,
               COUNT(CASE WHEN s.stream = 'account' THEN 1 END) as changes,
               COUNT(CASE WHEN s.stream = 'request' THEN 1 END) as reads
        FROM signals.signals s
        WHERE s.created_at > NOW() - INTERVAL '2 hours'
          AND s.instance_id = $1
        GROUP BY s.caller_id
        HAVING COUNT(CASE WHEN s.stream = 'account' THEN 1 END) >= 2
           AND COUNT(CASE WHEN s.stream = 'request' THEN 1 END) >= 50
      schedule: "@every 30m"
      alert:
        severity: medium
        name: privilege_escalation
        message: "User {{.caller_id}}: {{.changes}} account changes + {{.reads}} API reads in 2h"
```

### LLM-Enhanced Pattern Detection

For patterns that are hard to express as SQL thresholds, the background job can send
aggregated data to the LLM for classification:

```yaml
    - id: anomalous-behavior-llm
      description: "LLM-based anomaly detection on daily activity summaries"
      query: |
        SELECT caller_id,
               COUNT(*) as total, COUNT(DISTINCT ip) as ips,
               COUNT(DISTINCT country) as countries,
               COUNT(CASE WHEN outcome = 'failure' THEN 1 END) as failures,
               array_agg(DISTINCT operation) as operations
        FROM signals.signals
        WHERE created_at > NOW() - INTERVAL '24 hours'
          AND instance_id = $1
        GROUP BY caller_id
        HAVING COUNT(*) >= 20
      schedule: "@daily"
      engine: llm
      context_template: |
        User {{.caller_id}} activity summary (24h):
        Total signals: {{.total}}, IPs: {{.ips}}, Countries: {{.countries}}
        Failures: {{.failures}}, Operations: {{.operations}}
        Assess risk level and explain any anomalies.
      alert:
        severity: auto  # determined by LLM classification
```

---

## 26. Future: Alert Storage and External Integrations

### Motivation

Findings from real-time detection rules, background pattern scans, and LLM forensic
analysis need to be:
1. **Persisted** — for audit trail and incident timeline reconstruction.
2. **Forwarded** — to external incident management and SIEM systems.
3. **Queryable** — by admins through the Console and API.

Today, findings are logged (via structured logging to stdout) and recorded inline on
signal rows, but there is no dedicated alert storage or forwarding mechanism.

### Alert Model

```go
type Alert struct {
    ID          string          // unique alert ID (ULID or UUID)
    InstanceID  string
    CreatedAt   time.Time

    // Source
    Source      AlertSource     // "rule", "pattern", "forensic", "manual"
    SourceID    string          // rule ID, pattern ID, or forensic analysis ID
    RuleID      string          // when Source == "rule"

    // Subject
    CallerID    string          // affected user/service account
    SessionID   string          // when applicable
    IP          string          // source IP when applicable

    // Classification
    Severity    AlertSeverity   // critical, high, medium, low, info
    Name        string          // machine-readable alert type (e.g., "credential_stuffing")
    Message     string          // human-readable description
    Confidence  float64         // 0-1, from LLM or rule

    // State
    Status      AlertStatus     // open, acknowledged, resolved, false_positive
    ResolvedBy  string          // user ID who resolved
    ResolvedAt  *time.Time

    // Context
    Findings    []Finding       // underlying findings that triggered the alert
    Metadata    map[string]any  // extensible context (signal counts, IPs, etc.)
}

type AlertSource string
const (
    AlertSourceRule      AlertSource = "rule"
    AlertSourcePattern   AlertSource = "pattern"
    AlertSourceForensic  AlertSource = "forensic"
    AlertSourceManual    AlertSource = "manual"
)

type AlertSeverity string
const (
    AlertSeverityCritical AlertSeverity = "critical"
    AlertSeverityHigh     AlertSeverity = "high"
    AlertSeverityMedium   AlertSeverity = "medium"
    AlertSeverityLow      AlertSeverity = "low"
    AlertSeverityInfo     AlertSeverity = "info"
)

type AlertStatus string
const (
    AlertStatusOpen          AlertStatus = "open"
    AlertStatusAcknowledged  AlertStatus = "acknowledged"
    AlertStatusResolved      AlertStatus = "resolved"
    AlertStatusFalsePositive AlertStatus = "false_positive"
)
```

### Storage

Alerts are stored in a regular PostgreSQL table (not UNLOGGED — alerts are important):

```sql
CREATE TABLE IF NOT EXISTS projections.alerts (
    id              TEXT        PRIMARY KEY,
    instance_id     TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    source          TEXT        NOT NULL,
    source_id       TEXT,
    rule_id         TEXT,

    caller_id       TEXT,
    session_id      TEXT,
    ip              INET,

    severity        TEXT        NOT NULL,
    name            TEXT        NOT NULL,
    message         TEXT,
    confidence      FLOAT,

    status          TEXT        NOT NULL DEFAULT 'open',
    resolved_by     TEXT,
    resolved_at     TIMESTAMPTZ,

    findings        JSONB,
    metadata        JSONB
);

CREATE INDEX idx_alerts_instance_status
    ON projections.alerts (instance_id, status, created_at DESC);
CREATE INDEX idx_alerts_instance_caller
    ON projections.alerts (instance_id, caller_id, created_at DESC);
CREATE INDEX idx_alerts_instance_severity
    ON projections.alerts (instance_id, severity, created_at DESC);
```

### Alert API (v2)

```protobuf
service AlertService {
  rpc ListAlerts(ListAlertsRequest) returns (ListAlertsResponse);
  rpc GetAlert(GetAlertRequest) returns (GetAlertResponse);
  rpc UpdateAlertStatus(UpdateAlertStatusRequest) returns (UpdateAlertStatusResponse);
}
```

### External Integrations (Webhook / SIEM Forwarding)

Alerts are forwarded to external systems via a pluggable `AlertSink` interface:

```go
type AlertSink interface {
    Send(ctx context.Context, alert Alert) error
}
```

#### Planned Integrations

| Integration | Transport | Notes |
|-------------|-----------|-------|
| **Webhook** | HTTP POST (JSON) | Generic. Configurable URL, headers, retry policy. Works with any system that accepts webhooks. |
| **incident.io** | REST API | Create incidents with severity mapping. Auto-resolve when alert status changes. |
| **Datadog** | Events API or Log Ingestion | Forward as Datadog events with tags (instance, severity, rule). Queryable in Datadog SIEM. |
| **PagerDuty** | Events API v2 | Create alerts/incidents for critical/high severity. Auto-resolve. |
| **Splunk** | HEC (HTTP Event Collector) | Forward as structured events. Index by instance, severity. |
| **Slack / Teams** | Webhook | Notification-only. Link back to ZITADEL Console alert detail page. |
| **Syslog** | RFC 5424 (TCP/UDP) | Forward to on-prem SIEM (QRadar, ArcSight, etc.). |
| **OTel Logs** | OTLP exporter | Emit alerts as OTel log records into the existing instrumentation pipeline. Zero-config if OTel is already configured. |

#### Configuration

```yaml
Risk:
  Alerts:
    Enabled: true
    Sinks:
      - type: webhook
        url: "https://hooks.example.com/zitadel-alerts"
        headers:
          Authorization: "Bearer ${ALERT_WEBHOOK_TOKEN}"
        retry:
          maxAttempts: 3
          backoff: 5s
        filter:
          minSeverity: medium  # only forward medium+ alerts

      - type: datadog
        apiKey: "${DD_API_KEY}"
        site: "datadoghq.eu"
        tags: ["env:production", "service:zitadel"]
        filter:
          minSeverity: high

      - type: otel
        # Uses existing OTel exporter config — no additional settings.
        filter:
          minSeverity: low  # forward everything
```

### Alert Lifecycle

```
Detection Rule / Pattern Scan / Forensic Analysis
     │
     ▼
┌─────────────┐
│ Create Alert │ → INSERT INTO projections.alerts
└──────┬──────┘
       │
       ├─── Fan-out to AlertSinks (async, fire-and-forget)
       │       ├── Webhook → HTTP POST
       │       ├── Datadog → Events API
       │       ├── OTel → OTLP exporter
       │       └── ...
       │
       ├─── Console notification (via existing notification channel)
       │
       └─── Alert visible in Console UI:
              ├── Alert list (filterable by severity, status, user)
              ├── Alert detail (findings, signals timeline, metadata)
              ├── Actions: Acknowledge / Resolve / Mark False Positive
              └── "Analyze" button → triggers LLM forensics (§24)
```

### Deduplication

To avoid alert floods, alerts are deduplicated by `(instance_id, name, caller_id)`
within a configurable window (default: 1 hour). If an alert with the same key already
exists in `open` status, the existing alert's `metadata` is updated with the new
finding count, but no new alert row is created and no external notification is sent.

---

## 27. Implemented: FindingSink Contract

### Current State (POC)

The detection system produces `Finding` values during rule evaluation. Currently,
findings flow through two internal paths:

1. **FindingRecorder** — persists findings as JSON on the originating signal row
   (via `AppendFindings` on the DuckLake store). This makes findings queryable
   through DuckDB JSON functions.

2. **Signal emission** — each non-LLM rule match emits a detection-stream signal
   via the `signalEmitter` interface, creating cross-stream correlation entries.

### FindingSink Interface

For future external integrations (webhook, SIEM, OpenTelemetry), the `FindingSink`
interface defines the forwarding contract:

```go
type FindingSink interface {
    Forward(ctx context.Context, signal signals.Signal, findings []Finding) error
}
```

Design constraints:
- Sinks MUST NOT block the detection evaluation path
- Errors are logged but do not affect the detection decision
- `MultiSink` fans out to multiple sinks; first error is returned but all sinks run
- The originating signal is provided for correlation context (instance, user, session, trace)

### Relationship to §26 (Alert Storage)

The `FindingSink` is the forwarding layer — it pushes findings outward. §26's
proposed `Alert` model is a persistence layer — it stores findings as durable,
lifecycle-managed entities. When §26 is implemented:

- `AlertSink` would implement `FindingSink`, converting findings into alert rows
- External sinks (webhook, SIEM) would also implement `FindingSink`
- `MultiSink` would compose them: `NewMultiSink(alertSink, webhookSink, siemSink)`

### Vocabulary Alignment

| Term | Role | Layer |
|------|------|-------|
| **Finding** | Analytic result from a rule evaluation | Domain |
| **FindingRecorder** | Persists findings on signal rows (internal) | Storage |
| **FindingSink** | Forwards findings to external systems | Integration |
| **Alert** (§26, future) | Durable, lifecycle-managed finding | Persistence |

### What Actions Are NOT

Actions (`block`, `rate_limit`, `llm`, `log`, `captcha`) are rule-local
processing stages — they determine what happens when a rule matches. They are
not sinks. The distinction:

- **Action**: synchronous, inline, per-rule — "what to do when this rule fires"
- **Sink**: asynchronous, post-evaluation, cross-cutting — "where to send all findings"

This separation allows actions to focus on enforcement (block a request, rate-limit
an IP) while sinks handle observability and integration (alert admin, push to SIEM).

---
