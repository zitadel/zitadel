package model

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type MailTextSearchRequest iam_model.MailTextSearchRequest
type MailTextSearchQuery iam_model.MailTextSearchQuery
type MailTextSearchKey iam_model.MailTextSearchKey

func (req MailTextSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req MailTextSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req MailTextSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.MailTextSearchKeyUnspecified {
		return nil
	}
	return MailTextSearchKey(req.SortingColumn)
}

func (req MailTextSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req MailTextSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = MailTextSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req MailTextSearchQuery) GetKey() repository.ColumnKey {
	return MailTextSearchKey(req.Key)
}

func (req MailTextSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req MailTextSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key MailTextSearchKey) ToColumnName() string {
	switch iam_model.MailTextSearchKey(key) {
	case iam_model.MailTextSearchKeyAggregateID:
		return MailTextKeyAggregateID
	case iam_model.MailTextSearchKeyMailTextType:
		return MailTextKeyMailTextType
	case iam_model.MailTextSearchKeyLanguage:
		return MailTextKeyLanguage
	default:
		return ""
	}
}
