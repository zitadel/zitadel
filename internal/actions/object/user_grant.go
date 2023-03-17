package object

import (
	"time"

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

type userGrantList struct {
	Count     uint64
	Sequence  uint64
	Timestamp time.Time
	Grants    []*userGrant
}

type userGrant struct {
	Id                         string
	ProjectGrantId             string
	State                      domain.UserGrantState
	UserGrantResourceOwner     string
	UserGrantResourceOwnerName string

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64

	UserId            string
	UserResourceOwner string
	Roles             []string

	ProjectId   string
	ProjectName string
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

func UserGrantsFromQuery(c *actions.FieldConfig, userGrants *query.UserGrants) goja.Value {
	grantList := &userGrantList{
		Count:     userGrants.Count,
		Sequence:  userGrants.Sequence,
		Timestamp: userGrants.Timestamp,
		Grants:    make([]*userGrant, len(userGrants.UserGrants)),
	}

	for i, grant := range userGrants.UserGrants {
		grantList.Grants[i] = &userGrant{
			Id:                         grant.ID,
			ProjectGrantId:             grant.GrantID,
			State:                      grant.State,
			CreationDate:               grant.CreationDate,
			ChangeDate:                 grant.ChangeDate,
			Sequence:                   grant.Sequence,
			UserId:                     grant.UserID,
			Roles:                      grant.Roles,
			UserResourceOwner:          grant.UserResourceOwner,
			UserGrantResourceOwner:     grant.ResourceOwner,
			UserGrantResourceOwnerName: grant.OrgName,
			ProjectId:                  grant.ProjectID,
			ProjectName:                grant.ProjectName,
		}
	}

	return c.Runtime.ToValue(grantList)
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
