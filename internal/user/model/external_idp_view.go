package model

import (
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
)

type ExternalIDPSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ExternalIDPSearchKey
	Asc           bool
	Queries       []*ExternalIDPSearchQuery
}

type ExternalIDPSearchKey int32

const (
	ExternalIDPSearchKeyUnspecified ExternalIDPSearchKey = iota
	ExternalIDPSearchKeyExternalUserID
	ExternalIDPSearchKeyUserID
	ExternalIDPSearchKeyIdpConfigID
	ExternalIDPSearchKeyResourceOwner
	ExternalIDPSearchKeyInstanceID
)

type ExternalIDPSearchQuery struct {
	Key    ExternalIDPSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

func (r *ExternalIDPSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-3n8fM", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *ExternalIDPSearchRequest) AppendUserQuery(userID string) {
	r.Queries = append(r.Queries, &ExternalIDPSearchQuery{Key: ExternalIDPSearchKeyUserID, Method: domain.SearchMethodEquals, Value: userID})
}
