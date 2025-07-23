package policy

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
)

func ModelLoginPolicyToPb(policy *query.LoginPolicy) *policy_pb.LoginPolicy {
	return &policy_pb.LoginPolicy{
		IsDefault:                  policy.IsDefault,
		AllowUsernamePassword:      policy.AllowUsernamePassword,
		AllowRegister:              policy.AllowRegister,
		AllowExternalIdp:           policy.AllowExternalIDPs,
		ForceMfa:                   policy.ForceMFA,
		ForceMfaLocalOnly:          policy.ForceMFALocalOnly,
		PasswordlessType:           ModelPasswordlessTypeToPb(policy.PasswordlessType),
		HidePasswordReset:          policy.HidePasswordReset,
		IgnoreUnknownUsernames:     policy.IgnoreUnknownUsernames,
		AllowDomainDiscovery:       policy.AllowDomainDiscovery,
		DisableLoginWithEmail:      policy.DisableLoginWithEmail,
		DisableLoginWithPhone:      policy.DisableLoginWithPhone,
		DefaultRedirectUri:         policy.DefaultRedirectURI,
		PasswordCheckLifetime:      durationpb.New(time.Duration(policy.PasswordCheckLifetime)),
		ExternalLoginCheckLifetime: durationpb.New(time.Duration(policy.ExternalLoginCheckLifetime)),
		MfaInitSkipLifetime:        durationpb.New(time.Duration(policy.MFAInitSkipLifetime)),
		SecondFactorCheckLifetime:  durationpb.New(time.Duration(policy.SecondFactorCheckLifetime)),
		MultiFactorCheckLifetime:   durationpb.New(time.Duration(policy.MultiFactorCheckLifetime)),
		EnableRegistrationCaptcha:  policy.EnableRegistrationCaptcha,
		EnableLoginCaptcha:         policy.EnableLoginCaptcha,
		CaptchaType:                ModelCaptchaTypeToPb(policy.CaptchaType),
		CaptchaSiteKey:             policy.CaptchaSiteKey,
		CaptchaSecretKey:           policy.CaptchaSecretKey,
		SecondFactors:              ModelSecondFactorTypesToPb(policy.SecondFactors),
		MultiFactors:               ModelMultiFactorTypesToPb(policy.MultiFactors),
		Idps:                       idp_grpc.IDPLoginPolicyLinksToPb(policy.IDPLinks),
		Details: &object.ObjectDetails{
			Sequence:      policy.Sequence,
			CreationDate:  timestamppb.New(policy.CreationDate),
			ChangeDate:    timestamppb.New(policy.ChangeDate),
			ResourceOwner: policy.OrgID,
		},
	}
}

func PasswordlessTypeToDomain(passwordlessType policy_pb.PasswordlessType) domain.PasswordlessType {
	switch passwordlessType {
	case policy_pb.PasswordlessType_PASSWORDLESS_TYPE_ALLOWED:
		return domain.PasswordlessTypeAllowed
	case policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED:
		return domain.PasswordlessTypeNotAllowed
	default:
		return -1
	}
}

func ModelPasswordlessTypeToPb(passwordlessType domain.PasswordlessType) policy_pb.PasswordlessType {
	switch passwordlessType {
	case domain.PasswordlessTypeAllowed:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_ALLOWED
	case domain.PasswordlessTypeNotAllowed:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	default:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	}
}

func CaptchaTypeToDomain(captchaType policy_pb.CaptchaType) domain.CaptchaType {
	switch captchaType {
	case policy_pb.CaptchaType_CAPTCHA_TYPE_DISABLED:
		return domain.CaptchaTypeDisabled
	case policy_pb.CaptchaType_CAPTCHA_TYPE_RECAPTCHA:
		return domain.CaptchaTypeReCaptcha
	default:
		return -1
	}
}

func ModelCaptchaTypeToPb(captchaType domain.CaptchaType) policy_pb.CaptchaType {
	switch captchaType {
	case domain.CaptchaTypeDisabled:
		return policy_pb.CaptchaType_CAPTCHA_TYPE_DISABLED
	case domain.CaptchaTypeReCaptcha:
		return policy_pb.CaptchaType_CAPTCHA_TYPE_RECAPTCHA
	default:
		return policy_pb.CaptchaType_CAPTCHA_TYPE_DISABLED
	}
}
