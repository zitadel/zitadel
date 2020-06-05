package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type OrgDomainView struct {
	OrgID        string
	CreationDate time.Time
	ChangeDate   time.Time
	Domain       string
	Primary      bool
	Verified     bool
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
	ORGDOMAINSEARCHKEY_UNSPECIFIED OrgDomainSearchKey = iota
	ORGDOMAINSEARCHKEY_DOMAIN
	ORGDOMAINSEARCHKEY_ORG_ID
)

type OrgDomainSearchQuery struct {
	Key    OrgDomainSearchKey
	Method model.SearchMethod
	Value  string
}

type OrgDomainSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgDomainView
}

func (r *OrgDomainSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
