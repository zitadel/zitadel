package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view"
)

type UserSearchRequest usr_model.UserSearchRequest
type UserSearchQuery usr_model.UserSearchQuery
type UserSearchKey usr_model.UserSearchKey

func (req UserSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req UserSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req UserSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == usr_model.USERSEARCHKEY_UNSPECIFIED {
		return nil
	}
	return UserSearchKey(req.SortingColumn)
}

func (req UserSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req UserSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = UserSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req UserSearchQuery) GetKey() view.ColumnKey {
	return UserSearchKey(req.Key)
}

func (req UserSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req UserSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key UserSearchKey) ToColumnName() string {
	switch usr_model.UserSearchKey(key) {
	case usr_model.USERSEARCHKEY_USER_ID:
		return UserKeyUserID
	case usr_model.USERSEARCHKEY_USER_NAME:
		return UserKeyUserName
	case usr_model.USERSEARCHKEY_FIRST_NAME:
		return UserKeyFirstName
	case usr_model.USERSEARCHKEY_LAST_NAME:
		return UserKeyLastName
	case usr_model.USERSEARCHKEY_DISPLAY_NAME:
		return UserKeyDisplayName
	case usr_model.USERSEARCHKEY_NICK_NAME:
		return UserKeyNickName
	case usr_model.USERSEARCHKEY_EMAIL:
		return UserKeyEmail
	case usr_model.USERSEARCHKEY_STATE:
		return UserKeyState
	case usr_model.USERSEARCHKEY_RESOURCEOWNER:
		return UserKeyResourceOwner
	default:
		return ""
	}
}
