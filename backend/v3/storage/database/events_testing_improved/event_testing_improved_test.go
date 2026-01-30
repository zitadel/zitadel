package eventstestingimproved

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
	// Trigger migration registration
	_ "github.com/zitadel/zitadel/cmd/initialise"
	_ "github.com/zitadel/zitadel/cmd/setup"
	es "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var (
	pool   database.PoolTest
	pusher *es.Eventstore
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	var stop func()
	var err error
	ctx := context.Background()
	pool, stop, err = newEmbeddedDB(ctx)
	if err != nil {
		log.Printf("error with embedded postgres database: %v", err)
		os.Exit(1)
	}
	defer stop()

	// q, err := queue.NewQueueWithDB(pool.RawDB())
	pusher = eventstore.NewEventstore(&eventstore.Config{
		PushTimeout: 0,
		MaxRetries:  0,
		Pusher:      es.NewEventstoreFromPool(pool),
		Querier:     nil,
		Searcher:    nil,
		Queue:       nil,
	})


	return m.Run()
}

func newEmbeddedDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	connector, stop, err := embedded.StartEmbedded()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
	}

	dummyPool, err := connector.Connect(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to embedded postgres: %w", err)
	}

	pool = dummyPool.(database.PoolTest)

	err = pool.MigrateTest(ctx, es.RegisterEventstoreTypes)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to migrate database: %w", err)
	}

	return pool, stop, err
}

func TestMyTest(t *testing.T) {
	require.Nil(t, pool.Ping(t.Context()))
	tstCmd := instance.NewDomainAddedEvent(t.Context(), &instance.NewAggregate("123").Aggregate, "test-domain.com", false)
	evt, err := pusher.Push(t.Context(), tstCmd)

	require.NoError(t, err)
	assert.NotEmpty(t, evt)
}
