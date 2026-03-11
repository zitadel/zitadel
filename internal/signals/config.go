package signals

import "time"

type SnapshotConfig struct {
	HistoryWindow        time.Duration
	ContextChangeWindow  time.Duration
	MaxSignalsPerUser    int
	MaxSignalsPerSession int
}

type SignalStoreConfig struct {
	// Enabled activates the persistent signal store.
	Enabled bool
	// ChannelSize is the buffer size for the fire-and-forget emission channel.
	// Signals are dropped (and counted) when the channel is full.
	ChannelSize int
	// Debounce controls batching of signal writes.
	Debounce DebouncerConfig
	// DuckLake configures the DuckLake-based signal store (Parquet + PG catalog).
	DuckLake DuckLakeConfig
	// Streams controls which signal streams are captured.
	// Each stream defaults to enabled when the store is enabled.
	Streams StreamsConfig
}

// StreamsConfig controls which request pathways emit signals.
type StreamsConfig struct {
	// API enables signal emission for every connectRPC/gRPC API call (stream=request).
	// Covers all /zitadel.* service methods.
	API bool
	// HTTPAccess enables signal emission for every raw HTTP request (stream=request).
	// Covers OIDC endpoints, SAML, login UI, and any non-gRPC path.
	HTTPAccess bool
}

// WithDefaults returns the StreamsConfig with both streams enabled if none are
// explicitly configured. This preserves backward compatibility: an empty struct
// (zero value from config) enables all streams.
func (s StreamsConfig) WithDefaults() StreamsConfig {
	if !s.API && !s.HTTPAccess {
		return StreamsConfig{API: true, HTTPAccess: true}
	}
	return s
}

// DebouncerConfig controls how signals are batched before writing.
type DebouncerConfig struct {
	// MinFrequency is the maximum time between flushes.
	MinFrequency time.Duration
	// MaxBulkSize is the maximum batch size before a flush is triggered.
	MaxBulkSize uint
}

func (c SignalStoreConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	return c.DuckLake.Validate()
}
