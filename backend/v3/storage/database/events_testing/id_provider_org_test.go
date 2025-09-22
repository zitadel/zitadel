//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	durationpb "google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	zitadel_internal_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	idp_grpc "github.com/zitadel/zitadel/pkg/grpc/idp"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

func TestServer_TestIDProviderOrgReduces(t *testing.T) {
	instanceID := Instance.ID()

	orgID := Instance.DefaultOrg.Id

	t.Run("test iam idp add reduces", func(t *testing.T) {
		name := gofakeit.Name()

		before := time.Now()
		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(IAMCTX,
				idpRepo.NameCondition(name),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.config.added
			assert.Equal(t, instanceID, idp.InstanceID)
			assert.Equal(t, orgID, *idp.OrgID)
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateActive, idp.State)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, true, idp.AutoRegister)
			assert.Equal(t, true, idp.AllowCreation)
			assert.Equal(t, false, idp.AllowAutoUpdate)
			assert.Equal(t, true, idp.AllowLinking)
			assert.Nil(t, idp.AllowAutoLinkingField)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), *idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, before, after)
			assert.WithinRange(t, idp.CreatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test iam idp update reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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

		before := time.Now()
		_, err = MgmtClient.UpdateOrgIDP(IAMCTX, &management.UpdateOrgIDPRequest{
			IdpId:        addOIDC.IdpId,
			Name:         name,
			StylingType:  idp_grpc.IDPStylingType_STYLING_TYPE_UNSPECIFIED,
			AutoRegister: false,
		})
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(IAMCTX,
				idpRepo.NameCondition(name),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.config.changed
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, false, idp.AutoRegister)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_UNSPECIFIED), *idp.StylingType)
			assert.WithinRange(t, idp.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test iam idp deactivate reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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
		before := time.Now()
		_, err = MgmtClient.DeactivateOrgIDP(IAMCTX, &management.DeactivateOrgIDPRequest{
			IdpId: addOIDC.IdpId,
		})
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.config.deactivated
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateInactive, idp.State)
			assert.WithinRange(t, idp.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test iam idp reactivate reduces", func(t *testing.T) {
		name := gofakeit.Name()

		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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
		_, err = MgmtClient.DeactivateOrgIDP(IAMCTX, &management.DeactivateOrgIDPRequest{
			IdpId: addOIDC.IdpId,
		})
		require.NoError(t, err)
		// wait for idp to be deactivated
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateInactive, idp.State)
		}, retryDuration, tick)

		// reactivate idp
		before := time.Now()
		_, err = MgmtClient.ReactivateOrgIDP(IAMCTX, &management.ReactivateOrgIDPRequest{
			IdpId: addOIDC.IdpId,
		})
		after := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.config.reactivated
			assert.Equal(t, addOIDC.IdpId, idp.ID)
			assert.Equal(t, domain.IDPStateActive, idp.State)
			assert.WithinRange(t, idp.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test iam idp remove reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add idp
		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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
		_, err = MgmtClient.RemoveOrgIDP(IAMCTX, &management.RemoveOrgIDPRequest{
			IdpId: addOIDC.IdpId,
		})
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := idpRepo.Get(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)

			// event org.idp.config.remove
			require.ErrorIs(t, &database.NoRowFoundError{}, err)
		}, retryDuration, tick)
	})

	t.Run("test iam idp oidc added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oidc
		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err := idpRepo.GetOIDC(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.oidc.config.added
			// idp
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Equal(t, orgID, *oidc.OrgID)
			assert.Equal(t, name, oidc.Name)
			assert.Equal(t, addOIDC.IdpId, oidc.ID)
			assert.Equal(t, domain.IDPTypeOIDC, domain.IDPType(*oidc.Type))

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
		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(IAMCTX, idpRepo.IDCondition(addOIDC.IdpId), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addOIDC.IdpId, oidc.ID)
		}, retryDuration, tick)

		before := time.Now()
		_, err = MgmtClient.UpdateOrgIDPOIDCConfig(IAMCTX, &management.UpdateOrgIDPOIDCConfigRequest{
			IdpId:              addOIDC.IdpId,
			ClientId:           "new_clientID",
			ClientSecret:       "new_clientSecret",
			Issuer:             "new_issuer",
			Scopes:             []string{"new_scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
		})
		after := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateOIDC, err := idpRepo.GetOIDC(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.oidc.config.changed
			// idp
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Equal(t, orgID, *oidc.OrgID)
			assert.Equal(t, name, oidc.Name)
			assert.Equal(t, addOIDC.IdpId, updateOIDC.ID)
			assert.Equal(t, domain.IDPTypeOIDC, domain.IDPType(*updateOIDC.Type))
			assert.WithinRange(t, updateOIDC.UpdatedAt, before, after)

			// oidc
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Equal(t, orgID, *oidc.OrgID)
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
		addJWT, err := MgmtClient.AddOrgJWTIDP(IAMCTX, &management.AddOrgJWTIDPRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			jwt, err := idpRepo.GetJWT(IAMCTX,
				idpRepo.IDCondition(addJWT.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.jwt.config.added
			// idp
			assert.Equal(t, instanceID, jwt.InstanceID)
			assert.Equal(t, orgID, *jwt.OrgID)
			assert.Equal(t, name, jwt.Name)
			assert.Equal(t, addJWT.IdpId, jwt.ID)
			assert.Equal(t, domain.IDPTypeJWT, domain.IDPType(*jwt.Type))
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
		addJWT, err := MgmtClient.AddOrgJWTIDP(IAMCTX, &management.AddOrgJWTIDPRequest{
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

		before := time.Now()
		_, err = MgmtClient.UpdateOrgIDPJWTConfig(IAMCTX, &management.UpdateOrgIDPJWTConfigRequest{
			IdpId:        addJWT.IdpId,
			JwtEndpoint:  "new_jwtEndpoint",
			Issuer:       "new_issuer",
			KeysEndpoint: "new_keyEndpoint",
			HeaderName:   "new_headerName",
		})
		after := time.Now()
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateJWT, err := idpRepo.GetJWT(IAMCTX,
				idpRepo.IDCondition(addJWT.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.jwt.config.changed
			// idp
			assert.Equal(t, addJWT.IdpId, updateJWT.ID)
			assert.Equal(t, orgID, *updateJWT.OrgID)
			assert.Equal(t, domain.IDPTypeJWT, domain.IDPType(*updateJWT.Type))
			assert.WithinRange(t, updateJWT.UpdatedAt, before, after)

			// jwt
			assert.Equal(t, "new_jwtEndpoint", updateJWT.JWTEndpoint)
			assert.Equal(t, "new_issuer", updateJWT.Issuer)
			assert.Equal(t, "new_keyEndpoint", updateJWT.KeysEndpoint)
			assert.Equal(t, "new_headerName", updateJWT.HeaderName)
		}, retryDuration, tick)
	})

	t.Run("test org idp oauth added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oauth
		before := time.Now()
		addOAuth, err := MgmtClient.AddGenericOAuthProvider(IAMCTX, &management.AddGenericOAuthProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for oauth
		var oauth *domain.IDPOAuth
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oauth, err = idpRepo.GetOAuth(IAMCTX, idpRepo.IDCondition(addOAuth.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.oauth.added
			// idp
			assert.Equal(t, instanceID, oauth.InstanceID)
			assert.Equal(t, orgID, *oauth.OrgID)
			assert.Equal(t, addOAuth.Id, oauth.ID)
			assert.Equal(t, name, oauth.Name)
			assert.Equal(t, domain.IDPTypeOAuth, domain.IDPType(*oauth.Type))
			assert.Equal(t, false, oauth.AllowLinking)
			assert.Equal(t, false, oauth.AllowCreation)
			assert.Equal(t, false, oauth.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*oauth.AllowAutoLinkingField))
			assert.WithinRange(t, oauth.CreatedAt, before, after)
			assert.WithinRange(t, oauth.UpdatedAt, before, after)

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

	t.Run("test org idp oauth changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oauth
		addOAuth, err := MgmtClient.AddGenericOAuthProvider(IAMCTX, &management.AddGenericOAuthProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oauth, err = idpRepo.GetOAuth(IAMCTX, idpRepo.IDCondition(addOAuth.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addOAuth.Id, oauth.ID)
		}, retryDuration, tick)

		name = "new_" + name
		before := time.Now()
		_, err = MgmtClient.UpdateGenericOAuthProvider(IAMCTX, &management.UpdateGenericOAuthProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateOauth, err := idpRepo.GetOAuth(IAMCTX,
				idpRepo.IDCondition(addOAuth.Id),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.oauth.changed
			// idp
			assert.Equal(t, instanceID, updateOauth.InstanceID)
			assert.Equal(t, orgID, *updateOauth.OrgID)
			assert.Equal(t, addOAuth.Id, updateOauth.ID)
			assert.Equal(t, name, updateOauth.Name)
			assert.Equal(t, domain.IDPTypeOAuth, domain.IDPType(*updateOauth.Type))
			assert.Equal(t, true, updateOauth.AllowLinking)
			assert.Equal(t, true, updateOauth.AllowCreation)
			assert.Equal(t, true, updateOauth.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateOauth.AllowAutoLinkingField))
			assert.Equal(t, true, updateOauth.UsePKCE)
			assert.WithinRange(t, updateOauth.UpdatedAt, before, after)

			// oauth
			assert.Equal(t, "new_clientId", updateOauth.ClientID)
			assert.NotEqual(t, oauth.ClientSecret, updateOauth.ClientSecret)
			assert.Equal(t, "new_authoizationEndpoint", updateOauth.AuthorizationEndpoint)
			assert.Equal(t, "new_tokenEndpoint", updateOauth.TokenEndpoint)
			assert.Equal(t, "new_userEndpoint", updateOauth.UserEndpoint)
			assert.Equal(t, []string{"new_scope"}, updateOauth.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp oidc added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add oidc
		before := time.Now()
		addOIDC, err := MgmtClient.AddGenericOIDCProvider(IAMCTX, &management.AddGenericOIDCProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for oidc
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err := idpRepo.GetOIDC(IAMCTX, idpRepo.IDCondition(addOIDC.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.oidc added
			// idp
			assert.Equal(t, instanceID, oidc.InstanceID)
			assert.Equal(t, orgID, *oidc.OrgID)
			assert.Equal(t, addOIDC.Id, oidc.ID)
			assert.Equal(t, name, oidc.Name)
			assert.Equal(t, domain.IDPTypeOIDC, domain.IDPType(*oidc.Type))
			assert.Equal(t, false, oidc.AllowLinking)
			assert.Equal(t, false, oidc.AllowCreation)
			assert.Equal(t, false, oidc.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*oidc.AllowAutoLinkingField))
			assert.WithinRange(t, oidc.CreatedAt, before, after)
			assert.WithinRange(t, oidc.UpdatedAt, before, after)

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

		addOIDC, err := MgmtClient.AddGenericOIDCProvider(IAMCTX, &management.AddGenericOIDCProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(IAMCTX, idpRepo.IDCondition(addOIDC.Id), instanceID, &orgID)
			require.NoError(t, err)
		}, retryDuration, tick)

		name = "new_" + name
		before := time.Now()
		_, err = MgmtClient.UpdateGenericOIDCProvider(IAMCTX, &management.UpdateGenericOIDCProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateOIDC, err := idpRepo.GetOIDC(IAMCTX,
				idpRepo.IDCondition(addOIDC.Id),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)

			// event org.idp.oidc.changed
			// idp
			assert.Equal(t, instanceID, updateOIDC.InstanceID)
			assert.Equal(t, orgID, *updateOIDC.OrgID)
			assert.Equal(t, addOIDC.Id, updateOIDC.ID)
			assert.Equal(t, name, updateOIDC.Name)
			assert.Equal(t, domain.IDPTypeOIDC, domain.IDPType(*updateOIDC.Type))
			assert.Equal(t, true, updateOIDC.AllowLinking)
			assert.Equal(t, true, updateOIDC.AllowCreation)
			assert.Equal(t, true, updateOIDC.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateOIDC.AllowAutoLinkingField))
			assert.WithinRange(t, updateOIDC.UpdatedAt, before, after)

			// oidc
			assert.Equal(t, "new_clientId", updateOIDC.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, updateOIDC.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateOIDC.Scopes)
			assert.Equal(t, true, updateOIDC.IsIDTokenMapping)
			assert.Equal(t, true, updateOIDC.UsePKCE)
		}, retryDuration, tick)
	})

	t.Run("test org idp oidc migrated azure migration reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// create OIDC
		addOIDC, err := MgmtClient.AddGenericOIDCProvider(IAMCTX, &management.AddGenericOIDCProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(IAMCTX, idpRepo.IDCondition(addOIDC.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, domain.IDPTypeOIDC, domain.IDPType(*oidc.Type))
		}, retryDuration, tick)

		before := time.Now()
		_, err = MgmtClient.MigrateGenericOIDCProvider(IAMCTX, &management.MigrateGenericOIDCProviderRequest{
			Id: addOIDC.Id,
			Template: &management.MigrateGenericOIDCProviderRequest_Azure{
				Azure: &management.AddAzureADProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			azure, err := idpRepo.GetAzureAD(IAMCTX, idpRepo.IDCondition(addOIDC.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.oidc.migrated.azure
			// idp
			assert.Equal(t, instanceID, azure.InstanceID)
			assert.Equal(t, orgID, *azure.OrgID)
			assert.Equal(t, addOIDC.Id, azure.ID)
			assert.Equal(t, name, azure.Name)
			// type = azure
			assert.Equal(t, domain.IDPTypeAzure, domain.IDPType(*azure.Type))
			assert.Equal(t, true, azure.AllowLinking)
			assert.Equal(t, true, azure.AllowCreation)
			assert.Equal(t, true, azure.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*azure.AllowAutoLinkingField))
			assert.WithinRange(t, azure.UpdatedAt, before, after)

			// oidc
			assert.Equal(t, "new_clientId", azure.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, azure.ClientSecret)
			assert.Equal(t, domain.AzureTenantTypeOrganizations, azure.Tenant)
			assert.Equal(t, true, azure.IsEmailVerified)
			assert.Equal(t, []string{"new_scope"}, azure.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp oidc migrated google migration reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// create OIDC
		addOIDC, err := MgmtClient.AddGenericOIDCProvider(IAMCTX, &management.AddGenericOIDCProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			oidc, err = idpRepo.GetOIDC(IAMCTX, idpRepo.IDCondition(addOIDC.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, domain.IDPTypeOIDC, domain.IDPType(*oidc.Type))
		}, retryDuration, tick)

		before := time.Now()
		_, err = MgmtClient.MigrateGenericOIDCProvider(IAMCTX, &management.MigrateGenericOIDCProviderRequest{
			Id: addOIDC.Id,
			Template: &management.MigrateGenericOIDCProviderRequest_Google{
				Google: &management.AddGoogleProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			google, err := idpRepo.GetGoogle(IAMCTX, idpRepo.IDCondition(addOIDC.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.oidc.migrated.google
			// idp
			assert.Equal(t, instanceID, google.InstanceID)
			assert.Equal(t, orgID, *google.OrgID)
			assert.Equal(t, addOIDC.Id, google.ID)
			assert.Equal(t, name, google.Name)
			// type = google
			assert.Equal(t, domain.IDPTypeGoogle, domain.IDPType(*google.Type))
			assert.Equal(t, true, google.AllowLinking)
			assert.Equal(t, true, google.AllowCreation)
			assert.Equal(t, true, google.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*google.AllowAutoLinkingField))
			assert.WithinRange(t, google.UpdatedAt, before, after)

			// oidc
			assert.Equal(t, "new_clientId", google.ClientID)
			assert.NotEqual(t, oidc.ClientSecret, google.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, google.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp jwt added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add jwt
		before := time.Now()
		addJWT, err := MgmtClient.AddJWTProvider(IAMCTX, &management.AddJWTProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for jwt
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			jwt, err := idpRepo.GetJWT(IAMCTX, idpRepo.IDCondition(addJWT.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.jwt.added
			// idp
			assert.Equal(t, instanceID, jwt.InstanceID)
			assert.Equal(t, orgID, *jwt.OrgID)
			assert.Equal(t, addJWT.Id, jwt.ID)
			assert.Equal(t, name, jwt.Name)
			assert.Equal(t, domain.IDPTypeJWT, domain.IDPType(*jwt.Type))
			assert.Equal(t, false, jwt.AllowLinking)
			assert.Equal(t, false, jwt.AllowCreation)
			assert.Equal(t, false, jwt.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*jwt.AllowAutoLinkingField))
			assert.WithinRange(t, jwt.CreatedAt, before, after)
			assert.WithinRange(t, jwt.UpdatedAt, before, after)

			// jwt
			assert.Equal(t, "jwtEndpoint", jwt.JWTEndpoint)
			assert.Equal(t, "issuer", jwt.Issuer)
			assert.Equal(t, "keyEndpoint", jwt.KeysEndpoint)
			assert.Equal(t, "headerName", jwt.HeaderName)
		}, retryDuration, tick)
	})

	t.Run("test org idp jwt changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add jwt
		addJWT, err := MgmtClient.AddJWTProvider(IAMCTX, &management.AddJWTProviderRequest{
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
		before := time.Now()
		_, err = MgmtClient.UpdateJWTProvider(IAMCTX, &management.UpdateJWTProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for jwt
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateJwt, err := idpRepo.GetJWT(IAMCTX, idpRepo.IDCondition(addJWT.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.jwt.added
			// idp
			assert.Equal(t, instanceID, updateJwt.InstanceID)
			assert.Equal(t, orgID, *updateJwt.OrgID)
			assert.Equal(t, addJWT.Id, updateJwt.ID)
			assert.Equal(t, name, updateJwt.Name)
			assert.Equal(t, domain.IDPTypeJWT, domain.IDPType(*updateJwt.Type))
			assert.Equal(t, true, updateJwt.AllowLinking)
			assert.Equal(t, true, updateJwt.AllowCreation)
			assert.Equal(t, true, updateJwt.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateJwt.AllowAutoLinkingField))
			assert.WithinRange(t, updateJwt.UpdatedAt, before, after)

			// jwt
			assert.Equal(t, "new_jwtEndpoint", updateJwt.JWTEndpoint)
			assert.Equal(t, "new_issuer", updateJwt.Issuer)
			assert.Equal(t, "new_keyEndpoint", updateJwt.KeysEndpoint)
			assert.Equal(t, "new_headerName", updateJwt.HeaderName)
		}, retryDuration, tick)
	})

	t.Run("test org idp azure added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add azure
		before := time.Now()
		addAzure, err := MgmtClient.AddAzureADProvider(IAMCTX, &management.AddAzureADProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for azure
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			azure, err := idpRepo.GetAzureAD(IAMCTX, idpRepo.IDCondition(addAzure.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.azure.added
			// idp
			assert.Equal(t, instanceID, azure.InstanceID)
			assert.Equal(t, orgID, *azure.OrgID)
			assert.Equal(t, addAzure.Id, azure.ID)
			assert.Equal(t, name, azure.Name)
			assert.Equal(t, domain.IDPTypeAzure, domain.IDPType(*azure.Type))
			assert.Equal(t, true, azure.AllowLinking)
			assert.Equal(t, true, azure.AllowCreation)
			assert.Equal(t, true, azure.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*azure.AllowAutoLinkingField))
			assert.WithinRange(t, azure.UpdatedAt, before, after)

			// azure
			assert.Equal(t, "clientId", azure.ClientID)
			assert.NotNil(t, azure.ClientSecret)
			assert.Equal(t, domain.AzureTenantTypeOrganizations, azure.Tenant)
			assert.Equal(t, true, azure.IsEmailVerified)
			assert.Equal(t, []string{"scope"}, azure.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp azure changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add azure
		addAzure, err := MgmtClient.AddAzureADProvider(IAMCTX, &management.AddAzureADProviderRequest{
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

		var azure *domain.IDPAzureAD
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			azure, err = idpRepo.GetAzureAD(IAMCTX, idpRepo.IDCondition(addAzure.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addAzure.Id, azure.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change azure
		before := time.Now()
		_, err = MgmtClient.UpdateAzureADProvider(IAMCTX, &management.UpdateAzureADProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		// check values for azure
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateAzure, err := idpRepo.GetAzureAD(IAMCTX, idpRepo.IDCondition(addAzure.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.azure.changed
			// idp
			assert.Equal(t, instanceID, updateAzure.InstanceID)
			assert.Equal(t, orgID, *updateAzure.OrgID)
			assert.Equal(t, addAzure.Id, updateAzure.ID)
			assert.Equal(t, name, updateAzure.Name)
			assert.Equal(t, domain.IDPTypeAzure, domain.IDPType(*updateAzure.Type))
			assert.Equal(t, true, updateAzure.AllowLinking)
			assert.Equal(t, true, updateAzure.AllowCreation)
			assert.Equal(t, true, updateAzure.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*updateAzure.AllowAutoLinkingField))
			assert.WithinRange(t, updateAzure.UpdatedAt, before, after)

			// azure
			assert.Equal(t, "new_clientId", updateAzure.ClientID)
			assert.NotEqual(t, azure.ClientSecret, updateAzure.ClientSecret)
			assert.Equal(t, domain.AzureTenantTypeConsumers, updateAzure.Tenant)
			assert.Equal(t, true, updateAzure.IsEmailVerified)
			assert.Equal(t, []string{"new_scope"}, updateAzure.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp github added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github
		before := time.Now()
		addGithub, err := MgmtClient.AddGitHubProvider(IAMCTX, &management.AddGitHubProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for github
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			github, err := idpRepo.GetGithub(IAMCTX, idpRepo.IDCondition(addGithub.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.github.added
			// idp
			assert.Equal(t, instanceID, github.InstanceID)
			assert.Equal(t, orgID, *github.OrgID)
			assert.Equal(t, addGithub.Id, github.ID)
			assert.Equal(t, name, github.Name)
			assert.Equal(t, domain.IDPTypeGitHub, domain.IDPType(*github.Type))
			assert.Equal(t, false, github.AllowLinking)
			assert.Equal(t, false, github.AllowCreation)
			assert.Equal(t, false, github.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*github.AllowAutoLinkingField))
			assert.WithinRange(t, github.UpdatedAt, before, after)

			assert.Equal(t, "clientId", github.ClientID)
			assert.NotNil(t, github.ClientSecret)
			assert.Equal(t, []string{"scope"}, github.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp github changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github
		addGithub, err := MgmtClient.AddGitHubProvider(IAMCTX, &management.AddGitHubProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			github, err = idpRepo.GetGithub(IAMCTX, idpRepo.IDCondition(addGithub.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addGithub.Id, github.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change github
		before := time.Now()
		_, err = MgmtClient.UpdateGitHubProvider(IAMCTX, &management.UpdateGitHubProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		// check values for azure
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGithub, err := idpRepo.GetGithub(IAMCTX, idpRepo.IDCondition(addGithub.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.github.changed
			// idp
			assert.Equal(t, instanceID, updateGithub.InstanceID)
			assert.Equal(t, orgID, *updateGithub.OrgID)
			assert.Equal(t, addGithub.Id, updateGithub.ID)
			assert.Equal(t, name, updateGithub.Name)
			assert.Equal(t, domain.IDPTypeGitHub, domain.IDPType(*updateGithub.Type))
			assert.Equal(t, true, updateGithub.AllowLinking)
			assert.Equal(t, true, updateGithub.AllowCreation)
			assert.Equal(t, true, updateGithub.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateGithub.AllowAutoLinkingField))
			assert.WithinRange(t, updateGithub.UpdatedAt, before, after)

			// github
			assert.Equal(t, "new_clientId", updateGithub.ClientID)
			assert.NotEqual(t, github.ClientSecret, updateGithub.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateGithub.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp github enterprise added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github enterprise
		before := time.Now()
		addGithubEnterprise, err := MgmtClient.AddGitHubEnterpriseServerProvider(IAMCTX, &management.AddGitHubEnterpriseServerProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for github enterprise
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			githubEnterprise, err := idpRepo.GetGithubEnterprise(IAMCTX, idpRepo.IDCondition(addGithubEnterprise.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.github_enterprise.added
			// idp
			assert.Equal(t, instanceID, githubEnterprise.InstanceID)
			assert.Equal(t, orgID, *githubEnterprise.OrgID)
			assert.Equal(t, addGithubEnterprise.Id, githubEnterprise.ID)
			assert.Equal(t, name, githubEnterprise.Name)
			assert.Equal(t, domain.IDPTypeGitHubEnterprise, domain.IDPType(*githubEnterprise.Type))
			assert.Equal(t, false, githubEnterprise.AllowLinking)
			assert.Equal(t, false, githubEnterprise.AllowCreation)
			assert.Equal(t, false, githubEnterprise.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*githubEnterprise.AllowAutoLinkingField))
			assert.WithinRange(t, githubEnterprise.CreatedAt, before, after)
			assert.WithinRange(t, githubEnterprise.UpdatedAt, before, after)

			// github enterprise
			assert.Equal(t, "clientId", githubEnterprise.ClientID)
			assert.NotNil(t, githubEnterprise.ClientSecret)
			assert.Equal(t, "authoizationEndpoint", githubEnterprise.AuthorizationEndpoint)
			assert.Equal(t, "tokenEndpoint", githubEnterprise.TokenEndpoint)
			assert.Equal(t, "userEndpoint", githubEnterprise.UserEndpoint)
			assert.Equal(t, []string{"scope"}, githubEnterprise.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp github enterprise changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add github enterprise
		addGithubEnterprise, err := MgmtClient.AddGitHubEnterpriseServerProvider(IAMCTX, &management.AddGitHubEnterpriseServerProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			githubEnterprise, err = idpRepo.GetGithubEnterprise(IAMCTX, idpRepo.IDCondition(addGithubEnterprise.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addGithubEnterprise.Id, githubEnterprise.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change github enterprise
		before := time.Now()
		_, err = MgmtClient.UpdateGitHubEnterpriseServerProvider(IAMCTX, &management.UpdateGitHubEnterpriseServerProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		// check values for azure
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGithubEnterprise, err := idpRepo.GetGithubEnterprise(IAMCTX, idpRepo.IDCondition(addGithubEnterprise.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.github_enterprise.changed
			// idp
			assert.Equal(t, instanceID, githubEnterprise.InstanceID)
			assert.Equal(t, orgID, *githubEnterprise.OrgID)
			assert.Equal(t, addGithubEnterprise.Id, updateGithubEnterprise.ID)
			assert.Equal(t, name, updateGithubEnterprise.Name)
			assert.Equal(t, domain.IDPTypeGitHubEnterprise, domain.IDPType(*updateGithubEnterprise.Type))
			assert.Equal(t, false, updateGithubEnterprise.AllowLinking)
			assert.Equal(t, false, updateGithubEnterprise.AllowCreation)
			assert.Equal(t, false, updateGithubEnterprise.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*updateGithubEnterprise.AllowAutoLinkingField))
			assert.WithinRange(t, updateGithubEnterprise.UpdatedAt, before, after)

			// github enterprise
			assert.Equal(t, "new_clientId", updateGithubEnterprise.ClientID)
			assert.NotNil(t, updateGithubEnterprise.ClientSecret)
			assert.Equal(t, "new_authoizationEndpoint", updateGithubEnterprise.AuthorizationEndpoint)
			assert.Equal(t, "new_tokenEndpoint", updateGithubEnterprise.TokenEndpoint)
			assert.Equal(t, "new_userEndpoint", updateGithubEnterprise.UserEndpoint)
			assert.Equal(t, []string{"new_scope"}, updateGithubEnterprise.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp gitlab added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab
		before := time.Now()
		addGithub, err := MgmtClient.AddGitLabProvider(IAMCTX, &management.AddGitLabProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for gitlab
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gitlab, err := idpRepo.GetGitlab(IAMCTX, idpRepo.IDCondition(addGithub.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.gitlab.added
			// idp
			assert.Equal(t, instanceID, gitlab.InstanceID)
			assert.Equal(t, orgID, *gitlab.OrgID)
			assert.Equal(t, addGithub.Id, gitlab.ID)
			assert.Equal(t, name, gitlab.Name)
			assert.Equal(t, domain.IDPTypeGitLab, domain.IDPType(*gitlab.Type))
			assert.Equal(t, false, gitlab.AllowLinking)
			assert.Equal(t, false, gitlab.AllowCreation)
			assert.Equal(t, false, gitlab.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*gitlab.AllowAutoLinkingField))
			assert.WithinRange(t, gitlab.CreatedAt, before, after)
			assert.WithinRange(t, gitlab.UpdatedAt, before, after)

			// gitlab
			assert.Equal(t, "clientId", gitlab.ClientID)
			assert.NotNil(t, gitlab.ClientSecret)
			assert.Equal(t, []string{"scope"}, gitlab.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp gitlab changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab
		addGitlab, err := MgmtClient.AddGitLabProvider(IAMCTX, &management.AddGitLabProviderRequest{
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

		var gitlab *domain.IDPGitlab
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gitlab, err = idpRepo.GetGitlab(IAMCTX, idpRepo.IDCondition(addGitlab.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addGitlab.Id, gitlab.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change gitlab
		before := time.Now()
		_, err = MgmtClient.UpdateGitLabProvider(IAMCTX, &management.UpdateGitLabProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		// check values for gitlab
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGitlab, err := idpRepo.GetGitlab(IAMCTX, idpRepo.IDCondition(addGitlab.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.gitlab.changed
			// idp
			assert.Equal(t, instanceID, updateGitlab.InstanceID)
			assert.Equal(t, orgID, *updateGitlab.OrgID)
			assert.Equal(t, addGitlab.Id, updateGitlab.ID)
			assert.Equal(t, name, updateGitlab.Name)
			assert.Equal(t, true, updateGitlab.AllowLinking)
			assert.Equal(t, true, updateGitlab.AllowCreation)
			assert.Equal(t, true, updateGitlab.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateGitlab.AllowAutoLinkingField))
			assert.WithinRange(t, updateGitlab.UpdatedAt, before, after)

			// gitlab
			assert.Equal(t, "new_clientId", updateGitlab.ClientID)
			assert.NotEqual(t, gitlab.ClientSecret, updateGitlab.ClientSecret)
			assert.Equal(t, domain.IDPTypeGitLab, domain.IDPType(*updateGitlab.Type))
			assert.Equal(t, []string{"new_scope"}, updateGitlab.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp gitlab self hosted added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab self hosted
		before := time.Now()
		addGitlabSelfHosted, err := MgmtClient.AddGitLabSelfHostedProvider(IAMCTX, &management.AddGitLabSelfHostedProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for gitlab self hosted
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gitlabSelfHosted, err := idpRepo.GetGitlabSelfHosting(IAMCTX, idpRepo.IDCondition(addGitlabSelfHosted.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.gitlab_self_hosted.added
			// idp
			assert.Equal(t, instanceID, gitlabSelfHosted.InstanceID)
			assert.Equal(t, orgID, *gitlabSelfHosted.OrgID)
			assert.Equal(t, addGitlabSelfHosted.Id, gitlabSelfHosted.ID)
			assert.Equal(t, name, gitlabSelfHosted.Name)
			assert.Equal(t, domain.IDPTypeGitLabSelfHosted, domain.IDPType(*gitlabSelfHosted.Type))
			assert.Equal(t, false, gitlabSelfHosted.AllowLinking)
			assert.Equal(t, false, gitlabSelfHosted.AllowCreation)
			assert.Equal(t, false, gitlabSelfHosted.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*gitlabSelfHosted.AllowAutoLinkingField))
			assert.WithinRange(t, gitlabSelfHosted.CreatedAt, before, after)
			assert.WithinRange(t, gitlabSelfHosted.UpdatedAt, before, after)

			// gitlab self hosted
			assert.Equal(t, "clientId", gitlabSelfHosted.ClientID)
			assert.Equal(t, "issuer", gitlabSelfHosted.Issuer)
			assert.NotNil(t, gitlabSelfHosted.ClientSecret)
			assert.Equal(t, []string{"scope"}, gitlabSelfHosted.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp gitlab self hosted changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add gitlab self hosted
		addGitlabSelfHosted, err := MgmtClient.AddGitLabSelfHostedProvider(IAMCTX, &management.AddGitLabSelfHostedProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			githlabSelfHosted, err = idpRepo.GetGitlabSelfHosting(IAMCTX, idpRepo.IDCondition(addGitlabSelfHosted.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addGitlabSelfHosted.Id, githlabSelfHosted.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change gitlab self hosted
		before := time.Now()
		_, err = MgmtClient.UpdateGitLabSelfHostedProvider(IAMCTX, &management.UpdateGitLabSelfHostedProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		// check values for gitlab self hosted
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGitlabSelfHosted, err := idpRepo.GetGitlabSelfHosting(IAMCTX, idpRepo.IDCondition(addGitlabSelfHosted.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.gitlab_self_hosted.changed
			// idp
			assert.Equal(t, instanceID, updateGitlabSelfHosted.InstanceID)
			assert.Equal(t, orgID, *updateGitlabSelfHosted.OrgID)
			assert.Equal(t, addGitlabSelfHosted.Id, updateGitlabSelfHosted.ID)
			assert.Equal(t, name, updateGitlabSelfHosted.Name)
			assert.Equal(t, domain.IDPTypeGitLabSelfHosted, domain.IDPType(*updateGitlabSelfHosted.Type))
			assert.Equal(t, true, updateGitlabSelfHosted.AllowLinking)
			assert.Equal(t, true, updateGitlabSelfHosted.AllowCreation)
			assert.Equal(t, true, updateGitlabSelfHosted.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateGitlabSelfHosted.AllowAutoLinkingField))
			assert.WithinRange(t, updateGitlabSelfHosted.UpdatedAt, before, after)

			// gitlab self hosted
			assert.Equal(t, "new_clientId", updateGitlabSelfHosted.ClientID)
			assert.Equal(t, "new_issuer", updateGitlabSelfHosted.Issuer)
			assert.NotEqual(t, githlabSelfHosted.ClientSecret, updateGitlabSelfHosted.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateGitlabSelfHosted.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp google added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add google
		before := time.Now()
		addGoogle, err := MgmtClient.AddGoogleProvider(IAMCTX, &management.AddGoogleProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		// check values for google
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			google, err := idpRepo.GetGoogle(IAMCTX, idpRepo.IDCondition(addGoogle.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.google.added
			// idp
			assert.Equal(t, instanceID, google.InstanceID)
			assert.Equal(t, orgID, *google.OrgID)
			assert.Equal(t, addGoogle.Id, google.ID)
			assert.Equal(t, name, google.Name)
			assert.Equal(t, domain.IDPTypeGoogle, domain.IDPType(*google.Type))
			assert.Equal(t, false, google.AllowLinking)
			assert.Equal(t, false, google.AllowCreation)
			assert.Equal(t, false, google.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*google.AllowAutoLinkingField))
			assert.WithinRange(t, google.CreatedAt, before, after)
			assert.WithinRange(t, google.UpdatedAt, before, after)

			// google
			assert.Equal(t, "clientId", google.ClientID)
			assert.NotNil(t, google.ClientSecret)
			assert.Equal(t, []string{"scope"}, google.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org idp google changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add google
		addGoogle, err := MgmtClient.AddGoogleProvider(IAMCTX, &management.AddGoogleProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			google, err = idpRepo.GetGoogle(IAMCTX, idpRepo.IDCondition(addGoogle.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addGoogle.Id, google.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change google
		before := time.Now()
		_, err = MgmtClient.UpdateGoogleProvider(IAMCTX, &management.UpdateGoogleProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		// check values for google
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateGoogle, err := idpRepo.GetGoogle(IAMCTX, idpRepo.IDCondition(addGoogle.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.google.changed
			// idp
			assert.Equal(t, instanceID, updateGoogle.InstanceID)
			assert.Equal(t, orgID, *updateGoogle.OrgID)
			assert.Equal(t, addGoogle.Id, updateGoogle.ID)
			assert.Equal(t, name, updateGoogle.Name)
			assert.Equal(t, domain.IDPTypeGoogle, domain.IDPType(*updateGoogle.Type))
			assert.Equal(t, true, updateGoogle.AllowLinking)
			assert.Equal(t, true, updateGoogle.AllowCreation)
			assert.Equal(t, true, updateGoogle.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateGoogle.AllowAutoLinkingField))
			assert.WithinRange(t, updateGoogle.UpdatedAt, before, after)

			// google
			assert.Equal(t, "new_clientId", updateGoogle.ClientID)
			assert.NotEqual(t, google.ClientSecret, updateGoogle.ClientSecret)
			assert.Equal(t, []string{"new_scope"}, updateGoogle.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org ldap added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add ldap
		before := time.Now()
		addLdap, err := MgmtClient.AddLDAPProvider(IAMCTX, &management.AddLDAPProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			ldap, err := idpRepo.GetLDAP(IAMCTX, idpRepo.IDCondition(addLdap.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.ldap.v2.added
			// idp
			assert.Equal(t, instanceID, ldap.InstanceID)
			assert.Equal(t, orgID, *ldap.OrgID)
			assert.Equal(t, addLdap.Id, ldap.ID)
			assert.Equal(t, name, ldap.Name)
			assert.Equal(t, domain.IDPTypeLDAP, domain.IDPType(*ldap.Type))
			assert.Equal(t, false, ldap.AllowLinking)
			assert.Equal(t, false, ldap.AllowCreation)
			assert.Equal(t, false, ldap.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*ldap.AllowAutoLinkingField))
			assert.WithinRange(t, ldap.CreatedAt, before, after)
			assert.WithinRange(t, ldap.UpdatedAt, before, after)

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

	t.Run("test org ldap changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add ldap
		addLdap, err := MgmtClient.AddLDAPProvider(IAMCTX, &management.AddLDAPProviderRequest{
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			ldap, err = idpRepo.GetLDAP(IAMCTX, idpRepo.IDCondition(addLdap.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addLdap.Id, ldap.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change ldap
		before := time.Now()
		_, err = MgmtClient.UpdateLDAPProvider(IAMCTX, &management.UpdateLDAPProviderRequest{
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
		after := time.Now()
		require.NoError(t, err)

		// check values for ldap
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateLdap, err := idpRepo.GetLDAP(IAMCTX, idpRepo.IDCondition(addLdap.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.ldap.v2.changed
			// idp
			assert.Equal(t, instanceID, updateLdap.InstanceID)
			assert.Equal(t, orgID, *updateLdap.OrgID)
			assert.Equal(t, addLdap.Id, updateLdap.ID)
			assert.Equal(t, name, updateLdap.Name)
			assert.Equal(t, domain.IDPTypeLDAP, domain.IDPType(*updateLdap.Type))
			assert.Equal(t, true, updateLdap.AllowLinking)
			assert.Equal(t, true, updateLdap.AllowCreation)
			assert.Equal(t, true, updateLdap.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateLdap.AllowAutoLinkingField))
			assert.WithinRange(t, updateLdap.UpdatedAt, before, after)

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

	t.Run("test org apple added reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add apple
		before := time.Now()
		addApple, err := MgmtClient.AddAppleProvider(IAMCTX, &management.AddAppleProviderRequest{
			Name:       name,
			ClientId:   "clientID",
			TeamId:     "teamIDteam",
			KeyId:      "keyIDKeyId",
			PrivateKey: []byte("privateKey"),
			Scopes:     []string{"scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
		})
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			apple, err := idpRepo.GetApple(IAMCTX, idpRepo.IDCondition(addApple.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.apple.added
			// idp
			assert.Equal(t, instanceID, apple.InstanceID)
			assert.Equal(t, orgID, *apple.OrgID)
			assert.Equal(t, addApple.Id, apple.ID)
			assert.Equal(t, name, apple.Name)
			assert.Equal(t, domain.IDPTypeApple, domain.IDPType(*apple.Type))
			assert.Equal(t, false, apple.AllowLinking)
			assert.Equal(t, false, apple.AllowCreation)
			assert.Equal(t, false, apple.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*apple.AllowAutoLinkingField))
			assert.WithinRange(t, apple.CreatedAt, before, after)
			assert.WithinRange(t, apple.UpdatedAt, before, after)

			// apple
			assert.Equal(t, "clientID", apple.ClientID)
			assert.Equal(t, "teamIDteam", apple.TeamID)
			assert.Equal(t, "keyIDKeyId", apple.KeyID)
			assert.NotNil(t, apple.PrivateKey)
			assert.Equal(t, []string{"scope"}, apple.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org apple changed reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add apple
		addApple, err := MgmtClient.AddAppleProvider(IAMCTX, &management.AddAppleProviderRequest{
			Name:       name,
			ClientId:   "clientID",
			TeamId:     "teamIDteam",
			KeyId:      "keyIDKeyId",
			PrivateKey: []byte("privateKey"),
			Scopes:     []string{"scope"},
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

		var apple *domain.IDPApple
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			apple, err = idpRepo.GetApple(IAMCTX, idpRepo.IDCondition(addApple.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addApple.Id, apple.ID)
		}, retryDuration, tick)

		name = "new_" + name
		// change apple
		before := time.Now()
		_, err = MgmtClient.UpdateAppleProvider(IAMCTX, &management.UpdateAppleProviderRequest{
			Id:         addApple.Id,
			Name:       name,
			ClientId:   "new_clientID",
			TeamId:     "new_teamID",
			KeyId:      "new_kKeyId",
			PrivateKey: []byte("new_privateKey"),
			Scopes:     []string{"new_scope"},
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
		})
		after := time.Now()
		require.NoError(t, err)

		// check values for apple
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateApple, err := idpRepo.GetApple(IAMCTX, idpRepo.IDCondition(addApple.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event nstance.idp.apple.changed
			// idp
			assert.Equal(t, instanceID, updateApple.InstanceID)
			assert.Equal(t, orgID, *updateApple.OrgID)
			assert.Equal(t, addApple.Id, updateApple.ID)
			assert.Equal(t, name, updateApple.Name)
			assert.Equal(t, domain.IDPTypeApple, domain.IDPType(*updateApple.Type))
			assert.Equal(t, true, updateApple.AllowLinking)
			assert.Equal(t, true, updateApple.AllowCreation)
			assert.Equal(t, true, updateApple.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateApple.AllowAutoLinkingField))
			assert.WithinRange(t, updateApple.UpdatedAt, before, after)

			// apple
			assert.Equal(t, "new_clientID", updateApple.ClientID)
			assert.Equal(t, "new_teamID", updateApple.TeamID)
			assert.Equal(t, "new_kKeyId", updateApple.KeyID)
			assert.NotEqual(t, apple.PrivateKey, updateApple.PrivateKey)
			assert.Equal(t, []string{"new_scope"}, updateApple.Scopes)
		}, retryDuration, tick)
	})

	t.Run("test org saml added reduces", func(t *testing.T) {
		name := gofakeit.Name()
		federatedLogoutEnabled := false

		// add saml
		before := time.Now()
		addSAML, err := MgmtClient.AddSAMLProvider(IAMCTX, &management.AddSAMLProviderRequest{
			Name: name,
			Metadata: &management.AddSAMLProviderRequest_MetadataXml{
				MetadataXml: validSAMLMetadata1,
			},
			Binding:                       idp.SAMLBinding_SAML_BINDING_POST,
			WithSignedRequest:             false,
			TransientMappingAttributeName: &name,
			FederatedLogoutEnabled:        &federatedLogoutEnabled,
			NameIdFormat:                  idp.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_TRANSIENT.Enum(),
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			SignatureAlgorithm: idp.SAMLSignatureAlgorithm_SAML_SIGNATURE_RSA_SHA1,
		})
		after := time.Now()
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			saml, err := idpRepo.GetSAML(IAMCTX, idpRepo.IDCondition(addSAML.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.saml.added
			// idp
			assert.Equal(t, instanceID, saml.InstanceID)
			assert.Equal(t, orgID, *saml.OrgID)
			assert.Equal(t, addSAML.Id, saml.ID)
			assert.Equal(t, name, saml.Name)
			assert.Equal(t, domain.IDPTypeSAML, domain.IDPType(*saml.Type))
			assert.Equal(t, false, saml.AllowLinking)
			assert.Equal(t, false, saml.AllowCreation)
			assert.Equal(t, false, saml.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldEmail, domain.IDPAutoLinkingField(*saml.AllowAutoLinkingField))
			assert.WithinRange(t, saml.CreatedAt, before, after)
			assert.WithinRange(t, saml.UpdatedAt, before, after)

			// saml
			assert.Equal(t, validSAMLMetadata1, saml.Metadata)
			assert.NotNil(t, saml.Key)
			assert.NotNil(t, saml.Certificate)
			assert.NotNil(t, saml.Binding)
			assert.Equal(t, false, saml.WithSignedRequest)
			assert.Equal(t, zitadel_internal_domain.SAMLNameIDFormatTransient, *saml.NameIDFormat)
			assert.Equal(t, name, saml.TransientMappingAttributeName)
			assert.Equal(t, false, saml.FederatedLogoutEnabled)
			assert.Equal(t, "http://www.w3.org/2000/09/xmldsig#rsa-sha1", saml.SignatureAlgorithm)
		}, retryDuration, tick)
	})

	t.Run("test org saml changed reduces", func(t *testing.T) {
		name := gofakeit.Name()
		federatedLogoutEnabled := false

		// add saml
		addSAML, err := MgmtClient.AddSAMLProvider(IAMCTX, &management.AddSAMLProviderRequest{
			Name: name,
			Metadata: &management.AddSAMLProviderRequest_MetadataXml{
				MetadataXml: validSAMLMetadata1,
			},
			Binding:                       idp.SAMLBinding_SAML_BINDING_POST,
			WithSignedRequest:             false,
			TransientMappingAttributeName: &name,
			FederatedLogoutEnabled:        &federatedLogoutEnabled,
			NameIdFormat:                  idp.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_TRANSIENT.Enum(),
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  false,
				IsCreationAllowed: false,
				IsAutoCreation:    false,
				IsAutoUpdate:      false,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL,
			},
			SignatureAlgorithm: idp.SAMLSignatureAlgorithm_SAML_SIGNATURE_RSA_SHA1,
		})
		require.NoError(t, err)

		idpRepo := repository.IDProviderRepository(pool)

		var saml *domain.IDPSAML
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			saml, err = idpRepo.GetSAML(IAMCTX, idpRepo.IDCondition(addSAML.Id), instanceID, &orgID)
			require.NoError(t, err)
			assert.Equal(t, addSAML.Id, saml.ID)
		}, retryDuration, tick)

		name = "new_" + name
		federatedLogoutEnabled = true
		// change saml
		before := time.Now()
		_, err = MgmtClient.UpdateSAMLProvider(IAMCTX, &management.UpdateSAMLProviderRequest{
			Id:   addSAML.Id,
			Name: name,
			Metadata: &management.UpdateSAMLProviderRequest_MetadataXml{
				MetadataXml: validSAMLMetadata2,
			},
			Binding:                       idp.SAMLBinding_SAML_BINDING_ARTIFACT,
			WithSignedRequest:             true,
			TransientMappingAttributeName: &name,
			FederatedLogoutEnabled:        &federatedLogoutEnabled,
			NameIdFormat:                  idp.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_EMAIL_ADDRESS.Enum(),
			ProviderOptions: &idp_grpc.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
				AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
			},
			SignatureAlgorithm: idp.SAMLSignatureAlgorithm_SAML_SIGNATURE_RSA_SHA256,
		})
		after := time.Now()
		require.NoError(t, err)

		// check values for apple
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			updateSAML, err := idpRepo.GetSAML(IAMCTX, idpRepo.IDCondition(addSAML.Id), instanceID, &orgID)
			require.NoError(t, err)

			// event org.idp.saml.changed
			// idp
			assert.Equal(t, instanceID, updateSAML.InstanceID)
			assert.Equal(t, orgID, *updateSAML.OrgID)
			assert.Equal(t, addSAML.Id, updateSAML.ID)
			assert.Equal(t, name, updateSAML.Name)
			assert.Equal(t, domain.IDPTypeSAML, domain.IDPType(*updateSAML.Type))
			assert.Equal(t, true, updateSAML.AllowLinking)
			assert.Equal(t, true, updateSAML.AllowCreation)
			assert.Equal(t, true, updateSAML.AllowAutoUpdate)
			assert.Equal(t, domain.IDPAutoLinkingFieldUserName, domain.IDPAutoLinkingField(*updateSAML.AllowAutoLinkingField))
			assert.WithinRange(t, updateSAML.UpdatedAt, before, after)

			// saml
			assert.Equal(t, validSAMLMetadata2, updateSAML.Metadata)
			assert.NotNil(t, updateSAML.Key)
			assert.NotNil(t, updateSAML.Certificate)
			assert.NotNil(t, updateSAML.Binding)
			assert.NotEqual(t, saml.Binding, updateSAML.Binding)
			assert.Equal(t, true, updateSAML.WithSignedRequest)
			assert.Equal(t, zitadel_internal_domain.SAMLNameIDFormatEmailAddress, *updateSAML.NameIDFormat)
			assert.Equal(t, name, updateSAML.TransientMappingAttributeName)
			assert.Equal(t, true, updateSAML.FederatedLogoutEnabled)
			assert.Equal(t, "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256", updateSAML.SignatureAlgorithm)
		}, retryDuration, tick)
	})

	t.Run("test org iam remove reduces", func(t *testing.T) {
		name := gofakeit.Name()

		// add idp
		addOIDC, err := MgmtClient.AddOrgOIDCIDP(IAMCTX, &management.AddOrgOIDCIDPRequest{
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

		// check idp exists
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := idpRepo.Get(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)
			require.NoError(t, err)
		}, retryDuration, tick)

		// remove idp
		_, err = MgmtClient.DeleteProvider(IAMCTX, &management.DeleteProviderRequest{
			Id: addOIDC.IdpId,
		})
		require.NoError(t, err)

		// check idp is removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := idpRepo.Get(IAMCTX,
				idpRepo.IDCondition(addOIDC.IdpId),
				instanceID,
				&orgID,
			)
			require.ErrorIs(t, &database.NoRowFoundError{}, err)
		}, retryDuration, tick)
	})
}
