package query

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
)

func (q *Queries) GetIAMMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range q.zitadelRoles {
		if strings.HasPrefix(roleMap.Role, "IAM") {
			roles = append(roles, roleMap.Role)
		}
	}
	return roles
}

func (q *Queries) GetOrgMemberRoles(isGlobal bool) []string {
	roles := make([]string, 0)
	for _, roleMap := range q.zitadelRoles {
		if strings.HasPrefix(roleMap.Role, "ORG") {
			roles = append(roles, roleMap.Role)
		}
	}
	if isGlobal {
		roles = append(roles, domain.RoleSelfManagementGlobal)
	}
	return roles
}

func (q *Queries) GetProjectMemberRoles(ctx context.Context) ([]string, error) {
	instance, err := q.Instance(ctx, false)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0)
	defaultOrg := authz.GetCtxData(ctx).OrgID == instance.DefaultOrgID
	for _, roleMap := range q.zitadelRoles {
		if strings.HasPrefix(roleMap.Role, "PROJECT") && !strings.HasPrefix(roleMap.Role, "PROJECT_GRANT") {
			if defaultOrg && !strings.HasSuffix(roleMap.Role, "GLOBAL") {
				continue
			}
			roles = append(roles, roleMap.Role)
		}
	}
	return roles, nil
}

func (q *Queries) GetProjectGrantMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range q.zitadelRoles {
		if strings.HasPrefix(roleMap.Role, "PROJECT_GRANT") {
			roles = append(roles, roleMap.Role)
		}
	}
	return roles
}
