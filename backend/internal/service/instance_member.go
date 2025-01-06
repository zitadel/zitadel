package service

import (
	"context"

	"github.com/zitadel/zitadel/backend/internal/port"
)

type InstanceMember struct {
	Member
	Roles []InstanceMemberRole
}

type InstanceMemberRole uint8

const (
	InstanceMemberRoleUnspecified InstanceMemberRole = iota
	InstanceMemberRoleOwner
	InstanceMemberRoleAdmin
)

type InstanceMemberRepository interface {
	// CreateInstanceMember creates a new instance member
	CreateInstanceMember(ctx context.Context, executor port.Executor, member *InstanceMember) error
}

type InviteInstanceMemberRepository interface {
	// InviteInstanceMember creates a new invite for an instance admin
	InviteInstanceMember(ctx context.Context, executor port.Executor, roles ...InstanceMemberRole) (code string, err error)
}
