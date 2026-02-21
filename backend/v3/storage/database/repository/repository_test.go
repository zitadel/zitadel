package repository_test

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
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
	defer stop()
	if err != nil {
		log.Printf("error with embedded postgres database: %v", err)
		return 1
	}
	defer func() {
		r := recover()
		pool.Close(ctx)
		stop()
		if r != nil {
			panic(r)
		}
	}()

	return m.Run()
}

func newEmbeddedDB(ctx context.Context) (pool database.PoolTest, stop func(), err error) {
	var connector database.Connector
	if url := os.Getenv("ZITADEL_TEST_POSTGRES_URL"); url != "" {
		log.Println("using database provided by env")
		connector, err = postgres.DecodeConfig(url)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to connect to provided postgres: %w", err)
		}
		stop = func() {}
	} else {
		connector, stop, err = embedded.StartEmbedded()
		if err != nil {
			return nil, nil, fmt.Errorf("unable to start embedded postgres: %w", err)
		}
	}

	pool_, err := connector.Connect(ctx)
	if err != nil {
		return nil, stop, fmt.Errorf("unable to connect to embedded postgres: %w", err)
	}
	pool = pool_.(database.PoolTest)

	err = pool.MigrateTest(ctx)
	if err != nil {
		return nil, stop, fmt.Errorf("unable to migrate database: %w", err)
	}
	return pool, stop, err
}

func transactionForRollback(t *testing.T) (tx database.Transaction, rollback func()) {
	t.Helper()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	return tx, func() {
		// context.Background to ensure rollback does not return an error if test is already done
		err := tx.Rollback(context.Background())
		require.NoError(t, err)
	}
}

func savepointForRollback(t *testing.T, tx database.Transaction) (savepoint database.Transaction, rollback func()) {
	t.Helper()
	savepoint, err := tx.Begin(t.Context())
	require.NoError(t, err)
	return savepoint, func() {
		// context.Background to ensure rollback does not return an error if test is already done
		err := savepoint.Rollback(context.Background())
		require.NoError(t, err)
	}
}

func createInstance(t *testing.T, tx database.QueryExecutor) (instanceID string) {
	t.Helper()
	instance := domain.Instance{
		ID:              gofakeit.UUID(),
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "managementConsoleClient",
		ConsoleAppID:    "managementConsoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	instanceRepo := repository.InstanceRepository()
	err := instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	return instance.ID
}

func createOrganization(t *testing.T, tx database.QueryExecutor, instanceID string) (orgID string) {
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

func createProject(t *testing.T, tx database.QueryExecutor, instanceID, orgID string) (projectID string) {
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

func createProjectRole(t *testing.T, tx database.QueryExecutor, instanceID, orgID, projectID, key string) string {
	t.Helper()
	if key == "" {
		key = integration.RoleKey()
	}
	projectRole := domain.ProjectRole{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ProjectID:      projectID,
		Key:            key,
		DisplayName:    integration.RoleDisplayName(),
	}
	projectRoleRepo := repository.ProjectRepository().Role()
	err := projectRoleRepo.Create(t.Context(), tx, &projectRole)
	require.NoError(t, err)

	return projectRole.Key
}

func createProjectGrant(t *testing.T, tx database.Transaction, instanceID, grantingOrgID, grantedOrgID, projectID string, roleKeys []string) string {
	t.Helper()
	projectGrant := domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     gofakeit.UUID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  grantedOrgID,
		State:                  domain.ProjectGrantStateActive,
		RoleKeys:               roleKeys,
	}
	projectGrantRepo := repository.ProjectGrantRepository()
	err := projectGrantRepo.Create(t.Context(), tx, &projectGrant)
	require.NoError(t, err)

	return projectGrant.ID
}

func createIdentityProvider(t *testing.T, tx database.Transaction, instanceID, orgID string) string {
	t.Helper()
	idp := domain.IdentityProvider{
		InstanceID:        instanceID,
		OrgID:             &orgID,
		ID:                gofakeit.UUID(),
		State:             domain.IDPStateActive,
		Name:              gofakeit.Name(),
		Type:              gu.Ptr(domain.IDPTypeOIDC),
		AllowCreation:     true,
		AllowAutoCreation: true,
		AllowAutoUpdate:   true,
		AllowLinking:      true,
		StylingType:       &stylingType,
		Payload:           []byte("{}"),
	}
	idpRepo := repository.IDProviderRepository()
	err := idpRepo.Create(t.Context(), tx, &idp)
	require.NoError(t, err)

	return idp.ID
}

func createIDPIntent(t *testing.T, tx database.Transaction, instanceID, idpID string) string {
	t.Helper()
	successURL, err := url.Parse("https://example.com/success")
	require.NoError(t, err)
	failURL, err := url.Parse("https://example.com/fail")
	require.NoError(t, err)

	intent := domain.IDPIntent{
		ID:           gofakeit.UUID(),
		InstanceID:   instanceID,
		SuccessURL:   successURL,
		FailureURL:   failURL,
		IDPID:        idpID,
		IDPArguments: map[string]any{"arg1": map[string]any{"k1": 1, "k2": "v2"}},
		CreatedAt:    time.Now(),
	}
	idpIntentRepo := repository.IDPIntentRepository()
	err = idpIntentRepo.Create(t.Context(), tx, &intent)
	require.NoError(t, err)

	return intent.ID
}

func createMachineUser(t *testing.T, tx database.QueryExecutor, instanceID, orgID string) (userID string) {
	t.Helper()
	user := domain.User{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ID:             gofakeit.UUID(),
		Username:       gofakeit.Username(),
		State:          domain.UserStateActive,
		Machine: &domain.MachineUser{
			Name: gofakeit.Name(),
		},
	}
	userRepo := repository.UserRepository()
	err := userRepo.Create(t.Context(), tx, &user)
	require.NoError(t, err)
	return user.ID
}

func createHumanUser(t *testing.T, tx database.QueryExecutor, instanceID, orgID string) string {
	t.Helper()
	user := &domain.User{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ID:             gofakeit.UUID(),
		Username:       gofakeit.Username(),
		State:          domain.UserStateActive,
		Human: &domain.HumanUser{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email: domain.HumanEmail{
				Address: gofakeit.Email(),
			},
			Password: domain.HumanPassword{
				Hash: "hashed-password",
			},
		},
	}
	userRepo := repository.UserRepository()
	err := userRepo.Create(t.Context(), tx, user)
	require.NoError(t, err)
	return user.ID
}
