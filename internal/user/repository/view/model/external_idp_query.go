package model

import (
	"github.com/caos/zitadel/internal/domain"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type ExternalIDPSearchRequest usr_model.ExternalIDPSearchRequest
type ExternalIDPSearchQuery usr_model.ExternalIDPSearchQuery
type ExternalIDPSearchKey usr_model.ExternalIDPSearchKey

func (req ExternalIDPSearchRequest) GetLimit() uint64 {
	return req.Limit
}

func (req ExternalIDPSearchRequest) GetOffset() uint64 {
	return req.Offset
}

func (req ExternalIDPSearchRequest) GetSortingColumn() repository.ColumnKey {
	if req.SortingColumn == usr_model.ExternalIDPSearchKeyUnspecified {
		return nil
	}
	return ExternalIDPSearchKey(req.SortingColumn)
}

func (req ExternalIDPSearchRequest) GetAsc() bool {
	return req.Asc
}

func (req ExternalIDPSearchRequest) GetQueries() []repository.SearchQuery {
	result := make([]repository.SearchQuery, len(req.Queries))
	for i, q := range req.Queries {
		result[i] = ExternalIDPSearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
	}
	return result
}

func (req ExternalIDPSearchQuery) GetKey() repository.ColumnKey {
	return ExternalIDPSearchKey(req.Key)
}

func (req ExternalIDPSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req ExternalIDPSearchQuery) GetValue() interface{} {
	return req.Value
}

func (key ExternalIDPSearchKey) ToColumnName() string {
	switch usr_model.ExternalIDPSearchKey(key) {
	case usr_model.ExternalIDPSearchKeyExternalUserID:
		return ExternalIDPKeyExternalUserID
	case usr_model.ExternalIDPSearchKeyUserID:
		return ExternalIDPKeyUserID
	case usr_model.ExternalIDPSearchKeyIdpConfigID:
		return ExternalIDPKeyIDPConfigID
	case usr_model.ExternalIDPSearchKeyResourceOwner:
		return ExternalIDPKeyResourceOwner
	default:
		return ""
	}
}
