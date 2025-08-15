package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// CreateGroup creates a new user group in an organization
func (c *Commands) CreateGroup(ctx context.Context, group *domain.Group) (details *domain.ObjectDetails, err error) {
	return nil, zerrors.ThrowUnimplemented(nil, "GRP-Dmif", "Not implemented")
}

// UpdateGroup updates a user group
func (c *Commands) UpdateGroup(ctx context.Context, group *domain.Group) (details *domain.ObjectDetails, err error) {
	return nil, zerrors.ThrowUnimplemented(nil, "GRP-Dmif", "Not implemented")
}

// DeleteGroup deletes a user group from an organization
func (c *Commands) DeleteGroup(ctx context.Context, groupID string) (details *domain.ObjectDetails, err error) {
	return nil, zerrors.ThrowUnimplemented(nil, "GRP-Dmif", "Not implemented")
}
