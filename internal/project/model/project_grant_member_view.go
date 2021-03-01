package model

import (
	"github.com/caos/zitadel/internal/domain"
	"time"
)

type ProjectGrantMemberView struct {
	UserID       string
	GrantID      string
	ProjectID    string
	UserName     string
	Email        string
	FirstName    string
	LastName     string
	DisplayName  string
	Roles        []string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type ProjectGrantMemberSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ProjectGrantMemberSearchKey
	Asc           bool
	Queries       []*ProjectGrantMemberSearchQuery
}

type ProjectGrantMemberSearchKey int32

const (
	ProjectGrantMemberSearchKeyUnspecified ProjectGrantMemberSearchKey = iota
	ProjectGrantMemberSearchKeyUserName
	ProjectGrantMemberSearchKeyEmail
	ProjectGrantMemberSearchKeyFirstName
	ProjectGrantMemberSearchKeyLastName
	ProjectGrantMemberSearchKeyGrantID
	ProjectGrantMemberSearchKeyUserID
	ProjectGrantMemberSearchKeyProjectID
)

type ProjectGrantMemberSearchQuery struct {
	Key    ProjectGrantMemberSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type ProjectGrantMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectGrantMemberView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *ProjectGrantMemberSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
