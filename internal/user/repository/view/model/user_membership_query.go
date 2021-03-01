package model

import (
	"github.com/caos/zitadel/internal/domain"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type UserMembershipSearchRequest usr_model.UserMembershipSearchRequest
type UserMembershipSearchQuery usr_model.UserMembershipSearchQuery
type UserMembershipSearchKey usr_model.UserMembershipSearchKey

func (req UserMembershipSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req UserMembershipSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req UserMembershipSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == usr_model.UserMembershipSearchKeyUnspecified {
		return nil
	}
	return UserMembershipSearchKey(req.SortingColumn)
}

func (req UserMembershipSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req UserMembershipSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = UserMembershipSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req UserMembershipSearchQuery) GetKey() repository.ColumnKey {
	return UserMembershipSearchKey(req.Key)
}

func (req UserMembershipSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req UserMembershipSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key UserMembershipSearchKey) ToColumnName() string {
	switch usr_model.UserMembershipSearchKey(key) {
	case usr_model.UserMembershipSearchKeyUserID:
		return UserMembershipKeyUserID
	case usr_model.UserMembershipSearchKeyResourceOwner:
		return UserMembershipKeyResourceOwner
	case usr_model.UserMembershipSearchKeyMemberType:
		return UserMembershipKeyMemberType
	case usr_model.UserMembershipSearchKeyAggregateID:
		return UserMembershipKeyAggregateID
	case usr_model.UserMembershipSearchKeyObjectID:
		return UserMembershipKeyObjectID
	default:
		return ""
	}
}
