package redis

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/zitadel/zitadel/internal/cache"
)

var (
	//go:embed _select.lua
	selectComponent string
	//go:embed _util.lua
	utilComponent string
	//go:embed _remove.lua
	removeComponent string
	//go:embed set.lua
	setScript string
	//go:embed get.lua
	getScript string
	//go:embed invalidate.lua
	invalidateScript string

	// Don't mind the creative "import"
	setParsed        = redis.NewScript(strings.Join([]string{selectComponent, utilComponent, setScript}, "\n"))
	getParsed        = redis.NewScript(strings.Join([]string{selectComponent, utilComponent, removeComponent, getScript}, "\n"))
	invalidateParsed = redis.NewScript(strings.Join([]string{selectComponent, utilComponent, removeComponent, invalidateScript}, "\n"))
)

type redisCache[I, K comparable, V cache.Entry[I, K]] struct {
	db      int
	config  *cache.CacheConfig
	indices []I
	client  *redis.Client
	logger  *slog.Logger
}

// NewCache returns a cache that does nothing
func NewCache[I, K comparable, V cache.Entry[I, K]](config cache.CacheConfig, client *redis.Client, db int, indices []I) cache.Cache[I, K, V] {
	return &redisCache[I, K, V]{
		config:  &config,
		db:      db,
		indices: indices,
		client:  client,
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		})),
	}
}

func (c *redisCache[I, K, V]) Close() error {
	return c.client.Close()
}

func (c *redisCache[I, K, V]) Set(ctx context.Context, value V) {
	if _, err := c.set(ctx, value); err != nil {
		c.logger.ErrorContext(ctx, "redis cache set", "err", err)
	}
}

func (c *redisCache[I, K, V]) set(ctx context.Context, value V) (objectID string, err error) {
	// Internal ID used for the object
	objectID = uuid.NewString()
	keys := []string{objectID}
	// flatten the secondary keys
	for _, index := range c.indices {
		keys = append(keys, c.redisIndexKeys(index, value.Keys(index)...)...)
	}
	var buf strings.Builder
	err = json.NewEncoder(&buf).Encode(value)
	if err != nil {
		return "", err
	}
	err = setParsed.Run(ctx, c.client, keys,
		c.db,                                   // DB namespace
		buf.String(),                           // object
		int64(c.config.LastUseAge/time.Second), // usage_lifetime
		int64(c.config.MaxAge/time.Second),     // max_age,
	).Err()
	// redis.Nil is always returned because the script doesn't have a return value.
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return objectID, nil
}

func (c *redisCache[I, K, V]) Get(ctx context.Context, index I, key K) (value V, ok bool) {
	logger := c.logger.With("index", index, "key", key)
	obj, err := getParsed.Run(ctx, c.client, c.redisIndexKeys(index, key), c.db).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.ErrorContext(ctx, "redis cache get", "err", err)
		return value, false
	}
	data, ok := obj.(string)
	if !ok {
		logger.With("err", err).InfoContext(ctx, "redis cache miss")
		return value, false
	}
	err = json.NewDecoder(strings.NewReader(data)).Decode(&value)
	if err != nil {
		c.logger.ErrorContext(ctx, "redis cache get", "err", fmt.Errorf("decode: %w", err))
		return value, false
	}
	return value, true

}

func (c *redisCache[I, K, V]) Invalidate(ctx context.Context, index I, key ...K) (err error) {
	if len(key) == 0 {
		return nil
	}
	err = invalidateParsed.Run(ctx, c.client, c.redisIndexKeys(index, key...), c.db).Err()
	// redis.Nil is always returned because the script doesn't have a return value.
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}

func (c *redisCache[I, K, V]) Delete(ctx context.Context, index I, key ...K) (err error) {
	if len(key) == 0 {
		return nil
	}
	pipe := c.client.Pipeline()
	pipe.Select(ctx, c.db)
	pipe.Del(ctx, c.redisIndexKeys(index, key...)...)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *redisCache[I, K, V]) Truncate(ctx context.Context) (err error) {
	pipe := c.client.Pipeline()
	pipe.Select(ctx, c.db)
	pipe.FlushDB(ctx)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *redisCache[I, K, V]) redisIndexKeys(index I, keys ...K) []string {
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = fmt.Sprintf("%v:%v", index, k)
	}
	return out
}
