//go:build integration

package events_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/policy"
	user "github.com/zitadel/zitadel/pkg/grpc/user"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer_TestOrgLoginSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository(pool)

	t.Run("test adding login settings reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

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
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			// nil)
			require.NoError(t, err)

			// event org.policy.login.added
			// these values are found in default.yaml
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, false, setting.Settings.AllowRegister)
			assert.Equal(t, false, setting.Settings.AllowExternalSetting)
			assert.Equal(t, domain.PasswordlessTypeNotAllowed, setting.Settings.PasswordlessType)
			assert.Equal(t, false, setting.Settings.AllowDomainDiscovery)
			assert.Equal(t, false, setting.Settings.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Second*10), setting.Settings.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.Settings.ExternalLoginCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.Settings.MFAInitSkipLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.Settings.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), setting.Settings.MultiFactorCheckLifetime)
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
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event org.policy.login.changed
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, true, setting.Settings.AllowRegister)
			assert.Equal(t, true, setting.Settings.AllowExternalSetting)
			assert.Equal(t, true, setting.Settings.ForceMFA)
			assert.Equal(t, domain.PasswordlessTypeAllowed, setting.Settings.PasswordlessType)
			assert.Equal(t, true, setting.Settings.HidePasswordReset)
			assert.Equal(t, true, setting.Settings.IgnoreUnknownUsernames)
			assert.Equal(t, "http://www.new_example.com", setting.Settings.DefaultRedirectURI)
			assert.Equal(t, true, setting.Settings.AllowDomainDiscovery)
			assert.Equal(t, true, setting.Settings.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Second*5*20), setting.Settings.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*21), setting.Settings.ExternalLoginCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*22), setting.Settings.MFAInitSkipLifetime)
			assert.Equal(t, time.Duration(time.Second*5*23), setting.Settings.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*24), setting.Settings.MultiFactorCheckLifetime)
			assert.Equal(t, true, setting.Settings.DisableLoginWithEmail)
			assert.Equal(t, true, setting.Settings.DisableLoginWithPhone)
			assert.Equal(t, true, setting.Settings.ForceMFALocalOnly)
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

		// check inital MFAType value
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// assert.Equal(t, []domain.MultiFactorType{domain.MultiFactorType(policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION)}, setting.Settings.MFAType)
			assert.Nil(t, setting.Settings.MFAType)
		}, retryDuration, tick)

		// add MFAType
		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddMultiFactorToLoginPolicy(IAMCTX, &management.AddMultiFactorToLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event org.policy.login.multifactor.added
			assert.Equal(t, []domain.MultiFactorType{domain.MultiFactorType(policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION)}, setting.Settings.MFAType)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)

		// remove MFAType
		_, err = newInstance.Client.Mgmt.RemoveMultiFactorFromLoginPolicy(IAMCTX, &management.RemoveMultiFactorFromLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event org.policy.login.multifactor.remove
			assert.Equal(t, []domain.MultiFactorType{}, setting.Settings.MFAType)
		}, retryDuration, tick)
	})

	// TODO check this
	t.Run("test added/removed second multifactor reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// create org + login policy
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			secondFactorTypes = setting.Settings.SecondFactorTypes
		}, retryDuration, tick)

		// add new second factor type
		before = time.Now()
		_, err = newInstance.Client.Mgmt.AddSecondFactorToLoginPolicy(IAMCTX, &management.AddSecondFactorToLoginPolicyRequest{
			Type: policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		})
		require.NoError(t, err)
		after := time.Now()

		secondFactorTypes = append(secondFactorTypes, domain.SecondFactorType(policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS))

		// check new second factor type is added
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event org.policy.login.multifactor.secondfactor.added
			assert.Equal(t, secondFactorTypes, setting.Settings.SecondFactorTypes)
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
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event org.policy.login.multifactor.secondfactor.removed
			assert.Equal(t, secondFactorTypes, setting.Settings.SecondFactorTypes)
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
		require.NoError(t, err)
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

		_, err = newInstance.Client.Mgmt.ResetLoginPolicyToDefault(IAMCTX, &management.ResetLoginPolicyToDefaultRequest{
			// Id: organization.Id,
		})
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)

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
		require.NoError(t, err)
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

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLogin(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	// t.Run("test org remove reduces", func(t *testing.T) {
	// 	ctx := t.Context()
	// 	instanceID := Instance.ID()
	// 	orgName := gofakeit.Name()

	// 	// create org
	// 	organization, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
	// 		Name: orgName,
	// 	})
	// 	require.NoError(t, err)

	// 	// delete org
	// 	_, err = OrgClient.DeleteOrganization(CTX, &v2beta_org.DeleteOrganizationRequest{
	// 		Id: organization.Id,
	// 	})
	// 	require.NoError(t, err)

	// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		_, err := settingsRepo.GetLogin(
	// 			ctx,
	// 			instanceID,
	// 			&organization.Id)

	// 		// event org.remove
	// 		require.ErrorIs(t, err, new(database.NoRowFoundError))
	// 	}, retryDuration, tick)
	// })
}

func TestServer_TestOrgLabelSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository(pool)

	t.Run("test adding label settings reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		fmt.Printf("[DEBUGPRINT] [settings_org_test.go:1] organization.Id = %+v\n", organization.Id)

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
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event instance.policy.label.added
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, "#055090", setting.Settings.PrimaryColor)
			assert.Equal(t, "#055090", setting.Settings.BackgroundColor)
			assert.Equal(t, "#055090", setting.Settings.WarnColor)
			assert.Equal(t, "#055090", setting.Settings.FontColor)
			assert.Equal(t, "#055090", setting.Settings.PrimaryColorDark)
			assert.Equal(t, "#055090", setting.Settings.BackgroundColorDark)
			assert.Equal(t, "#055090", setting.Settings.WarnColorDark)
			assert.Equal(t, "#055090", setting.Settings.WarnColorDark)
			assert.Equal(t, "#055090", setting.Settings.FontColorDark)
			assert.Equal(t, false, setting.Settings.HideLoginNameSuffix)
			assert.Equal(t, false, setting.Settings.ErrorMsgPopup)
			assert.Equal(t, false, setting.Settings.DisableWatermark)
			assert.Equal(t, domain.LabelPolicyThemeDark, setting.Settings.ThemeMode)
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
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

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
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event instance.policy.label.change
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, "#055000", setting.Settings.PrimaryColor)
			assert.Equal(t, "#055000", setting.Settings.BackgroundColor)
			assert.Equal(t, "#055000", setting.Settings.WarnColor)
			assert.Equal(t, "#055000", setting.Settings.FontColor)
			assert.Equal(t, "#055000", setting.Settings.PrimaryColorDark)
			assert.Equal(t, "#055000", setting.Settings.BackgroundColorDark)
			assert.Equal(t, "#055000", setting.Settings.WarnColorDark)
			assert.Equal(t, "#055000", setting.Settings.WarnColorDark)
			assert.Equal(t, "#055000", setting.Settings.FontColorDark)
			assert.Equal(t, true, setting.Settings.HideLoginNameSuffix)
			assert.Equal(t, false, setting.Settings.ErrorMsgPopup)
			assert.Equal(t, true, setting.Settings.DisableWatermark)
			assert.Equal(t, domain.LabelPolicyThemeLight, setting.Settings.ThemeMode)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	// TODO
	t.Run("test policy label logo added", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		// organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
		// 	Name: gofakeit.Name(),
		// })
		// require.NoError(t, err)
		// IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// token := newInstance.Users.Get(integration.UserTypeIAMOwner).Token
		// token := integration.SystemToken

		// usr, err := newInstance.Client.Mgmt.AddMachineUser(IAMCTX, &management.AddMachineUserRequest{
		usr, err := newInstance.Client.Mgmt.AddMachineUser(IAMCTX, &management.AddMachineUserRequest{
			UserName:        "service_user",
			Name:            "service_user",
			AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
		})
		require.NoError(t, err)

		_, err = newInstance.Client.Mgmt.SetUserMetadata(IAMCTX, &management.SetUserMetadataRequest{
			Id:    usr.UserId,
			Key:   "key",
			Value: []byte("value"),
		})
		require.NoError(t, err)

		tkn, err := newInstance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{
			UserId:         usr.UserId,
			ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour)),
		})
		require.NoError(t, err)

		token := tkn.Token
		fmt.Printf("[DEBUGPRINT] [settings_org_test.go:1] token = %+v\n", token)

		time.Sleep(time.Second * 5)

		client := resty.New()
		// _, err = client.R().SetAuthToken(token).
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/logo")
		require.NoError(t, err)
		fmt.Printf("[DEBUGPRINT] [settings_org_test.go:1] out = %+v\n", out)
		fmt.Printf("[DEBUGPRINT] [settings_org_test.go:1] err = %+v\n", err)

		// before := time.Now()
		_, err = newInstance.Client.Admin.UpdateLabelPolicy(IAMCTX, &admin.UpdateLabelPolicyRequest{})
		require.NoError(t, err)
		// after := time.Now()

		// retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		// assert.EventuallyWithT(t, func(t *assert.CollectT) {
		// 	setting, err := settingsRepo.GetLabel(
		// 		ctx,
		// 		newInstance.ID(),
		// 		nil)
		// 	require.NoError(t, err)

		// 	// event instance.policy.label.logo.added
		// 	assert.Equal(t, domain.LabelPolicyThemeLight, setting.Settings.LabelPolicyLightLogoURL)
		// 	assert.WithinRange(t, setting.UpdatedAt, before, after)
		// }, retryDuration, tick)
	})

	t.Run("test label settings remove reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove label policy delete org
		_, err = newInstance.Client.Mgmt.ResetLabelPolicyToDefault(IAMCTX, &management.ResetLabelPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check label label settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLabel(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	// // TODO activated

	// t.Run("test label settings logo added reduces", func(t *testing.T) {
	// 	ctx := t.Context()

	// 	newInstance := integration.NewInstance(t.Context())
	// 	IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
	// 	organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
	// 		Name: gofakeit.Name(),
	// 	})
	// 	require.NoError(t, err)
	// 	IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

	// 	_, err = newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
	// 		PrimaryColor:        "#055090",
	// 		HideLoginNameSuffix: false,
	// 		WarnColor:           "#055090",
	// 		BackgroundColor:     "#055090",
	// 		FontColor:           "#055090",
	// 		PrimaryColorDark:    "#055090",
	// 		BackgroundColorDark: "#055090",
	// 		WarnColorDark:       "#055090",
	// 		FontColorDark:       "#055090",
	// 		DisableWatermark:    false,
	// 		ThemeMode:           policy.ThemeMode_THEME_MODE_DARK,
	// 	})
	// 	require.NoError(t, err)

	// 	// remove label policy delete org
	// 	_, err = newInstance.Client.Admin.label(IAMCTX, &management.AddCustomLoginPolicyRequest{

	// 	}
	// 	require.NoError(t, err)
	// })
}

func TestServer_TestOrgPasswordComplexitySettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository(pool)

	t.Run("test password complexity added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddCustomPasswordComplexityPolicy(IAMCTX, &management.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    10,
			HasUppercase: false,
			HasLowercase: false,
			HasNumber:    false,
			HasSymbol:    false,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordComplexity(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event org.policy.password.complexity.added
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, uint64(10), setting.Settings.MinLength)
			assert.Equal(t, false, setting.Settings.HasUppercase)
			assert.Equal(t, false, setting.Settings.HasLowercase)
			assert.Equal(t, false, setting.Settings.HasNumber)
			assert.Equal(t, false, setting.Settings.HasSymbol)
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
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

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
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordComplexity(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event instance.policy.password.complexity.changed
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, uint64(5), setting.Settings.MinLength)
			assert.Equal(t, true, setting.Settings.HasUppercase)
			assert.Equal(t, true, setting.Settings.HasLowercase)
			assert.Equal(t, true, setting.Settings.HasNumber)
			assert.Equal(t, true, setting.Settings.HasSymbol)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
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
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordComplexity(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetPasswordComplexity(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgPasswordPolicySettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository(pool)

	t.Run("test password policy added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		// add add password policy
		before := time.Now()
		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     0,
			ExpireWarnDays: 0,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordExpiry(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event instance.policy.password.age.added
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, uint64(0), setting.Settings.ExpireWarnDays)
			assert.Equal(t, uint64(0), setting.Settings.MaxAgeDays)
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
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     0,
			ExpireWarnDays: 0,
		})
		require.NoError(t, err)

		before := time.Now()
		_, err = newInstance.Client.Mgmt.UpdateCustomPasswordAgePolicy(IAMCTX, &management.UpdateCustomPasswordAgePolicyRequest{
			MaxAgeDays:     30,
			ExpireWarnDays: 30,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordExpiry(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event instance.policy.password.age.changed
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, uint64(30), setting.Settings.ExpireWarnDays)
			assert.Equal(t, uint64(30), setting.Settings.MaxAgeDays)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password poilcy removed reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     0,
			ExpireWarnDays: 0,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordExpiry(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove password policy delete org
		_, err = newInstance.Client.Mgmt.ResetPasswordAgePolicyToDefault(IAMCTX, &management.ResetPasswordAgePolicyToDefaultRequest{})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetPasswordExpiry(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	// t.Run("test delete org reduces", func(t *testing.T) {
	// 	ctx := t.Context()

	// 	newInstance := integration.NewInstance(t.Context())
	// 	IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
	// 	organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
	// 		Name: gofakeit.Name(),
	// 	})
	// 	require.NoError(t, err)
	// 	IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

	// 	_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
	// 		MaxAgeDays:     0,
	// 		ExpireWarnDays: 0,
	// 	})
	// 	require.NoError(t, err)

	// 	// check login settings exist
	// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		setting, err := settingsRepo.GetPasswordExpiry(
	// 			ctx,
	// 			newInstance.ID(),
	// 			&organization.Id)
	// 		require.NoError(t, err)

	// 		require.NotNil(t, setting)
	// 	}, retryDuration, tick)

	// 	// add delete org
	// 	_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
	// 		Id: organization.Id,
	// 	})
	// 	require.NoError(t, err)

	// 	// check password complexity settings removed
	// 	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		_, err := settingsRepo.GetPasswordExpiry(
	// 			ctx,
	// 			newInstance.ID(),
	// 			&organization.Id)

	// 		// event instance.removed
	// 		require.ErrorIs(t, err, new(database.NoRowFoundError))
	// 	}, retryDuration, tick)
	// })

	t.Run("test delete org reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     0,
			ExpireWarnDays: 0,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordExpiry(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetPasswordExpiry(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgLockoutSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository(pool)

	// t.Run("test lockout policy added", func(t *testing.T) {
	// 	ctx := t.Context()

	// 	newInstance := integration.NewInstance(t.Context())
	// 	IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
	// 	organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
	// 		Name: gofakeit.Name(),
	// 	})
	// 	require.NoError(t, err)
	// 	IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

	// 	before := time.Now()
	// 	_, err = newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
	// 		MaxPasswordAttempts: 1,
	// 		MaxOtpAttempts:      1,
	// 	})
	// 	require.NoError(t, err)
	// 	after := time.Now()

	// 	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
	// 	assert.EventuallyWithT(t, func(t *assert.CollectT) {
	// 		setting, err := settingsRepo.GetLockout(
	// 			ctx,
	// 			newInstance.ID(),
	// 			&organization.Id)
	// 		require.NoError(t, err)
	// 		fmt.Printf("[DEBUGPRINT] [settings_org_test.go:1] setting = %+v\n", setting)

	// 		// event instance.policy.lockout.added
	// 		assert.Equal(t, true, setting.Settings.IsDefault)
	// 		assert.Equal(t, uint64(1), setting.Settings.MaxOTPAttempts)
	// 		assert.Equal(t, uint64(1), setting.Settings.MaxPasswordAttempts)
	// 		assert.Equal(t, false, setting.Settings.ShowLockOutFailures)
	// 		assert.WithinRange(t, setting.CreatedAt, before, after)
	// 		assert.WithinRange(t, setting.UpdatedAt, before, after)
	// 	}, retryDuration, tick)
	// })

	t.Run("test lockout policy changed", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

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
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLockout(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event instance.policy.lockout.changed
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, uint64(5), setting.Settings.MaxOTPAttempts)
			assert.Equal(t, uint64(5), setting.Settings.MaxPasswordAttempts)
			assert.WithinRange(t, setting.CreatedAt, before, after)
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
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 1,
			MaxOtpAttempts:      1,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLockout(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove lockout policy org
		_, err = newInstance.Client.Mgmt.ResetLockoutPolicyToDefault(IAMCTX, &management.ResetLockoutPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLockout(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		_, err = newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 1,
			MaxOtpAttempts:      1,
		})
		require.NoError(t, err)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLockout(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLockout(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestOrgDomainSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository(pool)

	t.Run("test domain policy added", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
		// IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

		before := time.Now()
		_, err = newInstance.Client.Admin.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
			OrgId:                                  organization.Id,
			UserLoginMustBeDomain:                  false,
			ValidateOrgDomains:                     false,
			SmtpSenderAddressMatchesInstanceDomain: false,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetDomain(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)
			fmt.Printf("[DEBUGPRINT] [:1] setting = %+v\n", setting)

			// event instance.policy.domain.added
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, false, setting.Settings.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, false, setting.Settings.UserLoginMustBeDomain)
			assert.Equal(t, false, setting.Settings.ValidateOrgDomains)
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
		// IAMCTX = integration.SetOrgID(IAMCTX, organization.Id)

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
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetDomain(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			// event instance.policy.lockout.changed
			assert.Equal(t, false, setting.Settings.IsDefault)
			assert.Equal(t, true, setting.Settings.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, true, setting.Settings.UserLoginMustBeDomain)
			assert.Equal(t, true, setting.Settings.ValidateOrgDomains)
			assert.WithinRange(t, setting.CreatedAt, before, after)
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetDomain(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// remove domain policy org
		_, err = newInstance.Client.Admin.ResetCustomDomainPolicyToDefault(IAMCTX, &admin.ResetCustomDomainPolicyToDefaultRequest{
			OrgId: organization.Id,
		})
		require.NoError(t, err)

		// check domain settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetDomain(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
			Name: gofakeit.Name(),
		})
		require.NoError(t, err)
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
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetDomain(
				ctx,
				newInstance.ID(),
				&organization.Id)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*5)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetDomain(
				ctx,
				newInstance.ID(),
				&organization.Id)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
