package signals

import (
	"context"
	"sync"
	"testing"
	"time"
)

// mockSink records calls to WriteBatch for test assertions.
type mockSink struct {
	mu      sync.Mutex
	batches [][]RecordedSignal
}

func (m *mockSink) WriteBatch(_ context.Context, signals []RecordedSignal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]RecordedSignal, len(signals))
	copy(cp, signals)
	m.batches = append(m.batches, cp)
	return nil
}

func (m *mockSink) totalSignals() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	for _, b := range m.batches {
		total += len(b)
	}
	return total
}

func TestEmitter_BasicEmitAndDrain(t *testing.T) {
	sink := &mockSink{}
	cfg := StoreConfig{
		ChannelSize: 100,
		Debounce: DebouncerConfig{
			MinFrequency: 50 * time.Millisecond,
			MaxBulkSize:  10,
		},
	}
	emitter := NewEmitter(cfg, sink, nil)

	ctx, cancel := context.WithCancel(context.Background())
	go emitter.Start(ctx)

	for i := 0; i < 25; i++ {
		emitter.Emit(Signal{
			InstanceID: "inst-1",
			Operation:  "test.op",
			Stream:     StreamRequests,
		})
	}

	// Wait for debouncer to flush
	time.Sleep(200 * time.Millisecond)
	cancel()
	<-emitter.Done()

	if got := sink.totalSignals(); got != 25 {
		t.Errorf("expected 25 signals written, got %d", got)
	}
	if emitter.Dropped() != 0 {
		t.Errorf("expected 0 dropped, got %d", emitter.Dropped())
	}
}

func TestEmitter_DropCounting(t *testing.T) {
	sink := &mockSink{}
	cfg := StoreConfig{
		ChannelSize: 2, // tiny channel
		Debounce: DebouncerConfig{
			MinFrequency: time.Hour, // never auto-flush
			MaxBulkSize:  1000,      // never batch-flush
		},
	}
	emitter := NewEmitter(cfg, sink, nil)
	// Don't start the drain loop — channel will fill up

	for i := 0; i < 10; i++ {
		emitter.Emit(Signal{InstanceID: "inst-1"})
	}

	dropped := emitter.Dropped()
	if dropped < 8 {
		t.Errorf("expected at least 8 dropped signals, got %d", dropped)
	}
}

func TestEmitter_GracefulShutdownDrains(t *testing.T) {
	sink := &mockSink{}
	cfg := StoreConfig{
		ChannelSize: 100,
		Debounce: DebouncerConfig{
			MinFrequency: time.Hour, // no timer flush
			MaxBulkSize:  1000,      // no batch flush
		},
	}
	emitter := NewEmitter(cfg, sink, nil)

	ctx, cancel := context.WithCancel(context.Background())
	go emitter.Start(ctx)

	// Emit some signals
	for i := 0; i < 5; i++ {
		emitter.Emit(Signal{InstanceID: "inst-1"})
	}
	time.Sleep(10 * time.Millisecond) // let them enter the channel

	// Cancel context — should drain remaining signals
	cancel()
	<-emitter.Done()

	if got := sink.totalSignals(); got != 5 {
		t.Errorf("expected 5 signals after graceful shutdown, got %d", got)
	}
}

func TestEmitter_BatchFlush(t *testing.T) {
	sink := &mockSink{}
	cfg := StoreConfig{
		ChannelSize: 100,
		Debounce: DebouncerConfig{
			MinFrequency: time.Hour, // no timer flush
			MaxBulkSize:  5,         // flush every 5 signals
		},
	}
	emitter := NewEmitter(cfg, sink, nil)

	ctx, cancel := context.WithCancel(context.Background())
	go emitter.Start(ctx)

	for i := 0; i < 10; i++ {
		emitter.Emit(Signal{InstanceID: "inst-1"})
	}

	time.Sleep(100 * time.Millisecond)
	cancel()
	<-emitter.Done()

	if got := sink.totalSignals(); got != 10 {
		t.Errorf("expected 10 signals, got %d", got)
	}
	// Should have been flushed in at least 2 batches of 5
	sink.mu.Lock()
	batchCount := len(sink.batches)
	sink.mu.Unlock()
	if batchCount < 2 {
		t.Errorf("expected at least 2 batch flushes, got %d", batchCount)
	}
}

func TestEmitter_DefaultChannelSize(t *testing.T) {
	sink := &mockSink{}
	cfg := StoreConfig{
		ChannelSize: 0, // should default to 4096
	}
	emitter := NewEmitter(cfg, sink, nil)
	if cap(emitter.ch) != 4096 {
		t.Errorf("expected default channel size 4096, got %d", cap(emitter.ch))
	}
}
