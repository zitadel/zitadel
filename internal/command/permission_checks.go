package command

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/v2/user"
)

func (c *Commands) CheckPermission(ctx context.Context, permission string, aggregateType eventstore.AggregateType) eventstore.PermissionCheck {
	return func(resourceOwner, aggregateID string) error {
		// For example if a write model didn't query any events, the resource owner is probably empty.
		// In this case, we have to query an event on the given aggregate to get the resource owner.
		if resourceOwner == "" {
			r := NewResourceOwnerModel(authz.GetInstance(ctx).InstanceID(), aggregateType, aggregateID)
			err := c.eventstore.FilterToQueryReducer(ctx, r)
			if err != nil {
				return err
			}
			resourceOwner = r.resourceOwner
		}
		return c.checkPermission(ctx, permission, resourceOwner, aggregateID)
	}
}

func (c *Commands) checkPermissionOnUser(ctx context.Context, permission string, allowSelf bool) eventstore.PermissionCheck {
	return func(resourceOwner, aggregateID string) error {
		if allowSelf && aggregateID != "" && aggregateID == authz.GetCtxData(ctx).UserID {
			return nil
		}
		return c.CheckPermission(ctx, permission, user.AggregateType)(resourceOwner, aggregateID)
	}
}

func (c *Commands) CheckPermissionUserWrite(ctx context.Context, allowSelf bool) eventstore.PermissionCheck {
	return c.checkPermissionOnUser(ctx, domain.PermissionUserWrite, allowSelf)
}

func (c *Commands) CheckPermissionUserDelete(ctx context.Context, allowSelf bool) eventstore.PermissionCheck {
	return c.checkPermissionOnUser(ctx, domain.PermissionUserDelete, allowSelf)
}

// Deprecated: use CheckPermissionUserWrite and set the returned PermissionCheck to the eventstore.WriteModel to safely protect an API.
func (c *Commands) checkPermissionUpdateUser(ctx context.Context, resourceOwner, userID string) error {
	return c.CheckPermissionUserWrite(ctx, true)(resourceOwner, userID)
}

// Deprecated: use CheckPermission, allow self and use the returned function as PermissionCheck in eventstore.WriteModel to protect an API
func (c *Commands) checkPermissionUpdateUserCredentials(ctx context.Context, resourceOwner, userID string) error {
	if userID != "" && userID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := c.checkPermission(ctx, domain.PermissionUserCredentialWrite, resourceOwner, userID); err != nil {
		return err
	}
	return nil
}
