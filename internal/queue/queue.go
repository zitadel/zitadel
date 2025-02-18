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
	driver riverdriver.Driver[pgx.Tx]
	client *river.Client[pgx.Tx]

	config *river.Config
}

type Config struct {
	Client *database.DB `mapstructure:"-"`
}

func NewQueue(config *Config) *Queue {
	return &Queue{
		driver: riverpgxv5.New(config.Client.Pool),
		config: &river.Config{
			Workers:    river.NewWorkers(),
			Queues:     make(map[string]river.QueueConfig),
			JobTimeout: -1,
		},
	}
}

func (q *Queue) Start(ctx context.Context) (err error) {
	ctx = WithQueue(ctx)
	q.client, err = river.NewClient(q.driver, q.config)
	if err != nil {
		return err
	}
	return q.client.Start(ctx)
}

func (q *Queue) AddWorkers(w ...Worker) {
	for _, worker := range w {
		worker.Register(q.config.Workers, q.config.Queues)
	}
}

type insertOpt func(*river.InsertOpts)

func WithMaxAttempts(maxAttempts uint8) insertOpt {
	return func(opts *river.InsertOpts) {
		opts.MaxAttempts = int(maxAttempts)
	}
}

func WithQueueName(name string) insertOpt {
	return func(opts *river.InsertOpts) {
		opts.Queue = name
	}
}

func (q *Queue) Insert(ctx context.Context, args river.JobArgs, opts ...insertOpt) error {
	options := new(river.InsertOpts)
	ctx = WithQueue(ctx)
	for _, opt := range opts {
		opt(options)
	}
	_, err := q.client.Insert(ctx, args, options)
	return err
}

type Worker interface {
	Register(workers *river.Workers, queues map[string]river.QueueConfig)
}
