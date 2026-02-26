//go:build integration

package infra

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisConfig struct {
	Addr string // host:port
}

// StartRedis starts a Redis 8 container matching the devcontainer's
// cache-api-integration service.
func StartRedis(ctx context.Context) (*redis.RedisContainer, *RedisConfig, error) {
	// The redis module already waits for "Ready to accept connections" by default.
	container, err := redis.Run(ctx, "redis:8",
		testcontainers.WithLogger(log.New(os.Stderr, "[testcontainers/redis] ", log.LstdFlags)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start redis container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get redis host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, nil, fmt.Errorf("get redis port: %w", err)
	}

	cfg := &RedisConfig{
		Addr: fmt.Sprintf("%s:%s", host, mappedPort.Port()),
	}

	return container, cfg, nil
}
