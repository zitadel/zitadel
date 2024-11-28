package pg

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/exp/slices"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	//go:embed create_partition.sql.tmpl
	createPartitionQuery string
	createPartitionTmpl  = template.Must(template.New("create_partition").Parse(createPartitionQuery))
	//go:embed set.sql
	setQuery string
	//go:embed get.sql
	getQuery string
	//go:embed invalidate.sql
	invalidateQuery string
	//go:embed delete.sql
	deleteQuery string
	//go:embed prune.sql
	pruneQuery string
	//go:embed truncate.sql
	truncateQuery string
)

type PGXPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type pgCache[I ~int, K ~string, V cache.Entry[I, K]] struct {
	purpose   cache.Purpose
	config    *cache.Config
	indices   []I
	connector *Connector
	logger    *slog.Logger
}

// NewCache returns a cache that stores and retrieves objects using PostgreSQL unlogged tables.
func NewCache[I ~int, K ~string, V cache.Entry[I, K]](ctx context.Context, purpose cache.Purpose, config cache.Config, indices []I, connector *Connector) (cache.PrunerCache[I, K, V], error) {
	c := &pgCache[I, K, V]{
		purpose:   purpose,
		config:    &config,
		indices:   indices,
		connector: connector,
		logger:    config.Log.Slog().With("cache_purpose", purpose),
	}
	c.logger.InfoContext(ctx, "pg cache logging enabled")

	if connector.Dialect == "postgres" {
		if err := c.createPartition(ctx); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *pgCache[I, K, V]) createPartition(ctx context.Context) error {
	var query strings.Builder
	if err := createPartitionTmpl.Execute(&query, c.purpose.String()); err != nil {
		return err
	}
	_, err := c.connector.Exec(ctx, query.String())
	return err
}

func (c *pgCache[I, K, V]) Set(ctx context.Context, entry V) {
	//nolint:errcheck
	c.set(ctx, entry)
}

func (c *pgCache[I, K, V]) set(ctx context.Context, entry V) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	keys := c.indexKeysFromEntry(entry)
	c.logger.DebugContext(ctx, "pg cache set", "index_key", keys)

	_, err = c.connector.Exec(ctx, setQuery, c.purpose.String(), keys, entry)
	if err != nil {
		c.logger.ErrorContext(ctx, "pg cache set", "err", err)
		return err
	}
	return nil
}

func (c *pgCache[I, K, V]) Get(ctx context.Context, index I, key K) (value V, ok bool) {
	value, err := c.get(ctx, index, key)
	if err == nil {
		c.logger.DebugContext(ctx, "pg cache get", "index", index, "key", key)
		return value, true
	}
	logger := c.logger.With("err", err, "index", index, "key", key)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.InfoContext(ctx, "pg cache miss")
		return value, false
	}
	logger.ErrorContext(ctx, "pg cache get", "err", err)
	return value, false
}

func (c *pgCache[I, K, V]) get(ctx context.Context, index I, key K) (value V, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !slices.Contains(c.indices, index) {
		return value, cache.NewIndexUnknownErr(index)
	}
	err = c.connector.QueryRow(ctx, getQuery, c.purpose.String(), index, key, c.config.MaxAge, c.config.LastUseAge).Scan(&value)
	return value, err
}

func (c *pgCache[I, K, V]) Invalidate(ctx context.Context, index I, keys ...K) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, err = c.connector.Exec(ctx, invalidateQuery, c.purpose.String(), index, keys)
	c.logger.DebugContext(ctx, "pg cache invalidate", "index", index, "keys", keys)
	return err
}

func (c *pgCache[I, K, V]) Delete(ctx context.Context, index I, keys ...K) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, err = c.connector.Exec(ctx, deleteQuery, c.purpose.String(), index, keys)
	c.logger.DebugContext(ctx, "pg cache delete", "index", index, "keys", keys)
	return err
}

func (c *pgCache[I, K, V]) Prune(ctx context.Context) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, err = c.connector.Exec(ctx, pruneQuery, c.purpose.String(), c.config.MaxAge, c.config.LastUseAge)
	c.logger.DebugContext(ctx, "pg cache prune")
	return err
}

func (c *pgCache[I, K, V]) Truncate(ctx context.Context) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, err = c.connector.Exec(ctx, truncateQuery, c.purpose.String())
	c.logger.DebugContext(ctx, "pg cache truncate")
	return err
}

type indexKey[I, K comparable] struct {
	IndexID  I `json:"index_id"`
	IndexKey K `json:"index_key"`
}

func (c *pgCache[I, K, V]) indexKeysFromEntry(entry V) []indexKey[I, K] {
	keys := make([]indexKey[I, K], 0, len(c.indices)*3) // naive assumption
	for _, index := range c.indices {
		for _, key := range entry.Keys(index) {
			keys = append(keys, indexKey[I, K]{
				IndexID:  index,
				IndexKey: key,
			})
		}
	}
	return keys
}
