package signals

// PREVIEW: Identity Signals is a preview feature. APIs, storage format,
// and configuration may change between releases without notice.

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// Emitter provides fire-and-forget signal emission with bounded buffering.
// Signals are batched via a debouncer and flushed to a [SignalSink].
// If the internal channel is full, the signal is dropped and counted.
type Emitter struct {
	ch      chan Signal
	sink    SignalSink
	cfg     StoreConfig
	metrics *SignalMetrics
	dropped atomic.Int64
	done    chan struct{}
}

// NewEmitter creates a new signal emitter. Call [Emitter.Start] to begin
// draining signals from the channel.
func NewEmitter(cfg StoreConfig, sink SignalSink, m *SignalMetrics) *Emitter {
	size := cfg.ChannelSize
	if size <= 0 {
		size = 4096
	}
	return &Emitter{
		ch:      make(chan Signal, size),
		sink:    sink,
		cfg:     cfg,
		metrics: m,
		done:    make(chan struct{}),
	}
}

// Emit enqueues a signal for asynchronous persistence. It never blocks;
// if the channel is full the signal is dropped and counted.
func (e *Emitter) Emit(signal Signal) {
	select {
	case e.ch <- signal:
	default:
		count := e.dropped.Add(1)
		e.metrics.RecordDropped(context.Background(), 1)
		if count%100 == 0 {
			slog.Warn("identity_signals.channel_full",
				slog.Int64("total_dropped", count),
				slog.Int("channel_cap", cap(e.ch)),
			)
		}
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
		ctx:     ctx,
		sink:    e.sink,
		cfg:     e.cfg.Debounce,
		metrics: e.metrics,
		dropped: &e.dropped,
		cache:   make([]Signal, 0, e.cfg.Debounce.MaxBulkSize),
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
			// Drain remaining signals from the channel. The parent context
			// is cancelled, but debouncer.flush() detects ctx.Err() != nil
			// and creates a fresh 5-second background context for the
			// final WriteBatch call (see flush()).
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

// Done returns a channel closed when the emitter has fully stopped.
func (e *Emitter) Done() <-chan struct{} {
	return e.done
}

// signalDebouncer accumulates signals and flushes them in batches.
type signalDebouncer struct {
	ctx     context.Context
	sink    SignalSink
	cfg     DebouncerConfig
	metrics *SignalMetrics
	dropped *atomic.Int64
	mu      sync.Mutex
	cache   []Signal
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
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	recorded := make([]RecordedSignal, len(batch))
	for i, sig := range batch {
		recorded[i] = RecordedSignal{Signal: sig}
	}

	start := time.Now()
	if err := d.sink.WriteBatch(ctx, recorded); err != nil {
		dropped := int64(len(batch))
		d.dropped.Add(dropped)
		d.metrics.RecordDropped(ctx, dropped)
		slog.ErrorContext(ctx, "identity_signals.batch_write_failed",
			slog.Int("batch_size", len(batch)),
			slog.String("error", err.Error()),
		)
	} else {
		d.metrics.RecordBatchWriteDuration(ctx, time.Since(start).Seconds())
		// Count ingested signals per stream for OTEL metrics.
		streamCounts := make(map[string]int64)
		for _, sig := range batch {
			streamCounts[string(sig.Stream)]++
		}
		for stream, count := range streamCounts {
			d.metrics.RecordIngested(ctx, stream, count)
		}
	}
}
