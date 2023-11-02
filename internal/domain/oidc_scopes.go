package domain

import "github.com/zitadel/oidc/v3/pkg/oidc"

const (
	ScopeOpenID        = oidc.ScopeOpenID
	ScopeProfile       = oidc.ScopeProfile
	ScopeEmail         = oidc.ScopeEmail
	ScopeAddress       = oidc.ScopeAddress
	ScopePhone         = oidc.ScopePhone
	ScopeOfflineAccess = oidc.ScopeOfflineAccess

	ScopeProjectRolePrefix  = "urn:zitadel:iam:org:project:role:"
	ScopeProjectsRoles      = "urn:zitadel:iam:org:projects:roles"
	ClaimProjectRoles       = "urn:zitadel:iam:org:project:roles"
	ClaimProjectRolesFormat = "urn:zitadel:iam:org:project:%s:roles"
	ScopeUserMetaData       = "urn:zitadel:iam:user:metadata"
	ClaimUserMetaData       = ScopeUserMetaData
	ScopeResourceOwner      = "urn:zitadel:iam:user:resourceowner"
	ClaimResourceOwner      = ScopeResourceOwner + ":"
	ClaimActionLogFormat    = "urn:zitadel:iam:action:%s:log"
)
