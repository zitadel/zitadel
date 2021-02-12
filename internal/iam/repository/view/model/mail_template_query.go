package model

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	global_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type MailTemplateSearchRequest iam_model.MailTemplateSearchRequest
type MailTemplateSearchQuery iam_model.MailTemplateSearchQuery
type MailTemplateSearchKey iam_model.MailTemplateSearchKey

func (req MailTemplateSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req MailTemplateSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req MailTemplateSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.MailTemplateSearchKeyUnspecified {
		return nil
	}
	return MailTemplateSearchKey(req.SortingColumn)
}

func (req MailTemplateSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req MailTemplateSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = MailTemplateSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req MailTemplateSearchQuery) GetKey() repository.ColumnKey {
	return MailTemplateSearchKey(req.Key)
}

func (req MailTemplateSearchQuery) GetMethod() global_model.SearchMethod {
	return req.Method
}

func (req MailTemplateSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key MailTemplateSearchKey) ToColumnName() string {
	switch iam_model.MailTemplateSearchKey(key) {
	case iam_model.MailTemplateSearchKeyAggregateID:
		return MailTemplateKeyAggregateID
	default:
		return ""
	}
}
