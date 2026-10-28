package queue

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_logMiddleware_InsertMany(t *testing.T) {
	tests := []struct {
		name            string
		results         []*rivertype.JobInsertResult
		wantJobDetails  int
		wantLogContains []string
	}{
		{
			name: "fast insert, nil results",
			// river's fast insert path returns a slice of nil results with the length of the inserted jobs
			results:         make([]*rivertype.JobInsertResult, 2),
			wantJobDetails:  0,
			wantLogContains: []string{`"msg":"jobs inserted"`, `"count":2`},
		},
		{
			name: "results with jobs",
			results: []*rivertype.JobInsertResult{
				{Job: &rivertype.JobRow{ID: 1, Kind: "kind1", Queue: "queue1"}},
				{Job: &rivertype.JobRow{ID: 2, Kind: "kind2", Queue: "queue2"}},
			},
			wantJobDetails:  2,
			wantLogContains: []string{`"msg":"jobs inserted"`, `"count":2`, "kind1", "kind2"},
		},
		{
			name: "result without job",
			// not returned by river, covers the defensive check
			results:         []*rivertype.JobInsertResult{{}},
			wantJobDetails:  0,
			wantLogContains: []string{`"msg":"jobs inserted"`, `"count":1`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			m := &logMiddleware{
				logger: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})),
			}
			var got []*rivertype.JobInsertResult
			require.NotPanics(t, func() {
				var err error
				got, err = m.InsertMany(context.Background(), nil, func(context.Context) ([]*rivertype.JobInsertResult, error) {
					return tt.results, nil
				})
				require.NoError(t, err)
			})
			assert.Equal(t, tt.results, got)
			logs := buf.String()
			assert.Equal(t, tt.wantJobDetails, strings.Count(logs, "inserted job details"))
			for _, want := range tt.wantLogContains {
				assert.Contains(t, logs, want)
			}
		})
	}
}
