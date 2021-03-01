package model

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type TokenSearchRequest model.TokenSearchRequest
type TokenSearchQuery model.TokenSearchQuery
type TokenSearchKey model.TokenSearchKey

func (req TokenSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req TokenSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req TokenSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == model.TokenSearchKeyUnspecified {
		return nil
	}
	return TokenSearchKey(req.SortingColumn)
}

func (req TokenSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req TokenSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = TokenSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req TokenSearchQuery) GetKey() repository.ColumnKey {
	return TokenSearchKey(req.Key)
}

func (req TokenSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req TokenSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key TokenSearchKey) ToColumnName() string {
	switch model.TokenSearchKey(key) {
	case model.TokenSearchKeyTokenID:
		return TokenKeyTokenID
	case model.TokenSearchKeyUserAgentID:
		return TokenKeyUserAgentID
	case model.TokenSearchKeyUserID:
		return TokenKeyUserID
	case model.TokenSearchKeyApplicationID:
		return TokenKeyApplicationID
	case model.TokenSearchKeyExpiration:
		return TokenKeyExpiration
	case model.TokenSearchKeyResourceOwner:
		return TokenKeyResourceOwner
	default:
		return ""
	}
}
