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
	return nil
}

// RemoveMember implements [domain.MemberRepository].
func (o *orgMember) RemoveMember(ctx context.Context, orgID string, userID string) error {
	return nil
}

// SetMemberRoles implements [domain.MemberRepository].
func (o *orgMember) SetMemberRoles(ctx context.Context, orgID string, userID string, roles []string) error {
	return nil
}

var _ domain.MemberRepository = (*orgMember)(nil)
