package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

var pool database.PoolTest

func runTests(m *testing.M) int {
	var stop func()
	var err error
	ctx := context.Background()
	pool, stop, err = newEmbeddedDB(ctx)
	if err != nil {
		log.Printf("error with embedded postgres database: %v", err)
		return 1
	}
	defer stop()

	return m.Run()
}

func newEmbeddedDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	connector, stop, err := embedded.StartEmbedded()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
	}

	pool_, err := connector.Connect(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to embedded postgres: %w", err)
	}
	pool = pool_.(database.PoolTest)

	err = pool.MigrateTest(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to migrate database: %w", err)
	}
	return pool, stop, err
}

func transactionForRollback(t *testing.T) (tx database.Transaction, rollback func()) {
	t.Helper()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	return tx, func() {
		err := tx.Rollback(t.Context())
		require.NoError(t, err)
	}
}

func savepointForRollback(t *testing.T, tx database.Transaction) (savepoint database.Transaction, rollback func()) {
	t.Helper()
	savepoint, err := tx.Begin(t.Context())
	require.NoError(t, err)
	return savepoint, func() {
		err := savepoint.Rollback(t.Context())
		require.NoError(t, err)
	}
}

func createInstance(t *testing.T, tx database.Transaction) (instanceID string) {
	t.Helper()
	instance := domain.Instance{
		ID:              gofakeit.UUID(),
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "consoleClient",
		ConsoleAppID:    "consoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err := instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	return instance.ID
}

func createOrganization(t *testing.T, tx database.Transaction, instanceID string) (orgID string) {
	t.Helper()
	org := domain.Organization{
		InstanceID: instanceID,
		ID:         gofakeit.UUID(),
		Name:       gofakeit.Name(),
		State:      domain.OrgStateActive,
	}
	orgRepo := repository.OrganizationRepository()
	err := orgRepo.Create(t.Context(), tx, &org)
	require.NoError(t, err)

	return org.ID
}

func createProject(t *testing.T, tx database.Transaction, instanceID, orgID string) (projectID string) {
	t.Helper()
	project := domain.Project{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ID:             gofakeit.UUID(),
		Name:           gofakeit.Name(),
		State:          domain.ProjectStateActive,
	}
	projectRepo := repository.ProjectRepository()
	err := projectRepo.Create(t.Context(), tx, &project)
	require.NoError(t, err)

	return project.ID
}
