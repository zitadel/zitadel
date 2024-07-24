package settings

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
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
		PasswordCheckLifetime:      database.Duration(time.Hour),
		ExternalLoginCheckLifetime: database.Duration(time.Minute),
		MFAInitSkipLifetime:        database.Duration(time.Millisecond),
		SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
		MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
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

func Test_isDefaultToResourceOwnerTypePb(t *testing.T) {
	type args struct {
		isDefault bool
	}
	tests := []struct {
		args args
		want settings.ResourceOwnerType
	}{
		{
			args: args{false},
			want: settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_ORG,
		},
		{
			args: args{true},
			want: settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			got := isDefaultToResourceOwnerTypePb(tt.args.isDefault)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_passkeysTypeToPb(t *testing.T) {
	type args struct {
		passwordlessType domain.PasswordlessType
	}
	tests := []struct {
		args args
		want settings.PasskeysType
	}{
		{
			args: args{domain.PasswordlessTypeNotAllowed},
			want: settings.PasskeysType_PASSKEYS_TYPE_NOT_ALLOWED,
		},
		{
			args: args{domain.PasswordlessTypeAllowed},
			want: settings.PasskeysType_PASSKEYS_TYPE_ALLOWED,
		},
		{
			args: args{99},
			want: settings.PasskeysType_PASSKEYS_TYPE_NOT_ALLOWED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			got := passkeysTypeToPb(tt.args.passwordlessType)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_secondFactorTypeToPb(t *testing.T) {
	type args struct {
		secondFactorType domain.SecondFactorType
	}
	tests := []struct {
		args args
		want settings.SecondFactorType
	}{
		{
			args: args{domain.SecondFactorTypeTOTP},
			want: settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP,
		},
		{
			args: args{domain.SecondFactorTypeU2F},
			want: settings.SecondFactorType_SECOND_FACTOR_TYPE_U2F,
		},
		{
			args: args{domain.SecondFactorTypeOTPSMS},
			want: settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS,
		},
		{
			args: args{domain.SecondFactorTypeOTPEmail},
			want: settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP_EMAIL,
		},
		{
			args: args{domain.SecondFactorTypeUnspecified},
			want: settings.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED,
		},
		{
			args: args{99},
			want: settings.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			got := secondFactorTypeToPb(tt.args.secondFactorType)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_multiFactorTypeToPb(t *testing.T) {
	type args struct {
		typ domain.MultiFactorType
	}
	tests := []struct {
		args args
		want settings.MultiFactorType
	}{
		{
			args: args{domain.MultiFactorTypeU2FWithPIN},
			want: settings.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
		},
		{
			args: args{domain.MultiFactorTypeUnspecified},
			want: settings.MultiFactorType_MULTI_FACTOR_TYPE_UNSPECIFIED,
		},
		{
			args: args{99},
			want: settings.MultiFactorType_MULTI_FACTOR_TYPE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			got := multiFactorTypeToPb(tt.args.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_passwordComplexitySettingsToPb(t *testing.T) {
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

	got := passwordComplexitySettingsToPb(arg)
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("passwordComplexitySettingsToPb() =\n%v\nwant\n%v", got, want)
	}
}

func Test_passwordExpirySettingsToPb(t *testing.T) {
	arg := &query.PasswordAgePolicy{
		ExpireWarnDays: 80,
		MaxAgeDays:     90,
		IsDefault:      true,
	}
	want := &settings.PasswordExpirySettings{
		ExpireWarnDays:    80,
		MaxAgeDays:        90,
		ResourceOwnerType: settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE,
	}

	got := passwordExpirySettingsToPb(arg)
	grpc.AllFieldsSet(t, got.ProtoReflect(), ignoreTypes...)
	if !proto.Equal(got, want) {
		t.Errorf("passwordExpirySettingsToPb() =\n%v\nwant\n%v", got, want)
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
		ThemeMode:           domain.LabelPolicyThemeDark,
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
		ThemeMode:           settings.ThemeMode_THEME_MODE_DARK,
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
		TOSLink:        "http://example.com/tos",
		PrivacyLink:    "http://example.com/pricacy",
		HelpLink:       "http://example.com/help",
		SupportEmail:   "support@zitadel.com",
		IsDefault:      true,
		DocsLink:       "http://example.com/docs",
		CustomLink:     "http://example.com/custom",
		CustomLinkText: "Custom",
	}
	want := &settings.LegalAndSupportSettings{
		TosLink:           "http://example.com/tos",
		PrivacyPolicyLink: "http://example.com/pricacy",
		HelpLink:          "http://example.com/help",
		SupportEmail:      "support@zitadel.com",
		DocsLink:          "http://example.com/docs",
		CustomLink:        "http://example.com/custom",
		CustomLinkText:    "Custom",
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
		MaxOTPAttempts:      22,
		IsDefault:           true,
	}
	want := &settings.LockoutSettings{
		MaxPasswordAttempts: 22,
		MaxOtpAttempts:      22,
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

func Test_idpTypeToPb(t *testing.T) {
	type args struct {
		idpType domain.IDPType
	}
	tests := []struct {
		args args
		want settings.IdentityProviderType
	}{
		{
			args: args{domain.IDPTypeUnspecified},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_UNSPECIFIED,
		},
		{
			args: args{domain.IDPTypeOIDC},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_OIDC,
		},
		{
			args: args{domain.IDPTypeJWT},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_JWT,
		},
		{
			args: args{domain.IDPTypeOAuth},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_OAUTH,
		},
		{
			args: args{domain.IDPTypeLDAP},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_LDAP,
		},
		{
			args: args{domain.IDPTypeAzureAD},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_AZURE_AD,
		},
		{
			args: args{domain.IDPTypeGitHub},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITHUB,
		},
		{
			args: args{domain.IDPTypeGitHubEnterprise},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITHUB_ES,
		},
		{
			args: args{domain.IDPTypeGitLab},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITLAB,
		},
		{
			args: args{domain.IDPTypeGitLabSelfHosted},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED,
		},
		{
			args: args{domain.IDPTypeGoogle},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GOOGLE,
		},
		{
			args: args{domain.IDPTypeSAML},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_SAML,
		},
		{
			args: args{99},
			want: settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			if got := idpTypeToPb(tt.args.idpType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("idpTypeToPb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_securityPolicyToSettingsPb(t *testing.T) {
	want := &settings.SecuritySettings{
		EmbeddedIframe: &settings.EmbeddedIframeSettings{
			Enabled:        true,
			AllowedOrigins: []string{"foo", "bar"},
		},
		EnableImpersonation: true,
	}
	got := securityPolicyToSettingsPb(&query.SecurityPolicy{
		EnableIframeEmbedding: true,
		AllowedOrigins:        []string{"foo", "bar"},
		EnableImpersonation:   true,
	})
	assert.Equal(t, want, got)
}

func Test_securitySettingsToCommand(t *testing.T) {
	want := &command.SecurityPolicy{
		EnableIframeEmbedding: true,
		AllowedOrigins:        []string{"foo", "bar"},
		EnableImpersonation:   true,
	}
	got := securitySettingsToCommand(&settings.SetSecuritySettingsRequest{
		EmbeddedIframe: &settings.EmbeddedIframeSettings{
			Enabled:        true,
			AllowedOrigins: []string{"foo", "bar"},
		},
		EnableImpersonation: true,
	})
	assert.Equal(t, want, got)
}
