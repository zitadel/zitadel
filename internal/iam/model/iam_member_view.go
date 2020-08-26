package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type IamMemberView struct {
	UserID       string
	IamID        string
	UserName     string
	Email        string
	FirstName    string
	LastName     string
	DisplayName  string
	Roles        []string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Description  string
}

type IamMemberSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IamMemberSearchKey
	Asc           bool
	Queries       []*IamMemberSearchQuery
}

type IamMemberSearchKey int32

const (
	IamMemberSearchKeyUnspecified IamMemberSearchKey = iota
	IamMemberSearchKeyUserName
	IamMemberSearchKeyEmail
	IamMemberSearchKeyFirstName
	IamMemberSearchKeyLastName
	IamMemberSearchKeyIamID
	IamMemberSearchKeyUserID
)

type IamMemberSearchQuery struct {
	Key    IamMemberSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type IamMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IamMemberView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IamMemberSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
