package object

import (
	"github.com/dop251/goja"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

type UserGrants struct {
	UserGrants []UserGrant
}

type UserGrant struct {
	ProjectID      string
	ProjectGrantID string
	Roles          []string
}

func AppendGrantFunc(userGrants *UserGrants) func(c *actions.FieldConfig) func(call goja.FunctionCall) goja.Value {
	return func(c *actions.FieldConfig) func(call goja.FunctionCall) goja.Value {
		return func(call goja.FunctionCall) goja.Value {
			firstArg := objectFromFirstArgument(call, c.Runtime)
			grant := UserGrant{}
			mapObjectToGrant(firstArg, &grant)
			userGrants.UserGrants = append(userGrants.UserGrants, grant)
			return nil
		}
	}
}

func UserGrantsFromQuery(c *actions.FieldConfig,userGrants *query.UserGrants) goja.Value{
	return c.Runtime.ToValue(&)
}

func UserGrantsToDomain(userID string, actionUserGrants []UserGrant) []*domain.UserGrant {
	if actionUserGrants == nil {
		return nil
	}
	userGrants := make([]*domain.UserGrant, len(actionUserGrants))
	for i, grant := range actionUserGrants {
		userGrants[i] = &domain.UserGrant{
			UserID:         userID,
			ProjectID:      grant.ProjectID,
			ProjectGrantID: grant.ProjectGrantID,
			RoleKeys:       grant.Roles,
		}
	}
	return userGrants
}

func mapObjectToGrant(object *goja.Object, grant *UserGrant) {
	for _, key := range object.Keys() {
		switch key {
		case "projectId":
			grant.ProjectID = object.Get(key).String()
		case "projectGrantId":
			grant.ProjectGrantID = object.Get(key).String()
		case "roles":
			if roles, ok := object.Get(key).Export().([]interface{}); ok {
				for _, role := range roles {
					if r, ok := role.(string); ok {
						grant.Roles = append(grant.Roles, r)
					}
				}
			}
		}
	}
	if grant.ProjectID == "" {
		panic("projectId not set")
	}
}
