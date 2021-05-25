package model

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type RefreshTokenSearchRequest model.RefreshTokenSearchRequest
type RefreshTokenSearchQuery model.RefreshTokenSearchQuery
type RefreshTokenSearchKey model.RefreshTokenSearchKey

func (req RefreshTokenSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req RefreshTokenSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req RefreshTokenSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == model.RefreshTokenSearchKeyUnspecified {
		return nil
	}
	return RefreshTokenSearchKey(req.SortingColumn)
}

func (req RefreshTokenSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req RefreshTokenSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = RefreshTokenSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req RefreshTokenSearchQuery) GetKey() repository.ColumnKey {
	return RefreshTokenSearchKey(req.Key)
}

func (req RefreshTokenSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req RefreshTokenSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key RefreshTokenSearchKey) ToColumnName() string {
	switch model.RefreshTokenSearchKey(key) {
	case model.RefreshTokenSearchKeyRefreshTokenID:
		return RefreshTokenKeyTokenID
	case model.RefreshTokenSearchKeyUserAgentID:
		return RefreshTokenKeyUserAgentID
	case model.RefreshTokenSearchKeyUserID:
		return RefreshTokenKeyUserID
	case model.RefreshTokenSearchKeyApplicationID:
		return RefreshTokenKeyApplicationID
	case model.RefreshTokenSearchKeyExpiration:
		return RefreshTokenKeyExpiration
	case model.RefreshTokenSearchKeyResourceOwner:
		return RefreshTokenKeyResourceOwner
	default:
		return ""
	}
}
