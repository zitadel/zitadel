package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type GrantedProjectView struct {
	ProjectID     string
	Name          string
	CreationDate  time.Time
	ChangeDate    time.Time
	State         ProjectState
	Type          ProjectType
	ResourceOwner string
	OrgID         string
	OrgName       string
	OrgDomain     string
	Sequence      uint64
	GrantID       string
}

type ProjectType int32

const (
	PROJECTTYPE_OWNED ProjectType = iota
	PROJECTTYPE_GRANTED
)

type GrantedProjectSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn GrantedProjectSearchKey
	Asc           bool
	Queries       []*GrantedProjectSearchQuery
}

type GrantedProjectSearchKey int32

const (
	GRANTEDPROJECTSEARCHKEY_UNSPECIFIED GrantedProjectSearchKey = iota
	GRANTEDPROJECTSEARCHKEY_NAME
	GRANTEDPROJECTSEARCHKEY_PROJECTID
	GRANTEDPROJECTSEARCHKEY_GRANTID
	GRANTEDPROJECTSEARCHKEY_ORGID
)

type GrantedProjectSearchQuery struct {
	Key    GrantedProjectSearchKey
	Method model.SearchMethod
	Value  string
}

type GrantedProjectSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*GrantedProjectView
}
