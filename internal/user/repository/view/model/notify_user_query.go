package model

import (
	"github.com/caos/zitadel/internal/domain"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type NotifyUserSearchRequest usr_model.NotifyUserSearchRequest
type NotifyUserSearchQuery usr_model.NotifyUserSearchQuery
type NotifyUserSearchKey usr_model.NotifyUserSearchKey

func (req NotifyUserSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req NotifyUserSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req NotifyUserSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == usr_model.NotifyUserSearchKeyUnspecified {
		return nil
	}
	return NotifyUserSearchKey(req.SortingColumn)
}

func (req NotifyUserSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req NotifyUserSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = NotifyUserSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req NotifyUserSearchQuery) GetKey() repository.ColumnKey {
	return NotifyUserSearchKey(req.Key)
}

func (req NotifyUserSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req NotifyUserSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key NotifyUserSearchKey) ToColumnName() string {
	switch usr_model.NotifyUserSearchKey(key) {
	case usr_model.NotifyUserSearchKeyUserID:
		return NotifyUserKeyUserID
	case usr_model.NotifyUserSearchKeyResourceOwner:
		return NotifyUserKeyResourceOwner
	default:
		return ""
	}
}
