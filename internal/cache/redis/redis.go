package redis

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/zitadel/zitadel/internal/cache"
)

var (
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
	setParsed        = redis.NewScript(strings.Join([]string{utilComponent, setScript}, "\n"))
	getParsed        = redis.NewScript(strings.Join([]string{utilComponent, removeComponent, getScript}, "\n"))
	invalidateParsed = redis.NewScript(strings.Join([]string{utilComponent, removeComponent, invalidateScript}, "\n"))
)

type redisCache[I, K comparable, V cache.Entry[I, K]] struct {
	name    string
	config  *cache.CacheConfig
	indices []I
	client  *redis.Client
	logger  *slog.Logger
}

// NewCache returns a cache that does nothing
func NewCache[I, K comparable, V cache.Entry[I, K]]() cache.Cache[I, K, V] {
	return &redisCache[I, K, V]{
		client: redis.NewClient(nil),
	}
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
	argv := []any{
		buf.String(),                             // object
		int64(c.config.LastUseAge / time.Second), // usage_lifetime
		int64(c.config.MaxAge / time.Second),     // max_age
	}
	err = setParsed.Run(ctx, c.client, keys, argv...).Err()
	// redis.Nil is always returned because the script doesn't have a return value.
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return objectID, nil
}

func (c *redisCache[I, K, V]) Get(ctx context.Context, index I, key K) (value V, ok bool) {
	logger := c.logger.With("index", index, "key", key)
	obj, err := getParsed.Run(ctx, c.client, c.redisIndexKeys(index, key)).Result()
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
	err = invalidateParsed.Run(ctx, c.client, c.redisIndexKeys(index, key...)).Err()
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
	return c.client.Del(ctx, c.redisIndexKeys(index, key...)...).Err()
}
func (c *redisCache[I, K, V]) Truncate(ctx context.Context) (err error) {
	return c.client.FlushDB(ctx).Err()
}

func (c *redisCache[I, K, V]) redisIndexKeys(index I, keys ...K) []string {
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = fmt.Sprintf("%s:%v:%v", c.name, index, k)
	}
	return out
}
