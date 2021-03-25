package domain

import (
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
)

const (
	IAMRolePrefix          = "IAM"
	OrgRolePrefix          = "ORG"
	ProjectRolePrefix      = "PROJECT"
	ProjectGrantRolePrefix = "PROJECT_GRANT"
	RoleOrgOwner           = "ORG_OWNER"
	RoleOrgProjectCreator  = "ORG_PROJECT_CREATOR"
	RoleIAMOwner           = "IAM_OWNER"
	RoleProjectOwner       = "PROJECT_OWNER"
	RoleProjectOwnerGlobal = "PROJECT_OWNER_GLOBAL"
)

func CheckForInvalidRoles(roles []string, rolePrefix string, validRoles []authz.RoleMapping) []string {
	invalidRoles := make([]string, 0)
	for _, role := range roles {
		if !containsRole(role, rolePrefix, validRoles) {
			invalidRoles = append(invalidRoles, role)
		}
	}
	return invalidRoles
}

func containsRole(role, rolePrefix string, validRoles []authz.RoleMapping) bool {
	for _, validRole := range validRoles {
		if role == validRole.Role && strings.HasPrefix(role, rolePrefix) {
			return true
		}
	}
	return false
}
