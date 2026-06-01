package handlers

import (
	"context"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/internal/queue"
)

//go:generate mockgen -typed -package mock -destination ./mock/queue.mock.go . Queue
type Queue interface {
	Insert(ctx context.Context, args river.JobArgs, opts ...queue.InsertOpt) error
}
