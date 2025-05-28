//go:build integration

package instance_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

const ConnString = "host=localhost port=5432 user=zitadel dbname=zitadel sslmode=disable"

var (
	dbPool       *pgxpool.Pool
	CTX          context.Context
	SystemCTX    context.Context
	Instance     *integration.Instance
	SystemClient system.SystemServiceClient
)

var pool database.Pool

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// SystemCTX = integration.WithSystemAuthorization(ctx)
		SystemClient = integration.SystemClient()

		var err error
		dbPool, err = pgxpool.New(context.Background(), ConnString)
		if err != nil {
			panic(err)
		}

		pool = postgres.PGxPool(dbPool)

		return m.Run()
	}())
}

func TestServer_TestInstanceAddReduces(t *testing.T) {
	instanceName := "newInstance"
	_, err := SystemClient.CreateInstance(CTX, &system.CreateInstanceRequest{
		InstanceName: instanceName,
		Owner: &system.CreateInstanceRequest_Machine_{
			Machine: &system.CreateInstanceRequest_Machine{
				UserName:            "owner",
				Name:                "owner",
				PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
			},
		},
	})

	require.NoError(t, err)

	instanceRepo := repository.InstanceRepository(pool)
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
		instance, err := instanceRepo.Get(CTX,
			database.WithCondition(
				instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
			),
		)
		require.NoError(ttt, err)
		require.Equal(ttt, instanceName, instance.Name)
	}, retryDuration, tick)
}

func TestServer_TestInstanceUpdateNameReduces(t *testing.T) {
	instanceName := gofakeit.Name()
	res, err := SystemClient.CreateInstance(CTX, &system.CreateInstanceRequest{
		InstanceName: instanceName,
		Owner: &system.CreateInstanceRequest_Machine_{
			Machine: &system.CreateInstanceRequest_Machine{
				UserName:            "owner",
				Name:                "owner",
				PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
			},
		},
	})
	require.NoError(t, err)

	instanceName += "new"
	_, err = SystemClient.UpdateInstance(CTX, &system.UpdateInstanceRequest{
		InstanceId:   res.InstanceId,
		InstanceName: instanceName,
	})
	require.NoError(t, err)

	instanceRepo := repository.InstanceRepository(pool)
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
		instance, err := instanceRepo.Get(CTX,
			database.WithCondition(
				instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
			),
		)
		require.NoError(ttt, err)
		require.Equal(ttt, instanceName, instance.Name)
	}, retryDuration, tick)
}
