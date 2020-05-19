package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type OrgView struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	State         OrgState
	ResourceOwner string
	Sequence      uint64

	Name   string
	Domain string
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
	ORGSEARCHKEY_UNSPECIFIED OrgSearchKey = iota
	ORGSEARCHKEY_ORG_ID
	ORGSEARCHKEY_ORG_NAME
	ORGSEARCHKEY_ORG_DOMAIN
	ORGSEARCHKEY_STATE
	ORGSEARCHKEY_RESOURCEOWNER
)

type OrgSearchQuery struct {
	Key    OrgSearchKey
	Method model.SearchMethod
	Value  string
}

type OrgSearchResult struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgView
}

func (r *OrgSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
