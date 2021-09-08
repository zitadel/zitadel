package actions

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
)

type UserGrant struct {
	ProjectID      string
	ProjectGrantID string
	Roles          []string
}

func appendUserGrant(list *[]UserGrant) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		userGrantMap := call.Argument(0).Export()
		userGrant, _ := userGrantFromMap(userGrantMap)
		fmt.Println(userGrant)
		*list = append(*list, userGrant)
		return nil
	}
}

func userGrantFromMap(grantMap interface{}) (UserGrant, error) {
	m, ok := grantMap.(map[string]interface{})
	if !ok {
		return UserGrant{}, errors.New("invalid")
	}
	projectID, ok := m["projectID"].(string)
	if !ok {
		return UserGrant{}, errors.New("invalid")
	}
	var projectGrantID string
	if id, ok := m["projectGrantID"]; ok {
		projectGrantID, ok = id.(string)
		if !ok {
			return UserGrant{}, errors.New("invalid")
		}
	}
	var roles []string
	if r := m["roles"]; r != nil {
		rs, ok := r.([]interface{})
		if !ok {
			return UserGrant{}, errors.New("invalid")
		}
		for _, role := range rs {
			roles = append(roles, role.(string))
		}
	}
	return UserGrant{
		ProjectID:      projectID,
		ProjectGrantID: projectGrantID,
		Roles:          roles,
	}, nil
}
