package queue

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"github.com/riverqueue/rivercontrib/otelriver"
	"github.com/robfig/cron/v3"
	"github.com/zitadel/logging"

	new_db "github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	new_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

type Client interface {
	Start(ctx context.Context) error
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	InsertManyFastTx(ctx context.Context, tx new_db.Transaction, params []river.InsertManyParams) (int, error)
	PeriodicJobs() *river.PeriodicJobBundle
}

type pgxClient struct {
	*river.Client[pgx.Tx]
}

func (c *pgxClient) InsertManyFastTx(ctx context.Context, tx new_db.Transaction, params []river.InsertManyParams) (int, error) {
	return c.Client.InsertManyFastTx(ctx, tx.(*postgres.Tx).Tx, params)
}

type sqlClient struct {
	*river.Client[*sql.Tx]
}

func (c *sqlClient) InsertManyFastTx(ctx context.Context, tx new_db.Transaction, params []river.InsertManyParams) (int, error) {
	return c.Client.InsertManyFastTx(ctx, tx.(*new_sql.Tx).Tx, params)
}

// Queue abstracts the underlying queuing library
// For more information see github.com/riverqueue/river
type Queue struct {
	db *database.DB

	client      Client
	config      *river.Config
	shouldStart bool
}

type Config struct {
	Client *database.DB `mapstructure:"-"` // mapstructure is needed if we would like to use viper to configure the queue
}

func NewQueue(config *Config) (_ *Queue, err error) {
	middleware := []rivertype.Middleware{otelriver.NewMiddleware(&otelriver.MiddlewareConfig{
		MeterProvider: metrics.GetMetricsProvider(),
		DurationUnit:  "ms",
	})}
	return &Queue{
		// driver: driver,
		db: config.Client,
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

	switch pool := q.db.DB.(type) {
	case *new_sql.Pool:
		client, err := river.NewClient(riverdatabasesql.New(pool.DB), q.config)
		q.client = &sqlClient{client}
		if err != nil {
			return err
		}
	case *postgres.Pool:
		client, err := river.NewClient(riverpgxv5.New(pool.Pool), q.config)
		q.client = &pgxClient{client}
		if err != nil {
			return err
		}
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
	_, err := q.client.Insert(ctx, args, applyInsertOpts(opts))
	return err
}

// InsertManyFastTx wraps [river.Client.InsertManyFastTx] to insert all jobs in
// a single `COPY FROM` execution, within the existing transaction.
//
// Opts are applied to each job before sending them to river.
func (q *Queue) InsertManyFastTx(ctx context.Context, tx new_db.Transaction, args []river.JobArgs, opts ...InsertOpt) error {
	params := make([]river.InsertManyParams, len(args))
	for i, arg := range args {
		params[i] = river.InsertManyParams{
			Args:       arg,
			InsertOpts: applyInsertOpts(opts),
		}
	}

	_, err := q.client.InsertManyFastTx(ctx, tx, params)
	return err
}

func applyInsertOpts(opts []InsertOpt) *river.InsertOpts {
	options := new(river.InsertOpts)
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type Worker interface {
	Register(workers *river.Workers, queues map[string]river.QueueConfig)
}
