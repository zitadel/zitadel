package command

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) CheckPermission(ctx context.Context, permission string, resourceID string, allowSelf bool) (string, error) {
	// TODO: I assume we need to pass an eventstore.AggregateType to NewResourceOwnerModel for performance reasons.
	// Is there a way to change this, so we can get rid of the mapping in permissionAggregateType?
	aggregateType, err := c.permissionAggregateType(permission)
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
	if allowSelf && resourceID != "" && resourceID == authz.GetCtxData(ctx).UserID {
		return r.resourceOwner, nil
	}
	return r.resourceOwner, c.checkPermission(ctx, permission, r.resourceOwner, resourceID)
}

func (c *Commands) permissionAggregateType(permission string) (eventstore.AggregateType, error) {
	switch permission {
	case domain.PermissionUserWrite:
		return user.AggregateType, nil
	default:
		return "", zerrors.ThrowInternalf(nil, "COMMAND-4g3g2", "Permission %s not supported", permission)
	}
}
