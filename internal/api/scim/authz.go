package scim

import (
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
)

var AuthMapping = authz.MethodMapping{
	"POST:/scim/v2/" + http.OrgIdInPathVariable + "/Users": {
		Permission: domain.PermissionUserWrite,
	},
	"POST:/scim/v2/" + http.OrgIdInPathVariable + "/Users/.search": {
		Permission: domain.PermissionUserRead,
	},
	"GET:/scim/v2/" + http.OrgIdInPathVariable + "/Users": {
		Permission: domain.PermissionUserRead,
	},
	"GET:/scim/v2/" + http.OrgIdInPathVariable + "/Users/{id}": {
		Permission: domain.PermissionUserRead,
	},
	"PUT:/scim/v2/" + http.OrgIdInPathVariable + "/Users/{id}": {
		Permission: domain.PermissionUserWrite,
	},
	"PATCH:/scim/v2/" + http.OrgIdInPathVariable + "/Users/{id}": {
		Permission: domain.PermissionUserWrite,
	},
	"DELETE:/scim/v2/" + http.OrgIdInPathVariable + "/Users/{id}": {
		Permission: domain.PermissionUserDelete,
	},
	"POST:/scim/v2/" + http.OrgIdInPathVariable + "/Groups": {
		Permission: domain.PermissionGroupCreate,
	},
	"POST:/scim/v2/" + http.OrgIdInPathVariable + "/Groups/.search": {
		Permission: domain.PermissionGroupRead,
	},
	"GET:/scim/v2/" + http.OrgIdInPathVariable + "/Groups": {
		Permission: domain.PermissionGroupRead,
	},
	"GET:/scim/v2/" + http.OrgIdInPathVariable + "/Groups/{id}": {
		Permission: domain.PermissionGroupRead,
	},
	"PUT:/scim/v2/" + http.OrgIdInPathVariable + "/Groups/{id}": {
		Permission: domain.PermissionGroupWrite,
	},
	"PATCH:/scim/v2/" + http.OrgIdInPathVariable + "/Groups/{id}": {
		Permission: domain.PermissionGroupWrite,
	},
	"DELETE:/scim/v2/" + http.OrgIdInPathVariable + "/Groups/{id}": {
		Permission: domain.PermissionGroupDelete,
	},
	"POST:/scim/v2/" + http.OrgIdInPathVariable + "/Bulk": {
		Permission: "authenticated",
	},
}
