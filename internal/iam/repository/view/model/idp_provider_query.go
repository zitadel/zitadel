package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type IDPProviderSearchRequest iam_model.IDPProviderSearchRequest
type IDPProviderSearchQuery iam_model.IDPProviderSearchQuery
type IDPProviderSearchKey iam_model.IDPProviderSearchKey

func (req IDPProviderSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req IDPProviderSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req IDPProviderSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.IDPProviderSearchKeyUnspecified {
		return nil
	}
	return IDPProviderSearchKey(req.SortingColumn)
}

func (req IDPProviderSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req IDPProviderSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = IDPProviderSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req IDPProviderSearchQuery) GetKey() repository.ColumnKey {
	return IDPProviderSearchKey(req.Key)
}

func (req IDPProviderSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req IDPProviderSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key IDPProviderSearchKey) ToColumnName() string {
	switch iam_model.IDPProviderSearchKey(key) {
	case iam_model.IDPProviderSearchKeyAggregateID:
		return IDPProviderKeyAggregateID
	case iam_model.IDPProviderSearchKeyIdpConfigID:
		return IDPProviderKeyIdpConfigID
	case iam_model.IDPProviderSearchKeyState:
		return IDPProviderKeyState
	default:
		return ""
	}
}
