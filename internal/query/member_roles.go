package query

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
)

func (q *Queries) GetIAMMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range q.roles {
		if strings.HasPrefix(roleMap, "IAM") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (q *Queries) GetOrgMemberRoles(isGlobal bool) []string {
	roles := make([]string, 0)
	for _, roleMap := range q.roles {
		if strings.HasPrefix(roleMap, "ORG") {
			roles = append(roles, roleMap)
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
	for _, roleMap := range q.roles {
		if strings.HasPrefix(roleMap, "PROJECT") && !strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			if global && !strings.HasSuffix(roleMap, "GLOBAL") {
				continue
			}
			roles = append(roles, roleMap)
		}
	}
	return roles, nil
}

func (q *Queries) GetProjectGrantMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range q.roles {
		if strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}
