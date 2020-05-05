package model

import (
	"github.com/caos/zitadel/internal/model"
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
	PROJECTGRANTMEMBERSEARCHKEY_UNSPECIFIED ProjectGrantMemberSearchKey = iota
	PROJECTGRANTMEMBERSEARCHKEY_USER_NAME
	PROJECTGRANTMEMBERSEARCHKEY_EMAIL
	PROJECTGRANTMEMBERSEARCHKEY_FIRST_NAME
	PROJECTGRANTMEMBERSEARCHKEY_LAST_NAME
	PROJECTGRANTMEMBERSEARCHKEY_GRANT_ID
	PROJECTGRANTMEMBERSEARCHKEY_USER_ID
)

type ProjectGrantMemberSearchQuery struct {
	Key    ProjectGrantMemberSearchKey
	Method model.SearchMethod
	Value  string
}

type ProjectGrantMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectGrantMemberView
}
