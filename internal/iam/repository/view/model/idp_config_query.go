package model

import (
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type IDPConfigSearchRequest iam_model.IDPConfigSearchRequest
type IDPConfigSearchQuery iam_model.IDPConfigSearchQuery
type IDPConfigSearchKey iam_model.IDPConfigSearchKey

func (req IDPConfigSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req IDPConfigSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req IDPConfigSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == iam_model.IDPConfigSearchKeyUnspecified {
		return nil
	}
	return IDPConfigSearchKey(req.SortingColumn)
}

func (req IDPConfigSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req IDPConfigSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = IDPConfigSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req IDPConfigSearchQuery) GetKey() repository.ColumnKey {
	return IDPConfigSearchKey(req.Key)
}

func (req IDPConfigSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req IDPConfigSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key IDPConfigSearchKey) ToColumnName() string {
	switch iam_model.IDPConfigSearchKey(key) {
	case iam_model.IDPConfigSearchKeyAggregateID:
		return IDPConfigKeyAggregateID
	case iam_model.IDPConfigSearchKeyIdpConfigID:
		return IDPConfigKeyIdpConfigID
	case iam_model.IDPConfigSearchKeyName:
		return IDPConfigKeyName
	case iam_model.IDPConfigSearchKeyIdpProviderType:
		return IDPConfigKeyProviderType
	case iam_model.IDPConfigSearchKeyAuthConnectorMachineID:
		return IDPConfigKeyAuthConnectorMachineID
	default:
		return ""
	}
}
