package handler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/mock"
)

func TestHandler_Trigger_skipsInsideLifecycleTransaction(t *testing.T) {
	m := mock.NewSQLMock(t, mock.ExpectBegin(nil))
	defer m.Assert(t)
	tx, err := m.DB.Begin()
	require.NoError(t, err)
	h := &Handler{client: &database.DB{DB: m.DB}}
	ctx, cancel := context.WithTimeout(database.ContextWithTransaction(context.Background(), tx), time.Second)
	defer cancel()
	// Without the guard Trigger blocks on the worker queue, which ignores ctx.
	done := make(chan error, 1)
	go func() {
		_, err := h.Trigger(ctx)
		done <- err
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-ctx.Done():
		t.Fatal("Trigger did not return inside the lifecycle transaction")
	}
}
