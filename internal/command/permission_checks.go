package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/v2/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type PermissionCheck func(resourceOwner, aggregateID string) error

type UserGrantPermissionCheck func(projectID, projectGrantID string) PermissionCheck

func (c *Commands) newPermissionCheck(ctx context.Context, permission string, aggregateType eventstore.AggregateType) PermissionCheck {
	return func(resourceOwner, aggregateID string) error {
		if aggregateID == "" {
			return zerrors.ThrowInternal(nil, "COMMAND-ulBlS", "Errors.IDMissing")
		}
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
		if resourceOwner == "" {
			return zerrors.ThrowNotFound(nil, "COMMAND-4g3xq", "Errors.NotFound")
		}
		return c.checkPermission(ctx, permission, resourceOwner, aggregateID)
	}
}

func (c *Commands) checkPermissionOnUser(ctx context.Context, permission string) PermissionCheck {
	return func(resourceOwner, aggregateID string) error {
		if aggregateID != "" && aggregateID == authz.GetCtxData(ctx).UserID {
			return nil
		}
		return c.newPermissionCheck(ctx, permission, user.AggregateType)(resourceOwner, aggregateID)
	}
}

func (c *Commands) NewPermissionCheckUserWrite(ctx context.Context) PermissionCheck {
	return c.checkPermissionOnUser(ctx, domain.PermissionUserWrite)
}

func (c *Commands) checkPermissionDeleteUser(ctx context.Context, resourceOwner, userID string) error {
	return c.checkPermissionOnUser(ctx, domain.PermissionUserDelete)(resourceOwner, userID)
}

func (c *Commands) checkPermissionUpdateUser(ctx context.Context, resourceOwner, userID string) error {
	return c.NewPermissionCheckUserWrite(ctx)(resourceOwner, userID)
}

func (c *Commands) checkPermissionUpdateUserCredentials(ctx context.Context, resourceOwner, userID string) error {
	return c.checkPermissionOnUser(ctx, domain.PermissionUserCredentialWrite)(resourceOwner, userID)
}

func (c *Commands) checkPermissionCreateProject(ctx context.Context, resourceOwner, projectID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionProjectCreate, project.AggregateType)(resourceOwner, projectID)
}

func (c *Commands) checkPermissionDeleteProject(ctx context.Context, resourceOwner, projectID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionProjectDelete, project.AggregateType)(resourceOwner, projectID)
}

func (c *Commands) checkPermissionUpdateProject(ctx context.Context, resourceOwner, projectID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionProjectWrite, project.AggregateType)(resourceOwner, projectID)
}

func (c *Commands) checkPermissionUpdateProjectGrant(ctx context.Context, resourceOwner, projectID, projectGrantID string) (err error) {
	if err := c.newPermissionCheck(ctx, domain.PermissionProjectGrantWrite, project.AggregateType)(resourceOwner, projectGrantID); err != nil {
		if err := c.newPermissionCheck(ctx, domain.PermissionProjectGrantWrite, project.AggregateType)(resourceOwner, projectID); err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) checkPermissionDeleteProjectGrant(ctx context.Context, resourceOwner, projectID, projectGrantID string) (err error) {
	if err := c.newPermissionCheck(ctx, domain.PermissionProjectGrantDelete, project.AggregateType)(resourceOwner, projectGrantID); err != nil {
		if err := c.newPermissionCheck(ctx, domain.PermissionProjectGrantDelete, project.AggregateType)(resourceOwner, projectID); err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) checkPermissionUpdateApplication(ctx context.Context, resourceOwner, appID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionProjectAppWrite, project.AggregateType)(resourceOwner, appID)
}

func (c *Commands) checkPermissionDeleteApp(ctx context.Context, resourceOwner, appID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionProjectAppDelete, project.AggregateType)(resourceOwner, appID)
}

func (c *Commands) checkPermissionUpdateInstanceMember(ctx context.Context, instanceID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionInstanceMemberWrite, instance.AggregateType)(instanceID, instanceID)
}

func (c *Commands) checkPermissionDeleteInstanceMember(ctx context.Context, instanceID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionInstanceMemberDelete, instance.AggregateType)(instanceID, instanceID)
}

func (c *Commands) checkPermissionUpdateOrgMember(ctx context.Context, instanceID, orgID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionOrgMemberWrite, org.AggregateType)(instanceID, orgID)
}
func (c *Commands) checkPermissionDeleteOrgMember(ctx context.Context, instanceID, orgID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionOrgMemberDelete, org.AggregateType)(instanceID, orgID)
}

func (c *Commands) checkPermissionUpdateProjectMember(ctx context.Context, resourceOwner, projectID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionProjectMemberWrite, project.AggregateType)(resourceOwner, projectID)
}

func (c *Commands) checkPermissionDeleteProjectMember(ctx context.Context, resourceOwner, projectID string) error {
	return c.newPermissionCheck(ctx, domain.PermissionProjectMemberDelete, project.AggregateType)(resourceOwner, projectID)
}

func (c *Commands) checkPermissionUpdateProjectGrantMember(ctx context.Context, grantedOrgID, projectGrantID string) (err error) {
	// TODO: add permission check for project grant owners
	//if err := c.newPermissionCheck(ctx, domain.PermissionProjectGrantMemberWrite, project.AggregateType)(resourceOwner, projectGrantID); err != nil {
	return c.newPermissionCheck(ctx, domain.PermissionProjectGrantMemberWrite, project.AggregateType)(grantedOrgID, projectGrantID)
	//}
	//return nil
}

func (c *Commands) checkPermissionDeleteProjectGrantMember(ctx context.Context, grantedOrgID, projectGrantID string) (err error) {
	// TODO: add permission check for project grant owners
	//if err := c.newPermissionCheck(ctx, domain.PermissionProjectGrantMemberDelete, project.AggregateType)(resourceOwner, projectGrantID); err != nil {
	return c.newPermissionCheck(ctx, domain.PermissionProjectGrantMemberDelete, project.AggregateType)(grantedOrgID, projectGrantID)
	//}
	//return nil
}

func (c *Commands) newUserGrantPermissionCheck(ctx context.Context, permission string) UserGrantPermissionCheck {
	check := c.newPermissionCheck(ctx, permission, project.AggregateType)
	return func(projectID, projectGrantID string) PermissionCheck {
		return func(resourceOwner, _ string) error {
			if projectGrantID != "" {
				return check(resourceOwner, projectGrantID)
			}
			return check(resourceOwner, projectID)
		}
	}
}

func (c *Commands) NewPermissionCheckUserGrantWrite(ctx context.Context) UserGrantPermissionCheck {
	return c.newUserGrantPermissionCheck(ctx, domain.PermissionUserGrantWrite)
}

func (c *Commands) NewPermissionCheckUserGrantDelete(ctx context.Context) UserGrantPermissionCheck {
	return c.newUserGrantPermissionCheck(ctx, domain.PermissionUserGrantDelete)
}
