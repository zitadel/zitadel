package systemapi

import (
	"encoding/json"
	"reflect"

	"github.com/zitadel/zitadel/internal/api/authz"
)

type Users map[string]*authz.SystemAPIUser

func UsersDecodeHook(from, to reflect.Value) (any, error) {
	if to.Type() != reflect.TypeOf(Users{}) {
		return from.Interface(), nil
	}

	data, ok := from.Interface().(string)
	if !ok {
		return from.Interface(), nil
	}
	users := make(Users)
	err := json.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
