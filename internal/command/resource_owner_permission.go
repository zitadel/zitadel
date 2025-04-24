package command

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) CheckAggregatePermission(ctx context.Context, permission string, resourceID string) (string, error) {
	aggregateType, checkPermissionFunc, err := c.checkPermissionFunc(permission)
	if err != nil {
		return "", err
	}
	r := NewResourceOwnerModel(authz.GetInstance(ctx).InstanceID(), aggregateType, resourceID)
	err = c.eventstore.FilterToQueryReducer(ctx, r)
	if err != nil {
		return "", err
	}
	if r.resourceOwner == "" {
		return "", zerrors.ThrowPermissionDenied(nil, "COMMAND-4g3g2", "Errors.PermissionDenied")
	}
	return r.resourceOwner, checkPermissionFunc(ctx, resourceID, r.resourceOwner)
}

func (c *Commands) checkPermissionFunc(permission string) (eventstore.AggregateType, func(context.Context, string, string) error, error) {
	switch permission {
	case domain.PermissionUserWrite:
		return user.AggregateType, func(ctx context.Context, resourceID, resourceOwner string) error {
			return c.checkPermissionUpdateUser(ctx, resourceOwner, resourceID)
		}, nil
	default:
		return "", nil, zerrors.ThrowInternalf(nil, "COMMAND-4g3g2", "Permission %s not supported", permission)
	}
}
