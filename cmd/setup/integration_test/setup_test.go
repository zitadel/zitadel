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
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const ConnString = "host=localhost port=5432 user=zitadel dbname=zitadel sslmode=disable"

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
