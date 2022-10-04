package model

import (
	"github.com/zitadel/zitadel/internal/domain"

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
	UserMembershipSearchKeyInstanceID
)

type UserMembershipSearchQuery struct {
	Key    UserMembershipSearchKey
	Method domain.SearchMethod
	Value  interface{}
}
