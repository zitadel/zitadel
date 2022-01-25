package query

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
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
	iam, err := q.IAMByID(ctx, domain.IAMID)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0)
	global := authz.GetCtxData(ctx).OrgID == iam.GlobalOrgID
	for _, roleMap := range q.zitadelRoles {
		if strings.HasPrefix(roleMap.Role, "PROJECT") && !strings.HasPrefix(roleMap.Role, "PROJECT_GRANT") {
			if global && !strings.HasSuffix(roleMap.Role, "GLOBAL") {
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
