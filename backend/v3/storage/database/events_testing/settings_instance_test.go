//go:build integration

package events_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"

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
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/policy"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

//go:embed picture.png
var picture []byte

//go:embed font.otf
var font []byte

func uploadInstanceAsset(ctx context.Context, t *testing.T, instance *integration.Instance, path string, asset io.Reader) {
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
	req.Header.Set("authorization", md.Get("authorization")[0])
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServer_TestInstanceLoginSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.LoginSettingsRepository()
	before := time.Now()
	newInstance := integration.NewInstance(t.Context())
	after := time.Now()

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test adding login settings reduces", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.login.added
			// these values are found in default.yaml
			assert.Equal(collect, true, *setting.AllowRegister)
			assert.Equal(collect, true, *setting.AllowExternalIDP)
			assert.Equal(collect, domain.PasswordlessTypeAllowed, *setting.PasswordlessType)
			assert.Equal(collect, true, *setting.AllowDomainDiscovery)
			assert.Equal(collect, true, *setting.AllowUserNamePassword)
			assert.Equal(collect, time.Duration(time.Hour*240), *setting.PasswordCheckLifetime)
			assert.Equal(collect, time.Duration(time.Hour*12), *setting.MultiFactorCheckLifetime)
			assert.Equal(collect, time.Duration(time.Hour*18), *setting.SecondFactorCheckLifetime)
			assert.Equal(collect, time.Duration(time.Hour*240), *setting.ExternalLoginCheckLifetime)
			assert.WithinRange(collect, setting.CreatedAt, before, after)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test login settings change reduces", func(t *testing.T) {
		_, err := newInstance.Client.Admin.UpdateLoginPolicy(IAMCTX, &admin.UpdateLoginPolicyRequest{
			AllowUsernamePassword:      false,
			AllowRegister:              false,
			AllowExternalIdp:           true,
			ForceMfa:                   true,
			PasswordlessType:           policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED,
			HidePasswordReset:          true,
			IgnoreUnknownUsernames:     true,
			DefaultRedirectUri:         "http://www.example.com",
			PasswordCheckLifetime:      durationpb.New(time.Second * 20 * 20),
			ExternalLoginCheckLifetime: durationpb.New(time.Second * 20 * 21),
			MfaInitSkipLifetime:        durationpb.New(time.Second * 20 * 22),
			SecondFactorCheckLifetime:  durationpb.New(time.Second * 20 * 23),
			MultiFactorCheckLifetime:   durationpb.New(time.Second * 20 * 24),
			AllowDomainDiscovery:       false,
			DisableLoginWithEmail:      true,
			DisableLoginWithPhone:      true,
			ForceMfaLocalOnly:          true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(collect, setting.ForceMFA)

			// event instance.policy.login.changed
			assert.Equal(collect, false, *setting.AllowRegister)
			assert.Equal(collect, true, *setting.AllowExternalIDP)
			assert.Equal(collect, true, *setting.ForceMFA)
			assert.Equal(collect, domain.PasswordlessTypeNotAllowed, *setting.PasswordlessType)
			assert.Equal(collect, true, *setting.HidePasswordReset)
			assert.Equal(collect, true, *setting.IgnoreUnknownUsernames)
			assert.Equal(collect, "http://www.example.com", *setting.DefaultRedirectURI)
			assert.Equal(collect, false, *setting.AllowDomainDiscovery)
			assert.Equal(collect, false, *setting.AllowUserNamePassword)
			assert.Equal(collect, time.Duration(time.Second*20*20), *setting.PasswordCheckLifetime)
			assert.Equal(collect, time.Duration(time.Second*20*21), *setting.ExternalLoginCheckLifetime)
			assert.Equal(collect, time.Duration(time.Second*20*22), *setting.MFAInitSkipLifetime)
			assert.Equal(collect, time.Duration(time.Second*20*23), *setting.SecondFactorCheckLifetime)
			assert.Equal(collect, time.Duration(time.Second*20*24), *setting.MultiFactorCheckLifetime)
			assert.Equal(collect, true, *setting.DisableLoginWithEmail)
			assert.Equal(collect, true, *setting.DisableLoginWithPhone)
			assert.Equal(collect, true, *setting.ForceMFALocalOnly)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test added/remove login multifactor type reduces", func(t *testing.T) {
		// check inital MFAType value
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			assert.Equal(collect, []domain.MultiFactorType{domain.MultiFactorType(policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION)}, setting.MFAType)
		}, retryDuration, tick)

		// remove MFAType
		_, err := newInstance.Client.Admin.RemoveMultiFactorFromLoginPolicy(IAMCTX, &admin.RemoveMultiFactorFromLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.login.multifactor.remove
			assert.Equal(collect, []domain.MultiFactorType{}, setting.MFAType)
		}, retryDuration, tick)

		before := time.Now()
		_, err = newInstance.Client.Admin.AddMultiFactorToLoginPolicy(IAMCTX, &admin.AddMultiFactorToLoginPolicyRequest{
			Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		})
		require.NoError(t, err)
		after := time.Now()

		// add MFAType
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.login.multifactor.added
			assert.Equal(collect, []domain.MultiFactorType{domain.MultiFactorType(policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION)}, setting.MFAType)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test added/removed second multifactor reduces", func(t *testing.T) {
		before := time.Now()
		// get current second factor types
		var secondFactorTypes []domain.SecondFactorType
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			secondFactorTypes = setting.SecondFactorTypes
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
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.login.multifactor.secondfactor.added
			assert.Equal(collect, secondFactorTypes, setting.SecondFactorTypes)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
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
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.login.multifactor.secondfactor.removed
			assert.Equal(collect, secondFactorTypes, setting.SecondFactorTypes)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(collect, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check login settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLogin, domain.SettingStateActive),
				),
			)

			// event instance.removed
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceLabelSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.BrandingSettingsRepository()

	before := time.Now().Add(-time.Second)
	newInstance := integration.NewInstance(t.Context())
	after := time.Now().Add(time.Second)

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test adding label settings reduces", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.added
			// these values are found in default.yaml
			assert.Equal(collect, "#5469d4", *setting.PrimaryColorLight)
			assert.Equal(collect, "#fafafa", *setting.BackgroundColorLight)
			assert.Equal(collect, "#cd3d56", *setting.WarnColorLight)
			assert.Equal(collect, "#000000", *setting.FontColorLight)
			assert.Equal(collect, "#2073c4", *setting.PrimaryColorDark)
			assert.Equal(collect, "#111827", *setting.BackgroundColorDark)
			assert.Equal(collect, "#ff3b5b", *setting.WarnColorDark)
			assert.Equal(collect, "#ff3b5b", *setting.WarnColorDark)
			assert.Equal(collect, "#ffffff", *setting.FontColorDark)
			assert.Equal(collect, false, *setting.HideLoginNameSuffix)
			assert.Equal(collect, false, *setting.ErrorMsgPopup)
			assert.Equal(collect, false, *setting.DisableWatermark)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.CreatedAt, before, after)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label change", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdateLabelPolicy(IAMCTX, &admin.UpdateLabelPolicyRequest{
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
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(

				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.change
			assert.Equal(collect, "#055000", *setting.PrimaryColorLight)
			assert.Equal(collect, "#055000", *setting.BackgroundColorLight)
			assert.Equal(collect, "#055000", *setting.WarnColorLight)
			assert.Equal(collect, "#055000", *setting.FontColorLight)
			assert.Equal(collect, "#055000", *setting.PrimaryColorDark)
			assert.Equal(collect, "#055000", *setting.BackgroundColorDark)
			assert.Equal(collect, "#055000", *setting.WarnColorDark)
			assert.Equal(collect, "#055000", *setting.WarnColorDark)
			assert.Equal(collect, "#055000", *setting.FontColorDark)
			assert.Equal(collect, true, *setting.HideLoginNameSuffix)
			assert.Equal(collect, false, *setting.ErrorMsgPopup)
			assert.Equal(collect, true, *setting.DisableWatermark)
			assert.Equal(collect, domain.BrandingPolicyThemeLight, *setting.ThemeMode)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test label settings activated", func(t *testing.T) {
		// activate label
		before := time.Now().Add(-time.Second)
		_, err := newInstance.Client.Admin.ActivateLabelPolicy(IAMCTX, &admin.ActivateLabelPolicyRequest{})
		require.NoError(t, err)
		after := time.Now().Add(time.Second)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.activated
			assert.Equal(collect, domain.SettingStateActive, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo light added", func(t *testing.T) {
		// set logo light
		before := time.Now().Add(-time.Second)
		uploadInstanceAsset(IAMCTX, t, newInstance, "/instance/policy/label/logo", bytes.NewReader(picture))
		after := time.Now().Add(time.Second)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.logo.added
			assert.NotNil(collect, setting.LogoURLLight)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo light removed", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.RemoveLabelPolicyLogo(IAMCTX, &admin.RemoveLabelPolicyLogoRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check light logo removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.logo.removed
			assert.Equal(collect, url.URL{}, *setting.LogoURLLight)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark added", func(t *testing.T) {
		before := time.Now()
		uploadInstanceAsset(IAMCTX, t, newInstance, "/instance/policy/label/logo/dark", bytes.NewReader(picture))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.logo.dark.added
			assert.NotNil(collect, setting.LogoURLDark)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label logo dark removed", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.RemoveLabelPolicyLogoDark(IAMCTX, &admin.RemoveLabelPolicyLogoDarkRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check dark logo removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.logo.dark.removed
			assert.Equal(collect, url.URL{}, *setting.LogoURLDark)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light added", func(t *testing.T) {
		before := time.Now()
		uploadInstanceAsset(IAMCTX, t, newInstance, "/instance/policy/label/icon", bytes.NewReader(picture))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.icon.added
			assert.NotNil(collect, setting.IconURLLight)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark added", func(t *testing.T) {
		before := time.Now()
		uploadInstanceAsset(IAMCTX, t, newInstance, "/instance/policy/label/icon/dark", bytes.NewReader(picture))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.icon.dark.added
			assert.NotNil(collect, setting.IconURLDark)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon light removed", func(t *testing.T) {
		// remote logo light
		before := time.Now()
		_, err := newInstance.Client.Admin.RemoveLabelPolicyIcon(IAMCTX, &admin.RemoveLabelPolicyIconRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check light icon removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.icon.removed
			assert.Equal(collect, url.URL{}, *setting.IconURLLight)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy label icon dark removed", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.RemoveLabelPolicyIconDark(IAMCTX, &admin.RemoveLabelPolicyIconDarkRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check dark icon removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.icon.dark.removed
			assert.Equal(collect, url.URL{}, *setting.IconURLDark)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy font added", func(t *testing.T) {
		before := time.Now()
		uploadInstanceAsset(IAMCTX, t, newInstance, "/instance/policy/label/font", bytes.NewReader(font))
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.font.added
			assert.NotNil(collect, setting.FontURL)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test policy font removed", func(t *testing.T) {
		// remote logo light
		before := time.Now()
		_, err := newInstance.Client.Admin.RemoveLabelPolicyFont(IAMCTX, &admin.RemoveLabelPolicyFontRequest{})
		require.NoError(t, err)
		after := time.Now()

		// check font policy removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStatePreview),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.label.font.removed
			assert.Equal(collect, url.URL{}, *setting.FontURL)
			assert.Equal(collect, domain.SettingStatePreview, setting.State)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		// check label preview settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(collect, setting)
		}, retryDuration, tick)

		// check label activated settings exist
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*50)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(collect, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})

		// check label preview settings deleted
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)

			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)

		// check label activated settings deleted
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeBranding, domain.SettingStateActive),
				),
			)

			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
		require.NoError(t, err)
	})
}

func TestServer_TestInstancePasswordComplexitySettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.PasswordComplexitySettingsRepository()

	before := time.Now()
	newInstance := integration.NewInstance(t.Context())
	after := time.Now()

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)
	t.Run("test password complexity added", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.password.complexity.added
			// these values are found in default.yaml
			assert.Equal(collect, uint64(8), *setting.MinLength)
			assert.Equal(collect, true, *setting.HasUppercase)
			assert.Equal(collect, true, *setting.HasLowercase)
			assert.Equal(collect, true, *setting.HasNumber)
			assert.Equal(collect, true, *setting.HasSymbol)
			assert.WithinRange(collect, setting.CreatedAt, before, after)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password complexity change", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdatePasswordComplexityPolicy(IAMCTX, &admin.UpdatePasswordComplexityPolicyRequest{
			MinLength:    5,
			HasUppercase: false,
			HasLowercase: false,
			HasNumber:    false,
			HasSymbol:    false,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.password.complexity.changed
			assert.Equal(collect, uint64(5), *setting.MinLength)
			assert.Equal(collect, false, *setting.HasUppercase)
			assert.Equal(collect, false, *setting.HasLowercase)
			assert.Equal(collect, false, *setting.HasNumber)
			assert.Equal(collect, false, *setting.HasSymbol)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordComplexity, domain.SettingStateActive),
				),
			)

			// event instance.removed
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstancePasswordPolicySettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.PasswordExpirySettingsRepository()
	before := time.Now()
	newInstance := integration.NewInstance(t.Context())
	after := time.Now()

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test password policy added", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.password.age.added
			assert.Equal(collect, uint64(0), *setting.ExpireWarnDays)
			assert.Equal(collect, uint64(0), *setting.MaxAgeDays)
			assert.WithinRange(collect, setting.CreatedAt, before, after)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy changed", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdatePasswordAgePolicy(IAMCTX, &admin.UpdatePasswordAgePolicyRequest{
			MaxAgeDays:     30,
			ExpireWarnDays: 30,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.password.age.changed
			assert.Equal(collect, uint64(30), *setting.ExpireWarnDays)
			assert.Equal(collect, uint64(30), *setting.MaxAgeDays)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check password expiry settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypePasswordExpiry, domain.SettingStateActive),
				),
			)

			// event instance.removed
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceDomainSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.DomainSettingsRepository()
	before := time.Now()
	newInstance := integration.NewInstance(t.Context())
	after := time.Now()

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test domain policy added", func(t *testing.T) {

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.domain.added
			assert.Equal(collect, false, *setting.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(collect, false, *setting.LoginNameIncludesDomain)
			assert.Equal(collect, false, *setting.RequireOrgDomainVerification)
			assert.WithinRange(collect, setting.CreatedAt, before, after)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test domain policy changed", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdateDomainPolicy(IAMCTX, &admin.UpdateDomainPolicyRequest{
			UserLoginMustBeDomain:                  true,
			ValidateOrgDomains:                     true,
			SmtpSenderAddressMatchesInstanceDomain: true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.changed
			assert.Equal(collect, true, *setting.SMTPSenderAddressMatchesInstanceDomain)
			assert.Equal(collect, true, *setting.LoginNameIncludesDomain)
			assert.Equal(collect, true, *setting.RequireOrgDomainVerification)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check domain settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeDomain, domain.SettingStateActive),
				),
			)

			// event instance.removed
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceLockoutSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.LockoutSettingsRepository()

	before := time.Now()
	newInstance := integration.NewInstance(t.Context())
	after := time.Now()

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test lockout policy added", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.lockout.added
			assert.Equal(collect, uint64(0), *setting.MaxOTPAttempts)
			assert.Equal(collect, uint64(0), *setting.MaxPasswordAttempts)
			assert.Equal(collect, true, *setting.ShowLockOutFailures)
			assert.WithinRange(collect, setting.CreatedAt, before, after)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test password policy changed", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdateLockoutPolicy(IAMCTX, &admin.UpdateLockoutPolicyRequest{
			MaxPasswordAttempts: 5,
			MaxOtpAttempts:      5,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.lockout.changed
			assert.Equal(collect, uint64(5), *setting.MaxOTPAttempts)
			assert.Equal(collect, uint64(5), *setting.MaxPasswordAttempts)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		// check login settings exist
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		// delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// check password complexity settings removed
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLockout, domain.SettingStateActive),
				),
			)

			// event instance.removed
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceSecuritySettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.SecuritySettingsRepository()

	newInstance := integration.NewInstance(t.Context())
	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test security policy set", func(t *testing.T) {
		// 1. set security policy
		before := time.Now()
		_, err := newInstance.Client.Admin.SetSecurityPolicy(IAMCTX, &admin.SetSecurityPolicyRequest{
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"value"},
			EnableImpersonation:   true,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeSecurity, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.security.set
			assert.Equal(collect, true, *setting.EnableIframeEmbedding)
			assert.Equal(collect, []string{"value"}, setting.AllowedOrigins)
			assert.Equal(collect, true, *setting.EnableImpersonation)
			assert.WithinRange(collect, setting.CreatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test security policy re-set", func(t *testing.T) {
		// 2. re-set security policy
		before := time.Now()
		_, err := newInstance.Client.Admin.SetSecurityPolicy(IAMCTX, &admin.SetSecurityPolicyRequest{
			EnableIframeEmbedding: false,
			AllowedOrigins:        []string{"new_value"},
			EnableImpersonation:   false,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeSecurity, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			// event instance.policy.security.set
			assert.Equal(collect, false, *setting.EnableIframeEmbedding)
			assert.Equal(collect, []string{"new_value"}, setting.AllowedOrigins)
			assert.Equal(collect, false, *setting.EnableImpersonation)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		// 3. delete instance
		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		// 4. check security settings removed
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*10)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeSecurity, domain.SettingStateActive),
				),
			)

			// event instance.removed
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceNotificationSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.NotificationSettingsRepository()

	before := time.Now()
	newInstance := integration.NewInstance(t.Context())
	after := time.Now()
	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test add notification settings set", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeNotification, domain.SettingStateActive),
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
		_, err := newInstance.Client.Admin.UpdateNotificationPolicy(IAMCTX, &admin.UpdateNotificationPolicyRequest{
			PasswordChange: false,
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeNotification, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, false, *setting.PasswordChange)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeNotification, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeNotification, domain.SettingStateActive),
				),
			)

			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceLegalAndSupportSettingsReduces(t *testing.T) {
	t.Parallel()

	settingsRepo := repository.LegalAndSupportSettingsRepository()

	before := time.Now()
	newInstance := integration.NewInstance(t.Context())
	after := time.Now()
	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test add legal and support settings set", func(t *testing.T) {

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, "", *setting.TOSLink)
			assert.Equal(t, "", *setting.PrivacyPolicyLink)
			assert.Equal(t, "", *setting.HelpLink)
			assert.Equal(t, "", *setting.SupportEmail)
			assert.Equal(t, "https://zitadel.com/docs", *setting.DocsLink)
			assert.Equal(t, "", *setting.CustomLink)
			assert.Equal(t, "", *setting.CustomLinkText)
			assert.WithinRange(t, setting.CreatedAt, before, after)
			assert.WithinRange(t, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test add legal and support settings re-set", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdatePrivacyPolicy(IAMCTX, &admin.UpdatePrivacyPolicyRequest{
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
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
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

	t.Run("test delete instance reduces", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)

			require.NotNil(t, setting)
		}, retryDuration, tick)

		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeLegalAndSupport, domain.SettingStateActive),
				),
			)

			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}

func TestServer_TestInstanceSecretGeneratorSettingsReduces(t *testing.T) {
	settingsRepo := repository.SecretGeneratorSettingsRepository()
	newInstance := integration.NewInstance(t.Context())

	IAMCTX := newInstance.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	t.Run("test secret generator settings set", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)
			// ClientSecret
			assert.Equal(collect, uint(64), *setting.ClientSecret.Length)
			assert.Equal(collect, true, *setting.ClientSecret.IncludeLowerLetters)
			assert.Equal(collect, true, *setting.ClientSecret.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.ClientSecret.IncludeDigits)
			assert.Equal(collect, false, *setting.ClientSecret.IncludeSymbols)

			// InitializeUserCode
			assert.Equal(collect, uint(6), *setting.InitializeUserCode.Length)
			assert.Equal(collect, 72*time.Hour, *setting.InitializeUserCode.Expiry)
			assert.Equal(collect, false, *setting.InitializeUserCode.IncludeLowerLetters)
			assert.Equal(collect, true, *setting.InitializeUserCode.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.InitializeUserCode.IncludeDigits)
			assert.Equal(collect, false, *setting.InitializeUserCode.IncludeSymbols)

			// EmailVerificationCode
			assert.Equal(collect, uint(6), *setting.EmailVerificationCode.Length)
			assert.Equal(collect, 1*time.Hour, *setting.EmailVerificationCode.Expiry)
			assert.Equal(collect, false, *setting.EmailVerificationCode.IncludeLowerLetters)
			assert.Equal(collect, true, *setting.EmailVerificationCode.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.EmailVerificationCode.IncludeDigits)
			assert.Equal(collect, false, *setting.EmailVerificationCode.IncludeSymbols)

			// PhoneVerificationCode
			assert.Equal(collect, uint(6), *setting.PhoneVerificationCode.Length)
			assert.Equal(collect, 1*time.Hour, *setting.PhoneVerificationCode.Expiry)
			assert.Equal(collect, false, *setting.PhoneVerificationCode.IncludeLowerLetters)
			assert.Equal(collect, true, *setting.PhoneVerificationCode.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.PhoneVerificationCode.IncludeDigits)
			assert.Equal(collect, false, *setting.PhoneVerificationCode.IncludeSymbols)

			// PasswordVerificationCode
			assert.Equal(collect, uint(6), *setting.PasswordVerificationCode.Length)
			assert.Equal(collect, 1*time.Hour, *setting.PasswordVerificationCode.Expiry)
			assert.Equal(collect, false, *setting.PasswordVerificationCode.IncludeLowerLetters)
			assert.Equal(collect, true, *setting.PasswordVerificationCode.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.PasswordVerificationCode.IncludeDigits)
			assert.Equal(collect, false, *setting.PasswordVerificationCode.IncludeSymbols)

			// PasswordlessInitCode
			assert.Equal(collect, uint(12), *setting.PasswordlessInitCode.Length)
			assert.Equal(collect, 1*time.Hour, *setting.PasswordlessInitCode.Expiry)
			assert.Equal(collect, true, *setting.PasswordlessInitCode.IncludeLowerLetters)
			assert.Equal(collect, true, *setting.PasswordlessInitCode.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.PasswordlessInitCode.IncludeDigits)
			assert.Equal(collect, false, *setting.PasswordlessInitCode.IncludeSymbols)

			// DomainVerification
			assert.Equal(collect, uint(32), *setting.DomainVerification.Length)
			assert.Equal(collect, true, *setting.DomainVerification.IncludeLowerLetters)
			assert.Equal(collect, true, *setting.DomainVerification.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.DomainVerification.IncludeDigits)
			assert.Equal(collect, false, *setting.DomainVerification.IncludeSymbols)

			// OTPSMS
			assert.Equal(collect, uint(8), *setting.OTPSMS.Length)
			assert.Equal(collect, 5*time.Minute, *setting.OTPSMS.Expiry)
			assert.Equal(collect, false, *setting.OTPSMS.IncludeLowerLetters)
			assert.Equal(collect, false, *setting.OTPSMS.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.OTPSMS.IncludeDigits)
			assert.Equal(collect, false, *setting.OTPSMS.IncludeSymbols)

			// OTPEmail
			assert.Equal(collect, uint(8), *setting.OTPEmail.Length)
			assert.Equal(collect, 5*time.Minute, *setting.OTPEmail.Expiry)
			assert.Equal(collect, false, *setting.OTPEmail.IncludeLowerLetters)
			assert.Equal(collect, false, *setting.OTPEmail.IncludeUpperLetters)
			assert.Equal(collect, true, *setting.OTPEmail.IncludeDigits)
			assert.Equal(collect, false, *setting.OTPEmail.IncludeSymbols)
		}, retryDuration, tick)
	})

	t.Run("test secret generator settings update", func(t *testing.T) {
		before := time.Now()
		_, err := newInstance.Client.Admin.UpdateSecretGenerator(IAMCTX, &admin.UpdateSecretGeneratorRequest{
			GeneratorType:       settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_INIT_CODE,
			IncludeLowerLetters: true,
			Expiry:              durationpb.New(24 * time.Hour),
			Length:              uint32(8),
		})
		require.NoError(t, err)
		after := time.Now()

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
				),
			)
			require.NoError(t, err)
			assert.Equal(collect, uint(8), *setting.InitializeUserCode.Length)
			assert.Equal(collect, 24*time.Hour, *setting.InitializeUserCode.Expiry)
			assert.Equal(collect, true, *setting.InitializeUserCode.IncludeLowerLetters)
			assert.WithinRange(collect, setting.UpdatedAt, before, after)
		}, retryDuration, tick)
	})

	t.Run("test delete instance reduces", func(t *testing.T) {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			setting, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
				),
			)
			require.NoError(collect, err)
			require.NotNil(t, setting)
		}, retryDuration, tick)

		_, err := newInstance.Client.InstanceV2Beta.DeleteInstance(CTX, &instance.DeleteInstanceRequest{
			InstanceId: newInstance.ID(),
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(IAMCTX, time.Second*20)
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := settingsRepo.Get(
				IAMCTX, pool,
				database.WithCondition(
					settingsRepo.UniqueCondition(newInstance.ID(), nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
				),
			)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
