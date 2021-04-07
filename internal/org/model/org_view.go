package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type OrgView struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	State         OrgState
	ResourceOwner string
	Sequence      uint64

	Name string
}

type OrgSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn OrgSearchKey
	Asc           bool
	Queries       []*OrgSearchQuery
}

type OrgSearchKey int32

const (
	OrgSearchKeyUnspecified OrgSearchKey = iota
	OrgSearchKeyOrgID
	OrgSearchKeyOrgName
	OrgSearchKeyOrgDomain
	OrgSearchKeyState
	OrgSearchKeyResourceOwner
	OrgSearchKeyOrgNameLower //used for lowercase search
)

type OrgSearchQuery struct {
	Key    OrgSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type OrgSearchResult struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *OrgSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-200ds", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func OrgViewToOrg(o *OrgView) *Org {
	return &Org{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   o.ID,
			ChangeDate:    o.ChangeDate,
			CreationDate:  o.CreationDate,
			ResourceOwner: o.ResourceOwner,
			Sequence:      o.Sequence,
		},
		Name:  o.Name,
		State: o.State,
	}
}
