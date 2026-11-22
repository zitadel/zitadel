package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
)

func TestUserNotifier_Trigger_withoutWorkerPool(t *testing.T) {
	tests := []struct {
		name string
		new  func(context.Context, handler.Config) *handler.Handler
	}{
		{
			name: "legacy enabled",
			new: func(ctx context.Context, config handler.Config) *handler.Handler {
				return NewUserNotifier(ctx, config, nil, nil, nil, nil, WorkerConfig{LegacyEnabled: true}, nil)
			},
		},
		{
			name: "legacy disabled",
			new: func(ctx context.Context, config handler.Config) *handler.Handler {
				return NewUserNotifier(ctx, config, nil, nil, nil, nil, WorkerConfig{}, nil)
			},
		},
		{
			name: "direct legacy constructor",
			new: func(ctx context.Context, config handler.Config) *handler.Handler {
				return NewUserNotifierLegacy(ctx, config, nil, nil, nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { db.Close() })

			beginErr := errors.New("begin transaction failed")
			mock.ExpectBegin().WillReturnError(beginErr)
			ctx, cancel := context.WithTimeout(authz.WithInstanceID(t.Context(), "instance-1"), time.Second)
			defer cancel()
			notifier := tt.new(ctx, handler.Config{Client: &database.DB{DB: db}})

			// No projection worker pool is started in this package. The notifier must
			// reach the database in its calling goroutine instead of waiting for a worker.
			done := make(chan error, 1)
			go func() {
				_, err := notifier.Trigger(ctx)
				done <- err
			}()
			select {
			case err := <-done:
				require.ErrorIs(t, err, beginErr)
			case <-ctx.Done():
				t.Fatal("notifier waited for a projection worker")
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
