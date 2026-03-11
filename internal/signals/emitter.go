package signals

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

// SignalSink is the interface for batch-writing signals to persistent storage.
type SignalSink interface {
	WriteBatch(ctx context.Context, signals []RecordedSignal) error
}

// EnrichFunc is an optional callback that enriches raw signals with detection
// findings before they are persisted. When set on the [Emitter], the debouncer
// calls it during flush to attach any applicable findings. Signals for which
// enrichment is not relevant (e.g. no UserID) should be returned unchanged.
type EnrichFunc func(ctx context.Context, signals []Signal) []RecordedSignal

// Emitter provides fire-and-forget signal emission with bounded buffering.
// Signals are batched via a debouncer and flushed to a [SignalSink].
// If the internal channel is full, signals are dropped and counted.
type Emitter struct {
	ch      chan Signal
	sink    SignalSink
	cfg     SignalStoreConfig
	enrich  EnrichFunc
	dropped atomic.Int64
	done    chan struct{}
}

// NewEmitter creates a new signal emitter. Call [Emitter.Start] to begin
// draining signals from the channel.
// NewEmitter creates a new signal emitter. Call [Emitter.Start] to begin
// draining signals from the channel. An optional [EnrichFunc] can be set
// via [Emitter.SetEnrichFunc] to attach detection findings before persistence.
func NewEmitter(cfg SignalStoreConfig, sink SignalSink) *Emitter {
	size := cfg.ChannelSize
	if size <= 0 {
		size = 4096
	}
	return &Emitter{
		ch:   make(chan Signal, size),
		sink: sink,
		cfg:  cfg,
		done: make(chan struct{}),
	}
}

// SetEnrichFunc sets the enrichment callback. Must be called before [Emitter.Start].
func (e *Emitter) SetEnrichFunc(fn EnrichFunc) {
	e.enrich = fn
}

// Emit enqueues a signal for asynchronous persistence. It never blocks;
// if the channel is full the signal is dropped and the drop counter is
// incremented.
func (e *Emitter) Emit(signal Signal) {
	select {
	case e.ch <- signal:
	default:
		e.dropped.Add(1)
	}
}

// Dropped returns the number of signals dropped since the emitter was created.
func (e *Emitter) Dropped() int64 {
	return e.dropped.Load()
}

// Start begins the background drain loop. It blocks until ctx is cancelled,
// at which point it flushes any remaining buffered signals and closes the
// done channel. Call this in a goroutine.
func (e *Emitter) Start(ctx context.Context) {
	defer close(e.done)

	d := &signalDebouncer{
		ctx:    ctx,
		sink:   e.sink,
		enrich: e.enrich,
		cfg:    e.cfg.Debounce,
		cache:  make([]Signal, 0, e.cfg.Debounce.MaxBulkSize),
	}

	var ticker *time.Ticker
	var tickC <-chan time.Time
	if e.cfg.Debounce.MinFrequency > 0 {
		ticker = time.NewTicker(e.cfg.Debounce.MinFrequency)
		tickC = ticker.C
		defer ticker.Stop()
	}

	for {
		select {
		case sig, ok := <-e.ch:
			if !ok {
				d.flush()
				return
			}
			d.add(sig)
			if d.shouldFlush() {
				d.flush()
				if ticker != nil {
					ticker.Reset(e.cfg.Debounce.MinFrequency)
				}
			}
		case <-tickC:
			d.flush()
		case <-ctx.Done():
			// Drain remaining signals from the channel.
			for {
				select {
				case sig := <-e.ch:
					d.add(sig)
				default:
					d.flush()
					return
				}
			}
		}
	}
}

// Done returns a channel that is closed when the emitter has fully stopped
// (after context cancellation and final flush).
func (e *Emitter) Done() <-chan struct{} {
	return e.done
}

// signalDebouncer accumulates signals and flushes them in batches.
type signalDebouncer struct {
	ctx    context.Context
	sink   SignalSink
	enrich EnrichFunc
	cfg    DebouncerConfig
	mu     sync.Mutex
	cache  []Signal
}

func (d *signalDebouncer) add(sig Signal) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache = append(d.cache, sig)
}

func (d *signalDebouncer) shouldFlush() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.cfg.MaxBulkSize > 0 && uint(len(d.cache)) >= d.cfg.MaxBulkSize
}

func (d *signalDebouncer) flush() {
	d.mu.Lock()
	batch := d.cache
	d.cache = make([]Signal, 0, d.cfg.MaxBulkSize)
	d.mu.Unlock()

	if len(batch) == 0 {
		return
	}

	ctx := d.ctx
	if ctx.Err() != nil {
		// Use a short-lived context for the final flush on shutdown.
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	recorded := d.toRecorded(ctx, batch)
	if err := d.sink.WriteBatch(ctx, recorded); err != nil {
		if instrumentation.IsStreamEnabled(instrumentation.StreamRisk) {
			logging.WithError(ctx, err).Error("signal_store.batch_write_failed",
				slog.Int("batch_size", len(batch)),
			)
		}
	}
}

// toRecorded converts a batch of raw signals to RecordedSignals. When an
// EnrichFunc is configured, it is called to attach detection findings.
// Otherwise signals are wrapped with empty findings.
func (d *signalDebouncer) toRecorded(ctx context.Context, batch []Signal) []RecordedSignal {
	if d.enrich != nil {
		return d.enrich(ctx, batch)
	}
	recorded := make([]RecordedSignal, len(batch))
	for i, sig := range batch {
		recorded[i] = RecordedSignal{Signal: sig}
	}
	return recorded
}
