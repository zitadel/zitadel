package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type OrgDomainView struct {
	OrgID          string
	CreationDate   time.Time
	ChangeDate     time.Time
	Domain         string
	Primary        bool
	Verified       bool
	ValidationType OrgDomainValidationType
}

type OrgDomainSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn OrgDomainSearchKey
	Asc           bool
	Queries       []*OrgDomainSearchQuery
}

type OrgDomainSearchKey int32

const (
	OrgDomainSearchKeyUnspecified OrgDomainSearchKey = iota
	OrgDomainSearchKeyDomain
	OrgDomainSearchKeyOrgID
	OrgDomainSearchKeyVerified
	OrgDomainSearchKeyPrimary
)

type OrgDomainSearchQuery struct {
	Key    OrgDomainSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type OrgDomainSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgDomainView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *OrgDomainSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-8fn7f", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}
