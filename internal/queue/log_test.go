package queue

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/require"
)

// Regression for #12225: River InsertManyFast / InsertManyFastTx return a
// count-sized slice of nil *JobInsertResult (no job rows). The queue log
// middleware must not panic when detailing inserts under debug.
func TestLogMiddleware_InsertMany_NilFastResultsDoNotPanic(t *testing.T) {
	fastShapedNils := func(context.Context) ([]*rivertype.JobInsertResult, error) {
		return make([]*rivertype.JobInsertResult, 3), nil
	}
	mixedResults := func(context.Context) ([]*rivertype.JobInsertResult, error) {
		return []*rivertype.JobInsertResult{
			nil,
			{Job: nil},
			{Job: &rivertype.JobRow{
				ID:          1,
				Kind:        "test",
				Queue:       "default",
				State:       rivertype.JobStateAvailable,
				CreatedAt:   time.Now(),
				ScheduledAt: time.Now(),
			}},
		}, nil
	}

	tests := []struct {
		name     string
		level    slog.Level
		doInner  func(context.Context) ([]*rivertype.JobInsertResult, error)
		wantLen  int
	}{
		{
			name:    "debug, InsertManyFast-shaped nils",
			level:   slog.LevelDebug,
			doInner: fastShapedNils,
			wantLen: 3,
		},
		{
			name:    "info, InsertManyFast-shaped nils",
			level:   slog.LevelInfo,
			doInner: fastShapedNils,
			wantLen: 3,
		},
		{
			name:    "debug, mixed nil and populated results",
			level:   slog.LevelDebug,
			doInner: mixedResults,
			wantLen: 3,
		},
		{
			name:    "info, mixed nil and populated results",
			level:   slog.LevelInfo,
			doInner: mixedResults,
			wantLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &logMiddleware{
				logger: slog.New(slog.NewTextHandler(t.Output(), &slog.HandlerOptions{Level: tt.level})),
			}
			require.NotPanics(t, func() {
				results, err := m.InsertMany(context.Background(), nil, tt.doInner)
				require.NoError(t, err)
				require.Len(t, results, tt.wantLen)
			})
		})
	}
}
