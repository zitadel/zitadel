package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var globalLock sync.Mutex

type testLogEntry struct {
	Level string
	Msg   string
	Err   *testLogError `json:",omitempty"`
	Foo   string
	Group struct {
		V int
	}
}

type testLogError struct {
	Kind    string
	Parent  string
	Message string
	ID      string
}

// prepareDefaultLogger sets the global default logger to a JSON logger writing to a buffer.
// It returns a function that MUST be called to retrieve exactly one logged entry and release the global lock.
func prepareDefaultLogger() (done func() (*testLogEntry, error)) {
	globalLock.Lock()
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))
	slog.SetDefault(logger)

	return func() (*testLogEntry, error) {
		defer globalLock.Unlock()
		if buf.Len() == 0 {
			return nil, nil
		}

		entry := new(testLogEntry)
		decoder := json.NewDecoder(&buf)
		if err := decoder.Decode(entry); err != nil {
			return nil, err
		}
		return entry, nil
	}
}

func TestFromCtx(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
		slog.String("foo", "bar"),
	)
	FromCtx(ctx).Info("test message")
	want := &testLogEntry{
		Level: "INFO",
		Msg:   "test message",
		Foo:   "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestWith(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	ctx = With(ctx, slog.String("foo", "bar"))
	FromCtx(ctx).Info("test message")
	want := &testLogEntry{
		Level: "INFO",
		Msg:   "test message",
		Foo:   "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestWithGroup(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	ctx = WithGroup(ctx, "group")
	ctx = With(ctx, slog.Int("v", 42))
	FromCtx(ctx).Info("test message")
	want := &testLogEntry{
		Level: "INFO",
		Msg:   "test message",
		Group: struct{ V int }{V: 42},
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestLog(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	Log(ctx, slog.LevelInfo, "test message", slog.String("foo", "bar"))
	want := &testLogEntry{
		Level: "INFO",
		Msg:   "test message",
		Foo:   "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDebug(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	Debug(ctx, "test message", slog.String("foo", "bar"))
	want := &testLogEntry{
		Level: "DEBUG",
		Msg:   "test message",
		Foo:   "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestInfo(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	Info(ctx, "test message", slog.String("foo", "bar"))
	want := &testLogEntry{
		Level: "INFO",
		Msg:   "test message",
		Foo:   "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestWarn(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	Warn(ctx, "test message", slog.String("foo", "bar"))
	want := &testLogEntry{
		Level: "WARN",
		Msg:   "test message",
		Foo:   "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestError(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	Error(ctx, "test message", slog.String("foo", "bar"))
	want := &testLogEntry{
		Level: "ERROR",
		Msg:   "test message",
		Foo:   "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestWithError(t *testing.T) {
	done := prepareDefaultLogger()
	ctx := NewCtx(
		t.Context(),
		StreamRuntime,
	)
	err := errors.New("some error")
	logger := WithError(ctx, err)
	logger.Info("test message", slog.String("foo", "bar"))
	want := &testLogEntry{
		Level: "INFO",
		Msg:   "test message",
		Err: &testLogError{
			Kind:    "Unknown",
			Parent:  "some error",
			Message: "an unknown error occurred",
			ID:      "LOG-Ao5ch",
		},
		Foo: "bar",
	}
	got, err := done()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestOnError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *testLogEntry
	}{
		{
			name: "nil error",
			err:  nil,
		},
		{
			name: "non-zitadel error",
			err:  errors.New("some error"),
			want: &testLogEntry{
				Level: "INFO",
				Msg:   "test message",
				Err: &testLogError{
					Kind:    "Unknown",
					Parent:  "some error",
					Message: "an unknown error occurred",
					ID:      "LOG-ii6Pi",
				},
				Foo: "bar",
			},
		},
		{
			name: "zitadel error",
			err: zerrors.CreateZitadelError(
				zerrors.KindNotFound,
				errors.New("parent error"),
				"ZIT-404",
				"resource not found",
				0,
			),
			want: &testLogEntry{
				Level: "INFO",
				Msg:   "test message",
				Err: &testLogError{
					Kind:    "NotFound",
					Parent:  "parent error",
					Message: "resource not found",
					ID:      "ZIT-404",
				},
				Foo: "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := prepareDefaultLogger()
			ctx := NewCtx(
				t.Context(),
				StreamRuntime,
			)
			logger := OnError(ctx, tt.err)
			logger.Info("test message", slog.String("foo", "bar"))
			got, err := done()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
