package oidc

import (
	"context"
	"slices"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (s *Server) getUserInfoWithRoles(ctx context.Context, userID, projectID string, scope, roleAudience []string) (_ *oidc.UserInfo, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userInfoChan := make(chan *userInfoResult)
	go s.getUserInfo(ctx, userID, scope, roleAudience, userInfoChan)

	rolesChan := make(chan *assertRolesResult)
	go s.assertRoles(ctx, userID, projectID, scope, roleAudience, rolesChan)

	var (
		userInfoResult    *userInfoResult
		assertRolesResult *assertRolesResult
	)

	// make sure both channels are always read,
	// and cancel the context on first error
	for i := 0; i < 2; i++ {
		var resErr error

		select {
		case userInfoResult = <-userInfoChan:
			resErr = userInfoResult.err
		case assertRolesResult = <-rolesChan:
			resErr = assertRolesResult.err
		}

		if resErr == nil {
			continue
		}
		cancel()

		// we only care for the first error that occured,
		// as the next error is most probably a context error.
		if err == nil {
			err = resErr
		}
	}

	userInfo := userInfoToOIDC(userInfoResult.userInfo, scope)
	setUserInfoRoleClaims(userInfo, assertRolesResult.projectsRoles)

	return userInfo, nil
}

type userInfoResult struct {
	userInfo *query.OIDCUserInfo
	err      error
}

func (s *Server) getUserInfo(ctx context.Context, userID string, scope, roleAudience []string, rc chan<- *userInfoResult) {
	userInfo, err := s.storage.query.GetOIDCUserInfo(ctx, userID, scope, roleAudience)
	rc <- &userInfoResult{
		userInfo: userInfo,
		err:      err,
	}
}

type assertRolesResult struct {
	userGrants    *query.UserGrants
	projectsRoles *projectsRoles
	err           error
}

func (s *Server) assertRoles(ctx context.Context, userID, projectID string, scope, roleAudience []string, rc chan<- *assertRolesResult) {
	userGrands, projectsRoles, err := func() (*query.UserGrants, *projectsRoles, error) {
		// if all roles are requested take the audience for those from the scopes
		if slices.Contains(scope, domain.ScopeProjectsRoles) {
			roleAudience = domain.AddAudScopeToAudience(ctx, roleAudience, scope)
		}

		requestedRoles := make([]string, 0, len(scope))
		for _, s := range scope {
			if role, ok := strings.CutPrefix(s, ScopeProjectRolePrefix); ok {
				requestedRoles = append(requestedRoles, role)
			}
		}

		if len(requestedRoles) == 0 && len(roleAudience) == 0 {
			return nil, nil, nil
		}

		// ensure the projectID of the requesting is part of the roleAudience
		if projectID != "" {
			roleAudience = append(roleAudience, projectID)
		}
		queries := make([]query.SearchQuery, 0, 2)
		projectQuery, err := query.NewUserGrantProjectIDsSearchQuery(roleAudience)
		if err != nil {
			return nil, nil, err
		}
		queries = append(queries, projectQuery)
		userIDQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
		if err != nil {
			return nil, nil, err
		}
		queries = append(queries, userIDQuery)
		grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
			Queries: queries,
		}, false, false) // triggers disabled
		if err != nil {
			return nil, nil, err
		}
		roles := new(projectsRoles)
		// if specific roles where requested, check if they are granted and append them in the roles list
		if len(requestedRoles) > 0 {
			for _, requestedRole := range requestedRoles {
				for _, grant := range grants.UserGrants {
					checkGrantedRoles(roles, grant, requestedRole, grant.ProjectID == projectID)
				}
			}
			return grants, roles, nil
		}
		// no specific roles were requested, so convert any grants into roles
		for _, grant := range grants.UserGrants {
			for _, role := range grant.Roles {
				roles.Add(grant.ProjectID, role, grant.ResourceOwner, grant.OrgPrimaryDomain, grant.ProjectID == projectID)
			}
		}
		return grants, roles, nil
	}()

	rc <- &assertRolesResult{
		userGrants:    userGrands,
		projectsRoles: projectsRoles,
		err:           err,
	}
}

func userInfoToOIDC(user *query.OIDCUserInfo, scope []string) *oidc.UserInfo {
	out := new(oidc.UserInfo)
	for _, s := range scope {
		switch s {
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
			if strings.HasPrefix(s, domain.OrgDomainPrimaryScope) {
				out.AppendClaims(domain.OrgDomainPrimaryClaim, strings.TrimPrefix(s, domain.OrgDomainPrimaryScope))
			}
			if strings.HasPrefix(s, domain.OrgIDScope) {
				out.AppendClaims(domain.OrgIDClaim, strings.TrimPrefix(s, domain.OrgIDScope))
				setUserInfoOrgClaims(user, out)
			}
		}
	}

	return out
}

func userInfoEmailToOIDC(user *query.OIDCUserInfo) oidc.UserInfoEmail {
	return oidc.UserInfoEmail{
		Email:         string(user.Email),
		EmailVerified: oidc.Bool(user.IsEmailVerified),
	}
}

func userInfoProfileToOidc(user *query.OIDCUserInfo) oidc.UserInfoProfile {
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

func userInfoPhoneToOIDC(user *query.OIDCUserInfo) oidc.UserInfoPhone {
	return oidc.UserInfoPhone{
		PhoneNumber:         string(user.Phone),
		PhoneNumberVerified: user.IsPhoneVerified,
	}
}

func userInfoAddressToOIDC(user *query.OIDCUserInfo) *oidc.UserInfoAddress {
	return &oidc.UserInfoAddress{
		// Formatted: ??,
		StreetAddress: user.StreetAddress,
		Locality:      user.Locality,
		Region:        user.Region,
		PostalCode:    user.PostalCode,
		Country:       user.Country,
	}
}

func setUserInfoOrgClaims(user *query.OIDCUserInfo, out *oidc.UserInfo) {
	out.AppendClaims(ClaimResourceOwner+"id", user.OrgID)
	out.AppendClaims(ClaimResourceOwner+"name", user.OrgName)
	out.AppendClaims(ClaimResourceOwner+"primary_domain", user.OrgPrimaryDomain)
}
