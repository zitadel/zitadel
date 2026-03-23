package signals

// PREVIEW: Identity Signals is a preview feature. APIs, storage format,
// and configuration may change between releases without notice.

import (
	"fmt"
	"strings"
	"time"
)

// IdentitySignalsConfig is the top-level configuration for the Identity Signals
// subsystem. It sits alongside Eventstore and Projections as a core domain feature.
type IdentitySignalsConfig struct {
	// Enabled activates identity signal collection and storage.
	Enabled bool
	// GeoCountryHeader is the HTTP header used to extract the client's country
	// code (e.g. "CF-IPCountry", "X-Vercel-IP-Country").
	GeoCountryHeader string
	// Store configures signal buffering and persistence.
	Store StoreConfig
	// Streams configures per-stream enablement and retention.
	Streams StreamsConfig
	// Retention configures the background pruning worker.
	Retention RetentionConfig
}

// Validate checks the configuration for consistency.
func (c IdentitySignalsConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if !c.Store.DuckLake.Enabled {
		return fmt.Errorf("identity signals requires Store.DuckLake.Enabled=true when IdentitySignals.Enabled=true")
	}
	return c.Store.Validate()
}

// StoreConfig configures signal buffering and the DuckLake backend.
type StoreConfig struct {
	// ChannelSize is the buffer size for the fire-and-forget emission channel.
	// Signals are dropped when the channel is full.
	ChannelSize int
	// Debounce controls batching of signal writes.
	Debounce DebouncerConfig
	// DuckLake configures the DuckLake-based signal store.
	DuckLake DuckLakeConfig
}

// Validate checks the store configuration.
func (c StoreConfig) Validate() error {
	if c.Debounce.MinFrequency <= 0 {
		return fmt.Errorf("store.debounce.min_frequency must be > 0")
	}
	if c.Debounce.MaxBulkSize == 0 {
		return fmt.Errorf("store.debounce.max_bulk_size must be > 0")
	}
	return c.DuckLake.Validate()
}

// DebouncerConfig controls how signals are batched before writing.
type DebouncerConfig struct {
	// MinFrequency is the maximum time between flushes.
	MinFrequency time.Duration
	// MaxBulkSize is the maximum batch size before a flush is triggered.
	MaxBulkSize uint
}

// StreamsConfig configures per-stream enablement and retention policies.
type StreamsConfig struct {
	Requests StreamConfig
	Events   StreamConfig
}

// StreamConfig controls a single signal stream.
type StreamConfig struct {
	// Enabled activates signal collection for this stream.
	Enabled bool
	// Retention is how long signals are kept before the pruning worker deletes them.
	// Zero means signals are kept indefinitely.
	Retention time.Duration
}

// RetentionConfig configures the background pruning worker.
type RetentionConfig struct {
	// PruneInterval is how often the pruning worker runs.
	// Default: 6h.
	PruneInterval time.Duration
}

// SnapshotConfig controls how signal snapshots are built for risk evaluation.
// Kept for forward-compatibility with the detection engine.
type SnapshotConfig struct {
	HistoryWindow        time.Duration
	ContextChangeWindow  time.Duration
	MaxSignalsPerUser    int
	MaxSignalsPerSession int
}

// RetentionForStream returns the retention duration for the given stream.
// Returns zero (keep forever) if the stream has no explicit retention.
func (c StreamsConfig) RetentionForStream(stream SignalStream) time.Duration {
	switch stream {
	case StreamRequests:
		return c.Requests.Retention
	case StreamEvents:
		return c.Events.Retention
	default:
		return 0
	}
}

// EnabledStreams returns the list of streams that are enabled for collection.
func (c StreamsConfig) EnabledStreams() []SignalStream {
	var streams []SignalStream
	if c.Requests.Enabled {
		streams = append(streams, StreamRequests)
	}
	if c.Events.Enabled {
		streams = append(streams, StreamEvents)
	}
	return streams
}

// DuckLakeConfig configures the DuckLake-based signal store.
// When enabled, signals are stored as Parquet files with catalog metadata
// managed by DuckLake (PostgreSQL catalog backend).
type DuckLakeConfig struct {
	// Enabled activates the DuckLake signal store.
	Enabled bool
	// DataPath is the root path for Parquet data files.
	DataPath string
	// ExtensionDirectory is the directory where DuckDB extensions are cached.
	// When set, DuckDB uses this directory for INSTALL/LOAD instead of the
	// default (~/.duckdb). In container images the extensions are pre-downloaded
	// here so no internet access is needed at runtime.
	// Leave empty to use the DuckDB default.
	ExtensionDirectory string
	// MetadataSchema is the PostgreSQL schema where DuckLake stores its catalog
	// tables (ducklake_metadata, ducklake_snapshots, etc.). Defaults to "signals"
	// which is created by 'zitadel init'. Using a dedicated schema avoids
	// polluting the public schema and PG15+ permission issues.
	MetadataSchema string
	// Backend selects the storage backend: "fs" (default) or "s3".
	Backend ArchiveBackend
	// S3 configures S3-compatible storage when Backend is "s3".
	S3 ArchiveS3Config
	// FlushInterval is how often the emitter flushes buffered signals.
	FlushInterval time.Duration
	// CompactionInterval is how often the compaction worker merges small Parquet files.
	CompactionInterval time.Duration
	// CompactionThreshold is the minimum number of Parquet files that must exist
	// before compaction is triggered. Set to 0 to use the default (10).
	CompactionThreshold int
}

// Validate checks the DuckLake configuration for consistency.
func (c DuckLakeConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.DataPath == "" {
		return fmt.Errorf("identity signals ducklake data_path must not be empty")
	}
	if err := validateSafePath("data_path", c.DataPath); err != nil {
		return err
	}
	if c.ExtensionDirectory != "" {
		if err := validateSafePath("extension_directory", c.ExtensionDirectory); err != nil {
			return err
		}
	}
	if c.MetadataSchema != "" {
		if err := validateSafeIdentifier("metadata_schema", c.MetadataSchema); err != nil {
			return err
		}
	}
	if c.Backend == ArchiveBackendS3 && c.S3.Bucket == "" {
		return fmt.Errorf("identity signals ducklake s3 bucket must not be empty")
	}
	return nil
}

// validateSafePath checks that a filesystem or S3 path contains no control
// characters or SQL-significant characters that could be exploited if the
// value is used inside a DuckDB SQL string literal.
func validateSafePath(field, val string) error {
	for _, ch := range val {
		if ch < 32 || ch == 127 {
			return fmt.Errorf("identity signals ducklake %s contains control characters", field)
		}
	}
	if strings.ContainsAny(val, ";`") {
		return fmt.Errorf("identity signals ducklake %s contains forbidden characters", field)
	}
	return nil
}

// validateSafeIdentifier checks that a SQL identifier (schema/table name)
// is alphanumeric with underscores only.
func validateSafeIdentifier(field, val string) error {
	for _, ch := range val {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return fmt.Errorf("identity signals ducklake %s must be alphanumeric/underscore only", field)
		}
	}
	return nil
}

// ArchiveBackend selects the storage backend for signal data.
type ArchiveBackend string

const (
	ArchiveBackendFS ArchiveBackend = "fs"
	ArchiveBackendS3 ArchiveBackend = "s3"
)

// ArchiveS3Config configures S3-compatible storage for signal data.
type ArchiveS3Config struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	UseSSL    bool
}
