package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type IAMMemberView struct {
	UserID             string
	IAMID              string
	UserName           string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	PreferredLoginName string
	Roles              []string
	CreationDate       time.Time
	ChangeDate         time.Time
	Sequence           uint64
}

type IAMMemberSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IAMMemberSearchKey
	Asc           bool
	Queries       []*IAMMemberSearchQuery
}

type IAMMemberSearchKey int32

const (
	IAMMemberSearchKeyUnspecified IAMMemberSearchKey = iota
	IAMMemberSearchKeyUserName
	IAMMemberSearchKeyEmail
	IAMMemberSearchKeyFirstName
	IAMMemberSearchKeyLastName
	IAMMemberSearchKeyIamID
	IAMMemberSearchKeyUserID
)

type IAMMemberSearchQuery struct {
	Key    IAMMemberSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type IAMMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IAMMemberView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IAMMemberSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-8fn7f", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
