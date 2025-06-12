//go:build integration

package events_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

// const ConnString = "host=localhost port=5432 user=zitadel dbname=zitadel sslmode=disable"

// var (
// 	dbPool       *pgxpool.Pool
// 	CTX          context.Context
// 	Organization     *integration.Organization
// 	SystemClient system.SystemServiceClient
// )

// var pool database.Pool

// func TestMain(m *testing.M) {
// 	os.Exit(func() int {
// 		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
// 		defer cancel()

// 		Organization = integration.NewOrganization(ctx)

// 		CTX = Organization.WithAuthorization(ctx, integration.UserTypeIAMOwner)
// 		SystemClient = integration.SystemClient()

// 		var err error
// 		dbPool, err = pgxpool.New(context.Background(), ConnString)
// 		if err != nil {
// 			panic(err)
// 		}

// 		pool = postgres.PGxPool(dbPool)

// 		return m.Run()
// 	}())
// }

func TestServer_TestOrganizationAddReduces(t *testing.T) {
	orgName := gofakeit.Name()
	// beforeCreate := time.Now()

	_, err := OrgClient.AddOrganization(CTX, &org.AddOrganizationRequest{
		Name: orgName,
	})
	require.NoError(t, err)
	// afterCreate := time.Now()

	orgRepo := repository.OrgRepository(pool)
	organization, err := orgRepo.Get(CTX,
		orgRepo.NameCondition(database.TextOperationEqual, orgName),
	)
	require.NoError(t, err)
	fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> organization = %+v\n", organization)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
		organization, err := orgRepo.Get(CTX,
			orgRepo.NameCondition(database.TextOperationEqual, orgName),
		)
		require.NoError(ttt, err)
		fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> organization = %+v\n", organization)
		// // event instance.added
		// require.Equal(ttt, instanceName, organization.Name)
		// // event instance.default.org.set
		// require.NotNil(t, organization.DefaultOrgID)
		// // event instance.iam.project.set
		// require.NotNil(t, organization.IAMProjectID)
		// // event instance.iam.console.set
		// require.NotNil(t, organization.ConsoleAppID)
		// // event instance.default.language.set
		// require.NotNil(t, organization.DefaultLanguage)
		// // event instance.added
		// assert.WithinRange(t, organization.CreatedAt, beforeCreate, afterCreate)
		// // event instance.added
		// assert.WithinRange(t, organization.UpdatedAt, beforeCreate, afterCreate)
		// require.Nil(t, organization.DeletedAt)
	}, retryDuration, tick)
}

// func TestServer_TestOrganizationUpdateNameReduces(t *testing.T) {
// 	instanceName := gofakeit.Name()
// 	res, err := SystemClient.CreateOrganization(CTX, &system.CreateOrganizationRequest{
// 		OrganizationName: instanceName,
// 		Owner: &system.CreateOrganizationRequest_Machine_{
// 			Machine: &system.CreateOrganizationRequest_Machine{
// 				UserName:            "owner",
// 				Name:                "owner",
// 				PersonalAccessToken: &system.CreateOrganizationRequest_PersonalAccessToken{},
// 			},
// 		},
// 	})
// 	require.NoError(t, err)

// 	instanceName += "new"
// 	_, err = SystemClient.UpdateOrganization(CTX, &system.UpdateOrganizationRequest{
// 		OrganizationId:   res.OrganizationId,
// 		OrganizationName: instanceName,
// 	})
// 	require.NoError(t, err)

// 	instanceRepo := repository.OrganizationRepository(pool)
// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
// 	assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
// 		instance, err := instanceRepo.Get(CTX,
// 			instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
// 		)
// 		require.NoError(ttt, err)
// 		// event instance.changed
// 		require.Equal(ttt, instanceName, instance.Name)
// 	}, retryDuration, tick)
// }

// func TestServer_TestOrganizationDeleteReduces(t *testing.T) {
// 	instanceName := gofakeit.Name()
// 	res, err := SystemClient.CreateOrganization(CTX, &system.CreateOrganizationRequest{
// 		OrganizationName: instanceName,
// 		Owner: &system.CreateOrganizationRequest_Machine_{
// 			Machine: &system.CreateOrganizationRequest_Machine{
// 				UserName:            "owner",
// 				Name:                "owner",
// 				PersonalAccessToken: &system.CreateOrganizationRequest_PersonalAccessToken{},
// 			},
// 		},
// 	})
// 	require.NoError(t, err)

// 	_, err = SystemClient.RemoveOrganization(CTX, &system.RemoveOrganizationRequest{
// 		OrganizationId: res.OrganizationId,
// 	})
// 	require.NoError(t, err)

// 	instanceRepo := repository.OrganizationRepository(pool)
// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
// 	assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
// 		instance, err := instanceRepo.Get(CTX,
// 			instanceRepo.NameCondition(database.TextOperationEqual, instanceName),
// 		)
// 		// event instance.removed
// 		require.Nil(t, instance)
// 		require.NoError(ttt, err)
// 	}, retryDuration, tick)
// }
