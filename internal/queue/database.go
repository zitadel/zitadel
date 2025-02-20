package queue

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"

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
