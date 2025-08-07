//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
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
		addOIDC, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.NameCondition(name),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event iam.idp.config.added
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, instanceID, idp.InstanceID)
			assert.Equal(t, domain.IDPStateActive.String(), idp.State)
			assert.Equal(t, true, idp.AutoRegister)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp update reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOIDC, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
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
			IdpId:        addOIDC.IdpId,
			Name:         name,
			StylingType:  idp_grpc.IDPStylingType_STYLING_TYPE_UNSPECIFIED,
			AutoRegister: false,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.NameCondition(name),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event iam.idp.config.changed
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, false, idp.AutoRegister)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_UNSPECIFIED), idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp deactivate reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOIDC, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
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
			IdpId: addOIDC.IdpId,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event iam.idp.config.deactivated
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateInactive.String(), idp.State)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp reactivate reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOIDC, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
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
			IdpId: addOIDC.IdpId,
		})
		require.NoError(t, err)
		// wait for idp to be deactivated
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateInactive.String(), idp.State)
		}, retryDuration, tick)

		// reactivate idp
		beforeCreate := time.Now()
		_, err = AdminClient.ReactivateIDP(CTX, &admin.ReactivateIDPRequest{
			IdpId: addOIDC.IdpId,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event iam.idp.config.reactivated
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateActive.String(), idp.State)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test idp remove reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add idp
		addOIDC, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
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
			IdpId: addOIDC.IdpId,
		})
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := idpRepo.Get(CTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				nil,
			)

			// event iam.idp.config.remove
			require.ErrorIs(t, &database.NoRowFoundError{}, err)
		}, retryDuration, tick)
	})

	t.Run("test idp oidc addded reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oidc
		addOIDC, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err := idpRepo.GetOIDC(CTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event org.idp.oidc.config.added
			// idp
			assert.Equal(t, addOIDC.IdpId, oidc.ID)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)

			// oidc
			assert.Equal(t, addOIDC.IdpId, oidc.IDPConfigID)
			assert.Equal(t, "issuer", oidc.Issuer)
			assert.Equal(t, "clientID", oidc.ClientID)
			assert.Equal(t, []string{"scope"}, oidc.Scopes)
			assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL), oidc.IDPDisplayNameMapping)
			assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL), oidc.UserNameMapping)
		}, retryDuration, tick)
	})

	t.Run("test idp oidc changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oidc
		addOIDC, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
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

		// check original values for OCID
		var oidc *domain.IDPOIDC
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(CTX, idpRepo.IDCondition(addOIDC.IdpId), instanceID, nil)
			require.NoError(t, err)
		}, retryDuration, tick)

		// idp
		assert.Equal(t, addOIDC.IdpId, oidc.ID)
		assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)

		// oidc
		assert.Equal(t, addOIDC.IdpId, oidc.IDPConfigID)
		assert.Equal(t, "issuer", oidc.Issuer)
		assert.Equal(t, "clientID", oidc.ClientID)
		assert.Equal(t, []string{"scope"}, oidc.Scopes)
		assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL), oidc.IDPDisplayNameMapping)
		assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL), oidc.UserNameMapping)

		beforeCreate := time.Now()
		_, err = AdminClient.UpdateIDPOIDCConfig(CTX, &admin.UpdateIDPOIDCConfigRequest{
			IdpId:              addOIDC.IdpId,
			ClientId:           "new_clientID",
			ClientSecret:       "new_clientSecret",
			Issuer:             "new_issuer",
			Scopes:             []string{"new_scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateOIDC, err := idpRepo.GetOIDC(CTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event org.idp.oidc.config.changed
			// idp
			assert.Equal(t, addOIDC.IdpId, updateOIDC.ID)
			assert.Equal(t, domain.IDPTypeOIDC.String(), updateOIDC.Type)
			assert.WithinRange(t, updateOIDC.UpdatedAt, beforeCreate, afterCreate)

			// oidc
			assert.Equal(t, addOIDC.IdpId, updateOIDC.IDPConfigID)
			assert.Equal(t, "new_issuer", updateOIDC.Issuer)
			assert.Equal(t, "new_clientID", updateOIDC.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, updateOIDC.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateOIDC.Scopes)
			assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME), updateOIDC.IDPDisplayNameMapping)
			assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME), updateOIDC.UserNameMapping)
		}, retryDuration, tick)
	})

	t.Run("test idp jwt addded reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add jwt
		addJWT, err := AdminClient.AddJWTIDP(CTX, &admin.AddJWTIDPRequest{
			Name:         name,
			StylingType:  idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			JwtEndpoint:  "jwtEndpoint",
			Issuer:       "issuer",
			KeysEndpoint: "keyEndpoint",
			HeaderName:   "headerName",
			AutoRegister: true,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			jwt, err := idpRepo.GetJWT(CTX,
				idpRepo.IDCondition(addJWT.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event org.idp.jwt.config.added
			// idp
			assert.Equal(t, addJWT.IdpId, jwt.ID)
			assert.Equal(t, domain.IDPTypeJWT.String(), jwt.Type)

			// jwt
			assert.Equal(t, addJWT.IdpId, jwt.IDPConfigID)
			assert.Equal(t, "jwtEndpoint", jwt.JWTEndpoint)
			assert.Equal(t, "issuer", jwt.Issuer)
			assert.Equal(t, "keyEndpoint", jwt.KeysEndpoint)
			assert.Equal(t, "headerName", jwt.HeaderName)
		}, retryDuration, tick)
	})

	t.Run("test idp jwt changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add jwt
		addJWT, err := AdminClient.AddJWTIDP(CTX, &admin.AddJWTIDPRequest{
			Name:         name,
			StylingType:  idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			JwtEndpoint:  "jwtEndpoint",
			Issuer:       "issuer",
			KeysEndpoint: "keyEndpoint",
			HeaderName:   "headerName",
			AutoRegister: true,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check original values for jwt
		var jwt *domain.IDPJWT
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			jwt, err = idpRepo.GetJWT(CTX, idpRepo.IDCondition(addJWT.IdpId), instanceID, nil)
			require.NoError(t, err)
		}, retryDuration, tick)

		// idp
		assert.Equal(t, addJWT.IdpId, jwt.ID)
		assert.Equal(t, domain.IDPTypeJWT.String(), jwt.Type)

		// jwt
		assert.Equal(t, addJWT.IdpId, jwt.IDPConfigID)
		assert.Equal(t, "jwtEndpoint", jwt.JWTEndpoint)
		assert.Equal(t, "issuer", jwt.Issuer)
		assert.Equal(t, "keyEndpoint", jwt.KeysEndpoint)
		assert.Equal(t, "headerName", jwt.HeaderName)

		beforeCreate := time.Now()
		_, err = AdminClient.UpdateIDPJWTConfig(CTX, &admin.UpdateIDPJWTConfigRequest{
			IdpId:        addJWT.IdpId,
			JwtEndpoint:  "new_jwtEndpoint",
			Issuer:       "new_issuer",
			KeysEndpoint: "new_keyEndpoint",
			HeaderName:   "new_headerName",
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateJWT, err := idpRepo.GetJWT(CTX,
				idpRepo.IDCondition(addJWT.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event org.idp.jwt.config.changed
			// idp
			assert.Equal(t, addJWT.IdpId, updateJWT.ID)
			assert.Equal(t, domain.IDPTypeJWT.String(), updateJWT.Type)
			assert.WithinRange(t, updateJWT.UpdatedAt, beforeCreate, afterCreate)

			// jwt
			assert.Equal(t, addJWT.IdpId, updateJWT.IDPConfigID)
			assert.Equal(t, "new_jwtEndpoint", updateJWT.JWTEndpoint)
			assert.Equal(t, "new_issuer", updateJWT.Issuer)
			assert.Equal(t, "new_keyEndpoint", updateJWT.KeysEndpoint)
		}, retryDuration, tick)
	})
}
