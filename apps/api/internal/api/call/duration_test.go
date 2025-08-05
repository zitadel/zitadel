package call

import (
	"context"
	"testing"
	"time"
)

func TestTook(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		args        args
		startIsZero bool
	}{
		{
			name: "no start",
			args: args{
				ctx: context.Background(),
			},
			startIsZero: true,
		},
		{
			name: "with start",
			args: args{
				ctx: WithTimestamp(context.Background()),
			},
			startIsZero: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Took(tt.args.ctx)
			if tt.startIsZero && got != 0 {
				t.Errorf("Duration should be 0 but was %v", got)
			}
			if !tt.startIsZero && got <= 0 {
				t.Errorf("Duration should be greater 0 but was %d", got)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		args   args
		isZero bool
	}{
		{
			name: "no start",
			args: args{
				ctx: context.Background(),
			},
			isZero: true,
		},
		{
			name: "with start",
			args: args{
				ctx: WithTimestamp(context.Background()),
			},
			isZero: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromContext(tt.args.ctx)
			if tt.isZero != got.IsZero() {
				t.Errorf("Time is zero should be %v but was %v", tt.isZero, got.IsZero())
			}
		})
	}
}

func TestWithTimestamp(t *testing.T) {
	start := time.Date(2019, 4, 29, 0, 0, 0, 0, time.UTC)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		noPrevious bool
	}{
		{
			name: "fresh context",
			args: args{
				ctx: context.WithValue(context.Background(), key, start),
			},
			noPrevious: true,
		},
		{
			name: "with start",
			args: args{
				ctx: WithTimestamp(context.Background()),
			},
			noPrevious: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithTimestamp(tt.args.ctx)
			val := got.Value(key).(time.Time)

			if !tt.noPrevious && val.Before(start) {
				t.Errorf("time should be now not %v", val)
			}
			if tt.noPrevious && val.After(start) {
				t.Errorf("time should be start not %v", val)
			}
		})
	}
}
