//go:build integration

package infra

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// ConnectionString returns a postgres:// URL for the container.
func (c PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
}

// StartPostgres starts a Postgres 18 container with tuning matching the
// devcontainer's db-api-integration service.
func StartPostgres(ctx context.Context, logw io.Writer) (*postgres.PostgresContainer, *PostgresConfig, error) {
	const (
		user     = "postgres"
		password = "postgres"
		dbname   = "zitadel"
	)

	container, err := postgres.Run(ctx, "postgres:18",
		testcontainers.WithReuseByName("zitadel-integration-postgres"),
		testcontainers.WithLogger(log.New(logw, "[testcontainers/postgres] ", log.LstdFlags)),
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		}),
		// Match the performance tuning from .devcontainer/docker-compose.yaml
		postgres.BasicWaitStrategies(),
		// pass postgres tuning flags; BasicWaitStrategies() above handles readiness
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Cmd: []string{
					"-c", "shared_preload_libraries=pg_stat_statements",
					"-c", "pg_stat_statements.track=all",
					"-c", "shared_buffers=256MB",
					"-c", "work_mem=16MB",
					"-c", "effective_io_concurrency=100",
					"-c", "wal_level=minimal",
					"-c", "archive_mode=off",
					"-c", "max_wal_senders=0",
				},
			},
		}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start postgres container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get postgres host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("get postgres port: %w", err)
	}

	cfg := &PostgresConfig{
		Host:     host,
		Port:     mappedPort.Port(),
		User:     user,
		Password: password,
		Database: dbname,
	}

	return container, cfg, nil
}
