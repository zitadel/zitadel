package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type MessageTextSearchRequest iam_model.MessageTextSearchRequest
type MessageTextSearchQuery iam_model.MessageTextSearchQuery
type MessageTextSearchKey iam_model.MessageTextSearchKey

func (req MessageTextSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req MessageTextSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req MessageTextSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.MessageTextSearchKeyUnspecified {
		return nil
	}
	return MessageTextSearchKey(req.SortingColumn)
}

func (req MessageTextSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req MessageTextSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = MessageTextSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req MessageTextSearchQuery) GetKey() repository.ColumnKey {
	return MessageTextSearchKey(req.Key)
}

func (req MessageTextSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req MessageTextSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key MessageTextSearchKey) ToColumnName() string {
	switch iam_model.MessageTextSearchKey(key) {
	case iam_model.MessageTextSearchKeyAggregateID:
		return MessageTextKeyAggregateID
	case iam_model.MessageTextSearchKeyMessageTextType:
		return MessageTextKeyMessageTextType
	case iam_model.MessageTextSearchKeyLanguage:
		return MessageTextKeyLanguage
	default:
		return ""
	}
}
