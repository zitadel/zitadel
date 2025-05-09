package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type poolHandler[T any] struct {
	pool database.Pool

	client database.QueryExecutor
}

func (h *poolHandler[T]) acquire(ctx context.Context, in T) (out T, _ func(context.Context, error) error, err error) {
	client, err := h.pool.Acquire(ctx)
	if err != nil {
		return in, nil, err
	}
	h.client = client

	return in, func(ctx context.Context, _ error) error { return client.Release(ctx) }, nil
}

func (h *poolHandler[T]) begin(ctx context.Context, in T) (out T, _ func(context.Context, error) error, err error) {
	var beginner database.Beginner = h.pool
	if h.client != nil {
		beginner = h.client.(database.Beginner)
	}
	previousClient := h.client
	tx, err := beginner.Begin(ctx, nil)
	if err != nil {
		return in, nil, err
	}
	h.client = tx

	return in, func(ctx context.Context, err error) error {
		err = tx.End(ctx, err)
		if err != nil {
			return err
		}
		h.client = previousClient
		return nil
	}, nil
}
