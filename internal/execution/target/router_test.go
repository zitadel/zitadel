package target

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	eventGlobalTarget           = Target{ExecutionID: "event", TargetID: "event_global"}
	eventGroupTarget            = Target{ExecutionID: "event/foo.*", TargetID: "event_group"}
	eventMatchTarget            = Target{ExecutionID: "event/foo.bar", TargetID: "event_specific"}
	functionCallTarget1         = Target{ExecutionID: "function/Call", TargetID: "function_call_1"}
	functionCallTarget2         = Target{ExecutionID: "function/Call", TargetID: "function_call_2"}
	requestGlobalTarget         = Target{ExecutionID: "request", TargetID: "request_global"}
	requestServiceTarget        = Target{ExecutionID: "request/zitadel.test.TestService", TargetID: "request_service"}
	requestServiceMethodTarget  = Target{ExecutionID: "request/zitadel.test.TestService/TestMethod", TargetID: "request_service_method"}
	responseGlobalTarget        = Target{ExecutionID: "response", TargetID: "response_global"}
	responseServiceTarget       = Target{ExecutionID: "response/zitadel.test.TestService", TargetID: "response_service"}
	responseServiceMethodTarget = Target{ExecutionID: "response/zitadel.test.TestService/TestMethod", TargetID: "response_service_method"}

	testTargets = []Target{
		eventGlobalTarget,
		eventGroupTarget,
		eventMatchTarget,
		functionCallTarget1,
		functionCallTarget2,
		requestGlobalTarget,
		requestServiceTarget,
		requestServiceMethodTarget,
		responseGlobalTarget,
		responseServiceTarget,
		responseServiceMethodTarget,
	}
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
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		targets     []Target
		args        args
		wantTargets []Target
		wantOk      bool
	}{

		{
			name:    "event global match",
			targets: testTargets,
			args: args{
				id: "event/bar.foo",
			},
			wantTargets: []Target{eventGlobalTarget},
			wantOk:      true,
		},
		{
			name:    "event group match",
			targets: testTargets,
			args: args{
				id: "event/foo.baz",
			},
			wantTargets: []Target{eventGroupTarget},
			wantOk:      true,
		},
		{
			name:    "event group match with specific match available",
			targets: testTargets,
			args: args{
				id: "event/foo.bar.baz",
			},
			wantTargets: []Target{eventMatchTarget},
			wantOk:      true,
		},
		{
			name:    "event match",
			targets: testTargets,
			args: args{
				id: "event/foo.bar",
			},
			wantTargets: []Target{eventMatchTarget},
			wantOk:      true,
		},
		{
			name:    "function match",
			targets: testTargets,
			args: args{
				id: "function/Call",
			},
			wantTargets: []Target{functionCallTarget1, functionCallTarget2},
			wantOk:      true,
		},
		{
			name:    "request global match",
			targets: testTargets,
			args: args{
				id: "request/zitadel.test.OtherService/OtherMethod",
			},
			wantTargets: []Target{requestGlobalTarget},
			wantOk:      true,
		},
		{
			name:    "request service match",
			targets: testTargets,
			args: args{
				id: "request/zitadel.test.TestService/RandomMethod",
			},
			wantTargets: []Target{requestServiceTarget},
			wantOk:      true,
		},
		{
			name:    "request service method match",
			targets: testTargets,
			args: args{
				id: "request/zitadel.test.TestService/TestMethod",
			},
			wantTargets: []Target{requestServiceMethodTarget},
			wantOk:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(tt.targets)
			got, ok := r.GetEventBestMatch(tt.args.id)
			assert.Equal(t, tt.wantTargets, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}
