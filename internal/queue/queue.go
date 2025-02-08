package queue

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river/riverdriver"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
)

const (
	schema          = "queue"
	applicationName = "zitadel_queue"
)

var conns = &sync.Map{}

type queueKey struct{}

func WithQueue(parent context.Context) context.Context {
	return context.WithValue(parent, queueKey{}, struct{}{})
}

func init() {
	dialect.RegisterBeforeAcquire(func(ctx context.Context, c *pgx.Conn) error {
		if _, ok := ctx.Value(queueKey{}).(struct{}); !ok {
			return nil
		}
		_, err := c.Exec(ctx, "SET search_path TO "+schema+"; SET application_name TO "+applicationName)
		if err != nil {
			return err
		}
		conns.Store(c, struct{}{})
		return nil
	})
	dialect.RegisterAfterRelease(func(c *pgx.Conn) error {
		_, ok := conns.LoadAndDelete(c)
		if !ok {
			return nil
		}
		_, err := c.Exec(context.Background(), "SET search_path TO DEFAULT; SET application_name TO "+dialect.DefaultAppName)
		return err
	})
}

// Queue abstracts the underlying queuing library
// For more information see github.com/riverqueue/river
// TODO(adlerhurst): maybe it makes more sense to split the effective queue from the migrator.
type Queue struct {
	driver riverdriver.Driver[pgx.Tx]
}

func New(client *database.DB) *Queue {
	return &Queue{driver: riverpgxv5.New(client.Pool)}
}

func (q *Queue) ExecuteMigrations(ctx context.Context) error {
	_, err := q.driver.GetExecutor().Exec(ctx, "CREATE SCHEMA IF NOT EXISTS queue")
	if err != nil {
		return err
	}

	migrator, err := rivermigrate.New(q.driver, nil)
	if err != nil {
		return err
	}
	ctx = WithQueue(ctx)
	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err
}
