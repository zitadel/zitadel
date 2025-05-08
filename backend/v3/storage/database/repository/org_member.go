package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
)

type orgMember struct {
	*org
}

// AddMember implements [domain.MemberRepository].
func (o *orgMember) AddMember(ctx context.Context, orgID string, userID string, roles []string) error {
	panic("unimplemented")
}

// RemoveMember implements [domain.MemberRepository].
func (o *orgMember) RemoveMember(ctx context.Context, orgID string, userID string) error {
	panic("unimplemented")
}

// SetMemberRoles implements [domain.MemberRepository].
func (o *orgMember) SetMemberRoles(ctx context.Context, orgID string, userID string, roles []string) error {
	panic("unimplemented")
}

var _ domain.MemberRepository = (*orgMember)(nil)
