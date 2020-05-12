package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type ProjectMemberView struct {
	UserID       string
	ProjectID    string
	UserName     string
	Email        string
	FirstName    string
	LastName     string
	Roles        []string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type ProjectMemberSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ProjectMemberSearchKey
	Asc           bool
	Queries       []*ProjectMemberSearchQuery
}

type ProjectMemberSearchKey int32

const (
	PROJECTMEMBERSEARCHKEY_UNSPECIFIED ProjectMemberSearchKey = iota
	PROJECTMEMBERSEARCHKEY_USER_NAME
	PROJECTMEMBERSEARCHKEY_EMAIL
	PROJECTMEMBERSEARCHKEY_FIRST_NAME
	PROJECTMEMBERSEARCHKEY_LAST_NAME
	PROJECTMEMBERSEARCHKEY_PROJECT_ID
	PROJECTMEMBERSEARCHKEY_USER_ID
)

type ProjectMemberSearchQuery struct {
	Key    ProjectMemberSearchKey
	Method model.SearchMethod
	Value  string
}

type ProjectMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectMemberView
}

func (r *ProjectMemberSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
