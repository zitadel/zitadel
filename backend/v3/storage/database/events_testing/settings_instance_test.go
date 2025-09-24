//go:build integration

package events_test

import (
	"bytes"
	_ "embed"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	durationpb "google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/policy"
)

//go:embed picture.png
var picture []byte

//go:embed font.otf
var font []byte

func TestServer_TestInstanceLoginSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository()

	t.Run("test adding login settings reduces", func(t *testing.T) {
		ctx := t.Context()
		before := time.Now()
		newInstance := integration.NewInstance(t.Context())
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.login.added
			// these values are found in default.yaml
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, true, setting.Settings.AllowRegister)
			assert.Equal(t, true, setting.Settings.AllowExternalSetting)
			assert.Equal(t, domain.PasswordlessTypeAllowed, setting.Settings.PasswordlessType)
			assert.Equal(t, true, setting.Settings.AllowDomainDiscovery)
			assert.Equal(t, true, setting.Settings.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Hour*240), setting.Settings.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Hour*12), setting.Settings.MultiFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Hour*18), setting.Settings.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Hour*240), setting.Settings.ExternalLoginCheckLifetime)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test login settings change reduces", func(t *testing.T) {
		ctx := t.Context()
		before := time.Now()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		_, err := newInstance.Client.Admin.UpdateLoginPolicy(IAMCTX, &admin.UpdateLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           true,
			ForceMfa:                   true,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          true,
			IgnoreUnknownUsernames:     true,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Minute * 20),
			ExternalLoginCheckLifetime: durationpb.New(time.Minute * 21),
			MfaInitSkipLifetime:        durationpb.New(time.Minute * 22),
			SecondFactorCheckLifetime:  durationpb.New(time.Minute * 23),
			MultiFactorCheckLifetime:   durationpb.New(time.Minute * 24),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      true,
			DisableLoginWithPhone:      true,
			ForceMfaLocalOnly:          true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.login.changed
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, false, setting.Settings.AllowRegister)
			assert.Equal(t, true, setting.Settings.AllowExternalSetting)
			assert.Equal(t, true, setting.Settings.ForceMFA)
			assert.Equal(t, domain.PasswordlessTypeNotAllowed, setting.Settings.PasswordlessType)
			assert.Equal(t, true, setting.Settings.HidePasswordReset)
			assert.Equal(t, true, setting.Settings.IgnoreUnknownUsernames)
			assert.Equal(t, "http://www.example.com", setting.Settings.DefaultRedirectURI)
			assert.Equal(t, false, setting.Settings.AllowDomainDiscovery)
			assert.Equal(t, false, setting.Settings.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Minute*20), setting.Settings.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Minute*21), setting.Settings.ExternalLoginCheckLifetime)
			assert.Equal(t, time.Duration(time.Minute*22), setting.Settings.MFAInitSkipLifetime)
			assert.Equal(t, time.Duration(time.Minute*23), setting.Settings.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Minute*24), setting.Settings.MultiFactorCheckLifetime)
			assert.Equal(t, true, setting.Settings.DisableLoginWithEmail)
			assert.Equal(t, true, setting.Settings.DisableLoginWithPhone)
			assert.Equal(t, true, setting.Settings.ForceMFALocalOnly)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test added/remove multifactor type reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		// check inital MFAType value
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			assert.Equal(t, []domain.MultiFactorType{domain.MultiFactorType(policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION)}, setting.Settings.MFAType)
		}, retryDuration, tick)

		// remove MFAType
		_, err := newInstance.Client.Admin.RemoveMultiFactorFromLoginPolicy(IAMCTX, &admin.RemoveMultiFactorFromLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.login.multifactor.remove
			assert.Equal(t, []domain.MultiFactorType{}, setting.Settings.MFAType)
		}, retryDuration, tick)

		before := time.Now()
		_, err = newInstance.Client.Admin.AddMultiFactorToLoginPolicy(IAMCTX, &admin.AddMultiFactorToLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)
		after := time.Now()

		// add MFAType
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.login.multifactor.added
			assert.Equal(t, []domain.MultiFactorType{domain.MultiFactorType(policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION)}, setting.Settings.MFAType)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test added/removed second multifactor reduces", func(t *testing.T) {
		ctx := t.Context()
		before := time.Now()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		// get current second factor types
		var secondFactorTypes []domain.SecondFactorType
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			secondFactorTypes = setting.Settings.SecondFactorTypes
		}, retryDuration, tick)

		// add new second factor type
		before = time.Now()
		_, err := newInstance.Client.Admin.AddSecondFactorToLoginPolicy(IAMCTX, &admin.AddSecondFactorToLoginPolicyRequest{
			Type: policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		})
		require.NoError(t, err)
		after := time.Now()

		secondFactorTypes = append(secondFactorTypes, domain.SecondFactorType(policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS))

		// check new second factor type is added
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.login.multifactor.secondfactor.added
			assert.Equal(t, secondFactorTypes, setting.Settings.SecondFactorTypes)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)

		// remove second factor type
		before = time.Now()
		_, err = newInstance.Client.Admin.RemoveSecondFactorFromLoginPolicy(IAMCTX, &admin.RemoveSecondFactorFromLoginPolicyRequest{
			Type: policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		})
		require.NoError(t, err)
		after = time.Now()

		secondFactorTypes = secondFactorTypes[0 : len(secondFactorTypes)-1]

		// check new second factor type is removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.login.multifactor.secondfactor.removed
			assert.Equal(t, secondFactorTypes, setting.Settings.SecondFactorTypes)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		SystemCTX := integration.WithSystemAuthorization(ctx)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(SystemCTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLogin(
				ctx, pool,
				newInstance.ID(),
				nil)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceLabelSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository()

	t.Run("test adding label settings reduces", func(t *testing.T) {
		ctx := t.Context()
		before := time.Now()
		newInstance := integration.NewInstance(t.Context())
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				newInstance.ID(),
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.added
			// these values are found in default.yaml
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, "#5469d4", setting.Settings.PrimaryColor)
			assert.Equal(t, "#fafafa", setting.Settings.BackgroundColor)
			assert.Equal(t, "#cd3d56", setting.Settings.WarnColor)
			assert.Equal(t, "#000000", setting.Settings.FontColor)
			assert.Equal(t, "#2073c4", setting.Settings.PrimaryColorDark)
			assert.Equal(t, "#111827", setting.Settings.BackgroundColorDark)
			assert.Equal(t, "#ff3b5b", setting.Settings.WarnColorDark)
			assert.Equal(t, "#ff3b5b", setting.Settings.WarnColorDark)
			assert.Equal(t, "#ffffff", setting.Settings.FontColorDark)
			assert.Equal(t, false, setting.Settings.HideLoginNameSuffix)
			assert.Equal(t, false, setting.Settings.ErrorMsgPopup)
			assert.Equal(t, false, setting.Settings.DisableWatermark)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label change", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		_, err := newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
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

		before := time.Now()
		_, err = newInstance.Client.Admin.UpdateLabelPolicy(IAMCTX, &admin.UpdateLabelPolicyRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				newInstance.ID(),
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.change
			assert.Equal(t, true, setting.Settings.IsDefault)
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
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test label settings activated", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		// activate label
		before := time.Now()
		_, err := newInstance.Client.Admin.ActivateLabelPolicy(IAMCTX, &admin.ActivateLabelPolicyRequest{})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				newInstance.ID(),
				nil,
				domain.LabelStateActivated)

			require.NoError(t, err)

			// event instance.policy.label.activated
			assert.Equal(t, domain.LabelStateActivated, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo light added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set logo light
		before := time.Now()
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/logo")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.logo.added
			assert.NotNil(t, setting.Settings.LabelPolicyLightLogoURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set logo dark
		before := time.Now()
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/logo/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.logo.dark.added
			assert.NotNil(t, setting.Settings.LabelPolicyDarkLogoURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo light removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set logo light
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/logo")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check light logo set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			assert.NotNil(t, setting.Settings.LabelPolicyLightLogoURL)
		}, retryDuration, tick)

		// remote logo light
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			Delete("http://localhost:8080" + "/admin/v1" + "/policies/label/logo")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		// check light logo removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.logo.removed
			assert.Nil(t, setting.Settings.LabelPolicyLightLogoURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set logo dark
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/logo/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check dark logo set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			assert.NotNil(t, setting.Settings.LabelPolicyDarkLogoURL)
		}, retryDuration, tick)

		// remote logo dark
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			Delete("http://localhost:8080" + "/admin/v1" + "/policies/label/logo_dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		// check dark logo removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.logo.dark.removed
			assert.Nil(t, setting.Settings.LabelPolicyDarkLogoURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set icon light
		before := time.Now()
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/icon")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.icon.added
			assert.NotNil(t, setting.Settings.LabelPolicyLightIconURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set icon dark
		before := time.Now()
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/icon/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.icon.dark.added
			assert.NotNil(t, setting.Settings.LabelPolicyDarkIconURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set icon light
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/icon")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check light icon set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			assert.NotNil(t, setting.Settings.LabelPolicyLightIconURL)
		}, retryDuration, tick)

		// remote icon light
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			Delete("http://localhost:8080" + "/admin/v1" + "/policies/label/icon")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		// check light icon removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.icon.removed
			assert.Nil(t, setting.Settings.LabelPolicyLightIconURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set icon dark
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(picture)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/icon/dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		// check dark icon set
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			assert.NotNil(t, setting.Settings.LabelPolicyDarkIconURL)
		}, retryDuration, tick)

		// remote icon dark
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			Delete("http://localhost:8080" + "/admin/v1" + "/policies/label/icon_dark")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		// check dark icon removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.icon.dark.removed
			assert.Nil(t, setting.Settings.LabelPolicyDarkIconURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy font added", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set logo light
		before := time.Now()
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(font)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/font")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.font.added
			assert.NotNil(t, setting.Settings.LabelPolicyFontURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy font removed", func(t *testing.T) {
		ctx := t.Context()
		token := integration.SystemToken

		instanceRepo := repository.InstanceRepository()
		instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.NameCondition(database.TextOperationEqual, "ZITADEL")))
		instanceID := instance.ID
		require.NoError(t, err)

		// set logo light
		client := resty.New()
		out, err := client.R().SetAuthToken(token).
			SetMultipartField("file", "filename", "image/png", bytes.NewReader(font)).
			Post("http://localhost:8080" + "/assets/v1" + "/instance/policy/label/font")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			assert.NotNil(t, setting.Settings.LabelPolicyFontURL)
		}, retryDuration, tick)

		// remote font policy
		before := time.Now()
		client = resty.New()
		out, err = client.R().SetAuthToken(token).
			Delete("http://localhost:8080" + "/admin/v1" + "/policies/label/font")
		require.NoError(t, err)
		require.Equal(t, 200, out.StatusCode())
		after := time.Now()

		// check font policy removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				instanceID,
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			// event instance.policy.label.font.removed
			assert.Nil(t, setting.Settings.LabelPolicyFontURL)
			assert.Equal(t, domain.LabelStatePreview, *setting.LabelState)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		SystemCTX := integration.WithSystemAuthorization(ctx)

		// check label preview settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				newInstance.ID(),
				nil,
				domain.LabelStatePreview)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// check label activated settings exist
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLabel(
				ctx, pool,
				newInstance.ID(),
				nil,
				domain.LabelStateActivated)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(SystemCTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})

		// check label preview settings deleted
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLabel(
				ctx, pool,
				newInstance.ID(),
				nil,
				domain.LabelStatePreview)

			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)

		// check label activated settings deleted
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLabel(
				ctx, pool,
				newInstance.ID(),
				nil,
				domain.LabelStateActivated)

			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
		require.NoError(t, err)
	})
}

func TestServer_TestPasswordComplexitySettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository()

	t.Run("test password complexity added", func(t *testing.T) {
		ctx := t.Context()

		before := time.Now()
		newInstance := integration.NewInstance(t.Context())
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordComplexity(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.password.complexity.added
			// these values are found in default.yaml
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, uint64(8), setting.Settings.MinLength)
			assert.Equal(t, true, setting.Settings.HasUppercase)
			assert.Equal(t, true, setting.Settings.HasLowercase)
			assert.Equal(t, true, setting.Settings.HasNumber)
			assert.Equal(t, true, setting.Settings.HasSymbol)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password complexity change", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		before := time.Now()
		_, err := newInstance.Client.Admin.UpdatePasswordComplexityPolicy(IAMCTX, &admin.UpdatePasswordComplexityPolicyRequest{
			MinLength:    5,
			HasUppercase: true,
			HasLowercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordComplexity(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.password.complexity.changed
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, uint64(5), setting.Settings.MinLength)
			assert.Equal(t, true, setting.Settings.HasUppercase)
			assert.Equal(t, true, setting.Settings.HasLowercase)
			assert.Equal(t, true, setting.Settings.HasNumber)
			assert.Equal(t, true, setting.Settings.HasSymbol)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		SystemCTX := integration.WithSystemAuthorization(ctx)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordComplexity(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(SystemCTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetPasswordComplexity(
				ctx, pool,
				newInstance.ID(),
				nil)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstancePasswordPolicySettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository()

	t.Run("test password policy added", func(t *testing.T) {
		ctx := t.Context()

		before := time.Now()
		newInstance := integration.NewInstance(t.Context())
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordExpiry(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.password.age.added
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, uint64(0), setting.Settings.ExpireWarnDays)
			assert.Equal(t, uint64(0), setting.Settings.MaxAgeDays)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy changed", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		before := time.Now()
		_, err := newInstance.Client.Admin.UpdatePasswordAgePolicy(IAMCTX, &admin.UpdatePasswordAgePolicyRequest{
			MaxAgeDays:     30,
			ExpireWarnDays: 30,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordExpiry(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.password.age.changed
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, uint64(30), setting.Settings.ExpireWarnDays)
			assert.Equal(t, uint64(30), setting.Settings.MaxAgeDays)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		SystemCTX := integration.WithSystemAuthorization(ctx)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetPasswordExpiry(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(SystemCTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetPasswordExpiry(
				ctx, pool,
				newInstance.ID(),
				nil)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestDomainSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository()

	t.Run("test domain policy added", func(t *testing.T) {
		ctx := t.Context()

		before := time.Now()
		newInstance := integration.NewInstance(t.Context())
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetDomain(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.domain.added
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, false, setting.Settings.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, false, setting.Settings.UserLoginMustBeDomain)
			assert.Equal(t, false, setting.Settings.ValidateOrgDomains)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test domain policy changed", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		before := time.Now()
		_, err := newInstance.Client.Admin.UpdateDomainPolicy(IAMCTX, &admin.UpdateDomainPolicyRequest{
			UserLoginMustBeDomain:                  true,
			ValidateOrgDomains:                     true,
			SmtpSenderAddressMatchesInstanceDomain: true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetDomain(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.changed
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, true, setting.Settings.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, true, setting.Settings.UserLoginMustBeDomain)
			assert.Equal(t, true, setting.Settings.ValidateOrgDomains)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		SystemCTX := integration.WithSystemAuthorization(ctx)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetDomain(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(SystemCTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check domain settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetDomain(
				ctx, pool,
				newInstance.ID(),
				nil)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestLockoutSettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository()

	t.Run("test lockout policy added", func(t *testing.T) {
		ctx := t.Context()

		before := time.Now()
		newInstance := integration.NewInstance(t.Context())
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLockout(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.lockout.added
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, uint64(0), setting.Settings.MaxOTPAttempts)
			assert.Equal(t, uint64(0), setting.Settings.MaxPasswordAttempts)
			assert.Equal(t, true, setting.Settings.ShowLockOutFailures)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy changed", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		before := time.Now()
		_, err := newInstance.Client.Admin.UpdateLockoutPolicy(IAMCTX, &admin.UpdateLockoutPolicyRequest{
			MaxPasswordAttempts: 5,
			MaxOtpAttempts:      5,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLockout(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.lockout.changed
			assert.Equal(t, true, setting.Settings.IsDefault)
			assert.Equal(t, uint64(5), setting.Settings.MaxOTPAttempts)
			assert.Equal(t, uint64(5), setting.Settings.MaxPasswordAttempts)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		SystemCTX := integration.WithSystemAuthorization(ctx)

		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetLockout(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(SystemCTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetLockout(
				ctx, pool,
				newInstance.ID(),
				nil)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestSecuritySettingsReduces(t *testing.T) {
	settingsRepo := repository.SettingsRepository()

	t.Run("test security policy set", func(t *testing.T) {
		ctx := t.Context()

		newInstance := integration.NewInstance(t.Context())
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		// 1. set security policy
		before := time.Now()
		_, err := newInstance.Client.Admin.SetSecurityPolicy(IAMCTX, &admin.SetSecurityPolicyRequest{
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"value"},
			EnableImpersonation:   true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetSecurity(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.security.set
			assert.Equal(t, true, setting.Settings.EnableIframeEmbedding)
			assert.Equal(t, []string{"value"}, setting.Settings.AllowedOrigins)
			assert.Equal(t, true, setting.Settings.EnableImpersonation)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)

		// 2. re-set security policy
		before = time.Now()
		_, err = newInstance.Client.Admin.SetSecurityPolicy(IAMCTX, &admin.SetSecurityPolicyRequest{
			EnableIframeEmbedding: false,
			AllowedOrigins:        []string{"no_value"},
			EnableImpersonation:   false,
		})
		require.NoError(t, err)
		after = time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetSecurity(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			// event instance.policy.security.set
			assert.Equal(t, false, setting.Settings.EnableIframeEmbedding)
			assert.Equal(t, []string{"no_value"}, setting.Settings.AllowedOrigins)
			assert.Equal(t, false, setting.Settings.EnableImpersonation)
			assert.WithinRange(t, *setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		ctx := t.Context()
		newInstance := integration.NewInstance(t.Context())

		SystemCTX := integration.WithSystemAuthorization(ctx)
		IAMCTX := newInstance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		// 1. set security policy
		_, err := newInstance.Client.Admin.SetSecurityPolicy(IAMCTX, &admin.SetSecurityPolicyRequest{
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"value"},
			EnableImpersonation:   true,
		})
		require.NoError(t, err)

		// 2. check security instance exists
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.GetSecurity(
				ctx, pool,
				newInstance.ID(),
				nil)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// 3. delete instance
		_, err = newInstance.Client.InstanceV2Beta.DeleteInstance(SystemCTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// 4. check security settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Second*10)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.GetSecurity(
				ctx, pool,
				newInstance.ID(),
				nil)

			// event instance.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
