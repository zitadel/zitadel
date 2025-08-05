package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
)

func AddAudScopeToAudience(ctx context.Context, audience, scopes []string) []string {
	for _, scope := range scopes {
		if !(strings.HasPrefix(scope, ProjectIDScope) && strings.HasSuffix(scope, AudSuffix)) {
			continue
		}
		projectID := strings.TrimSuffix(strings.TrimPrefix(scope, ProjectIDScope), AudSuffix)
		if projectID == ProjectIDScopeZITADEL {
			projectID = authz.GetInstance(ctx).ProjectID()
		}
		audience = addProjectID(audience, projectID)
	}
	return audience
}

// RoleOrgIDsFromScope parses orgIDs from [OrgRoleIDScope] prefixed scopes.
func RoleOrgIDsFromScope(scopes []string) (orgIDs []string) {
	for _, scope := range scopes {
		orgID, found := strings.CutPrefix(scope, OrgRoleIDScope)
		if found {
			orgIDs = append(orgIDs, orgID)
		}
	}
	return orgIDs
}

func addProjectID(audience []string, projectID string) []string {
	for _, a := range audience {
		if a == projectID {
			return audience
		}
	}
	return append(audience, projectID)
}

//go:generate enumer -type TokenReason -transform snake -trimprefix TokenReason -json
type TokenReason int

const (
	TokenReasonUnspecified TokenReason = iota
	TokenReasonAuthRequest
	TokenReasonRefresh
	TokenReasonJWTProfile
	TokenReasonClientCredentials
	TokenReasonExchange
	TokenReasonImpersonation
)

type TokenActor struct {
	Actor  *TokenActor `json:"actor,omitempty"`
	UserID string      `json:"user_id,omitempty"`
	Issuer string      `json:"issuer,omitempty"`
}
