package queue

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/zitadel/zitadel/internal/database"
)

// Queue abstracts the underlying queuing library
// For more information see github.com/riverqueue/river
type Queue struct {
	driver  riverdriver.Driver[pgx.Tx]
	client  *river.Client[pgx.Tx]
	workers *river.Workers
}

type Config struct {
	*river.Config
	Client *database.DB
}

func NewQueue(config *Config) (queue *Queue, err error) {
	queue = &Queue{
		driver:  riverpgxv5.New(config.Client.Pool),
		workers: river.NewWorkers(),
	}

	queue.client, err = river.NewClient(queue.driver, config.Config)
	if err != nil {
		return nil, err
	}

	return queue, nil
}

func (q *Queue) Start(ctx context.Context) error {
	return q.client.Start(ctx)
}

func (q *Queue) AddWorkers(w ...Worker) {
	for _, worker := range w {
		worker.Register(q.workers)
	}
}

func (q *Queue) Insert(ctx context.Context, args river.JobArgs) error {
	ctx = WithQueue(ctx)
	_, err := q.client.Insert(ctx, args, nil)
	return err
}

type Worker interface {
	Register(workers *river.Workers)
}
