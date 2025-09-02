//go:build integration

package events_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel/internal/integration"
)

func TestServer_TestLoginSettingsReduces(t *testing.T) {
	// instanceID := Instance.ID()

	// orgID := Instance.DefaultOrg.Id

	t.Run("test adding login settings reduces", func(t *testing.T) {
		newInstance := integration.NewInstance(t.Context())
		fmt.Printf("[DEBUGPRINT] [:1] newInstance = %+v\n", newInstance)

		// beforeCreate := time.Now()
		// addOIDC, err := MgmtClient.AddOrgOIDCIDP(CTX, &management.AddOrgOIDCIDPRequest{
		// 	Name:               name,
		// 	StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
		// 	ClientId:           "clientID",
		// 	ClientSecret:       "clientSecret",
		// 	Issuer:             "issuer",
		// 	Scopes:             []string{"scope"},
		// 	DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
		// 	UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
		// 	AutoRegister:       true,
		// })
		// afterCreate := time.Now()
		// require.NoError(t, err)

		// idpRepo := repository.IDProviderRepository(pool)

		// retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		// assert.EventuallyWithT(t, func(t *assert.CollectT) {
		// 	idp, err := idpRepo.Get(CTX,
		// 		idpRepo.NameCondition(name),
		// 		instanceID,
		// 		&orgID,
		// 	)
		// 	require.NoError(t, err)

		// 	// event org.idp.config.added
		// 	assert.Equal(t, instanceID, idp.InstanceID)
		// 	assert.Equal(t, orgID, *idp.OrgID)
		// 	assert.Equal(t, addOIDC.IdpId, idp.ID)
		// 	assert.Equal(t, domain.IDPStateActive.String(), idp.State)
		// 	assert.Equal(t, name, idp.Name)
		// 	assert.Equal(t, true, idp.AutoRegister)
		// 	assert.Equal(t, true, idp.AllowCreation)
		// 	assert.Equal(t, false, idp.AllowAutoUpdate)
		// 	assert.Equal(t, true, idp.AllowLinking)
		// 	assert.Equal(t, domain.IDPAutoLinkingOptionUnspecified.String(), idp.AllowAutoLinking)
		// 	assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), *idp.StylingType)
		// 	assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		// 	assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
		// }, retryDuration, tick)
	})
}
