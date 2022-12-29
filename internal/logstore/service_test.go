package logstore_test

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/thejerf/abtime"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/mock"
)

const (
	tick  = time.Second
	ticks = 60
)

func TestService(t *testing.T) {
	// tests should run on a single thread
	// important for deterministic results
	beforeProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(beforeProcs)
	type args struct {
		mainSink      *logstore.EmitterConfig
		secondarySink *logstore.EmitterConfig
	}
	type wantSink struct {
		err   error
		bulks []int
		len   int
	}
	type want struct {
		enabled       bool
		handleErr     error
		limitErr      error
		doLimit       bool
		remaining     *uint64
		mainSink      wantSink
		secondarySink wantSink
	}
	tests := []struct {
		name string
		args args
		want want
	}{ /*{
			name: "max and min debouncing works",
			args: args{
				mainSink: emitterConfig(withDebouncerConfig(&logstore.DebouncerConfig{
					MinFrequency: 1 * time.Minute,
					MaxBulkSize:  60,
				})),
				secondarySink: emitterConfig(),
			},
			want: want{
				enabled:   true,
				handleErr: nil,
				limitErr:  nil,
				doLimit:   false,
				remaining: nil,
				mainSink: wantSink{
					err:   nil,
					bulks: repeat(60, 1),
					len:   60,
				},
				secondarySink: wantSink{
					err:   nil,
					bulks: repeat(1, 60),
					len:   60,
				},
			},
		}, */{
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
			},
			want: want{
				enabled:   true,
				handleErr: nil,
				limitErr:  nil,
				doLimit:   false,
				remaining: nil,
				mainSink: wantSink{
					err:   nil,
					bulks: repeat(6, 10),
					len:   60,
				},
				secondarySink: wantSink{
					err:   nil,
					bulks: repeat(10, 6),
					len:   60,
				},
			},
		}, /*, {
			name: "when disabling main sink, secondary sink still works",
			args: args{
				mainSink:      emitterConfig(withDisabled()),
				secondarySink: emitterConfig(),
			},
			want: want{
				enabled:   true,
				handleErr: nil,
				limitErr:  nil,
				doLimit:   false,
				remaining: nil,
				mainSink: wantSink{
					err:   nil,
					bulks: repeat(99, 0),
					len:   0,
				},
				secondarySink: wantSink{
					err:   nil,
					bulks: repeat(1, 60),
					len:   60,
				},
			},
		}, {
			name: "cleanupping works",
			args: args{
				mainSink: emitterConfig(withCleanupping(17*time.Second, 30*time.Second)),
				secondarySink: emitterConfig(withDebouncerConfig(&logstore.DebouncerConfig{
					MinFrequency: 0,
					MaxBulkSize:  15,
				}), withCleanupping(5*time.Second, 47*time.Second)),
			},
			want: want{
				enabled:   true,
				handleErr: nil,
				limitErr:  nil,
				doLimit:   false,
				remaining: nil,
				mainSink: wantSink{
					err:   nil,
					bulks: repeat(1, 60),
					len:   16,
				},
				secondarySink: wantSink{
					err:   nil,
					bulks: repeat(15, 4),
					len:   16, // In-Memory cleanup is inclusive
				},
			},
		}*/}
	for _, ttt := range tests {
		t.Run("Given over a minute, each second a log record is emitted", func(tt *testing.T) {
			tt.Run(ttt.name, func(t *testing.T) {
				ctx := context.Background()
				clock := abtime.NewManual()
				mainStorage := mock.NewInMemoryStorage(clock)
				ttt.args.mainSink.Z_ManualTickerCleanupID = 1
				ttt.args.mainSink.Z_ManualTickerMinFrequencyID = 2
				mainEmitter, err := logstore.NewEmitter(ctx, clock, ttt.args.mainSink, mainStorage)
				if err != nil {
					if !errors.Is(err, ttt.want.mainSink.err) {
						t.Errorf("wantet err %v but got err %v", ttt.want.mainSink.err, err)
					}
					return
				}
				ttt.args.mainSink.Z_ManualTickerCleanupID = 3
				ttt.args.mainSink.Z_ManualTickerMinFrequencyID = 4
				secondaryStorage := mock.NewInMemoryStorage(clock)
				secondaryEmitter, err := logstore.NewEmitter(ctx, clock, ttt.args.secondarySink, secondaryStorage)
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
					return
				}

				svc := logstore.New(
					mainEmitter,
					nil,
					nil,
					secondaryEmitter)

				if svc.Enabled() != ttt.want.enabled {
					t.Errorf("wantet service enabled to be %t but is %t", ttt.want.enabled, svc.Enabled())
					return
				}

				now := clock.Now()
				for i := 0; i < ticks; i++ {
					err = svc.Handle(ctx, mock.NewRecord(clock))
					clock.Advance(tick)
					trigger(clock, now, ttt.args.mainSink.Debounce.MinFrequency, ttt.args.mainSink.Z_ManualTickerMinFrequencyID)
					trigger(clock, now, ttt.args.secondarySink.Debounce.MinFrequency, ttt.args.secondarySink.Z_ManualTickerMinFrequencyID)
					trigger(clock, now, ttt.args.mainSink.CleanupInterval, ttt.args.mainSink.Z_ManualTickerCleanupID)
					trigger(clock, now, ttt.args.secondarySink.CleanupInterval, ttt.args.secondarySink.Z_ManualTickerCleanupID)
				}

				if !errors.Is(err, ttt.want.handleErr) {
					t.Errorf("wantet err %v but got err %v", ttt.want.handleErr, err)
					return
				}
				err = nil

				mainBulks := mainStorage.Bulks()
				if !reflect.DeepEqual(ttt.want.mainSink.bulks, mainBulks) {
					t.Errorf("wanted main storage to have bulks %v, but got %v", ttt.want.mainSink.bulks, mainBulks)
				}

				mainLen := mainStorage.Len()
				if !reflect.DeepEqual(ttt.want.mainSink.len, mainLen) {
					t.Errorf("wanted main storage to have len %d, but got %d", ttt.want.mainSink.len, mainLen)
				}

				secondaryBulks := secondaryStorage.Bulks()
				if !reflect.DeepEqual(ttt.want.secondarySink.bulks, secondaryBulks) {
					t.Errorf("wanted secondary storage to have bulks %v, but got %v", ttt.want.secondarySink.bulks, secondaryBulks)
				}

				secondaryLen := secondaryStorage.Len()
				if !reflect.DeepEqual(ttt.want.secondarySink.len, secondaryLen) {
					t.Errorf("wanted secondary storage to have len %d, but got %d", ttt.want.secondarySink.len, secondaryLen)
				}

				doLimit, remaining, err := svc.Limit(ctx, "")
				if !errors.Is(err, ttt.want.limitErr) {
					t.Errorf("wantet err %v but got err %v", ttt.want.limitErr, err)
				}
				if doLimit != ttt.want.doLimit {
					t.Errorf("wantet limit %t but got %t", ttt.want.doLimit, doLimit)
				}

				if remaining == nil && ttt.want.remaining == nil {
					return
				}

				if remaining == nil && ttt.want.remaining != nil ||
					remaining != nil && ttt.want.remaining == nil {
					t.Errorf("wantet remaining nil %t but got %t", ttt.want.remaining == nil, remaining == nil)
				}
				if *remaining != *ttt.want.remaining {
					t.Errorf("wantet remaining %d but got %d", *ttt.want.remaining, *remaining)
				}
			})
		})
	}
}

func trigger(clock *abtime.ManualTime, since time.Time, interval time.Duration, id int) {
	if interval == 0 {
		return
	}
	if int64(clock.Now().Sub(since))%int64(interval) == 0 {
		clock.Trigger(id)
		// run all routines
		// important for deterministic results
		runtime.Gosched()
	}
}
