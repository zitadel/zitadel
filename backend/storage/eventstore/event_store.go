package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type Eventstore struct {
	executor database.Executor
}

func New(executor database.Executor) *Eventstore {
	return &Eventstore{executor: executor}
}

type Event interface{}

func (e *Eventstore) Push(ctx context.Context, events ...Event) error {
	return nil
}
