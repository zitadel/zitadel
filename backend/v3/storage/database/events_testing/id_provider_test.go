//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	idp_grpc "github.com/zitadel/zitadel/pkg/grpc/idp"
)

func TestServer_TestIDProviderReduces(t *testing.T) {
	instanceID := Instance.ID()

	t.Run("test idp add reduces", func(t *testing.T) {
		name := gofakeit.Name()

		beforeCreate := time.Now()
		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)
		afterCreate := time.Now()

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.NameCondition(name),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event iam.idp.config.added
			assert.Equal(t, addOCID.IdpId, idp.ID)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, instanceID, idp.InstanceID)
			assert.Equal(t, domain.IDPStateActive.String(), idp.State)
			assert.Equal(t, true, idp.AllowAutoCreation)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp update reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)

		name = "new_" + name

		beforeCreate := time.Now()
		_, err = AdminClient.UpdateIDP(CTX, &admin.UpdateIDPRequest{
			IdpId:        addOCID.IdpId,
			Name:         name,
			StylingType:  idp_grpc.IDPStylingType_STYLING_TYPE_UNSPECIFIED,
			AutoRegister: false,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.NameCondition(name),
				// idpRepo.IDCondition(addOCID.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event "iam.idp.config.changed"
			assert.Equal(t, addOCID.IdpId, idp.ID)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, false, idp.AllowAutoCreation)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_UNSPECIFIED), idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp deactivate reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)

		// deactivate idp
		beforeCreate := time.Now()
		_, err = AdminClient.DeactivateIDP(CTX, &admin.DeactivateIDPRequest{
			IdpId: addOCID.IdpId,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				// idpRepo.NameCondition(name),
				idpRepo.IDCondition(addOCID.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event "iam.idp.config.deactivated"
			assert.Equal(t, addOCID.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateInactive.String(), idp.State)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp reactivate reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// deactivate idp
		_, err = AdminClient.DeactivateIDP(CTX, &admin.DeactivateIDPRequest{
			IdpId: addOCID.IdpId,
		})
		require.NoError(t, err)
		// wait for idp to be deactivated
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.IDCondition(addOCID.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			assert.Equal(t, addOCID.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateInactive.String(), idp.State)
		}, retryDuration, tick)

		// reactivate idp
		// beforeCreate := time.Now().Add(-time.Second)
		beforeCreate := time.Now()
		_, err = AdminClient.ReactivateIDP(CTX, &admin.ReactivateIDPRequest{
			IdpId: addOCID.IdpId,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				// idpRepo.NameCondition(name),
				idpRepo.IDCondition(addOCID.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event "iam.idp.config.reactivated"
			assert.Equal(t, addOCID.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateActive.String(), idp.State)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp remove reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add idp
		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// remove idp
		_, err = AdminClient.RemoveIDP(CTX, &admin.RemoveIDPRequest{
			IdpId: addOCID.IdpId,
		})
		require.NoError(t, err)

		// retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Delete(CTX,
				// idpRepo.NameCondition(name),
				idpRepo.IDCondition(addOCID.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event "iam.idp.config.remove"
			assert.Nil(t, idp)
		}, retryDuration, tick)
	})

	t.Run("test idp oidc addded reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add idp
		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// remove idp
		_, err = AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)

		// retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Delete(CTX,
				// idpRepo.NameCondition(name),
				idpRepo.IDCondition(addOCID.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event "iam.idp.config.remove"
			assert.Nil(t, idp)
		}, retryDuration, tick)
	})
}
