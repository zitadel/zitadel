package memory

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/internal/v3/storage"
)

var _ storage.Client = (*Client)(nil)

type Client struct{}

func (c *Client) Begin(ctx context.Context) (storage.Transaction, error) {
	return new(Transaction), nil
}

var _ storage.Transaction = (*Transaction)(nil)

type Transaction struct {
	commitHooks   []func(ctx context.Context) error
	rollbackHooks []func(ctx context.Context) error
}

// Commit implements storage.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	for _, hook := range t.commitHooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

// End implements storage.Transaction.
func (t *Transaction) End(ctx context.Context, err error) error {
	if err != nil {
		rollbackErr := t.Rollback(ctx)
		slog.WarnContext(ctx, "Rollback failed", slog.Any("cause", rollbackErr))

		return err
	}
	return t.Commit(ctx)
}

// OnCommit implements storage.Transaction.
func (t *Transaction) OnCommit(hook func(ctx context.Context) error) {
	t.commitHooks = append(t.commitHooks, hook)
}

// OnRollback implements storage.Transaction.
func (t *Transaction) OnRollback(hook func(ctx context.Context) error) {
	t.rollbackHooks = append(t.rollbackHooks, hook)
}

// Rollback implements storage.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	for _, hook := range t.rollbackHooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}
