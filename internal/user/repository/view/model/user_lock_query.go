package model

import (
	"github.com/caos/zitadel/internal/domain"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type UserLockSearchRequest usr_model.UserLockSearchRequest
type UserLockSearchQuery usr_model.UserLockSearchQuery
type UserLockSearchKey usr_model.UserLockSearchKey

func (req UserLockSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req UserLockSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req UserLockSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == usr_model.UserLockSearchKeyUnspecified {
		return nil
	}
	return UserLockSearchKey(req.SortingColumn)
}

func (req UserLockSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req UserLockSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = UserLockSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req UserLockSearchQuery) GetKey() repository.ColumnKey {
	return UserLockSearchKey(req.Key)
}

func (req UserLockSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req UserLockSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key UserLockSearchKey) ToColumnName() string {
	switch usr_model.UserLockSearchKey(key) {
	case usr_model.UserLockSearchKeyUserID:
		return UserLockedKeyUserID
	default:
		return ""
	}
}
