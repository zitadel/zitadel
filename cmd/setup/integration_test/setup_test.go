// Package setup_test implements tests for procedural PostgreSQL functions,
// created in the database during Zitadel setup.
// Tests depend on `zitadel setup` being run first and therefore is run as integration tests.
// A PGX connection is used directly to the integration test database.
// This package assumes the database server available as per integration test defaults.
// See the [ConnString] constant.

//go:build integration

package setup_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var ConnString = fmt.Sprintf(
	"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	getEnv("ZITADEL_DATABASE_POSTGRES_HOST", "localhost"),
	getEnv("ZITADEL_DATABASE_POSTGRES_PORT", "5433"),
	getEnv("ZITADEL_DATABASE_POSTGRES_USER", "zitadel"),
	getEnv("ZITADEL_DATABASE_POSTGRES_PASSWORD", "zitadel"),
	getEnv("ZITADEL_DATABASE_POSTGRES_DATABASE", "zitadel"),
	getEnv("ZITADEL_DATABASE_POSTGRES_SSL_MODE", "disable"),
)

var (
	CTX    context.Context
	dbPool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	CTX, cancel = context.WithTimeout(context.Background(), time.Second*10)

	var err error
	dbPool, err = pgxpool.New(context.Background(), ConnString)
	if err != nil {
		panic(err)
	}
	exit := m.Run()
	cancel()
	dbPool.Close()
	os.Exit(exit)
}
