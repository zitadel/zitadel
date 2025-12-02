//go:build integration

package events_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/policy"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func createInstanceWithOrg(t *testing.T) (context.Context, *integration.Instance, string) {
	newInstance := integration.NewInstance(t.Context())

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	organization, err := newInstance.Client.OrgV2beta.CreateOrganization(IAMCTX, &v2beta_org.CreateOrganizationRequest{
		Name: integration.OrganizationName(),
	})
	require.NoError(t, err)
	orgId := organization.GetId()
	IAMCTX = integration.SetOrgID(IAMCTX, orgId)

	cleanupInstance(t, newInstance)
	return IAMCTX, newInstance, orgId
}

func cleanupInstance(t *testing.T, instance *integration.Instance) {
	t.Cleanup(func() {
		_, err := SystemClient.RemoveInstance(CTX, &system.RemoveInstanceRequest{
			InstanceId: instance.ID(),
		})
		if err != nil {
			t.Logf("Failed to delete instance on cleanup: %v", err)
		}
	})
}

func uploadOrganizationAsset(ctx context.Context, t *testing.T, instance *integration.Instance, path string, asset io.Reader) {
	url := http_util.BuildOrigin(instance.Host(), instance.Config.Secure) + "/assets/v1" + path

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "filename")
	require.NoError(t, err)

	_, err = io.Copy(part, asset)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	require.NoError(t, err)
	md, ok := metadata.FromOutgoingContext(ctx)
	require.True(t, ok)
	// context information has to be HTTP headers, as the asset API is only HTTP
	req.Header.Set("x-zitadel-orgid", md.Get("x-zitadel-orgid")[0])
	req.Header.Set("authorization", md.Get("authorization")[0])
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServer_TestOrgLoginSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.LoginSettingsRepository()
	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test adding login settings reduces", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddCustomLoginPolicy(IAMCTX, &management.AddCustomLoginPolicyRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.added
			// these values are found in default.yaml
			assert.Equal(t, false, *setting.AllowRegister)
			assert.Equal(t, false, *setting.AllowExternalIDP)
			assert.Equal(t, domain.PasswordlessTypeNotAllowed, *setting.PasswordlessType)
			assert.Equal(t, false, *setting.AllowDomainDiscovery)
			assert.Equal(t, false, *setting.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Second*10), *setting.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), *setting.ExternalLoginCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), *setting.MFAInitSkipLifetime)
			assert.Equal(t, time.Duration(time.Second*10), *setting.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*10), *setting.MultiFactorCheckLifetime)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test change login settings reduces", func(t *testing.T) {
		// update org policy
		before := time.Now()
		_, err := newInstance.Client.Mgmt.UpdateCustomLoginPolicy(IAMCTX, &management.UpdateCustomLoginPolicyRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.login.changed
			assert.Equal(t, true, *setting.AllowRegister)
			assert.Equal(t, true, *setting.AllowExternalIDP)
			assert.Equal(t, true, *setting.ForceMFA)
			assert.Equal(t, domain.PasswordlessTypeAllowed, *setting.PasswordlessType)
			assert.Equal(t, true, *setting.HidePasswordReset)
			assert.Equal(t, true, *setting.IgnoreUnknownUsernames)
			assert.Equal(t, "http://www.new_example.com", *setting.DefaultRedirectURI)
			assert.Equal(t, true, *setting.AllowDomainDiscovery)
			assert.Equal(t, true, *setting.AllowUserNamePassword)
			assert.Equal(t, time.Duration(time.Second*5*20), *setting.PasswordCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*21), *setting.ExternalLoginCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*22), *setting.MFAInitSkipLifetime)
			assert.Equal(t, time.Duration(time.Second*5*23), *setting.SecondFactorCheckLifetime)
			assert.Equal(t, time.Duration(time.Second*5*24), *setting.MultiFactorCheckLifetime)
			assert.Equal(t, true, *setting.DisableLoginWithEmail)
			assert.Equal(t, true, *setting.DisableLoginWithPhone)
			assert.Equal(t, true, *setting.ForceMFALocalOnly)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test added/remove multifactor type reduces", func(t *testing.T) {
		// check initial MFAType value
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Nil(t, setting.MFAType)
		}, retryDuration, tick)

		// add MFAType
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddMultiFactorToLoginPolicy(IAMCTX, &management.AddMultiFactorToLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
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

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
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
		before := time.Now()
		// get current second factor types
		var secondFactorTypes []domain.SecondFactorType
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			secondFactorTypes = setting.SecondFactorTypes
		}, retryDuration, tick)

		// add new second factor type
		before = time.Now()
		_, err := newInstance.Client.Mgmt.AddSecondFactorToLoginPolicy(IAMCTX, &management.AddSecondFactorToLoginPolicyRequest{
			Type: policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		})
		require.NoError(t, err)
		after := time.Now()

		secondFactorTypes = append(secondFactorTypes, domain.SecondFactorType(policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS))

		// check new second factor type is added
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
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
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
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
		_, err := newInstance.Client.Mgmt.ResetLoginPolicyToDefault(IAMCTX, &management.ResetLoginPolicyToDefaultRequest{})
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)

			// event org.policy.login.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		// add delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
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
	t.Parallel()

	settingsRepo := repository.BrandingSettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test adding label settings reduces", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddCustomLabelPolicy(IAMCTX, &management.AddCustomLabelPolicyRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.added
			assert.Equal(t, "#055090", *setting.PrimaryColorLight)
			assert.Equal(t, "#055090", *setting.BackgroundColorLight)
			assert.Equal(t, "#055090", *setting.WarnColorLight)
			assert.Equal(t, "#055090", *setting.FontColorLight)
			assert.Equal(t, "#055090", *setting.PrimaryColorDark)
			assert.Equal(t, "#055090", *setting.BackgroundColorDark)
			assert.Equal(t, "#055090", *setting.WarnColorDark)
			assert.Equal(t, "#055090", *setting.WarnColorDark)
			assert.Equal(t, "#055090", *setting.FontColorDark)
			assert.Equal(t, false, *setting.HideLoginNameSuffix)
			assert.Equal(t, false, *setting.ErrorMsgPopup)
			assert.Equal(t, false, *setting.DisableWatermark)
			assert.Equal(t, domain.BrandingPolicyThemeDark, *setting.ThemeMode)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label change", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.UpdateCustomLabelPolicy(IAMCTX, &management.UpdateCustomLabelPolicyRequest{
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

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.change
			assert.Equal(t, "#055000", *setting.PrimaryColorLight)
			assert.Equal(t, "#055000", *setting.BackgroundColorLight)
			assert.Equal(t, "#055000", *setting.WarnColorLight)
			assert.Equal(t, "#055000", *setting.FontColorLight)
			assert.Equal(t, "#055000", *setting.PrimaryColorDark)
			assert.Equal(t, "#055000", *setting.BackgroundColorDark)
			assert.Equal(t, "#055000", *setting.WarnColorDark)
			assert.Equal(t, "#055000", *setting.WarnColorDark)
			assert.Equal(t, "#055000", *setting.FontColorDark)
			assert.Equal(t, true, *setting.HideLoginNameSuffix)
			assert.Equal(t, false, *setting.ErrorMsgPopup)
			assert.Equal(t, true, *setting.DisableWatermark)
			assert.Equal(t, domain.BrandingPolicyThemeLight, *setting.ThemeMode)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test label settings activated", func(t *testing.T) {
		// activate label
		before := time.Now()
		_, err := newInstance.Client.Mgmt.ActivateCustomLabelPolicy(IAMCTX, &management.ActivateCustomLabelPolicyRequest{})
		require.NoError(t, err)
		after := time.Now()

		_, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.activated
			assert.Equal(t, domain.SettingStateActive, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, time.Minute*5, tick)
	})

	t.Run("test policy label logo light added", func(t *testing.T) {
		// set logo light
		before := time.Now()
		uploadOrganizationAsset(IAMCTX, t, newInstance, "/org/policy/label/logo", bytes.NewReader(picture))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.added
			assert.NotNil(t, setting.LogoURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo light removed", func(t *testing.T) {
		// remote logo light
		before := time.Now()
		_, err := newInstance.Client.Mgmt.RemoveCustomLabelPolicyLogo(IAMCTX, &management.RemoveCustomLabelPolicyLogoRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check light logo removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.removed
			assert.Equal(t, url.URL{}, *setting.LogoURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark added", func(t *testing.T) {
		// set logo dark
		before := time.Now()
		uploadOrganizationAsset(IAMCTX, t, newInstance, "/org/policy/label/logo/dark", bytes.NewReader(picture))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.dark.added
			assert.NotNil(t, setting.LogoURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark removed", func(t *testing.T) {
		// remove logo dark
		before := time.Now()
		_, err := newInstance.Client.Mgmt.RemoveCustomLabelPolicyLogoDark(IAMCTX, &management.RemoveCustomLabelPolicyLogoDarkRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check dark logo removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.logo.dark.removed
			assert.Equal(t, url.URL{}, *setting.LogoURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light added", func(t *testing.T) {
		// set icon light
		before := time.Now()
		uploadOrganizationAsset(IAMCTX, t, newInstance, "/org/policy/label/icon", bytes.NewReader(picture))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.added
			assert.NotNil(t, setting.IconURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light removed", func(t *testing.T) {
		// remote icon light
		before := time.Now()
		_, err := newInstance.Client.Mgmt.RemoveCustomLabelPolicyIcon(IAMCTX, &management.RemoveCustomLabelPolicyIconRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check light icon removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.removed
			assert.Equal(t, url.URL{}, *setting.IconURLLight)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark added", func(t *testing.T) {
		// set icon dark
		before := time.Now()
		uploadOrganizationAsset(IAMCTX, t, newInstance, "/org/policy/label/icon/dark", bytes.NewReader(picture))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.dark.added
			assert.NotNil(t, setting.IconURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark removed", func(t *testing.T) {
		// remove icon dark
		before := time.Now()
		_, err := newInstance.Client.Mgmt.RemoveCustomLabelPolicyIconDark(IAMCTX, &management.RemoveCustomLabelPolicyIconDarkRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check dark icon removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.icon.dark.removed
			assert.Equal(t, url.URL{}, *setting.IconURLDark)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy font added", func(t *testing.T) {
		// set font
		before := time.Now()
		uploadOrganizationAsset(IAMCTX, t, newInstance, "/org/policy/label/font", bytes.NewReader(font))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
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
		// remove font
		before := time.Now()
		_, err := newInstance.Client.Mgmt.RemoveCustomLabelPolicyFont(IAMCTX, &management.RemoveCustomLabelPolicyFontRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check font policy removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(t, err)

			// event org.policy.label.font.removed
			assert.Equal(t, url.URL{}, *setting.FontURL)
			assert.Equal(t, domain.SettingStatePreview, setting.State)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test label settings remove reduces", func(t *testing.T) {
		// remove label policy
		_, err := newInstance.Client.Mgmt.ResetLabelPolicyToDefault(IAMCTX, &management.ResetLabelPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check label settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)

			// event org.policy.label.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test org delete reduces", func(t *testing.T) {
		// delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check label settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
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
	t.Parallel()

	settingsRepo := repository.PasswordComplexitySettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test password complexity added", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddCustomPasswordComplexityPolicy(IAMCTX, &management.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    10,
			HasUppercase: true,
			HasLowercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.complexity.added
			assert.Equal(t, uint64(10), *setting.MinLength)
			assert.Equal(t, true, *setting.HasUppercase)
			assert.Equal(t, true, *setting.HasLowercase)
			assert.Equal(t, true, *setting.HasNumber)
			assert.Equal(t, true, *setting.HasSymbol)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password complexity change", func(t *testing.T) {
		// update password compexity
		before := time.Now()
		_, err := newInstance.Client.Mgmt.UpdateCustomPasswordComplexityPolicy(IAMCTX, &management.UpdateCustomPasswordComplexityPolicyRequest{
			MinLength:    5,
			HasUppercase: true,
			HasLowercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.complexity.changed
			assert.Equal(t, uint64(5), *setting.MinLength)
			assert.Equal(t, true, *setting.HasUppercase)
			assert.Equal(t, true, *setting.HasLowercase)
			assert.Equal(t, true, *setting.HasNumber)
			assert.Equal(t, true, *setting.HasSymbol)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password complexity removed", func(t *testing.T) {
		// delete password complexity policy
		_, err := newInstance.Client.Mgmt.ResetPasswordComplexityPolicyToDefault(IAMCTX, &management.ResetPasswordComplexityPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		// create password complexity
		_, err := newInstance.Client.Mgmt.AddCustomPasswordComplexityPolicy(IAMCTX, &management.AddCustomPasswordComplexityPolicyRequest{
			MinLength:    10,
			HasUppercase: false,
			HasLowercase: false,
			HasNumber:    false,
			HasSymbol:    false,
		})
		require.NoError(t, err)

		// check password complexity settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// add delete org
		_, err = newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
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
	t.Parallel()

	settingsRepo := repository.PasswordExpirySettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test password policy added", func(t *testing.T) {
		// add password policy
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddCustomPasswordAgePolicy(IAMCTX, &management.AddCustomPasswordAgePolicyRequest{
			MaxAgeDays:     10,
			ExpireWarnDays: 10,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.age.added
			assert.Equal(t, uint64(10), *setting.ExpireWarnDays)
			assert.Equal(t, uint64(10), *setting.MaxAgeDays)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy change reduces", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.UpdateCustomPasswordAgePolicy(IAMCTX, &management.UpdateCustomPasswordAgePolicyRequest{
			MaxAgeDays:     40,
			ExpireWarnDays: 40,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.password.age.changed
			assert.Equal(t, uint64(40), *setting.ExpireWarnDays)
			assert.Equal(t, uint64(40), *setting.MaxAgeDays)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy removed reduces", func(t *testing.T) {
		// remove password policy delete org
		_, err := newInstance.Client.Mgmt.ResetPasswordAgePolicyToDefault(IAMCTX, &management.ResetPasswordAgePolicyToDefaultRequest{})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)

			// event org.policy.password.age.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		// add delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
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
	t.Parallel()

	settingsRepo := repository.LockoutSettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test lockout policy added", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddCustomLockoutPolicy(IAMCTX, &management.AddCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 1,
			MaxOtpAttempts:      1,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.lockout.added
			assert.Equal(t, uint64(1), *setting.MaxOTPAttempts)
			assert.Equal(t, uint64(1), *setting.MaxPasswordAttempts)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test lockout policy changed", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.UpdateCustomLockoutPolicy(IAMCTX, &management.UpdateCustomLockoutPolicyRequest{
			MaxPasswordAttempts: 5,
			MaxOtpAttempts:      5,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.lockout.changed
			assert.Equal(t, uint64(5), *setting.MaxOTPAttempts)
			assert.Equal(t, uint64(5), *setting.MaxPasswordAttempts)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test remove lockout policy reduces", func(t *testing.T) {
		// remove lockout policy org
		_, err := newInstance.Client.Mgmt.ResetLockoutPolicyToDefault(IAMCTX, &management.ResetLockoutPolicyToDefaultRequest{})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		// add delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
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
	t.Parallel()

	settingsRepo := repository.DomainSettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test domain policy added", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.AddCustomDomainPolicy(IAMCTX, &admin.AddCustomDomainPolicyRequest{
			OrgId:                                  orgId,
			UserLoginMustBeDomain:                  false,
			ValidateOrgDomains:                     false,
			SmtpSenderAddressMatchesInstanceDomain: false,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.domain.added
			assert.Equal(t, false, *setting.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, false, *setting.LoginNameIncludesDomain)
			assert.Equal(t, false, *setting.RequireOrgDomainVerification)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test domain policy changed", func(t *testing.T) {
		// update domain policy
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdateCustomDomainPolicy(IAMCTX, &admin.UpdateCustomDomainPolicyRequest{
			OrgId:                                  orgId,
			UserLoginMustBeDomain:                  true,
			ValidateOrgDomains:                     true,
			SmtpSenderAddressMatchesInstanceDomain: true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event org.policy.domain.changed
			assert.Equal(t, true, *setting.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(t, true, *setting.LoginNameIncludesDomain)
			assert.Equal(t, true, *setting.RequireOrgDomainVerification)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test remove domain policy reduces", func(t *testing.T) {
		// remove domain policy org
		_, err := newInstance.Client.Admin.ResetCustomDomainPolicyToDefault(IAMCTX, &admin.ResetCustomDomainPolicyToDefaultRequest{
			OrgId: orgId,
		})
		require.NoError(t, err)

		// check domain settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)

			// event org.policy.domain.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test delete org reduces", func(t *testing.T) {
		// add delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
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
	t.Parallel()

	settingsRepo := repository.OrganizationSettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test add org settings set", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.SettingsV2beta.SetOrganizationSettings(IAMCTX, &settings.SetOrganizationSettingsRequest{
			OrganizationId:              orgId,
			OrganizationScopedUsernames: gu.Ptr(true),
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event settings.organization.set
			assert.Equal(t, true, *setting.OrganizationScopedUsernames)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test add org settings re-set", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.SettingsV2beta.SetOrganizationSettings(IAMCTX, &settings.SetOrganizationSettingsRequest{
			OrganizationId:              orgId,
			OrganizationScopedUsernames: gu.Ptr(false),
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			// event settings.organization.set
			assert.Equal(t, false, *setting.OrganizationScopedUsernames)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test organization settings removed reduces", func(t *testing.T) {
		// delete organization setting
		_, err := newInstance.Client.SettingsV2beta.DeleteOrganizationSettings(IAMCTX, &settings.DeleteOrganizationSettingsRequest{
			OrganizationId: orgId,
		})
		require.NoError(t, err)

		// check organization settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)

			// event settings.organization.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test organization removed reduces", func(t *testing.T) {
		// add delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check organization settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeOrganization, domain.SettingStateActive),
				),
			)

			// event org.removed
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestNotificationSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.NotificationSettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test add notification settings set", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddCustomNotificationPolicy(IAMCTX, &management.AddCustomNotificationPolicyRequest{
			PasswordChange: true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeNotification, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, true, *setting.PasswordChange)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test add notification settings re-set", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.UpdateCustomNotificationPolicy(IAMCTX, &management.UpdateCustomNotificationPolicyRequest{
			PasswordChange: false,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeNotification, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, false, *setting.PasswordChange)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test notification settings removed reduces", func(t *testing.T) {
		_, err := newInstance.Client.Mgmt.ResetNotificationPolicyToDefault(IAMCTX, &management.ResetNotificationPolicyToDefaultRequest{})
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeNotification, domain.SettingStateActive),
				),
			)

			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test organization removed reduces", func(t *testing.T) {
		// add delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check organization settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeNotification, domain.SettingStateActive),
				),
			)

			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestLegalAndSupportSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.LegalAndSupportSettingsRepository()

	IAMCTX, newInstance, orgId := createInstanceWithOrg(t)

	t.Run("test add legal and support settings set", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.AddCustomPrivacyPolicy(IAMCTX, &management.AddCustomPrivacyPolicyRequest{
			TosLink:        "https://tos.example.com",
			PrivacyLink:    "https://privacy.example.com",
			HelpLink:       "https://help.example.com",
			SupportEmail:   "support@example.com",
			DocsLink:       "https://docs.example.com",
			CustomLink:     "https://custom.example.com",
			CustomLinkText: "Custom link text",
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, "https://tos.example.com", *setting.TOSLink)
			assert.Equal(t, "https://privacy.example.com", *setting.PrivacyPolicyLink)
			assert.Equal(t, "https://help.example.com", *setting.HelpLink)
			assert.Equal(t, "support@example.com", *setting.SupportEmail)
			assert.Equal(t, "https://docs.example.com", *setting.DocsLink)
			assert.Equal(t, "https://custom.example.com", *setting.CustomLink)
			assert.Equal(t, "Custom link text", *setting.CustomLinkText)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test add legal and support settings re-set", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Mgmt.UpdateCustomPrivacyPolicy(IAMCTX, &management.UpdateCustomPrivacyPolicyRequest{
			TosLink:        "https://tos.example2.com",
			PrivacyLink:    "https://privacy.example2.com",
			HelpLink:       "https://help.example2.com",
			SupportEmail:   "support@example2.com",
			DocsLink:       "https://docs.example2.com",
			CustomLink:     "https://custom.example2.com",
			CustomLinkText: "Custom link text2",
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, "https://tos.example2.com", *setting.TOSLink)
			assert.Equal(t, "https://privacy.example2.com", *setting.PrivacyPolicyLink)
			assert.Equal(t, "https://help.example2.com", *setting.HelpLink)
			assert.Equal(t, "support@example2.com", *setting.SupportEmail)
			assert.Equal(t, "https://docs.example2.com", *setting.DocsLink)
			assert.Equal(t, "https://custom.example2.com", *setting.CustomLink)
			assert.Equal(t, "Custom link text2", *setting.CustomLinkText)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test legal and support settings removed reduces", func(t *testing.T) {
		_, err := newInstance.Client.Mgmt.ResetPrivacyPolicyToDefault(IAMCTX, &management.ResetPrivacyPolicyToDefaultRequest{})
		require.NoError(t, err)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
				),
			)

			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test organization removed reduces", func(t *testing.T) {
		// add delete org
		_, err := newInstance.Client.OrgV2beta.DeleteOrganization(IAMCTX, &v2beta_org.DeleteOrganizationRequest{
			Id: orgId,
		})
		require.NoError(t, err)

		// check organization settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), &orgId, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
				),
			)

			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
