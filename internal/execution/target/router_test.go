package target

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	eventGlobalTarget   = Target{ExecutionID: "event", TargetID: "event_global"}
	eventGroupTarget    = Target{ExecutionID: "event/foo.*", TargetID: "event_group"}
	eventMatchTarget    = Target{ExecutionID: "event/foo.bar", TargetID: "event_specific"}
	functionCallTarget1 = Target{ExecutionID: "function/Call", TargetID: "function_call_1"}
	functionCallTarget2 = Target{ExecutionID: "function/Call", TargetID: "function_call_2"}

	testTargets = []Target{eventGlobalTarget, eventGroupTarget, eventMatchTarget, functionCallTarget1, functionCallTarget2}
)

func TestBinarySearchRouter_Get(t *testing.T) {
	r := NewRouter(testTargets)
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		wantTargets []Target
		wantOk      bool
	}{
		{
			name: "event global does not match exactly",
			args: args{
				id: "event/bar.foo",
			},
			wantTargets: nil,
			wantOk:      false,
		},
		{
			name: "event group does not match exactly",
			args: args{
				id: "event/foo.bar.baz",
			},
			wantTargets: nil,
			wantOk:      false,
		},
		{
			name: "event match",
			args: args{
				id: "event/foo.bar",
			},
			wantTargets: []Target{eventMatchTarget},
			wantOk:      true,
		},
		{
			name: "function match",
			args: args{
				id: "function/Call",
			},
			wantTargets: []Target{functionCallTarget1, functionCallTarget2},
			wantOk:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := r.Get(tt.args.id)
			assert.Equal(t, tt.wantTargets, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestBinarySearchRouter_GetEventBestMatch(t *testing.T) {
	r := NewRouter(testTargets)
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		wantTargets []Target
		wantOk      bool
	}{
		{
			name: "event global match",
			args: args{
				id: "event/bar.foo",
			},
			wantTargets: []Target{eventGlobalTarget},
			wantOk:      true,
		},
		{
			name: "event group match",
			args: args{
				id: "event/foo.bar.baz",
			},
			wantTargets: []Target{eventGroupTarget},
			wantOk:      true,
		},
		{
			name: "event match",
			args: args{
				id: "event/foo.bar",
			},
			wantTargets: []Target{eventMatchTarget},
			wantOk:      true,
		},
		{
			name: "function match",
			args: args{
				id: "function/Call",
			},
			wantTargets: []Target{functionCallTarget1, functionCallTarget2},
			wantOk:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := r.GetEventBestMatch(tt.args.id)
			assert.Equal(t, tt.wantTargets, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}
