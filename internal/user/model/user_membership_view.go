package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type UserMembershipView struct {
	UserID      string
	MemberType  MemberType
	AggregateID string
	ObjectID    string

	Roles         []string
	DisplayName   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
}

type MemberType int32

const (
	MemberTypeUnspecified MemberType = iota
	MemberTypeOrganisation
	MemberTypeProject
	MemberTypeProjectGrant
)

type UserMembershipSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn UserMembershipSearchKey
	Asc           bool
	Queries       []*UserMembershipSearchQuery
}

type UserMembershipSearchKey int32

const (
	UserMembershipSearchKeyUnspecified UserMembershipSearchKey = iota
	UserMembershipSearchKeyUserID
	UserMembershipSearchKeyMemberType
	UserMembershipSearchKeyAggregateID
	UserMembershipSearchKeyObjectID
	UserMembershipSearchKeyResourceOwner
)

type UserMembershipSearchQuery struct {
	Key    UserMembershipSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type UserMembershipSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserMembershipView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *UserMembershipSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *UserMembershipSearchRequest) GetSearchQuery(key UserMembershipSearchKey) (int, *UserMembershipSearchQuery) {
	for i, q := range r.Queries {
		if q.Key == key {
			return i, q
		}
	}
	return -1, nil
}

func (r *UserMembershipSearchRequest) AppendResourceOwnerQuery(orgID string) {
	r.Queries = append(r.Queries, &UserMembershipSearchQuery{Key: UserMembershipSearchKeyResourceOwner, Method: model.SearchMethodEquals, Value: orgID})
}

func (r *UserMembershipSearchRequest) AppendUserIDQuery(userID string) {
	r.Queries = append(r.Queries, &UserMembershipSearchQuery{Key: UserMembershipSearchKeyUserID, Method: model.SearchMethodEquals, Value: userID})
}
