package domain

import "time"

type UserMembership struct {
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
