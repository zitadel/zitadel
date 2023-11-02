package oidc

import (
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func userinfoToOIDC(user *query.OIDCUserinfo, scopes []string) *oidc.UserInfo {
	out := new(oidc.UserInfo)
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			out.Subject = user.ID
		case oidc.ScopeEmail:
			out.UserInfoEmail = userInfoEmailToOIDC(user)
		case oidc.ScopeProfile:
			out.UserInfoProfile = userInfoProfileToOidc(user)
		case oidc.ScopePhone:
			out.UserInfoPhone = userInfoPhoneToOIDC(user)
		case oidc.ScopeAddress:
			out.Address = userInfoAddressToOIDC(user)
		case ScopeUserMetaData:
			if len(user.Metadata) > 0 {
				out.AppendClaims(ClaimUserMetaData, user.Metadata)
			}
		case ScopeResourceOwner:
			setUserInfoOrgClaims(user, out)
		default:
			if strings.HasPrefix(scope, domain.OrgDomainPrimaryScope) {
				out.AppendClaims(domain.OrgDomainPrimaryClaim, strings.TrimPrefix(scope, domain.OrgDomainPrimaryScope))
			}
			if strings.HasPrefix(scope, domain.OrgIDScope) {
				out.AppendClaims(domain.OrgIDClaim, strings.TrimPrefix(scope, domain.OrgIDScope))
				setUserInfoOrgClaims(user, out)
			}
		}
	}

	return out
}

func userInfoEmailToOIDC(user *query.OIDCUserinfo) oidc.UserInfoEmail {
	return oidc.UserInfoEmail{
		Email:         string(user.Email),
		EmailVerified: oidc.Bool(user.IsEmailVerified),
	}
}

func userInfoProfileToOidc(user *query.OIDCUserinfo) oidc.UserInfoProfile {
	return oidc.UserInfoProfile{
		Name:       user.Name,
		GivenName:  user.FirstName,
		FamilyName: user.LastName,
		Nickname:   user.NickName,
		// Picture:    domain.AvatarURL(o.assetAPIPrefix(ctx), user.ResourceOwner, user.Human.AvatarKey),
		Gender:    getGender(user.Gender),
		Locale:    oidc.NewLocale(user.PreferredLanguage),
		UpdatedAt: oidc.FromTime(user.UpdatedAt),
		// PreferredUsername: user.PreferredLoginName,
	}
}

func userInfoPhoneToOIDC(user *query.OIDCUserinfo) oidc.UserInfoPhone {
	return oidc.UserInfoPhone{
		PhoneNumber:         string(user.Phone),
		PhoneNumberVerified: user.IsPhoneVerified,
	}
}

func userInfoAddressToOIDC(user *query.OIDCUserinfo) *oidc.UserInfoAddress {
	return &oidc.UserInfoAddress{
		// Formatted: ??,
		StreetAddress: user.StreetAddress,
		Locality:      user.Locality,
		Region:        user.Region,
		PostalCode:    user.PostalCode,
		Country:       user.Country,
	}
}

func setUserInfoOrgClaims(user *query.OIDCUserinfo, out *oidc.UserInfo) {
	out.AppendClaims(ClaimResourceOwner+"id", user.OrgID)
	out.AppendClaims(ClaimResourceOwner+"name", user.OrgName)
	out.AppendClaims(ClaimResourceOwner+"primary_domain", user.OrgPrimaryDomain)
}
