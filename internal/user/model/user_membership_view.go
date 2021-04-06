package model

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"

	"time"
)

type UserMembershipView struct {
	UserID      string
	MemberType  MemberType
	AggregateID string
	//ObjectID differs from aggregate id if obejct is sub of an aggregate
	ObjectID string

	Roles             []string
	DisplayName       string
	CreationDate      time.Time
	ChangeDate        time.Time
	ResourceOwner     string
	ResourceOwnerName string
	Sequence          uint64
}

type MemberType int32

const (
	MemberTypeUnspecified MemberType = iota
	MemberTypeOrganisation
	MemberTypeProject
	MemberTypeProjectGrant
	MemberTypeIam
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
	Method domain.SearchMethod
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

func (r *UserMembershipSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return caos_errors.ThrowInvalidArgument(nil, "SEARCH-8fn7f", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *UserMembershipSearchRequest) GetSearchQuery(key UserMembershipSearchKey) (int, *UserMembershipSearchQuery) {
	for i, q := range r.Queries {
		if q.Key == key {
			return i, q
		}
	}
	return -1, nil
}

func (r *UserMembershipSearchRequest) AppendResourceOwnerAndIamQuery(orgID, iamID string) {
	r.Queries = append(r.Queries, &UserMembershipSearchQuery{Key: UserMembershipSearchKeyResourceOwner, Method: domain.SearchMethodIsOneOf, Value: []string{orgID, iamID}})
}

func (r *UserMembershipSearchRequest) AppendUserIDQuery(userID string) {
	r.Queries = append(r.Queries, &UserMembershipSearchQuery{Key: UserMembershipSearchKeyUserID, Method: domain.SearchMethodEquals, Value: userID})
}
