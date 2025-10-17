//go:build integration

package events_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta_project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/system"
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
	dbPool        *pgxpool.Pool
	CTX           context.Context
	IAMCTX        context.Context
	Instance      *integration.Instance
	SystemClient  system.SystemServiceClient
	OrgClient     v2beta_org.OrganizationServiceClient
	ProjectClient v2beta_project.ProjectServiceClient
	AdminClient   admin.AdminServiceClient
	MgmtClient    mgmt.ManagementServiceClient
)

var pool database.Pool

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		CTX = integration.WithSystemAuthorization(ctx)
		Instance = integration.NewInstance(CTX)

		IAMCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		SystemClient = integration.SystemClient()
		OrgClient = Instance.Client.OrgV2beta
		ProjectClient = Instance.Client.Projectv2Beta
		AdminClient = Instance.Client.Admin
		MgmtClient = Instance.Client.Mgmt

		defer func() {
			_, err := Instance.Client.InstanceV2Beta.DeleteInstance(CTX, &v2beta.DeleteInstanceRequest{
				InstanceId: Instance.Instance.Id,
			})
			if err != nil {
				log.Printf("Failed to delete instance on cleanup: %v\n", err)
			}
		}()

		var err error
		dbConfig, err := pgxpool.ParseConfig(ConnString)
		if err != nil {
			panic(err)
		}
		dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			orgState, err := conn.LoadType(ctx, "zitadel.organization_state")
			if err != nil {
				return err
			}
			conn.TypeMap().RegisterType(orgState)
			return nil
		}

		dbPool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err != nil {
			panic(err)
		}

		pool = postgres.PGxPool(dbPool)

		return m.Run()
	}())
}
