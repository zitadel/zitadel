package migrate

import (
	"encoding/json"
	"reflect"

	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
)

type SystemAPIUsers map[string]*internal_authz.SystemAPIUser

func systemAPIUsersDecodeHook(from, to reflect.Value) (any, error) {
	if to.Type() != reflect.TypeOf(SystemAPIUsers{}) {
		return from.Interface(), nil
	}

	data, ok := from.Interface().(string)
	if !ok {
		return from.Interface(), nil
	}
	users := make(SystemAPIUsers)
	err := json.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
