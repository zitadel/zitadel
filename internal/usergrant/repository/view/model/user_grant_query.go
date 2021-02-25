package model

import (
	"github.com/caos/zitadel/internal/domain"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type UserGrantSearchRequest grant_model.UserGrantSearchRequest
type UserGrantSearchQuery grant_model.UserGrantSearchQuery
type UserGrantSearchKey grant_model.UserGrantSearchKey

func (req UserGrantSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req UserGrantSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req UserGrantSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == grant_model.UserGrantSearchKeyUnspecified {
		return nil
	}
	return UserGrantSearchKey(req.SortingColumn)
}

func (req UserGrantSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req UserGrantSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = UserGrantSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req UserGrantSearchQuery) GetKey() repository.ColumnKey {
	return UserGrantSearchKey(req.Key)
}

func (req UserGrantSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req UserGrantSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key UserGrantSearchKey) ToColumnName() string {
	switch grant_model.UserGrantSearchKey(key) {
	case grant_model.UserGrantSearchKeyUserID:
		return UserGrantKeyUserID
	case grant_model.UserGrantSearchKeyProjectID:
		return UserGrantKeyProjectID
	case grant_model.UserGrantSearchKeyState:
		return UserGrantKeyState
	case grant_model.UserGrantSearchKeyResourceOwner:
		return UserGrantKeyResourceOwner
	case grant_model.UserGrantSearchKeyGrantID:
		return UserGrantKeyGrantID
	case grant_model.UserGrantSearchKeyOrgName:
		return UserGrantKeyOrgName
	case grant_model.UserGrantSearchKeyRoleKey:
		return UserGrantKeyRole
	case grant_model.UserGrantSearchKeyID:
		return UserGrantKeyID
	case grant_model.UserGrantSearchKeyUserName:
		return UserGrantKeyUserName
	case grant_model.UserGrantSearchKeyFirstName:
		return UserGrantKeyFirstName
	case grant_model.UserGrantSearchKeyLastName:
		return UserGrantKeyLastName
	case grant_model.UserGrantSearchKeyEmail:
		return UserGrantKeyEmail
	case grant_model.UserGrantSearchKeyOrgDomain:
		return UserGrantKeyOrgDomain
	case grant_model.UserGrantSearchKeyProjectName:
		return UserGrantKeyProjectName
	case grant_model.UserGrantSearchKeyDisplayName:
		return UserGrantKeyDisplayName
	default:
		return ""
	}
}
