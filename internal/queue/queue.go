package queue

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"github.com/riverqueue/rivercontrib/otelriver"
	"github.com/robfig/cron/v3"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

// Queue abstracts the underlying queuing library
// For more information see github.com/riverqueue/river
type Queue struct {
	driver riverdriver.Driver[pgx.Tx]
	client *river.Client[pgx.Tx]

	config      *river.Config
	shouldStart bool
}

type Config struct {
	Client *database.DB `mapstructure:"-"` // mapstructure is needed if we would like to use viper to configure the queue
}

func NewQueue(config *Config) (_ *Queue, err error) {
	middleware := []rivertype.Middleware{otelriver.NewMiddleware(&otelriver.MiddlewareConfig{
		MeterProvider: metrics.GetMetricsProvider(),
	})}
	return &Queue{
		driver: riverpgxv5.New(config.Client.Pool),
		config: &river.Config{
			Workers:    river.NewWorkers(),
			Queues:     make(map[string]river.QueueConfig),
			JobTimeout: -1,
			Middleware: middleware,
			Schema:     schema,
		},
	}, nil
}

func (q *Queue) ShouldStart() {
	if q == nil {
		return
	}
	q.shouldStart = true
}

func (q *Queue) Start(ctx context.Context) (err error) {
	if q == nil || !q.shouldStart {
		return nil
	}

	q.client, err = river.NewClient(q.driver, q.config)
	if err != nil {
		return err
	}

	return q.client.Start(ctx)
}

func (q *Queue) AddWorkers(w ...Worker) {
	if q == nil {
		logging.Info("skip adding workers because queue is not set")
		return
	}
	for _, worker := range w {
		worker.Register(q.config.Workers, q.config.Queues)
	}
}

func (q *Queue) AddPeriodicJob(schedule cron.Schedule, jobArgs river.JobArgs, opts ...InsertOpt) (handle rivertype.PeriodicJobHandle) {
	if q == nil {
		logging.Info("skip adding periodic job because queue is not set")
		return
	}
	options := new(river.InsertOpts)
	for _, opt := range opts {
		opt(options)
	}
	return q.client.PeriodicJobs().Add(
		river.NewPeriodicJob(
			schedule,
			func() (river.JobArgs, *river.InsertOpts) {
				return jobArgs, options
			},
			nil,
		),
	)
}

type InsertOpt func(*river.InsertOpts)

func WithMaxAttempts(maxAttempts uint8) InsertOpt {
	return func(opts *river.InsertOpts) {
		opts.MaxAttempts = int(maxAttempts)
	}
}

func WithQueueName(name string) InsertOpt {
	return func(opts *river.InsertOpts) {
		opts.Queue = name
	}
}

func (q *Queue) Insert(ctx context.Context, args river.JobArgs, opts ...InsertOpt) error {
	options := new(river.InsertOpts)
	for _, opt := range opts {
		opt(options)
	}
	_, err := q.client.Insert(ctx, args, options)
	return err
}

type Worker interface {
	Register(workers *river.Workers, queues map[string]river.QueueConfig)
}
