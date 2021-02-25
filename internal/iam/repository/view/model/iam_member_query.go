package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type IAMMemberSearchRequest iam_model.IAMMemberSearchRequest
type IAMMemberSearchQuery iam_model.IAMMemberSearchQuery
type IAMMemberSearchKey iam_model.IAMMemberSearchKey

func (req IAMMemberSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req IAMMemberSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req IAMMemberSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.IAMMemberSearchKeyUnspecified {
		return nil
	}
	return IAMMemberSearchKey(req.SortingColumn)
}

func (req IAMMemberSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req IAMMemberSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = IAMMemberSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req IAMMemberSearchQuery) GetKey() repository.ColumnKey {
	return IAMMemberSearchKey(req.Key)
}

func (req IAMMemberSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req IAMMemberSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key IAMMemberSearchKey) ToColumnName() string {
	switch iam_model.IAMMemberSearchKey(key) {
	case iam_model.IAMMemberSearchKeyEmail:
		return IAMMemberKeyEmail
	case iam_model.IAMMemberSearchKeyFirstName:
		return IAMMemberKeyFirstName
	case iam_model.IAMMemberSearchKeyLastName:
		return IAMMemberKeyLastName
	case iam_model.IAMMemberSearchKeyUserName:
		return IAMMemberKeyUserName
	case iam_model.IAMMemberSearchKeyUserID:
		return IAMMemberKeyUserID
	case iam_model.IAMMemberSearchKeyIamID:
		return IAMMemberKeyIamID
	default:
		return ""
	}
}
