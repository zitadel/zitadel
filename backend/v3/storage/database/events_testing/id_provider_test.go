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
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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
			assert.Equal(t, instanceID, idp.InstanceID)
			assert.Nil(t, idp.OrgID)
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateActive.String(), idp.State)
			assert.Equal(t, name, idp.Name)
			// assert.Equal(t, domain.IDPTypeUnspecified.String(), idp.Type)
			assert.Equal(t, true, idp.AutoRegister)
			assert.Equal(t, true, idp.AllowCreation)
			assert.Equal(t, false, idp.AllowAutoUpdate)
			assert.Equal(t, true, idp.AllowLinking)
			assert.Equal(t, domain.IDPAutoLinkingOptionUnspecified.String(), idp.AllowAutoLinking)
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
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

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
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

	t.Run("test iam idp oidc added reduces", func(t *testing.T) {
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
			AutoRegister:       false,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err := idpRepo.GetOIDC(CTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event org.idp.oidc.config.added
			// idp
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Nil(t, oidc.OrgID)
			assert.Equal(t, name, oidc.Name)
			assert.Equal(t, addOIDC.IdpId, oidc.ID)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)

			// oidc
			assert.Equal(t, "issuer", oidc.Issuer)
			assert.Equal(t, "clientID", oidc.ClientID)
			assert.Equal(t, []string{"scope"}, oidc.Scopes)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), *oidc.StylingType)
			assert.Equal(t, false, oidc.AutoRegister)
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(CTX, idpRepo.IDCondition(addOIDC.IdpId), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addOIDC.IdpId, oidc.ID)
		}, retryDuration, tick)

		// // idp
		// assert.Equal(t, addOIDC.IdpId, oidc.ID)
		// assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)

		// // oidc
		// assert.Equal(t, instanceID, oidc.InstanceID)
		// assert.Nil(t, oidc.OrgID)
		// assert.Equal(t, "issuer", oidc.Issuer)
		// assert.Equal(t, "clientID", oidc.ClientID)
		// assert.Equal(t, []string{"scope"}, oidc.Scopes)
		// assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL), oidc.IDPDisplayNameMapping)
		// assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL), oidc.UserNameMapping)

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
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Nil(t, oidc.OrgID)
			assert.Equal(t, name, oidc.Name)
			assert.Equal(t, addOIDC.IdpId, updateOIDC.ID)
			assert.Equal(t, domain.IDPTypeOIDC.String(), updateOIDC.Type)
			assert.WithinRange(t, updateOIDC.UpdatedAt, beforeCreate, afterCreate)

			// oidc
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Nil(t, oidc.OrgID)
			assert.Equal(t, "new_issuer", updateOIDC.Issuer)
			assert.Equal(t, "new_clientID", updateOIDC.ClientID)
			assert.NotNil(t, oidc.ClientSecret)
			assert.NotEqual(t, oidc.ClientSecret, updateOIDC.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateOIDC.Scopes)
			assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME), updateOIDC.IDPDisplayNameMapping)
			assert.Equal(t, domain.OIDCMappingField(idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME), updateOIDC.UserNameMapping)
		}, retryDuration, tick)
	})

	t.Run("test iam idp jwt added reduces", func(t *testing.T) {
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
			assert.Equal(t, instanceID, jwt.InstanceID)
			assert.Nil(t, jwt.OrgID)
			assert.Equal(t, name, jwt.Name)
			assert.Equal(t, addJWT.IdpId, jwt.ID)
			assert.Equal(t, domain.IDPTypeJWT.String(), jwt.Type)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), *jwt.StylingType)

			// jwt
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
		// var jwt *domain.IDPJWT
		// retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		// assert.EventuallyWithT(t, func(t *assert.CollectT) {
		// 	jwt, err = idpRepo.GetJWT(CTX, idpRepo.IDCondition(addJWT.IdpId), instanceID, nil)
		// 	require.NoError(t, err)
		// 	assert.Equal(t, addJWT.IdpId, jwt.ID)
		// }, retryDuration, tick)

		// // idp
		// assert.Equal(t, addJWT.IdpId, jwt.ID)
		// assert.Equal(t, domain.IDPTypeJWT.String(), jwt.Type)

		// // jwt
		// assert.Equal(t, "jwtEndpoint", jwt.JWTEndpoint)
		// assert.Equal(t, "issuer", jwt.Issuer)
		// assert.Equal(t, "keyEndpoint", jwt.KeysEndpoint)
		// assert.Equal(t, "headerName", jwt.HeaderName)

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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
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
			assert.Equal(t, "new_jwtEndpoint", updateJWT.JWTEndpoint)
			assert.Equal(t, "new_issuer", updateJWT.Issuer)
			assert.Equal(t, "new_keyEndpoint", updateJWT.KeysEndpoint)
			assert.Equal(t, "new_headerName", updateJWT.HeaderName)
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
			assert.Equal(t, instanceID, oauth.InstanceID)
			assert.Nil(t, oauth.OrgID)
			assert.Equal(t, addOAuth.Id, oauth.ID)
			assert.Equal(t, name, oauth.Name)
			assert.Equal(t, domain.IDPTypeOAuth.String(), oauth.Type)
			assert.Equal(t, false, oauth.AllowLinking)
			assert.Equal(t, false, oauth.AllowCreation)
			assert.Equal(t, false, oauth.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), oauth.AllowAutoLinking)
			assert.WithinRange(t, oauth.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, oauth.UpdatedAt, beforeCreate, afterCreate)

			// oauth
			assert.Equal(t, "clientId", oauth.ClientID)
			assert.NotNil(t, oauth.ClientSecret)
			assert.Equal(t, "authoizationEndpoint", oauth.AuthorizationEndpoint)
			assert.Equal(t, "tokenEndpoint", oauth.TokenEndpoint)
			assert.Equal(t, "userEndpoint", oauth.UserEndpoint)
			assert.Equal(t, []string{"scope"}, oauth.Scopes)
			assert.Equal(t, "idAttribute", oauth.IDAttribute)
			assert.Equal(t, false, oauth.UsePKCE)
		}, retryDuration, tick)
	})

	t.Run("test instance idp oauth changed reduces", func(t *testing.T) {
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
			assert.Equal(t, addOAuth.Id, oauth.ID)
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
			assert.Equal(t, instanceID, oauth.InstanceID)
			assert.Nil(t, oauth.OrgID)
			assert.Equal(t, addOAuth.Id, updateOauth.ID)
			assert.Equal(t, name, updateOauth.Name)
			assert.Equal(t, domain.IDPTypeOAuth.String(), oauth.Type)
			assert.Equal(t, true, updateOauth.AllowLinking)
			assert.Equal(t, true, updateOauth.AllowCreation)
			assert.Equal(t, true, updateOauth.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateOauth.AllowAutoLinking)
			assert.Equal(t, true, updateOauth.UsePKCE)
			assert.WithinRange(t, updateOauth.UpdatedAt, beforeCreate, afterCreate)

			// oauth
			assert.Equal(t, "new_clientId", updateOauth.ClientID)
			assert.NotEqual(t, oauth.ClientSecret, updateOauth.ClientSecret)
			assert.Equal(t, "new_authoizationEndpoint", updateOauth.AuthorizationEndpoint)
			assert.Equal(t, "new_tokenEndpoint", updateOauth.TokenEndpoint)
			assert.Equal(t, "new_userEndpoint", updateOauth.UserEndpoint)
			assert.Equal(t, []string{"new_scope"}, updateOauth.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp oidc added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oidc
		beforeCreate := time.Now()
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
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Nil(t, oidc.OrgID)
			assert.Equal(t, addOIDC.Id, oidc.ID)
			assert.Equal(t, name, oidc.Name)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)
			assert.Equal(t, false, oidc.AllowLinking)
			assert.Equal(t, false, oidc.AllowCreation)
			assert.Equal(t, false, oidc.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), oidc.AllowAutoLinking)
			assert.WithinRange(t, oidc.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, oidc.UpdatedAt, beforeCreate, afterCreate)

			// oidc
			assert.Equal(t, addOIDC.Id, oidc.ID)
			assert.Equal(t, "clientId", oidc.ClientID)
			assert.NotNil(t, oidc.ClientSecret)
			assert.Equal(t, []string{"scope"}, oidc.Scopes)
			assert.Equal(t, "issuer", oidc.Issuer)
			assert.Equal(t, false, oidc.IsIDTokenMapping)
			assert.Equal(t, false, oidc.UsePKCE)
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
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Nil(t, oidc.OrgID)
			assert.Equal(t, addOIDC.Id, oidc.ID)
			assert.Equal(t, name, updateOIDC.Name)
			assert.Equal(t, domain.IDPTypeOIDC.String(), oidc.Type)
			assert.Equal(t, true, updateOIDC.AllowLinking)
			assert.Equal(t, true, updateOIDC.AllowCreation)
			assert.Equal(t, true, updateOIDC.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateOIDC.AllowAutoLinking)
			assert.WithinRange(t, updateOIDC.UpdatedAt, beforeCreate, afterCreate)

			// oidc
			assert.Equal(t, addOIDC.Id, updateOIDC.ID)
			assert.Equal(t, "new_clientId", updateOIDC.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, updateOIDC.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateOIDC.Scopes)
			assert.Equal(t, true, updateOIDC.IsIDTokenMapping)
			assert.Equal(t, true, updateOIDC.UsePKCE)
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
			assert.Equal(t, instanceID, azure.InstanceID)
			assert.Nil(t, azure.OrgID)
			assert.Equal(t, addOIDC.Id, azure.ID)
			assert.Equal(t, name, azure.Name)
			// type = azure
			assert.Equal(t, domain.IDPTypeAzure.String(), azure.Type)
			assert.Equal(t, true, azure.AllowLinking)
			assert.Equal(t, true, azure.AllowCreation)
			assert.Equal(t, true, azure.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), azure.AllowAutoLinking)
			assert.WithinRange(t, azure.UpdatedAt, beforeCreate, afterCreate)

			// oidc
			assert.Equal(t, "new_clientId", azure.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, azure.ClientSecret)
			assert.Equal(t, domain.AzureTenantTypeOrganizations.String(), azure.Tenant)
			assert.Equal(t, true, azure.IsEmailVerified)
			assert.Equal(t, []string{"new_scope"}, azure.Scopes)
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
			assert.Equal(t, instanceID, google.InstanceID)
			assert.Nil(t, google.OrgID)
			assert.Equal(t, addOIDC.Id, google.ID)
			assert.Equal(t, name, google.Name)
			// type = google
			assert.Equal(t, domain.IDPTypeGoogle.String(), google.Type)
			assert.Equal(t, true, google.AllowLinking)
			assert.Equal(t, true, google.AllowCreation)
			assert.Equal(t, true, google.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), google.AllowAutoLinking)
			assert.WithinRange(t, google.UpdatedAt, beforeCreate, afterCreate)

			// oidc
			assert.Equal(t, "new_clientId", google.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, google.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, google.Scopes)
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
			assert.Equal(t, instanceID, jwt.InstanceID)
			assert.Nil(t, jwt.OrgID)
			assert.Equal(t, addJWT.Id, jwt.ID)
			assert.Equal(t, name, jwt.Name)
			assert.Equal(t, domain.IDPTypeJWT.String(), jwt.Type)
			assert.Equal(t, false, jwt.AllowLinking)
			assert.Equal(t, false, jwt.AllowCreation)
			assert.Equal(t, false, jwt.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), jwt.AllowAutoLinking)
			assert.WithinRange(t, jwt.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, jwt.UpdatedAt, beforeCreate, afterCreate)

			// jwt
			assert.Equal(t, "jwtEndpoint", jwt.JWTEndpoint)
			assert.Equal(t, "issuer", jwt.Issuer)
			assert.Equal(t, "keyEndpoint", jwt.KeysEndpoint)
			assert.Equal(t, "headerName", jwt.HeaderName)
		}, retryDuration, tick)
	})

	t.Run("test instance idp jwt changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add jwt
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
		require.NoError(t, err)

		name = "new_" + name
		// change jwt
		beforeCreate := time.Now().Add(-1 * time.Second)
		_, err = AdminClient.UpdateJWTProvider(CTX, &admin.UpdateJWTProviderRequest{
			Id:           addJWT.Id,
			Name:         name,
			Issuer:       "new_issuer",
			JwtEndpoint:  "new_jwtEndpoint",
			KeysEndpoint: "new_keyEndpoint",
			HeaderName:   "new_headerName",
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
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
			assert.Equal(t, instanceID, jwt.InstanceID)
			assert.Nil(t, jwt.OrgID)
			assert.Equal(t, addJWT.Id, jwt.ID)
			assert.Equal(t, name, jwt.Name)
			assert.Equal(t, domain.IDPTypeJWT.String(), jwt.Type)
			assert.Equal(t, true, jwt.AllowLinking)
			assert.Equal(t, true, jwt.AllowCreation)
			assert.Equal(t, true, jwt.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), jwt.AllowAutoLinking)
			assert.WithinRange(t, jwt.UpdatedAt, beforeCreate, afterCreate)

			// jwt
			assert.Equal(t, "new_jwtEndpoint", jwt.JWTEndpoint)
			assert.Equal(t, "new_issuer", jwt.Issuer)
			assert.Equal(t, "new_keyEndpoint", jwt.KeysEndpoint)
			assert.Equal(t, "new_headerName", jwt.HeaderName)
		}, retryDuration, tick)
	})

	t.Run("test instance idp azure added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add azure
		beforeCreate := time.Now()
		addAzure, err := AdminClient.AddAzureADProvider(CTX, &admin.AddAzureADProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Tenant: &idp_grpc.AzureADTenant{
				Type: &idp_grpc.AzureADTenant_TenantType{
					TenantType: idp.AzureADTenantType_AZURE_AD_TENANT_TYPE_ORGANISATIONS,
				},
			},
			EmailVerified: true,
			Scopes:        []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for azure
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			azure, err := idpRepo.GetOAzureAD(CTX, idpRepo.IDCondition(addAzure.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.azure.added
			// idp
			assert.Equal(t, instanceID, azure.InstanceID)
			assert.Nil(t, azure.OrgID)
			assert.Equal(t, addAzure.Id, azure.ID)
			assert.Equal(t, name, azure.Name)
			assert.Equal(t, domain.IDPTypeAzure.String(), azure.Type)
			assert.Equal(t, true, azure.AllowLinking)
			assert.Equal(t, true, azure.AllowCreation)
			assert.Equal(t, true, azure.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), azure.AllowAutoLinking)
			assert.WithinRange(t, azure.UpdatedAt, beforeCreate, afterCreate)

			// azure
			assert.Equal(t, "clientId", azure.ClientID)
			assert.NotNil(t, azure.ClientSecret)
			assert.Equal(t, domain.AzureTenantTypeOrganizations.String(), azure.Tenant)
			assert.Equal(t, true, azure.IsEmailVerified)
			assert.Equal(t, []string{"scope"}, azure.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp azure changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add azure
		addAzure, err := AdminClient.AddAzureADProvider(CTX, &admin.AddAzureADProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Tenant: &idp_grpc.AzureADTenant{
				Type: &idp_grpc.AzureADTenant_TenantType{
					TenantType: idp.AzureADTenantType_AZURE_AD_TENANT_TYPE_ORGANISATIONS,
				},
			},
			EmailVerified: false,
			Scopes:        []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var azure *domain.IDPOAzureAD
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			azure, err = idpRepo.GetOAzureAD(CTX, idpRepo.IDCondition(addAzure.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addAzure.Id, azure.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change azure
		beforeCreate := time.Now().Add(-1 * time.Second)
		_, err = AdminClient.UpdateAzureADProvider(CTX, &admin.UpdateAzureADProviderRequest{
			Id:           addAzure.Id,
			Name:         name,
			ClientId:     "new_clientId",
			ClientSecret: "new_clientSecret",
			Tenant: &idp_grpc.AzureADTenant{
				Type: &idp_grpc.AzureADTenant_TenantType{
					TenantType: idp.AzureADTenantType_AZURE_AD_TENANT_TYPE_CONSUMERS,
				},
			},
			EmailVerified: true,
			Scopes:        []string{"new_scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		// check values for azure
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateAzure, err := idpRepo.GetOAzureAD(CTX, idpRepo.IDCondition(addAzure.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.azure.changed
			// idp
			assert.Equal(t, instanceID, updateAzure.InstanceID)
			assert.Nil(t, updateAzure.OrgID)
			assert.Equal(t, addAzure.Id, updateAzure.ID)
			assert.Equal(t, name, updateAzure.Name)
			assert.Equal(t, domain.IDPTypeAzure.String(), updateAzure.Type)
			assert.Equal(t, true, updateAzure.AllowLinking)
			assert.Equal(t, true, updateAzure.AllowCreation)
			assert.Equal(t, true, updateAzure.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), updateAzure.AllowAutoLinking)
			assert.WithinRange(t, updateAzure.UpdatedAt, beforeCreate, afterCreate)

			// azure
			assert.Equal(t, "new_clientId", updateAzure.ClientID)
			assert.NotEqual(t, azure.ClientSecret, updateAzure.ClientSecret)
			assert.Equal(t, domain.AzureTenantTypeConsumers.String(), updateAzure.Tenant)
			assert.Equal(t, true, updateAzure.IsEmailVerified)
			assert.Equal(t, []string{"new_scope"}, updateAzure.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp github added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github
		beforeCreate := time.Now()
		addGithub, err := AdminClient.AddGitHubProvider(CTX, &admin.AddGitHubProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for github
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			github, err := idpRepo.GetGithub(CTX, idpRepo.IDCondition(addGithub.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.github.added
			// idp
			assert.Equal(t, instanceID, github.InstanceID)
			assert.Nil(t, github.OrgID)
			assert.Equal(t, addGithub.Id, github.ID)
			assert.Equal(t, name, github.Name)
			assert.Equal(t, domain.IDPTypeGitHub.String(), github.Type)
			assert.Equal(t, false, github.AllowLinking)
			assert.Equal(t, false, github.AllowCreation)
			assert.Equal(t, false, github.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), github.AllowAutoLinking)
			assert.WithinRange(t, github.UpdatedAt, beforeCreate, afterCreate)

			assert.Equal(t, "clientId", github.ClientID)
			assert.NotNil(t, github.ClientSecret)
			assert.Equal(t, []string{"scope"}, github.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp github changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github
		addGithub, err := AdminClient.AddGitHubProvider(CTX, &admin.AddGitHubProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var github *domain.IDPGithub
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			github, err = idpRepo.GetGithub(CTX, idpRepo.IDCondition(addGithub.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addGithub.Id, github.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change github
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateGitHubProvider(CTX, &admin.UpdateGitHubProviderRequest{
			Id:           addGithub.Id,
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
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		// check values for azure
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGithub, err := idpRepo.GetGithub(CTX, idpRepo.IDCondition(addGithub.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.github.changed
			// idp
			assert.Equal(t, instanceID, updateGithub.InstanceID)
			assert.Nil(t, updateGithub.OrgID)
			assert.Equal(t, addGithub.Id, updateGithub.ID)
			assert.Equal(t, name, updateGithub.Name)
			assert.Equal(t, domain.IDPTypeGitHub.String(), updateGithub.Type)
			assert.Equal(t, true, updateGithub.AllowLinking)
			assert.Equal(t, true, updateGithub.AllowCreation)
			assert.Equal(t, true, updateGithub.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateGithub.AllowAutoLinking)
			assert.WithinRange(t, updateGithub.UpdatedAt, beforeCreate, afterCreate)

			// github
			assert.Equal(t, "new_clientId", updateGithub.ClientID)
			assert.NotEqual(t, github.ClientSecret, updateGithub.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateGithub.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp github enterprise added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github enterprise
		beforeCreate := time.Now()
		addGithubEnterprise, err := AdminClient.AddGitHubEnterpriseServerProvider(CTX, &admin.AddGitHubEnterpriseServerProviderRequest{
			Name:                  name,
			ClientId:              "clientId",
			ClientSecret:          "clientSecret",
			AuthorizationEndpoint: "authoizationEndpoint",
			TokenEndpoint:         "tokenEndpoint",
			UserEndpoint:          "userEndpoint",
			Scopes:                []string{"scope"},
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

		// check values for github enterprise
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			githubEnterprise, err := idpRepo.GetGithubEnterprise(CTX, idpRepo.IDCondition(addGithubEnterprise.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.github_enterprise.added
			// idp
			assert.Equal(t, instanceID, githubEnterprise.InstanceID)
			assert.Nil(t, githubEnterprise.OrgID)
			assert.Equal(t, addGithubEnterprise.Id, githubEnterprise.ID)
			assert.Equal(t, name, githubEnterprise.Name)
			assert.Equal(t, domain.IDPTypeGitHubEnterprise.String(), githubEnterprise.Type)
			assert.Equal(t, false, githubEnterprise.AllowLinking)
			assert.Equal(t, false, githubEnterprise.AllowCreation)
			assert.Equal(t, false, githubEnterprise.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), githubEnterprise.AllowAutoLinking)
			assert.WithinRange(t, githubEnterprise.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, githubEnterprise.UpdatedAt, beforeCreate, afterCreate)

			// github enterprise
			assert.Equal(t, "clientId", githubEnterprise.ClientID)
			assert.NotNil(t, githubEnterprise.ClientSecret)
			assert.Equal(t, "authoizationEndpoint", githubEnterprise.AuthorizationEndpoint)
			assert.Equal(t, "tokenEndpoint", githubEnterprise.TokenEndpoint)
			assert.Equal(t, "userEndpoint", githubEnterprise.UserEndpoint)
			assert.Equal(t, []string{"scope"}, githubEnterprise.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp github enterprise changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github enterprise
		addGithubEnterprise, err := AdminClient.AddGitHubEnterpriseServerProvider(CTX, &admin.AddGitHubEnterpriseServerProviderRequest{
			Name:                  name,
			ClientId:              "clientId",
			ClientSecret:          "clientSecret",
			AuthorizationEndpoint: "authoizationEndpoint",
			TokenEndpoint:         "tokenEndpoint",
			UserEndpoint:          "userEndpoint",
			Scopes:                []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var githubEnterprise *domain.IDPGithubEnterprise
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			githubEnterprise, err = idpRepo.GetGithubEnterprise(CTX, idpRepo.IDCondition(addGithubEnterprise.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addGithubEnterprise.Id, githubEnterprise.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change github enterprise
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateGitHubEnterpriseServerProvider(CTX, &admin.UpdateGitHubEnterpriseServerProviderRequest{
			Id:                    addGithubEnterprise.Id,
			Name:                  name,
			ClientId:              "new_clientId",
			ClientSecret:          "new_clientSecret",
			AuthorizationEndpoint: "new_authoizationEndpoint",
			TokenEndpoint:         "new_tokenEndpoint",
			UserEndpoint:          "new_userEndpoint",
			Scopes:                []string{"new_scope"},
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

		// check values for azure
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGithubEnterprise, err := idpRepo.GetGithubEnterprise(CTX, idpRepo.IDCondition(addGithubEnterprise.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.github_enterprise.changed
			// idp
			assert.Equal(t, instanceID, githubEnterprise.InstanceID)
			assert.Nil(t, githubEnterprise.OrgID)
			assert.Equal(t, addGithubEnterprise.Id, updateGithubEnterprise.ID)
			assert.Equal(t, name, updateGithubEnterprise.Name)
			assert.Equal(t, domain.IDPTypeGitHubEnterprise.String(), updateGithubEnterprise.Type)
			assert.Equal(t, false, updateGithubEnterprise.AllowLinking)
			assert.Equal(t, false, updateGithubEnterprise.AllowCreation)
			assert.Equal(t, false, updateGithubEnterprise.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), updateGithubEnterprise.AllowAutoLinking)
			assert.WithinRange(t, updateGithubEnterprise.UpdatedAt, beforeCreate, afterCreate)

			// github enterprise
			assert.Equal(t, "new_clientId", updateGithubEnterprise.ClientID)
			assert.NotNil(t, updateGithubEnterprise.ClientSecret)
			assert.Equal(t, "new_authoizationEndpoint", updateGithubEnterprise.AuthorizationEndpoint)
			assert.Equal(t, "new_tokenEndpoint", updateGithubEnterprise.TokenEndpoint)
			assert.Equal(t, "new_userEndpoint", updateGithubEnterprise.UserEndpoint)
			assert.Equal(t, []string{"new_scope"}, updateGithubEnterprise.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp gitlab added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab
		beforeCreate := time.Now()
		addGithub, err := AdminClient.AddGitLabProvider(CTX, &admin.AddGitLabProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
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

		// check values for gitlab
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gitlab, err := idpRepo.GetGitlab(CTX, idpRepo.IDCondition(addGithub.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.gitlab.added
			// idp
			assert.Equal(t, instanceID, gitlab.InstanceID)
			assert.Nil(t, gitlab.OrgID)
			assert.Equal(t, addGithub.Id, gitlab.ID)
			assert.Equal(t, name, gitlab.Name)
			assert.Equal(t, domain.IDPTypeGitLab.String(), gitlab.Type)
			assert.Equal(t, false, gitlab.AllowLinking)
			assert.Equal(t, false, gitlab.AllowCreation)
			assert.Equal(t, false, gitlab.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), gitlab.AllowAutoLinking)
			assert.WithinRange(t, gitlab.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, gitlab.UpdatedAt, beforeCreate, afterCreate)

			// gitlab
			assert.Equal(t, "clientId", gitlab.ClientID)
			assert.NotNil(t, gitlab.ClientSecret)
			assert.Equal(t, []string{"scope"}, gitlab.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp gitlab changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab
		addGitlab, err := AdminClient.AddGitLabProvider(CTX, &admin.AddGitLabProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var githlab *domain.IDPGitlab
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			githlab, err = idpRepo.GetGitlab(CTX, idpRepo.IDCondition(addGitlab.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addGitlab.Id, githlab.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change gitlab
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateGitLabProvider(CTX, &admin.UpdateGitLabProviderRequest{
			Id:           addGitlab.Id,
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
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		// check values for gitlab
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGithlab, err := idpRepo.GetGitlab(CTX, idpRepo.IDCondition(addGitlab.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.gitlab.changed
			// idp
			assert.Equal(t, instanceID, updateGithlab.InstanceID)
			assert.Nil(t, updateGithlab.OrgID)
			assert.Equal(t, addGitlab.Id, updateGithlab.ID)
			assert.Equal(t, name, updateGithlab.Name)
			assert.Equal(t, true, updateGithlab.AllowLinking)
			assert.Equal(t, true, updateGithlab.AllowCreation)
			assert.Equal(t, true, updateGithlab.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateGithlab.AllowAutoLinking)
			assert.WithinRange(t, updateGithlab.UpdatedAt, beforeCreate, afterCreate)

			// gitlab
			assert.Equal(t, "new_clientId", updateGithlab.ClientID)
			assert.NotEqual(t, githlab.ClientSecret, updateGithlab.ClientSecret)
			assert.Equal(t, domain.IDPTypeGitLab.String(), updateGithlab.Type)
			assert.Equal(t, []string{"new_scope"}, updateGithlab.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp gitlab self hosted added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab self hosted
		beforeCreate := time.Now()
		addGitlabSelfHosted, err := AdminClient.AddGitLabSelfHostedProvider(CTX, &admin.AddGitLabSelfHostedProviderRequest{
			Name:         name,
			Issuer:       "issuer",
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
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

		// check values for gitlab self hosted
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gitlabSelfHosted, err := idpRepo.GetGitlabSelfHosting(CTX, idpRepo.IDCondition(addGitlabSelfHosted.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.gitlab_self_hosted.added
			// idp
			assert.Equal(t, instanceID, gitlabSelfHosted.InstanceID)
			assert.Nil(t, gitlabSelfHosted.OrgID)
			assert.Equal(t, addGitlabSelfHosted.Id, gitlabSelfHosted.ID)
			assert.Equal(t, name, gitlabSelfHosted.Name)
			assert.Equal(t, domain.IDPTypeGitLabSelfHosted.String(), gitlabSelfHosted.Type)
			assert.Equal(t, false, gitlabSelfHosted.AllowLinking)
			assert.Equal(t, false, gitlabSelfHosted.AllowCreation)
			assert.Equal(t, false, gitlabSelfHosted.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), gitlabSelfHosted.AllowAutoLinking)
			assert.WithinRange(t, gitlabSelfHosted.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, gitlabSelfHosted.UpdatedAt, beforeCreate, afterCreate)

			// gitlab self hosted
			assert.Equal(t, "clientId", gitlabSelfHosted.ClientID)
			assert.Equal(t, "issuer", gitlabSelfHosted.Issuer)
			assert.NotNil(t, gitlabSelfHosted.ClientSecret)
			assert.Equal(t, []string{"scope"}, gitlabSelfHosted.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp gitlab self hosted changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab self hosted
		addGitlabSelfHosted, err := AdminClient.AddGitLabSelfHostedProvider(CTX, &admin.AddGitLabSelfHostedProviderRequest{
			Name:         name,
			Issuer:       "issuer",
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var githlabSelfHosted *domain.IDPGitlabSelfHosting
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			githlabSelfHosted, err = idpRepo.GetGitlabSelfHosting(CTX, idpRepo.IDCondition(addGitlabSelfHosted.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addGitlabSelfHosted.Id, githlabSelfHosted.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change gitlab self hosted
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateGitLabSelfHostedProvider(CTX, &admin.UpdateGitLabSelfHostedProviderRequest{
			Id:           addGitlabSelfHosted.Id,
			Name:         name,
			ClientId:     "new_clientId",
			Issuer:       "new_issuer",
			ClientSecret: "new_clientSecret",
			Scopes:       []string{"new_scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		// check values for gitlab self hosted
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGithlabSelfHosted, err := idpRepo.GetGitlabSelfHosting(CTX, idpRepo.IDCondition(addGitlabSelfHosted.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.gitlab_self_hosted.changed
			// idp
			assert.Equal(t, instanceID, updateGithlabSelfHosted.InstanceID)
			assert.Nil(t, updateGithlabSelfHosted.OrgID)
			assert.Equal(t, addGitlabSelfHosted.Id, updateGithlabSelfHosted.ID)
			assert.Equal(t, name, updateGithlabSelfHosted.Name)
			assert.Equal(t, domain.IDPTypeGitLabSelfHosted.String(), updateGithlabSelfHosted.Type)
			assert.Equal(t, true, updateGithlabSelfHosted.AllowLinking)
			assert.Equal(t, true, updateGithlabSelfHosted.AllowCreation)
			assert.Equal(t, true, updateGithlabSelfHosted.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateGithlabSelfHosted.AllowAutoLinking)
			assert.WithinRange(t, updateGithlabSelfHosted.UpdatedAt, beforeCreate, afterCreate)

			// gitlab self hosted
			assert.Equal(t, "new_clientId", updateGithlabSelfHosted.ClientID)
			assert.Equal(t, "new_issuer", updateGithlabSelfHosted.Issuer)
			assert.NotEqual(t, githlabSelfHosted.ClientSecret, updateGithlabSelfHosted.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateGithlabSelfHosted.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp google added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add google
		beforeCreate := time.Now()
		addGoogle, err := AdminClient.AddGoogleProvider(CTX, &admin.AddGoogleProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
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

		// check values for google
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			google, err := idpRepo.GetGoogle(CTX, idpRepo.IDCondition(addGoogle.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.google.added
			// idp
			assert.Equal(t, instanceID, google.InstanceID)
			assert.Nil(t, google.OrgID)
			assert.Equal(t, addGoogle.Id, google.ID)
			assert.Equal(t, name, google.Name)
			assert.Equal(t, domain.IDPTypeGoogle.String(), google.Type)
			assert.Equal(t, false, google.AllowLinking)
			assert.Equal(t, false, google.AllowCreation)
			assert.Equal(t, false, google.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), google.AllowAutoLinking)
			assert.WithinRange(t, google.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, google.UpdatedAt, beforeCreate, afterCreate)

			// google
			assert.Equal(t, "clientId", google.ClientID)
			assert.NotNil(t, google.ClientSecret)
			assert.Equal(t, []string{"scope"}, google.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance idp google changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add google
		addGoogle, err := AdminClient.AddGoogleProvider(CTX, &admin.AddGoogleProviderRequest{
			Name:         name,
			ClientId:     "clientId",
			ClientSecret: "clientSecret",
			Scopes:       []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var google *domain.IDPGoogle
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			google, err = idpRepo.GetGoogle(CTX, idpRepo.IDCondition(addGoogle.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addGoogle.Id, google.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change google
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateGoogleProvider(CTX, &admin.UpdateGoogleProviderRequest{
			Id:           addGoogle.Id,
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
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		// check values for google
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGoogle, err := idpRepo.GetGoogle(CTX, idpRepo.IDCondition(addGoogle.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.google.changed
			// idp
			assert.Equal(t, instanceID, updateGoogle.InstanceID)
			assert.Nil(t, updateGoogle.OrgID)
			assert.Equal(t, addGoogle.Id, updateGoogle.ID)
			assert.Equal(t, name, updateGoogle.Name)
			assert.Equal(t, domain.IDPTypeGoogle.String(), updateGoogle.Type)
			assert.Equal(t, true, updateGoogle.AllowLinking)
			assert.Equal(t, true, updateGoogle.AllowCreation)
			assert.Equal(t, true, updateGoogle.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateGoogle.AllowAutoLinking)
			assert.WithinRange(t, updateGoogle.UpdatedAt, beforeCreate, afterCreate)

			// google
			assert.Equal(t, "new_clientId", updateGoogle.ClientID)
			assert.NotEqual(t, google.ClientSecret, updateGoogle.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateGoogle.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test instance ldap added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add ldap
		beforeCreate := time.Now()
		addLdap, err := AdminClient.AddLDAPProvider(CTX, &admin.AddLDAPProviderRequest{
			Name:              name,
			Servers:           []string{"servers"},
			StartTls:          true,
			BaseDn:            "baseDN",
			BindDn:            "bindND",
			BindPassword:      "bindPassword",
			UserBase:          "userBase",
			UserObjectClasses: []string{"userOhjectClasses"},
			UserFilters:       []string{"userFilters"},
			Timeout:           durationpb.New(time.Minute),
			Attributes: &idp_grpc.LDAPAttributes{
				IdAttribute:                "idAttribute",
				FirstNameAttribute:         "firstNameAttribute",
				LastNameAttribute:          "lastNameAttribute",
				DisplayNameAttribute:       "displayNameAttribute",
				NickNameAttribute:          "nickNameAttribute",
				PreferredUsernameAttribute: "preferredUsernameAttribute",
				EmailAttribute:             "emailAttribute",
				EmailVerifiedAttribute:     "emailVerifiedAttribute",
				PhoneAttribute:             "phoneAttribute",
				PhoneVerifiedAttribute:     "phoneVerifiedAttribute",
				PreferredLanguageAttribute: "preferredLanguageAttribute",
				AvatarUrlAttribute:         "avatarUrlAttribute",
				ProfileAttribute:           "profileAttribute",
			},
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			ldap, err := idpRepo.GetLDAP(CTX, idpRepo.IDCondition(addLdap.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.ldap.v2.added
			// idp
			assert.Equal(t, instanceID, ldap.InstanceID)
			assert.Nil(t, ldap.OrgID)
			assert.Equal(t, addLdap.Id, ldap.ID)
			assert.Equal(t, name, ldap.Name)
			assert.Equal(t, domain.IDPTypeLDAP.String(), ldap.Type)
			assert.Equal(t, false, ldap.AllowLinking)
			assert.Equal(t, false, ldap.AllowCreation)
			assert.Equal(t, false, ldap.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionEmail.String(), ldap.AllowAutoLinking)
			assert.WithinRange(t, ldap.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, ldap.UpdatedAt, beforeCreate, afterCreate)

			// ldap
			assert.Equal(t, []string{"servers"}, ldap.Servers)
			assert.Equal(t, true, ldap.StartTLS)
			assert.Equal(t, "baseDN", ldap.BaseDN)
			assert.Equal(t, "bindND", ldap.BindDN)
			assert.NotNil(t, ldap.BindPassword)
			assert.Equal(t, "userBase", ldap.UserBase)
			assert.Equal(t, []string{"userOhjectClasses"}, ldap.UserObjectClasses)
			assert.Equal(t, []string{"userFilters"}, ldap.UserFilters)
			assert.Equal(t, time.Minute, ldap.Timeout)
			assert.Equal(t, "idAttribute", ldap.IDAttribute)
			assert.Equal(t, "firstNameAttribute", ldap.FirstNameAttribute)
			assert.Equal(t, "lastNameAttribute", ldap.LastNameAttribute)
			assert.Equal(t, "displayNameAttribute", ldap.DisplayNameAttribute)
			assert.Equal(t, "nickNameAttribute", ldap.NickNameAttribute)
			assert.Equal(t, "preferredUsernameAttribute", ldap.PreferredUsernameAttribute)
			assert.Equal(t, "emailAttribute", ldap.EmailAttribute)
			assert.Equal(t, "emailVerifiedAttribute", ldap.EmailVerifiedAttribute)
			assert.Equal(t, "phoneAttribute", ldap.PhoneAttribute)
			assert.Equal(t, "phoneVerifiedAttribute", ldap.PhoneVerifiedAttribute)
			assert.Equal(t, "preferredLanguageAttribute", ldap.PreferredLanguageAttribute)
			assert.Equal(t, "avatarUrlAttribute", ldap.AvatarURLAttribute)
			assert.Equal(t, "profileAttribute", ldap.ProfileAttribute)
		}, retryDuration, tick)
	})

	t.Run("test instance ldap changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add ldap
		addLdap, err := AdminClient.AddLDAPProvider(CTX, &admin.AddLDAPProviderRequest{
			Name:              name,
			Servers:           []string{"servers"},
			StartTls:          true,
			BaseDn:            "baseDN",
			BindDn:            "bindND",
			BindPassword:      "bindPassword",
			UserBase:          "userBase",
			UserObjectClasses: []string{"userOhjectClasses"},
			UserFilters:       []string{"userFilters"},
			Timeout:           durationpb.New(time.Minute),
			Attributes: &idp_grpc.LDAPAttributes{
				IdAttribute:                "idAttribute",
				FirstNameAttribute:         "firstNameAttribute",
				LastNameAttribute:          "lastNameAttribute",
				DisplayNameAttribute:       "displayNameAttribute",
				NickNameAttribute:          "nickNameAttribute",
				PreferredUsernameAttribute: "preferredUsernameAttribute",
				EmailAttribute:             "emailAttribute",
				EmailVerifiedAttribute:     "emailVerifiedAttribute",
				PhoneAttribute:             "phoneAttribute",
				PhoneVerifiedAttribute:     "phoneVerifiedAttribute",
				PreferredLanguageAttribute: "preferredLanguageAttribute",
				AvatarUrlAttribute:         "avatarUrlAttribute",
				ProfileAttribute:           "profileAttribute",
			},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var ldap *domain.IDPLDAP
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			ldap, err = idpRepo.GetLDAP(CTX, idpRepo.IDCondition(addLdap.Id), instanceID, nil)
			require.NoError(t, err)
			assert.Equal(t, addLdap.Id, ldap.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change ldap
		beforeCreate := time.Now()
		_, err = AdminClient.UpdateLDAPProvider(CTX, &admin.UpdateLDAPProviderRequest{
			Id:                addLdap.Id,
			Name:              name,
			Servers:           []string{"new_servers"},
			StartTls:          false,
			BaseDn:            "new_baseDN",
			BindDn:            "new_bindND",
			BindPassword:      "new_bindPassword",
			UserBase:          "new_userBase",
			UserObjectClasses: []string{"new_userOhjectClasses"},
			UserFilters:       []string{"new_userFilters"},
			Timeout:           durationpb.New(time.Second),
			Attributes: &idp_grpc.LDAPAttributes{
				IdAttribute:                "new_idAttribute",
				FirstNameAttribute:         "new_firstNameAttribute",
				LastNameAttribute:          "new_lastNameAttribute",
				DisplayNameAttribute:       "new_displayNameAttribute",
				NickNameAttribute:          "new_nickNameAttribute",
				PreferredUsernameAttribute: "new_preferredUsernameAttribute",
				EmailAttribute:             "new_emailAttribute",
				EmailVerifiedAttribute:     "new_emailVerifiedAttribute",
				PhoneAttribute:             "new_phoneAttribute",
				PhoneVerifiedAttribute:     "new_phoneVerifiedAttribute",
				PreferredLanguageAttribute: "new_preferredLanguageAttribute",
				AvatarUrlAttribute:         "new_avatarUrlAttribute",
				ProfileAttribute:           "new_profileAttribute",
			},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		afterCreate := time.Now()
		require.NoError(t, err)

		// check values for ldap
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateLdap, err := idpRepo.GetLDAP(CTX, idpRepo.IDCondition(addLdap.Id), instanceID, nil)
			require.NoError(t, err)

			// event instance.idp.ldap.v2.changed
			// idp
			assert.Equal(t, instanceID, updateLdap.InstanceID)
			assert.Nil(t, updateLdap.OrgID)
			assert.Equal(t, addLdap.Id, updateLdap.ID)
			assert.Equal(t, name, updateLdap.Name)
			assert.Equal(t, domain.IDPTypeLDAP.String(), updateLdap.Type)
			assert.Equal(t, true, updateLdap.AllowLinking)
			assert.Equal(t, true, updateLdap.AllowCreation)
			assert.Equal(t, true, updateLdap.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingOptionUserName.String(), updateLdap.AllowAutoLinking)
			assert.WithinRange(t, updateLdap.UpdatedAt, beforeCreate, afterCreate)

			// ldap
			assert.Equal(t, []string{"new_servers"}, updateLdap.Servers)
			assert.Equal(t, false, updateLdap.StartTLS)
			assert.Equal(t, "new_baseDN", updateLdap.BaseDN)
			assert.Equal(t, "new_bindND", updateLdap.BindDN)
			assert.NotEqual(t, ldap.BindPassword, updateLdap.BindPassword)
			assert.Equal(t, "new_userBase", updateLdap.UserBase)
			assert.Equal(t, []string{"new_userOhjectClasses"}, updateLdap.UserObjectClasses)
			assert.Equal(t, []string{"new_userFilters"}, updateLdap.UserFilters)
			assert.Equal(t, time.Second, updateLdap.Timeout)
			assert.Equal(t, "new_idAttribute", updateLdap.IDAttribute)
			assert.Equal(t, "new_firstNameAttribute", updateLdap.FirstNameAttribute)
			assert.Equal(t, "new_lastNameAttribute", updateLdap.LastNameAttribute)
			assert.Equal(t, "new_displayNameAttribute", updateLdap.DisplayNameAttribute)
			assert.Equal(t, "new_nickNameAttribute", updateLdap.NickNameAttribute)
			assert.Equal(t, "new_preferredUsernameAttribute", updateLdap.PreferredUsernameAttribute)
			assert.Equal(t, "new_emailAttribute", updateLdap.EmailAttribute)
			assert.Equal(t, "new_emailVerifiedAttribute", updateLdap.EmailVerifiedAttribute)
			assert.Equal(t, "new_phoneAttribute", updateLdap.PhoneAttribute)
			assert.Equal(t, "new_phoneVerifiedAttribute", updateLdap.PhoneVerifiedAttribute)
			assert.Equal(t, "new_preferredLanguageAttribute", updateLdap.PreferredLanguageAttribute)
			assert.Equal(t, "new_avatarUrlAttribute", updateLdap.AvatarURLAttribute)
			assert.Equal(t, "new_profileAttribute", updateLdap.ProfileAttribute)
		}, retryDuration, tick)
	})
}
