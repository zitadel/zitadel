package model

import (
	"github.com/caos/zitadel/internal/domain"

	"time"
)

type OrgProjectMapping struct {
	OrgID     string
	ProjectID string
}

type OrgProjectMappingViewSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn OrgProjectMappingViewSearchKey
	Asc           bool
	Queries       []*OrgProjectMappingViewSearchQuery
}

type OrgProjectMappingViewSearchKey int32

const (
	OrgProjectMappingSearchKeyUnspecified OrgProjectMappingViewSearchKey = iota
	OrgProjectMappingSearchKeyProjectID
	OrgProjectMappingSearchKeyOrgID
)

type OrgProjectMappingViewSearchQuery struct {
	Key    OrgProjectMappingViewSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type OrgProjectMappingViewSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*OrgProjectMapping
	Sequence    uint64
	Timestamp   time.Time
}

func (r *OrgProjectMappingViewSearchRequest) GetSearchQuery(key OrgProjectMappingViewSearchKey) (int, *OrgProjectMappingViewSearchQuery) {
	for i, q := range r.Queries {
		if q.Key == key {
			return i, q
		}
	}
	return -1, nil
}
