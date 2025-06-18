//go:build integration

package events_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/internal/integration"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

const ConnString = "host=localhost port=5432 user=zitadel dbname=zitadel sslmode=disable"

var (
	dbPool       *pgxpool.Pool
	CTX          context.Context
	Instance     *integration.Instance
	SystemClient system.SystemServiceClient
	OrgClient    v2beta_org.OrganizationServiceClient
	MgmtClient   mgmt.ManagementServiceClient
)

var pool database.Pool

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		SystemClient = integration.SystemClient()
		OrgClient = Instance.Client.OrgV2beta
		MgmtClient = Instance.Client.Mgmt

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
