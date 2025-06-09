package queue

import (
	"testing"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/riverqueue/rivercontrib/otelriver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

func TestNewQueue(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Queue
		wantErr bool
	}{
		{
			name: "create queue with metrics disabled",
			args: args{
				config: &Config{
					Client:        &database.DB{},
					EnableMetrics: false,
				},
			},
			want: &Queue{
				config: &river.Config{
					Workers:    river.NewWorkers(),
					Queues:     make(map[string]river.QueueConfig),
					JobTimeout: -1,
					Middleware: []rivertype.Middleware{},
				},
			},
			wantErr: false,
		},
		{
			name: "create queue with metrics enabled",
			args: args{
				config: &Config{
					Client:        &database.DB{},
					EnableMetrics: true,
				},
			},
			want: &Queue{
				config: &river.Config{
					Workers:    river.NewWorkers(),
					Queues:     make(map[string]river.QueueConfig),
					JobTimeout: -1,
					Middleware: []rivertype.Middleware{
						otelriver.NewMiddleware(&otelriver.MiddlewareConfig{
							EnableSemanticMetrics: true,
							MeterProvider:         metrics.GetMetricsProvider(),
						}),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQueue(tt.args.config)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, len(tt.want.config.Middleware), len(got.config.Middleware))
			if tt.args.config.EnableMetrics {
				assert.NotNil(t, got.config.Middleware[0])
			}
		})
	}
}
