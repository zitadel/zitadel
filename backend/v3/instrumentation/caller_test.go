package instrumentation

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCaller(t *testing.T) {
	_, _, prevLine, _ := runtime.Caller(0)
	caller, ok := GetCaller(0)
	require.True(t, ok)
	want := Caller{
		Function: "github.com/zitadel/zitadel/backend/v3/instrumentation.TestGetCaller",
		File:     "backend/v3/instrumentation/caller_test.go",
		Line:     prevLine + 1,
	}
	assert.Equal(t, want.Function, caller.Function, "function")
	assert.Contains(t, caller.File, want.File, "file")
	assert.Equal(t, want.Line, caller.Line, "line")
}

func TestGetCallingFunc(t *testing.T) {
	tests := []struct {
		name string
		skip int
		want string
	}{
		{
			name: "test function",
			skip: 0,
			want: "github.com/zitadel/zitadel/backend/v3/instrumentation.TestGetCallingFunc.func1",
		},
		{
			name: "unknown caller",
			skip: 100,
			want: CallerUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callerFunc := GetCallingFunc(tt.skip)
			assert.Equal(t, tt.want, callerFunc, "function")
		})
	}
}
