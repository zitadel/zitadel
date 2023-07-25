package settings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2alpha"
)

var ignoreTypes = []protoreflect.FullName{"google.protobuf.Duration"}

func Test_loginSettingsToPb(t *testing.T) {
	arg := &query.LoginPolicy{
		AllowUsernamePassword:      true,
		AllowRegister:              true,
		AllowExternalIDPs:          true,
		ForceMFA:                   true,
		ForceMFALocalOnly:          true,
		PasswordlessType:           domain.PasswordlessTypeAllowed,
		HidePasswordReset:          true,
		IgnoreUnknownUsernames:     true,
		AllowDomainDiscovery:       true,
		DisableLoginWithEmail:      true,
		DisableLoginWithPhone:      true,
		DefaultRedirectURI:         "example.com",
		PasswordCheckLifetime:      time.Hour,
		ExternalLoginCheckLifetime: time.Minute,
		MFAInitSkipLifetime:        time.Millisecond,
		SecondFactorCheckLifetime:  time.Microsecond,
		MultiFactorCheckLifetime:   time.Nanosecond,
		SecondFactors: []domain.SecondFactorType{
			domain.SecondFactorTypeTOTP,
			domain.SecondFactorTypeU2F,
			domain.SecondFactorTypeOTPEmail,
			domain.SecondFactorTypeOTPSMS,
		},
		MultiFactors: []domain.MultiFactorType{
			domain.MultiFactorTypeU2FWithPIN,
		},
		IsDefault: true,
	}

	want := &settings.LoginSettings{
		AllowUsernamePassword:      true,
		AllowRegister:              true,
		AllowExternalIdp:           true,
		ForceMfa:                   true,
		ForceMfaLocalOnly:          true,
		PasskeysType:               settings.PasskeysType_PASSKEYS_TYPE_ALLOWED,
		HidePasswordReset:          true,
		IgnoreUnknownUsernames:     true,
		AllowDomainDiscovery:       true,
		DisableLoginWithEmail:      true,
		DisableLoginWithPhone:      true,
		DefaultRedirectUri:         "example.com",
		PasswordCheckLifetime:      durationpb.New(time.Hour),
		ExternalLoginCheckLifetime: durationpb.New(time.Minute),
		MfaInitSkipLifetime:        durationpb.New(time.Millisecond),
		SecondFactorCheckLifetime:  durationpb.New(time.Microsecond),
		MultiFactorCheckLifetime:   durationpb.New(time.Nanosecond),
		SecondFactors: []settings.SecondFactorType{
			settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP,
			settings.SecondFactorType_SECOND_FACTOR_TYPE_U2F,
			settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP_EMAIL,
			settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		},
		MultiFactors: []settings.MultiFactorType{
			settings.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		},
		ResourceOwnerType: settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
	}

	got := loginSettingsToPb(arg)
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("loginSettingsToPb() =\n%v\nwant\n%v", got, want)
	}
}

func Test_passwordSettingsToPb(t *testing.T) {
	arg := &query.PasswordComplexityPolicy{
		MinLength:    12,
		HasUppercase: true,
		HasLowercase: true,
		HasNumber:    true,
		HasSymbol:    true,
		IsDefault:    true,
	}
	want := &settings.PasswordComplexitySettings{
		MinLength:         12,
		RequiresUppercase: true,
		RequiresLowercase: true,
		RequiresNumber:    true,
		RequiresSymbol:    true,
		ResourceOwnerType: settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
	}

	got := passwordSettingsToPb(arg)
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("passwordSettingsToPb() =\n%v\nwant\n%v", got, want)
	}
}

func Test_brandingSettingsToPb(t *testing.T) {
	arg := &query.LabelPolicy{
		Light: query.Theme{
			PrimaryColor:    "red",
			WarnColor:       "white",
			BackgroundColor: "blue",
			FontColor:       "orange",
			LogoURL:         "light-logo",
			IconURL:         "light-icon",
		},
		Dark: query.Theme{
			PrimaryColor:    "magenta",
			WarnColor:       "pink",
			BackgroundColor: "black",
			FontColor:       "white",
			LogoURL:         "dark-logo",
			IconURL:         "dark-icon",
		},
		ResourceOwner:       "me",
		FontURL:             "fonts",
		WatermarkDisabled:   true,
		HideLoginNameSuffix: true,
		IsDefault:           true,
	}
	want := &settings.BrandingSettings{
		LightTheme: &settings.Theme{
			PrimaryColor:    "red",
			WarnColor:       "white",
			BackgroundColor: "blue",
			FontColor:       "orange",
			LogoUrl:         "http://example.com/me/light-logo",
			IconUrl:         "http://example.com/me/light-icon",
		},
		DarkTheme: &settings.Theme{
			PrimaryColor:    "magenta",
			WarnColor:       "pink",
			BackgroundColor: "black",
			FontColor:       "white",
			LogoUrl:         "http://example.com/me/dark-logo",
			IconUrl:         "http://example.com/me/dark-icon",
		},
		FontUrl:             "http://example.com/me/fonts",
		DisableWatermark:    true,
		HideLoginNameSuffix: true,
		ResourceOwnerType:   settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
	}

	got := brandingSettingsToPb(arg, "http://example.com")
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("brandingSettingsToPb() =\n%v\nwant\n%v", got, want)
	}
}

func Test_domainSettingsToPb(t *testing.T) {
	arg := &query.DomainPolicy{
		UserLoginMustBeDomain:                  true,
		ValidateOrgDomains:                     true,
		SMTPSenderAddressMatchesInstanceDomain: true,
		IsDefault:                              true,
	}
	want := &settings.DomainSettings{
		LoginNameIncludesDomain:                true,
		RequireOrgDomainVerification:           true,
		SmtpSenderAddressMatchesInstanceDomain: true,
		ResourceOwnerType:                      settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
	}
	got := domainSettingsToPb(arg)
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("domainSettingsToPb() =\n%v\nwant\n%v", got, want)
	}
}

func Test_legalSettingsToPb(t *testing.T) {
	arg := &query.PrivacyPolicy{
		TOSLink:      "http://example.com/tos",
		PrivacyLink:  "http://example.com/pricacy",
		HelpLink:     "http://example.com/help",
		SupportEmail: "support@zitadel.com",
		IsDefault:    true,
	}
	want := &settings.LegalAndSupportSettings{
		TosLink:           "http://example.com/tos",
		PrivacyPolicyLink: "http://example.com/pricacy",
		HelpLink:          "http://example.com/help",
		SupportEmail:      "support@zitadel.com",
		ResourceOwnerType: settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
	}
	got := legalAndSupportSettingsToPb(arg)
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("legalSettingsToPb() =\n%v\nwant\n%v", got, want)
	}
}

func Test_lockoutSettingsToPb(t *testing.T) {
	arg := &query.LockoutPolicy{
		MaxPasswordAttempts: 22,
		IsDefault:           true,
	}
	want := &settings.LockoutSettings{
		MaxPasswordAttempts: 22,
		ResourceOwnerType:   settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
	}
	got := lockoutSettingsToPb(arg)
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("lockoutSettingsToPb() =\n%v\nwant\n%v", got, want)
	}
}

func Test_identityProvidersToPb(t *testing.T) {
	arg := []*query.IDPLoginPolicyLink{
		{
			IDPID:   "1",
			IDPName: "foo",
			IDPType: domain.IDPTypeOIDC,
		},
		{
			IDPID:   "2",
			IDPName: "bar",
			IDPType: domain.IDPTypeGitHub,
		},
	}
	want := []*settings.IdentityProvider{
		{
			Id:   "1",
			Name: "foo",
			Type: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_OIDC,
		},
		{
			Id:   "2",
			Name: "bar",
			Type: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITHUB,
		},
	}
	got := identityProvidersToPb(arg)
	require.Len(t, got, len(got))
	for i, v := range got {
		grpc.AllFieldsSet(t, v.ProtoReflect(), ignoreTypes...)
		if !proto.Equal(v, want[i]) {
			t.Errorf("identityProvidersToPb() =\n%v\nwant\n%v", got, want)
		}
	}
}
