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

	t.Run("test iam idp add reduces", func(t *testing.T) {
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

			// event iam.idp.config.added
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, instanceID, idp.InstanceID)
			assert.Equal(t, domain.IDPStateActive.String(), idp.State)
			assert.Equal(t, true, idp.AutoRegister)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), *idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test iam idp update reduces", func(t *testing.T) {
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
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_UNSPECIFIED), *idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test iam idp deactivate reduces", func(t *testing.T) {
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

	t.Run("test iam idp reactivate reduces", func(t *testing.T) {
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

	t.Run("test iam idp remove reduces", func(t *testing.T) {
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

	t.Run("test iam idp oidc addded reduces", func(t *testing.T) {
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

	t.Run("test iam idp oidc changed reduces", func(t *testing.T) {
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

	t.Run("test iam idp jwt addded reduces", func(t *testing.T) {
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

			// event iam.idp.jwt.config.added
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

	t.Run("test iam idp jwt changed reduces", func(t *testing.T) {
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

			// event iam.idp.jwt.config.changed
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

	t.Run("test instance idp oauth added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oauth
		beforeCreate := time.Now().Add(-1 * time.Second)
		addOAuth, err := AdminClient.AddGenericOAuthProvider(CTX, &admin.AddGenericOAuthProviderRequest{
			Name:                  name,
			ClientId:              "clientId",
			ClientSecret:          "clientSecret",
			AuthorizationEndpoint: "authoizationEndpoint",
			TokenEndpoint:         "tokenEndpoint",
			UserEndpoint:          "userEndpoint",
			Scopes:                []string{"scope"},
			IdAttribute:           "idAttribute",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			UsePkce: false,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for oauth
		var oauth *domain.IDPOAuth
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oauth, err = idpRepo.GetOAuth(CTX, idpRepo.IDCondition(addOAuth.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.oauth.added
			// idp
			assert.Equal(t, addOAuth.Id, oauth.IdentityProvider.ID)
			assert.Equal(t, domain.IDPTypeOAuth.String(), oauth.Type)

			// oauth
			assert.Equal(t, addOAuth.Id, oauth.IdentityProvider.ID)
			assert.Equal(t, "clientId", oauth.ClientID)
			assert.NotNil(t, oauth.ClientSecret)
			assert.Equal(t, "authoizationEndpoint", oauth.AuthorizationEndpoint)
			assert.Equal(t, "authoizationEndpoint", oauth.AuthorizationEndpoint)
			assert.Equal(t, "tokenEndpoint", oauth.TokenEndpoint)
			assert.Equal(t, "userEndpoint", oauth.UserEndpoint)
			assert.Equal(t, "userEndpoint", oauth.UserEndpoint)
			assert.Equal(t, []string{"scope"}, oauth.Scopes)
			assert.Equal(t, false, oauth.AllowLinking)
			assert.Equal(t, false, oauth.AllowCreation)
			assert.Equal(t, false, oauth.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), oauth.AllowAutoLinking)
			assert.Equal(t, false, oauth.UsePKCE)
			assert.WithinRange(t, oauth.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, oauth.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test instanceidp oauth changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oauth
		addOAuth, err := AdminClient.AddGenericOAuthProvider(CTX, &admin.AddGenericOAuthProviderRequest{
			Name:                  name,
			ClientId:              "clientId",
			ClientSecret:          "clientSecret",
			AuthorizationEndpoint: "authoizationEndpoint",
			TokenEndpoint:         "tokenEndpoint",
			UserEndpoint:          "userEndpoint",
			Scopes:                []string{"scope"},
			IdAttribute:           "idAttribute",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			UsePkce: false,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for oauth
		var oauth *domain.IDPOAuth
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oauth, err = idpRepo.GetOAuth(CTX, idpRepo.IDCondition(addOAuth.Id), instanceID, nil)
			require.NoError(t, err)
		}, retryDuration, tick)

		name = "new_" + name
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateGenericOAuthProvider(CTX, &admin.UpdateGenericOAuthProviderRequest{
			Id:                    addOAuth.Id,
			Name:                  name,
			ClientId:              "new_clientId",
			ClientSecret:          "new_clientSecret",
			AuthorizationEndpoint: "new_authoizationEndpoint",
			TokenEndpoint:         "new_tokenEndpoint",
			UserEndpoint:          "new_userEndpoint",
			Scopes:                []string{"new_scope"},
			IdAttribute:           "new_idAttribute",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
			UsePkce: true,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateOauth, err := idpRepo.GetOAuth(CTX,
				idpRepo.IDCondition(addOAuth.Id),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event instance.idp.oauth.changed
			// idp
			assert.Equal(t, addOAuth.Id, oauth.IdentityProvider.ID)
			assert.Equal(t, domain.IDPTypeOAuth.String(), oauth.Type)

			// oauth
			assert.Equal(t, addOAuth.Id, updateOauth.IdentityProvider.ID)
			assert.Equal(t, "new_clientId", updateOauth.ClientID)
			assert.NotEqual(t, oauth.ClientSecret, updateOauth.ClientSecret)
			assert.Equal(t, "new_authoizationEndpoint", updateOauth.AuthorizationEndpoint)
			assert.Equal(t, "new_tokenEndpoint", updateOauth.TokenEndpoint)
			assert.Equal(t, "new_userEndpoint", updateOauth.UserEndpoint)
			assert.Equal(t, []string{"new_scope"}, updateOauth.Scopes)
			assert.Equal(t, true, updateOauth.AllowLinking)
			assert.Equal(t, true, updateOauth.AllowCreation)
			assert.Equal(t, true, updateOauth.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateOauth.AllowAutoLinking)
			assert.Equal(t, true, updateOauth.UsePKCE)
			assert.WithinRange(t, updateOauth.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test instance idp oidc added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oidc
		beforeCreate := time.Now().Add(-1 * time.Second)
		addOIDC, err := AdminClient.AddGenericOIDCProvider(CTX, &admin.AddGenericOIDCProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			Issuer:       "issuer",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			IsIdTokenMapping: false,
			UsePkce:          false,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for oidc
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err := idpRepo.GetOIDC(CTX, idpRepo.IDCondition(addOIDC.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.oidc added
			// idp
			assert.Equal(t, addOIDC.Id, oidc.ID)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)

			// oidc
			assert.Equal(t, addOIDC.Id, oidc.ID)
			assert.Equal(t, "clientId", oidc.ClientID)
			// assert.NotNil(t, oidc.ClientSecret)
			// assert.Equal(t, "authoizationEndpoint", oidc.AuthorizationEndpoint)
			// assert.Equal(t, "tokenEndpoint", oidc.TokenEndpoint)
			// assert.Equal(t, "userEndpoint", oidc.UserEndpoint)
			// assert.Equal(t, "userEndpoint", oidc.UserEndpoint)
			assert.Equal(t, []string{"scope"}, oidc.Scopes)
			assert.Equal(t, "issuer", oidc.Issuer)
			assert.Equal(t, false, oidc.IsIDTokenMapping)
			assert.Equal(t, false, oidc.AllowLinking)
			assert.Equal(t, false, oidc.AllowCreation)
			assert.Equal(t, false, oidc.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), oidc.AllowAutoLinking)
			assert.Equal(t, false, oidc.UsePKCE)
			assert.WithinRange(t, oidc.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, oidc.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test instanceidp oidc changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOIDC, err := AdminClient.AddGenericOIDCProvider(CTX, &admin.AddGenericOIDCProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			Issuer:       "issuer",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			IsIdTokenMapping: false,
			UsePkce:          false,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for oidc
		var oidc *domain.IDPOIDC
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(CTX, idpRepo.IDCondition(addOIDC.Id), instanceID, nil)
			require.NoError(t, err)
		}, retryDuration, tick)

		name = "new_" + name
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateGenericOIDCProvider(CTX, &admin.UpdateGenericOIDCProviderRequest{
			Id:           addOIDC.Id,
			Name:         name,
			Issuer:       "new_issuer",
			ClientId:     "new_clientId",
			ClientSecret: "new_clientSecret",
			Scopes:       []string{"new_scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
			IsIdTokenMapping: true,
			UsePkce:          true,
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateOIDC, err := idpRepo.GetOIDC(CTX,
				idpRepo.IDCondition(addOIDC.Id),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event instance.idp.oidc.changed
			// idp
			assert.Equal(t, addOIDC.Id, oidc.ID)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)

			// oidc
			assert.Equal(t, addOIDC.Id, updateOIDC.ID)
			assert.Equal(t, "new_clientId", updateOIDC.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, updateOIDC.ClientSecret)
			// assert.Equal(t, "new_authoizationEndpoint", updateOIDC.AuthorizationEndpoint)
			// assert.Equal(t, "new_tokenEndpoint", updateOIDC.TokenEndpoint)
			assert.Equal(t, []string{"new_scope"}, updateOIDC.Scopes)
			assert.Equal(t, true, updateOIDC.IsIDTokenMapping)
			assert.Equal(t, true, updateOIDC.AllowLinking)
			assert.Equal(t, true, updateOIDC.AllowCreation)
			assert.Equal(t, true, updateOIDC.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateOIDC.AllowAutoLinking)
			assert.Equal(t, true, updateOIDC.UsePKCE)
			assert.WithinRange(t, updateOIDC.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test instance idp oidc migrated azure migration reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// create OIDC
		addOIDC, err := AdminClient.AddGenericOIDCProvider(CTX, &admin.AddGenericOIDCProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			Issuer:       "issuer",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			IsIdTokenMapping: false,
			UsePkce:          false,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var oidc *domain.IDPOIDC
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(CTX, idpRepo.IDCondition(addOIDC.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)
		}, retryDuration, tick)

		beforeCreate := time.Now()
		_, err = AdminClient.MigrateGenericOIDCProvider(CTX, &admin.MigrateGenericOIDCProviderRequest{
			Id: addOIDC.Id,
			Template: &admin.MigrateGenericOIDCProviderRequest_Azure{
				Azure: &admin.AddAzureADProviderRequest{
					Name:         name,
					ClientId:     "new_clientId",
					ClientSecret: "new_clientSecret",
					Tenant: &idp_grpc.AzureADTenant{
						Type: &idp_grpc.AzureADTenant_TenantType{
							TenantType: idp.AzureADTenantType_AZURE_AD_TENANT_TYPE_ORGANISATIONS,
						},
					},
					EmailVerified: true,
					Scopes:        []string{"new_scope"},
					ProviderOptions: &idp_grpc.Options{
						IsLinkingAllowed:  true,
						IsCreationAllowed: true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
						AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
					},
				},
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			azure, err := idpRepo.GetOAzureAD(CTX, idpRepo.IDCondition(addOIDC.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.oidc.migrated.azure
			// idp
			assert.Equal(t, addOIDC.Id, azure.IdentityProvider.ID)
			assert.Equal(t, name, azure.IdentityProvider.Name)

			// oidc
			assert.Equal(t, "new_clientId", azure.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, azure.ClientSecret)
			// type = azure
			assert.Equal(t, domain.AzureTenantTypeOrganizations.String(), azure.Tenant)
			assert.Equal(t, domain.IDPTypeAzure.String(), azure.Type)
			assert.Equal(t, true, azure.IsEmailVerified)
			assert.Equal(t, []string{"new_scope"}, azure.Scopes)
			assert.Equal(t, true, azure.AllowLinking)
			assert.Equal(t, true, azure.AllowCreation)
			assert.Equal(t, true, azure.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), azure.AllowAutoLinking)
			assert.WithinRange(t, azure.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test instance idp oidc migrated google migration reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// create OIDC
		addOIDC, err := AdminClient.AddGenericOIDCProvider(CTX, &admin.AddGenericOIDCProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			Issuer:       "issuer",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			IsIdTokenMapping: false,
			UsePkce:          false,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var oidc *domain.IDPOIDC
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(CTX, idpRepo.IDCondition(addOIDC.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)
		}, retryDuration, tick)

		beforeCreate := time.Now()
		_, err = AdminClient.MigrateGenericOIDCProvider(CTX, &admin.MigrateGenericOIDCProviderRequest{
			Id: addOIDC.Id,
			Template: &admin.MigrateGenericOIDCProviderRequest_Google{
				Google: &admin.AddGoogleProviderRequest{
					Name:         name,
					ClientId:     "new_clientId",
					ClientSecret: "new_clientSecret",
					Scopes:       []string{"new_scope"},
					ProviderOptions: &idp_grpc.Options{
						IsLinkingAllowed:  true,
						IsCreationAllowed: true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
						AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
					},
				},
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			google, err := idpRepo.GetGoogle(CTX, idpRepo.IDCondition(addOIDC.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.oidc.migrated.google
			// idp
			assert.Equal(t, addOIDC.Id, google.IdentityProvider.ID)
			assert.Equal(t, name, google.IdentityProvider.Name)

			// oidc
			assert.Equal(t, "new_clientId", google.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, google.ClientSecret)
			// type = google
			assert.Equal(t, domain.IDPTypeGoogle.String(), google.Type)
			assert.Equal(t, []string{"new_scope"}, google.Scopes)
			assert.Equal(t, true, google.AllowLinking)
			assert.Equal(t, true, google.AllowCreation)
			assert.Equal(t, true, google.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), google.AllowAutoLinking)
			assert.WithinRange(t, google.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test instance idp jwt added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add jwt
		beforeCreate := time.Now().Add(-1 * time.Second)
		addJWT, err := AdminClient.AddJWTProvider(CTX, &admin.AddJWTProviderRequest{
			Name:         name,
			Issuer:       "issuer",
			JwtEndpoint:  "jwtEndpoint",
			KeysEndpoint: "keyEndpoint",
			HeaderName:   "headerName",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for jwt
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			jwt, err := idpRepo.GetJWT(CTX, idpRepo.IDCondition(addJWT.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.jwt.added
			// idp
			assert.Equal(t, addJWT.Id, jwt.ID)
			assert.Equal(t, domain.IDPTypeJWT.String(), jwt.Type)

			// jwt
			assert.Equal(t, addJWT.Id, jwt.ID)
			assert.Equal(t, "jwtEndpoint", jwt.JWTEndpoint)
			assert.Equal(t, "issuer", jwt.Issuer)
			assert.Equal(t, "keyEndpoint", jwt.KeysEndpoint)
			assert.Equal(t, "headerName", jwt.HeaderName)

			assert.Equal(t, false, jwt.AllowLinking)
			assert.Equal(t, false, jwt.AllowCreation)
			assert.Equal(t, false, jwt.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), jwt.AllowAutoLinking)
			assert.WithinRange(t, jwt.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, jwt.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})
}
