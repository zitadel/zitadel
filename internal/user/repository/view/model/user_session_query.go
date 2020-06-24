package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view"
)

type UserSessionSearchRequest usr_model.UserSessionSearchRequest
type UserSessionSearchQuery usr_model.UserSessionSearchQuery
type UserSessionSearchKey usr_model.UserSessionSearchKey

func (req UserSessionSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req UserSessionSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req UserSessionSearchRequest) GetSortingColumn() view.ColumnKey {
	if req.SortingColumn == usr_model.UserSessionSearchKeyUnspecified {
		return nil
	}
	return UserSessionSearchKey(req.SortingColumn)
}

func (req UserSessionSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req UserSessionSearchRequest) GetQueries() []view.SearchQuery {
	result := make([]view.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = UserSessionSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req UserSessionSearchQuery) GetKey() view.ColumnKey {
	return UserSessionSearchKey(req.Key)
}

func (req UserSessionSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req UserSessionSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key UserSessionSearchKey) ToColumnName() string {
	switch usr_model.UserSessionSearchKey(key) {
	case usr_model.UserSessionSearchKeyUserAgentID:
		return UserSessionKeyUserAgentID
	case usr_model.UserSessionSearchKeyUserID:
		return UserSessionKeyUserID
	case usr_model.UserSessionSearchKeyState:
		return UserSessionKeyState
	case usr_model.UserSessionSearchKeyResourceOwner:
		return UserSessionKeyResourceOwner
	default:
		return ""
	}
}
