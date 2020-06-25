package model

import (
	global_model "github.com/caos/zitadel/internal/model"
	token_model "github.com/caos/zitadel/internal/token/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type TokenSearchRequest token_model.TokenSearchRequest
type TokenSearchQuery token_model.TokenSearchQuery
type TokenSearchKey token_model.TokenSearchKey

func (req TokenSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req TokenSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req TokenSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == token_model.TokenSearchKeyUnspecified {
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

func (req TokenSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req TokenSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key TokenSearchKey) ToColumnName() string {
	switch token_model.TokenSearchKey(key) {
	case token_model.TokenSearchKeyTokenID:
		return TokenKeyTokenID
	case token_model.TokenSearchKeyUserAgentID:
		return TokenKeyUserAgentID
	case token_model.TokenSearchKeyUserID:
		return TokenKeyUserID
	case token_model.TokenSearchKeyApplicationID:
		return TokenKeyApplicationID
	case token_model.TokenSearchKeyExpiration:
		return TokenKeyExpiration
	case token_model.TokenSearchKeyResourceOwner:
		return TokenKeyResourceOwner
	default:
		return ""
	}
}
