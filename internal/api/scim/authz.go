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
	"GET:/scim/v2/" + http.OrgIdInPathVariable + "/Users/{id}": {
		Permission: domain.PermissionUserRead,
	},
	"PUT:/scim/v2/" + http.OrgIdInPathVariable + "/Users/{id}": {
		Permission: domain.PermissionUserWrite,
	},
	"DELETE:/scim/v2/" + http.OrgIdInPathVariable + "/Users/{id}": {
		Permission: domain.PermissionUserDelete,
	},
}
