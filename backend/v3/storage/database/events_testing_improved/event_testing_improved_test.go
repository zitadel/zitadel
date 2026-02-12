package eventstestingimproved

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	// Trigger migration registration
	_ "github.com/zitadel/zitadel/cmd/initialise"
	_ "github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	es_v3 "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var (
	pool database.PoolTest
	es   *eventstore.Eventstore
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
	pusher := es_v3.NewEventstoreFromPool(pool)
	es = eventstore.NewEventstore(&eventstore.Config{
		PushTimeout: 0,
		MaxRetries:  0,
		Pusher:      pusher,
		Querier:     old_es.NewPostgres(pool.InternalDB()),
		Searcher:    pusher,
		Queue:       nil,
	})

	projection.CreateRelational(ctx, pool.InternalDB(), es, projection.Config{})

	if err := projection.Start(ctx); err != nil {
		log.Printf("error starting projections: %v", err)
		return 1
	}

	return m.Run()
}

func newEmbeddedDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	connector, stop, err := embedded.StartEmbedded()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
	}

	dummyPool, err := connector.Connect(ctx, postgres.WithAfterConnectFunc(func(ctx context.Context, c *pgx.Conn) error {
		// TODO(IAM-Marco): I am not sure whether this is needed or not. I speculate it is used for the ES table to handle
		// the `position` column. Without it, I imagine that events might not be in the right order (because of position having low res).
		pgxdecimal.Register(c.TypeMap())
		return es_v3.RegisterEventstoreTypes(ctx, c)
	}))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to embedded postgres: %w", err)
	}

	pool = dummyPool.(database.PoolTest)
	err = pool.MigrateTest(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to migrate database: %w", err)
	}

	return pool, stop, err
}

func TestDomainAddedEvent(t *testing.T) {
	ctx := t.Context()
	instanceDomainRepo := repository.InstanceDomainRepository()
	instanceRepo := repository.InstanceRepository()

	err := instanceRepo.Create(ctx, pool, &domain.Instance{
		ID:              "123",
		Name:            "my instance",
		DefaultOrgID:    gofakeit.UUID(),
		IAMProjectID:    gofakeit.UUID(),
		ConsoleClientID: gofakeit.UUID(),
		ConsoleAppID:    gofakeit.UUID(),
		DefaultLanguage: gofakeit.LanguageAbbreviation(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	tstCmd := instance.NewDomainAddedEvent(ctx, &instance.NewAggregate("123").Aggregate, "test-domain.com", false)
	_, err = es.Push(ctx, tstCmd)
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*30)
	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		domain, err := instanceDomainRepo.Get(ctx, pool,
			database.WithCondition(
				database.And(
					instanceDomainRepo.InstanceIDCondition("123"),
					instanceDomainRepo.DomainCondition(database.TextOperationEqual, "test-domain.com"),
				),
			),
		)
		require.NoError(t, err)
		assert.Equal(t, "test-domain.com", domain.Domain)
		assert.Equal(t, "123", domain.InstanceID)
	}, retryDuration, tick)
}
