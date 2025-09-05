//go:build integration

package events_test

import (
	"context"
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
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

const ConnString = "host=localhost port=5432 user=zitadel password=zitadel dbname=zitadel sslmode=disable"

var (
	dbPool       *pgxpool.Pool
	CTX          context.Context
	IAMCTX       context.Context
	Instance     *integration.Instance
	SystemClient system.SystemServiceClient
	OrgClient    v2beta_org.OrganizationServiceClient
	AdminClient  admin.AdminServiceClient
	MgmtClient   mgmt.ManagementServiceClient
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
