//go:build integration

package events_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestServer_TestIDProviderReduces(t *testing.T) {
	// instanceID := Instance.ID()

	t.Run("test org add reduces", func(t *testing.T) {
		// beforeCreate := time.Now()
		orgName := gofakeit.Name()

		// create org
		_, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
			Name: orgName,
		})
		require.NoError(t, err)
		// afterCreate := time.Now()

		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               gofakeit.Name(),
			StylingType:        idp.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> addOCID = %+v\n", addOCID)
		fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> err = %+v\n", err)

		// idpRepo := repository.IDProviderRepository(pool)

		// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		// 	assert.EventuallyWithT(t, func(tt *assert.CollectT) {
		// 		organization, err := idpRepo.Get(CTX,
		// 			idpRepo.NameCondition(orgName),
		// 			instanceID,
		// 		)
		// 		require.NoError(tt, err)

		// 		// event org.added
		// 		assert.NotNil(t, organization.ID)
		// 		assert.Equal(t, orgName, organization.Name)
		// 		assert.NotNil(t, organization.InstanceID)
		// 		assert.Equal(t, domain.OrgStateActive.String(), organization.State)
		// 		assert.WithinRange(t, organization.CreatedAt, beforeCreate, afterCreate)
		// 		assert.WithinRange(t, organization.UpdatedAt, beforeCreate, afterCreate)
		// 	}, retryDuration, tick)
	})
}
