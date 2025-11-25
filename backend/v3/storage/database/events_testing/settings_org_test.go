//go:build integration

package events_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/policy"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
)

func TestServer_TestOrgLoginSettingsReduces(t *testing.T) {
	settingsRepo := repository.LoginSettingsRepository()

	t.Run("test adding login settings reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddCustomLoginPolicy(IAMCTX, &management.AddCustomLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           false,
			ForceMfa:                   false,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          false,
			IgnoreUnknownUsernames:     false,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 10),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 10),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 10),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 10),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 10),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      false,
			DisableLoginWithPhone:      false,
			ForceMfaLocalOnly:          false,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.added
			// these values are found in default.yaml
			assert.Equal(t, false, setting.AllowRegister)
			assert.Equal(t, false, setting.AllowExternalIDP)
			assert.Equal(t, domain.PasswordlessTypeNotAllowed, setting.PasswordlessType)
			assert.Equal(t, false, setting.AllowDomainDiscovery)
			assert.Equal(t, false, setting.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Second*10), setting.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.ExternalLoginCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.MFAInitSkipLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.MultiFactorCheckLifetime)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test change login settings reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// create org + login policy
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLoginPolicy(IAMCTX, &management.AddCustomLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           false,
			ForceMfa:                   false,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          false,
			IgnoreUnknownUsernames:     false,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 10),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 10),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 10),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 10),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 10),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      false,
			DisableLoginWithPhone:      false,
			ForceMfaLocalOnly:          false,
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// update org policy
		before := time.Now()
		_, err = newInstance.Client.Mgmt.UpdateCustomLoginPolicy(IAMCTX, &management.UpdateCustomLoginPolicyRequest{
			AllowUsernamePassword:      true,
			AllowRegister:              true,
			AllowExternalIdp:           true,
			ForceMfa:                   true,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_ALLOWED,
			HidePasswordReset:          true,
			IgnoreUnknownUsernames:     true,
			DefaultRedirectUri:         "http://www.new_example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 5 * 20),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 5 * 21),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 5 * 22),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 5 * 23),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 5 * 24),
			AllowDomainDiscovery:       true,
			DisableLoginWithEmail:      true,
			DisableLoginWithPhone:      true,
			ForceMfaLocalOnly:          true,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.changed
			assert.Equal(t, true, setting.AllowRegister)
			assert.Equal(t, true, setting.AllowExternalIDP)
			assert.Equal(t, true, setting.ForceMFA)
			assert.Equal(t, domain.PasswordlessTypeAllowed, setting.PasswordlessType)
			assert.Equal(t, true, setting.HidePasswordReset)
			assert.Equal(t, true, setting.IgnoreUnknownUsernames)
			assert.Equal(t, "http://www.new_example.com", setting.DefaultRedirectURI)
			assert.Equal(t, true, setting.AllowDomainDiscovery)
			assert.Equal(t, true, setting.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Second*5*20), setting.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*21), setting.ExternalLoginCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*22), setting.MFAInitSkipLifetime)
			assert.Equal(t, time.Duration(time.Second*5*23), setting.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*24), setting.MultiFactorCheckLifetime)
			assert.Equal(t, true, setting.DisableLoginWithEmail)
			assert.Equal(t, true, setting.DisableLoginWithPhone)
			assert.Equal(t, true, setting.ForceMFALocalOnly)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test added/remove multifactor type reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// create org + login policy
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.AddCustomLoginPolicy(IAMCTX, &management.AddCustomLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           false,
			ForceMfa:                   false,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          false,
			IgnoreUnknownUsernames:     false,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 10),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 10),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 10),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 10),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 10),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      false,
			DisableLoginWithPhone:      false,
			ForceMfaLocalOnly:          false,
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// check initial MFAType value
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Nil(t, setting.MFAType)
		}, retryDuration, tick)

		// add MFAType
		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddMultiFactorToLoginPolicy(IAMCTX, &management.AddMultiFactorToLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.multifactor.added
			assert.Equal(t, []domain.MultiFactorType{domain.MultiFactorType(policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION)}, setting.MFAType)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)

		// remove MFAType
		_, err = newInstance.Client.Mgmt.RemoveMultiFactorFromLoginPolicy(IAMCTX, &management.RemoveMultiFactorFromLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.multifactor.remove
			assert.Equal(t, []domain.MultiFactorType{}, setting.MFAType)
		}, retryDuration, tick)
	})

	t.Run("test added/removed second multifactor reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// create org + login policy
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.AddCustomLoginPolicy(IAMCTX, &management.AddCustomLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           false,
			ForceMfa:                   false,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          false,
			IgnoreUnknownUsernames:     false,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 10),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 10),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 10),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 10),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 10),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      false,
			DisableLoginWithPhone:      false,
			ForceMfaLocalOnly:          false,
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		before := time.Now()
		// get current second factor types
		var secondFactorTypes []domain.SecondFactorType
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			secondFactorTypes = setting.SecondFactorTypes
		}, retryDuration, tick)

		// add new second factor type
		before = time.Now()
		_, err = newInstance.Client.Mgmt.AddSecondFactorToLoginPolicy(IAMCTX, &management.AddSecondFactorToLoginPolicyRequest{
			Type: policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		secondFactorTypes = append(secondFactorTypes, domain.SecondFactorType(policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS))

		// check new second factor type is added
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.multifactor.secondfactor.added
			assert.Equal(t, secondFactorTypes, setting.SecondFactorTypes)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)

		// remove second factor type
		before = time.Now()
		_, err = newInstance.Client.Mgmt.RemoveSecondFactorFromLoginPolicy(IAMCTX, &management.RemoveSecondFactorFromLoginPolicyRequest{
			Type: policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		})
		require.NoError(t, err)
		after = time.Now()

		secondFactorTypes = secondFactorTypes[0 : len(secondFactorTypes)-1]

		// check new second factor type is removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.multifactor.secondfactor.removed
			assert.Equal(t, secondFactorTypes, setting.SecondFactorTypes)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test login settings removed reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// create org + login policy
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		orgId := organization.Id
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)
		_, err = newInstance.Client.Mgmt.AddCustomLoginPolicy(IAMCTX, &management.AddCustomLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           false,
			ForceMfa:                   false,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          false,
			IgnoreUnknownUsernames:     false,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 10),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 10),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 10),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 10),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 10),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      false,
			DisableLoginWithPhone:      false,
			ForceMfaLocalOnly:          false,
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.ResetLoginPolicyToDefault(IAMCTX, &management.ResetLoginPolicyToDefaultRequest{})
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)

			// event org.policy.login.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// create org + login policy
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLoginPolicy(IAMCTX, &management.AddCustomLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           false,
			ForceMfa:                   false,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          false,
			IgnoreUnknownUsernames:     false,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 10),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 10),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 10),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 10),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 10),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      false,
			DisableLoginWithPhone:      false,
			ForceMfaLocalOnly:          false,
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgLabelSettingsReduces(t *testing.T) {
	settingsRepo := repository.BrandingSettingsRepository()

	t.Run("test adding label settings reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.added
			assert.Equal(t, "#055090", setting.PrimaryColorLight)
			assert.Equal(t, "#055090", setting.BackgroundColorLight)
			assert.Equal(t, "#055090", setting.WarnColorLight)
			assert.Equal(t, "#055090", setting.FontColorLight)
			assert.Equal(t, "#055090", setting.PrimaryColorDark)
			assert.Equal(t, "#055090", setting.BackgroundColorDark)
			assert.Equal(t, "#055090", setting.WarnColorDark)
			assert.Equal(t, "#055090", setting.WarnColorDark)
			assert.Equal(t, "#055090", setting.FontColorDark)
			assert.Equal(t, false, setting.HideLoginNameSuffix)
			assert.Equal(t, false, setting.ErrorMsgPopup)
			assert.Equal(t, false, setting.DisableWatermark)
			assert.Equal(t, domain.BrandingPolicyThemeDark, setting.ThemeMode)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label change", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		})
		require.NoError(t, err)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.UpdateCustomLabelPolicy(IAMCTX, &management.UpdateCustomLabelPolicyRequest{
			PrimaryColor:        "#055000",
			HideLoginNameSuffix: true,
			WarnColor:           "#055000",
			BackgroundColor:     "#055000",
			FontColor:           "#055000",
			PrimaryColorDark:    "#055000",
			BackgroundColorDark: "#055000",
			WarnColorDark:       "#055000",
			FontColorDark:       "#055000",
			DisableWatermark:    true,
			ThemeMode:           policy.ThemeMode_THEME_MODE_LIGHT,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.change
			assert.Equal(t, "#055000", setting.PrimaryColorLight)
			assert.Equal(t, "#055000", setting.BackgroundColorLight)
			assert.Equal(t, "#055000", setting.WarnColorLight)
			assert.Equal(t, "#055000", setting.FontColorLight)
			assert.Equal(t, "#055000", setting.PrimaryColorDark)
			assert.Equal(t, "#055000", setting.BackgroundColorDark)
			assert.Equal(t, "#055000", setting.WarnColorDark)
			assert.Equal(t, "#055000", setting.WarnColorDark)
			assert.Equal(t, "#055000", setting.FontColorDark)
			assert.Equal(t, true, setting.HideLoginNameSuffix)
			assert.Equal(t, false, setting.ErrorMsgPopup)
			assert.Equal(t, true, setting.DisableWatermark)
			assert.Equal(t, domain.BrandingPolicyThemeLight, setting.ThemeMode)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test label settings activated", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		})
		require.NoError(t, err)

		// check label policy exists in preview
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, domain.SettingStatePreview, setting.State)
		}, retryDuration, tick)

		// check label policy does not exists in activated
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)

			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)

		// activate label
		before := time.Now()
		_, err = newInstance.Client.Mgmt.ActivateCustomLabelPolicy(IAMCTX, &management.ActivateCustomLabelPolicyRequest{})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.activated
			assert.Equal(t, domain.SettingStateActive, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo light added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set logo light
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/logo")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.added
			assert.NotNil(t, setting.LogoURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set logo dark
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/logo/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.dark.added
			assert.NotNil(t, setting.LogoURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo light removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set logo light
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/logo")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check light logo set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			assert.NotNil(t, setting.LogoURLLight)
		}, retryDuration, tick)

		// remote logo light
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label/logo")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now().Add(time.Second * 30)

		// check light logo removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.removed
			assert.Nil(t, setting.LogoURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set logo dark
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/logo/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check dark logo set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			assert.NotNil(t, setting.LogoURLDark)
		}, retryDuration, tick)

		// remote logo dark
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label/logo_dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now().Add(time.Second * 30)

		// check dark logo removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.dark.removed
			assert.Nil(t, setting.LogoURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set icon light
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/icon")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.added
			assert.NotNil(t, setting.IconURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set icon dark
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/icon/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.dark.added
			assert.NotNil(t, setting.IconURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set icon light
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/icon")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check light icon set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			assert.NotNil(t, setting.IconURLLight)
		}, retryDuration, tick)

		// remote icon light
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label/icon")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now().Add(time.Second * 30)

		// check light icon removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.removed
			assert.Nil(t, setting.IconURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set icon dark
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/icon/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check dark icon set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// assert.Equal(t, domain.LabelPolicyThemeLight, setting.LogoURLLight)
			assert.NotNil(t, setting.IconURLDark)
		}, retryDuration, tick)

		// remote icon dark
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label/icon_dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now().Add(time.Second * 30)

		// check dark icon removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.dark.removed
			assert.Nil(t, setting.IconURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy font added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set logo light
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(font)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/font")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.font.added
			assert.NotNil(t, setting.FontURL)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy font removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)
		orgRepo := repository.OrganizationRepository()
		org_, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(orgRepo.InstanceIDCondition(instanceID), orgRepo.NameCondition(database.TextOperationEqual, "ZITADEL"))))
		require.NoError(t, err)
		orgId := org_.ID

		// delete previous label settings
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// add label policy
		label := management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		}
		setLabelPolicyPayload, err := json.Marshal(&label)
		require.NoError(t, err)
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			SetBody(setLabelPolicyPayload).
			Post("http://localhost:8082" + "/management/v1" + "/policies/label")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// set logo light
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(font)).
			SetHeader("x-zitadel-orgid", orgId).
			Post("http://localhost:8082" + "/assets/v1" + "/org/policy/label/font")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			assert.NotNil(t, setting.FontURL)
		}, retryDuration, tick)

		// remote font policy
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			SetHeader("x-zitadel-orgid", orgId).
			Delete("http://localhost:8082" + "/management/v1" + "/policies/label/font")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now().Add(time.Second * 30)

		// check font policy removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(instanceID, &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.font.removed
			assert.Nil(t, setting.FontURL)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test label settings remove reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		})
		require.NoError(t, err)

		// check label settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove label policy delete org
		_, err = newInstance.Client.Mgmt.ResetLabelPolicyToDefault(IAMCTX, &management.ResetLabelPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check label label settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)

			// event org.policy.label.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test org delete reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
			PrimaryColor:        "#055090",
			HideLoginNameSuffix: false,
			WarnColor:           "#055090",
			BackgroundColor:     "#055090",
			FontColor:           "#055090",
			PrimaryColorDark:    "#055090",
			BackgroundColorDark: "#055090",
			WarnColorDark:       "#055090",
			FontColorDark:       "#055090",
			DisableWatermark:    false,
			ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
		})
		require.NoError(t, err)

		// check label settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check label label settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgPasswordComplexitySettingsReduces(t *testing.T) {
	settingsRepo := repository.PasswordComplexitySettingsRepository()

	t.Run("test password complexity added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddCustomPasswordComplexityPolicy(IAMCTX, &management.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    10,
			HasUppercase: true,
			HasLowercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.complexity.added
			assert.Equal(t, uint64(10), setting.MinLength)
			assert.Equal(t, true, setting.HasUppercase)
			assert.Equal(t, true, setting.HasLowercase)
			assert.Equal(t, true, setting.HasNumber)
			assert.Equal(t, true, setting.HasSymbol)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password complexity change", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		// create password complexity
		_, err = newInstance.Client.Mgmt.AddCustomPasswordComplexityPolicy(IAMCTX, &management.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    10,
			HasUppercase: false,
			HasLowercase: false,
			HasNumber:    false,
			HasSymbol:    false,
		})
		require.NoError(t, err)

		// update password compexity
		before := time.Now()
		_, err = newInstance.Client.Mgmt.UpdateCustomPasswordComplexityPolicy(IAMCTX, &management.UpdateCustomPasswordComplexityPolicyRequest{
			MinLength:    5,
			HasUppercase: true,
			HasLowercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.complexity.changed
			assert.Equal(t, uint64(5), setting.MinLength)
			assert.Equal(t, true, setting.HasUppercase)
			assert.Equal(t, true, setting.HasLowercase)
			assert.Equal(t, true, setting.HasNumber)
			assert.Equal(t, true, setting.HasSymbol)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password complexity removed", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		// create password complexity
		_, err = newInstance.Client.Mgmt.AddCustomPasswordComplexityPolicy(IAMCTX, &management.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    10,
			HasUppercase: false,
			HasLowercase: false,
			HasNumber:    false,
			HasSymbol:    false,
		})
		require.NoError(t, err)

		// check password complexity settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete password complexity policy
		_, err = newInstance.Client.Mgmt.ResetPasswordComplexityPolicyToDefault(IAMCTX, &management.ResetPasswordComplexityPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		// create password complexity
		_, err = newInstance.Client.Mgmt.AddCustomPasswordComplexityPolicy(IAMCTX, &management.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    10,
			HasUppercase: false,
			HasLowercase: false,
			HasNumber:    false,
			HasSymbol:    false,
		})
		require.NoError(t, err)

		// check password complexity settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgPasswordPolicySettingsReduces(t *testing.T) {
	settingsRepo := repository.PasswordExpirySettingsRepository()

	t.Run("test password policy added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		// add add password policy
		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     10,
			ExpireWarnDays: 10,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.age.added
			assert.Equal(t, uint64(10), setting.ExpireWarnDays)
			assert.Equal(t, uint64(10), setting.MaxAgeDays)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy change reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     10,
			ExpireWarnDays: 10,
		})
		require.NoError(t, err)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.UpdateCustomPasswordAgePolicy(IAMCTX, &management.UpdateCustomPasswordAgePolicyRequest{
			MaxAgeDays:     40,
			ExpireWarnDays: 40,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.age.changed
			assert.Equal(t, uint64(40), setting.ExpireWarnDays)
			assert.Equal(t, uint64(40), setting.MaxAgeDays)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy removed reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     0,
			ExpireWarnDays: 0,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove password policy delete org
		_, err = newInstance.Client.Mgmt.ResetPasswordAgePolicyToDefault(IAMCTX, &management.ResetPasswordAgePolicyToDefaultRequest{})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)

			// event org.policy.password.age.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     0,
			ExpireWarnDays: 0,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgLockoutSettingsReduces(t *testing.T) {
	settingsRepo := repository.LockoutSettingsRepository()

	t.Run("test lockout policy added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 1,
			MaxOtpAttempts:      1,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.lockout.added
			assert.Equal(t, uint64(1), setting.MaxOTPAttempts)
			assert.Equal(t, uint64(1), setting.MaxPasswordAttempts)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test lockout policy changed", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 1,
			MaxOtpAttempts:      1,
		})
		require.NoError(t, err)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.UpdateCustomLockoutPolicy(IAMCTX, &management.UpdateCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 5,
			MaxOtpAttempts:      5,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.lockout.changed
			assert.Equal(t, uint64(5), setting.MaxOTPAttempts)
			assert.Equal(t, uint64(5), setting.MaxPasswordAttempts)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test remove lockout policy reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 1,
			MaxOtpAttempts:      1,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove lockout policy org
		_, err = newInstance.Client.Mgmt.ResetLockoutPolicyToDefault(IAMCTX, &management.ResetLockoutPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, orgId)

		_, err = newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 1,
			MaxOtpAttempts:      1,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgDomainSettingsReduces(t *testing.T) {
	settingsRepo := repository.DomainSettingsRepository()

	t.Run("test domain policy added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id

		before := time.Now()
		_, err = newInstance.Client.Admin.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
			OrgId:                                  organization.Id,
			UserLoginMustBeDomain:                  false,
			ValidateOrgDomains:                     false,
			SmtpSenderAddressMatchesInstanceDomain: false,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.domain.added
			assert.Equal(t, false, setting.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, false, setting.LoginNameIncludesDomain)
			assert.Equal(t, false, setting.RequireOrgDomainVerification)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test domain policy changed", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id

		// add domain policy
		_, err = newInstance.Client.Admin.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
			OrgId:                                  organization.Id,
			UserLoginMustBeDomain:                  false,
			ValidateOrgDomains:                     false,
			SmtpSenderAddressMatchesInstanceDomain: false,
		})
		require.NoError(t, err)

		// update domain policy
		before := time.Now()
		_, err = newInstance.Client.Admin.UpdateCustomDomainPolicy(IAMCTX, &admin.UpdateCustomDomainPolicyRequest{
			OrgId:                                  organization.Id,
			UserLoginMustBeDomain:                  true,
			ValidateOrgDomains:                     true,
			SmtpSenderAddressMatchesInstanceDomain: true,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.domain.changed
			assert.Equal(t, true, setting.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, true, setting.LoginNameIncludesDomain)
			assert.Equal(t, true, setting.RequireOrgDomainVerification)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test remove domain policy reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		// IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// add domain policy
		_, err = newInstance.Client.Admin.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
			OrgId:                                  organization.Id,
			UserLoginMustBeDomain:                  false,
			ValidateOrgDomains:                     false,
			SmtpSenderAddressMatchesInstanceDomain: false,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove domain policy org
		_, err = newInstance.Client.Admin.ResetCustomDomainPolicyToDefault(IAMCTX, &admin.ResetCustomDomainPolicyToDefaultRequest{
			OrgId: organization.Id,
		})
		require.NoError(t, err)

		// check domain settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)

			// event org.policy.domain.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id

		// add domain policy
		_, err = newInstance.Client.Admin.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
			OrgId:                                  organization.Id,
			UserLoginMustBeDomain:                  false,
			ValidateOrgDomains:                     false,
			SmtpSenderAddressMatchesInstanceDomain: false,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgSettingsReduces(t *testing.T) {
	settingsRepo := repository.OrganizationSettingsRepository()

	t.Run("test add org settings added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		organizationScopedUsernames := true
		before := time.Now()
		_, err = newInstance.Client.SettingsV2beta.SetOrganizationSettings(IAMCTX, &settings.SetOrganizationSettingsRequest{
			OrganizationId:              organization.Id,
			OrganizationScopedUsernames: &organizationScopedUsernames,
		})
		require.NoError(t, err)
		after := time.Now().Add(time.Second * 30)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event settings.organization.set
			assert.Equal(t, organizationScopedUsernames, setting.OrganizationScopedUsernames)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test organization settings removed reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// create organization setting
		organizationScopedUsernames := true
		_, err = newInstance.Client.SettingsV2beta.SetOrganizationSettings(IAMCTX, &settings.SetOrganizationSettingsRequest{
			OrganizationId:              organization.Id,
			OrganizationScopedUsernames: &organizationScopedUsernames,
		})
		require.NoError(t, err)

		// check organization settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete organization setting
		_, err = newInstance.Client.SettingsV2beta.DeleteOrganizationSettings(IAMCTX, &settings.DeleteOrganizationSettingsRequest{
			OrganizationId: organization.Id,
		})
		require.NoError(t, err)

		// check organization settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)

			// event settings.organization.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test organization removed reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// create organization setting
		organizationScopedUsernames := true
		_, err = newInstance.Client.SettingsV2beta.SetOrganizationSettings(IAMCTX, &settings.SetOrganizationSettingsRequest{
			OrganizationId:              organization.Id,
			OrganizationScopedUsernames: &organizationScopedUsernames,
		})
		require.NoError(t, err)

		// check organization settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check organization settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		orgId := organization.Id
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// create organization setting
		organizationScopedUsernames := true
		_, err = newInstance.Client.SettingsV2beta.SetOrganizationSettings(IAMCTX, &settings.SetOrganizationSettingsRequest{
			OrganizationId:              organization.Id,
			OrganizationScopedUsernames: &organizationScopedUsernames,
		})
		require.NoError(t, err)

		// check organization settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check organization settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				ctx, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)

			// event settings.organization.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
