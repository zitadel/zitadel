// The library github.com/benbjohnson/clock fails when race is enabled
// https://github.com/benbjohnson/clock/issues/44
//go:build !race

package logstore_test

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/zitadel/zitadel/internal/logstore"
	emittermock "github.com/zitadel/zitadel/internal/logstore/emitters/mock"
	quotaqueriermock "github.com/zitadel/zitadel/internal/logstore/quotaqueriers/mock"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	tick  = time.Second
	ticks = 60
)

type args struct {
	mainSink      *logstore.EmitterConfig
	secondarySink *logstore.EmitterConfig
	config        quota.AddedEvent
}

type want struct {
	enabled       bool
	remaining     *uint64
	mainSink      wantSink
	secondarySink wantSink
}

type wantSink struct {
	bulks []int
	len   int
}

func TestService(t *testing.T) {
	// tests should run on a single thread
	// important for deterministic results
	beforeProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(beforeProcs)

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "max and min debouncing works",
		args: args{
			mainSink: emitterConfig(withDebouncerConfig(&logstore.DebouncerConfig{
				MinFrequency: 1 * time.Minute,
				MaxBulkSize:  60,
			})),
			secondarySink: emitterConfig(),
			config:        quotaConfig(),
		},
		want: want{
			enabled:   true,
			remaining: nil,
			mainSink: wantSink{
				bulks: repeat(60, 1),
				len:   60,
			},
			secondarySink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
		},
	}, {
		name: "mixed debouncing works",
		args: args{
			mainSink: emitterConfig(withDebouncerConfig(&logstore.DebouncerConfig{
				MinFrequency: 0,
				MaxBulkSize:  6,
			})),
			secondarySink: emitterConfig(withDebouncerConfig(&logstore.DebouncerConfig{
				MinFrequency: 10 * time.Second,
				MaxBulkSize:  0,
			})),
			config: quotaConfig(),
		},
		want: want{
			enabled:   true,
			remaining: nil,
			mainSink: wantSink{
				bulks: repeat(6, 10),
				len:   60,
			},
			secondarySink: wantSink{
				bulks: repeat(10, 6),
				len:   60,
			},
		},
	}, {
		name: "when disabling main sink, secondary sink still works",
		args: args{
			mainSink:      emitterConfig(withDisabled()),
			secondarySink: emitterConfig(),
			config:        quotaConfig(),
		},
		want: want{
			enabled:   true,
			remaining: nil,
			mainSink: wantSink{
				bulks: repeat(99, 0),
				len:   0,
			},
			secondarySink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
		},
	}, {
		name: "when all sink are disabled, the service is disabled",
		args: args{
			mainSink:      emitterConfig(withDisabled()),
			secondarySink: emitterConfig(withDisabled()),
			config:        quotaConfig(),
		},
		want: want{
			enabled:   false,
			remaining: nil,
			mainSink: wantSink{
				bulks: repeat(99, 0),
				len:   0,
			},
			secondarySink: wantSink{
				bulks: repeat(99, 0),
				len:   0,
			},
		},
	}, {
		name: "cleanupping works",
		args: args{
			mainSink: emitterConfig(withCleanupping(17*time.Second, 28*time.Second)),
			secondarySink: emitterConfig(withDebouncerConfig(&logstore.DebouncerConfig{
				MinFrequency: 0,
				MaxBulkSize:  15,
			}), withCleanupping(5*time.Second, 47*time.Second)),
			config: quotaConfig(),
		},
		want: want{
			enabled:   true,
			remaining: nil,
			mainSink: wantSink{
				bulks: repeat(1, 60),
				len:   21,
			},
			secondarySink: wantSink{
				bulks: repeat(15, 4),
				len:   18,
			},
		},
	}, {
		name: "when quota has a limit of 90, 30 are remaining",
		args: args{
			mainSink:      emitterConfig(),
			secondarySink: emitterConfig(),
			config:        quotaConfig(withLimiting()),
		},
		want: want{
			enabled:   true,
			remaining: uint64Ptr(30),
			mainSink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
			secondarySink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
		},
	}, {
		name: "when quota has a limit of 30, 0 are remaining",
		args: args{
			mainSink:      emitterConfig(),
			secondarySink: emitterConfig(),
			config:        quotaConfig(withLimiting(), withAmountAndInterval(30)),
		},
		want: want{
			enabled:   true,
			remaining: uint64Ptr(0),
			mainSink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
			secondarySink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
		},
	}, {
		name: "when quota has amount of 30 but is not limited, remaining is nil",
		args: args{
			mainSink:      emitterConfig(),
			secondarySink: emitterConfig(),
			config:        quotaConfig(withAmountAndInterval(30)),
		},
		want: want{
			enabled:   true,
			remaining: nil,
			mainSink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
			secondarySink: wantSink{
				bulks: repeat(1, 60),
				len:   60,
			},
		},
	}}
	for _, tt := range tests {
		runTest(t, tt.name, tt.args, tt.want)
	}
}

func runTest(t *testing.T, name string, args args, want want) bool {
	return t.Run("Given over a minute, each second a log record is emitted", func(tt *testing.T) {
		tt.Run(name, func(t *testing.T) {
			ctx, clock, mainStorage, secondaryStorage, svc := given(t, args, want)
			remaining := when(svc, ctx, clock)
			then(t, mainStorage, secondaryStorage, remaining, want)
		})
	})
}

func given(t *testing.T, args args, want want) (context.Context, *clock.Mock, *emittermock.InmemLogStorage, *emittermock.InmemLogStorage, *logstore.Service) {
	ctx := context.Background()
	clock := clock.NewMock()

	periodStart := time.Time{}
	clock.Set(args.config.From)

	mainStorage := emittermock.NewInMemoryStorage(clock)
	mainEmitter, err := logstore.NewEmitter(ctx, clock, args.mainSink, mainStorage)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	secondaryStorage := emittermock.NewInMemoryStorage(clock)
	secondaryEmitter, err := logstore.NewEmitter(ctx, clock, args.secondarySink, secondaryStorage)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	svc := logstore.New(
		quotaqueriermock.NewNoopQuerier(&args.config, periodStart),
		logstore.UsageReporterFunc(func(context.Context, []*quota.NotifiedEvent) error { return nil }),
		mainEmitter,
		secondaryEmitter)

	if svc.Enabled() != want.enabled {
		t.Errorf("wantet service enabled to be %t but is %t", want.enabled, svc.Enabled())
	}
	return ctx, clock, mainStorage, secondaryStorage, svc
}

func when(svc *logstore.Service, ctx context.Context, clock *clock.Mock) *uint64 {
	var remaining *uint64
	for i := 0; i < ticks; i++ {
		svc.Handle(ctx, emittermock.NewRecord(clock))
		runtime.Gosched()
		remaining = svc.Limit(ctx, "non-empty-instance-id")
		clock.Add(tick)
	}
	time.Sleep(time.Millisecond)
	runtime.Gosched()
	return remaining
}

func then(t *testing.T, mainStorage, secondaryStorage *emittermock.InmemLogStorage, remaining *uint64, want want) {
	mainBulks := mainStorage.Bulks()
	if !reflect.DeepEqual(want.mainSink.bulks, mainBulks) {
		t.Errorf("wanted main storage to have bulks %v, but got %v", want.mainSink.bulks, mainBulks)
	}

	mainLen := mainStorage.Len()
	if !reflect.DeepEqual(want.mainSink.len, mainLen) {
		t.Errorf("wanted main storage to have len %d, but got %d", want.mainSink.len, mainLen)
	}

	secondaryBulks := secondaryStorage.Bulks()
	if !reflect.DeepEqual(want.secondarySink.bulks, secondaryBulks) {
		t.Errorf("wanted secondary storage to have bulks %v, but got %v", want.secondarySink.bulks, secondaryBulks)
	}

	secondaryLen := secondaryStorage.Len()
	if !reflect.DeepEqual(want.secondarySink.len, secondaryLen) {
		t.Errorf("wanted secondary storage to have len %d, but got %d", want.secondarySink.len, secondaryLen)
	}

	if remaining == nil && want.remaining == nil {
		return
	}

	if remaining == nil && want.remaining != nil ||
		remaining != nil && want.remaining == nil {
		t.Errorf("wantet remaining nil %t but got %t", want.remaining == nil, remaining == nil)
		return
	}
	if *remaining != *want.remaining {
		t.Errorf("wantet remaining %d but got %d", *want.remaining, *remaining)
		return
	}
}
