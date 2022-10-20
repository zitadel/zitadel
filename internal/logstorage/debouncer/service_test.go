package debouncer_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/logstorage/debouncer"
)

type shipper struct {
	shipped []uint
}

func (s *shipper) Ship(items []any) {
	s.shipped = append(s.shipped, uint(len(items)))
}

type given struct {
	cfg   *debouncer.Config
	ticks uint
	delay time.Duration
}

func TestNew(t *testing.T) {

	tests := []struct {
		name  string
		given given
		want  []uint
	}{{
		name: "When an empty config is passed, no calls should be made to ship",
		given: given{
			ticks: 2,
		},
		want: nil,
	}, {
		name: "When MinFrequency is 0 seconds and MaxBulkSize is 5, calls should be made immediately",
		given: given{
			ticks: 3,
			cfg: &debouncer.Config{
				MinFrequency: 0,
				MaxBulkSize:  5,
			},
		},
		want: []uint{1, 1, 1},
	}, {
		name: "When MinFrequency is 2 seconds and MaxBulkSize is 0, calls should be made immediately",
		given: given{
			ticks: 3,
			cfg: &debouncer.Config{
				MinFrequency: 2 * time.Second,
				MaxBulkSize:  0,
			},
			delay: time.Second,
		},
		want: []uint{1, 1, 1},
	}, {
		name: "When MinFrequency is 2 second and MaxBulkSize is 4, one call should be made",
		given: given{
			ticks: 3,
			cfg: &debouncer.Config{
				MinFrequency: 2 * time.Second,
				MaxBulkSize:  4,
			},
			delay: time.Second,
		},
		want: []uint{3},
	}, {
		name: "When MinFrequency is 1 second and MaxBulkSize is 4, first two calls, then one call should be made",
		given: given{
			ticks: 3,
			cfg: &debouncer.Config{
				MinFrequency: 1 * time.Second,
				MaxBulkSize:  4,
			},
		},
		want: []uint{2, 1},
	}, {
		name: "When MinFrequency is 2 second and MaxBulkSize is 2, first two calls, then one call should be made",
		given: given{
			ticks: 3,
			cfg: &debouncer.Config{
				MinFrequency: 2 * time.Second,
				MaxBulkSize:  2,
			},
		},
		want: []uint{2, 1},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run(t, tt.given, tt.want)
		})
	}
}

func run(t *testing.T, in given, expect []uint) {
	mock := &shipper{}
	svc := debouncer.New(in.cfg, mock)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	var ticked uint
	for range ticker.C {
		if ticked >= in.ticks {
			break
		}
		ticked++
		svc.Add(1)
	}
	if ticked != in.ticks {
		t.Fatalf("Test setup is wrong. Wanted %d ticks, but broke with %d ticks", in.ticks, ticked)
	}

	if in.delay != 0 {
		timer := time.NewTimer(in.delay)
		defer timer.Stop()
		<-timer.C
	}

	if !reflect.DeepEqual(mock.shipped, expect) {
		t.Errorf("Got calls to Ship() %v, want %v", mock.shipped, expect)
	}
}
