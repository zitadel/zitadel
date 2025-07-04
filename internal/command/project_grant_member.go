package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddProjectGrantMember struct {
	ResourceOwner string
	UserID        string
	GrantID       string
	ProjectID     string
	Roles         []string
}

func (i *AddProjectGrantMember) IsValid(zitadelRoles []authz.RoleMapping) error {
	if i.ProjectID == "" || i.GrantID == "" || i.UserID == "" || len(i.Roles) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-8fi7G", "Errors.Project.Grant.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(i.Roles, domain.ProjectGrantRolePrefix, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-m9gKK", "Errors.Project.Grant.Member.Invalid")
	}
	return nil
}

func (c *Commands) AddProjectGrantMember(ctx context.Context, member *AddProjectGrantMember) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := member.IsValid(c.zitadelRoles); err != nil {
		return nil, err
	}
	_, err = c.checkUserExists(ctx, member.UserID, "")
	if err != nil {
		return nil, err
	}
	grantedOrgID, projectGrantResourceOwner, err := c.checkProjectGrantExists(ctx, member.GrantID, "", member.ProjectID, "")
	if err != nil {
		return nil, err
	}
	if member.ResourceOwner == "" {
		member.ResourceOwner = projectGrantResourceOwner
	}
	addedMember, err := c.projectGrantMemberWriteModelByID(ctx, member.ProjectID, member.UserID, member.GrantID, member.ResourceOwner)
	if err != nil {
		return nil, err
	}
	// TODO: change e2e tests to use correct resourceowner, wrong resource owner is corrected through aggregate
	// error if provided resourceowner is not equal to the resourceowner of the project grant
	//if projectGrantResourceOwner != addedMember.ResourceOwner {
	//	return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-0l10S9OmZV", "Errors.Project.Grant.Invalid")
	//}
	if addedMember.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-37fug", "Errors.AlreadyExists")
	}
	if err := c.checkPermissionUpdateProjectGrantMember(ctx, grantedOrgID, addedMember.GrantID); err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(
		ctx,
		project.NewProjectGrantMemberAddedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &addedMember.WriteModel),
			member.UserID,
			member.GrantID,
			member.Roles...,
		))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&addedMember.WriteModel), nil
}

type ChangeProjectGrantMember struct {
	UserID    string
	GrantID   string
	ProjectID string
	Roles     []string
}

func (i *ChangeProjectGrantMember) IsValid(zitadelRoles []authz.RoleMapping) error {
	if i.ProjectID == "" || i.GrantID == "" || i.UserID == "" || len(i.Roles) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-109fs", "Errors.Project.Grant.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(i.Roles, domain.ProjectGrantRolePrefix, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-m0sDf", "Errors.Project.Grant.Member.Invalid")
	}
	return nil
}

// ChangeProjectGrantMember updates an existing member
func (c *Commands) ChangeProjectGrantMember(ctx context.Context, member *ChangeProjectGrantMember) (*domain.ObjectDetails, error) {
	if err := member.IsValid(c.zitadelRoles); err != nil {
		return nil, err
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, member.GrantID, "", member.ProjectID, "")
	if err != nil {
		return nil, err
	}
	existingMember, err := c.projectGrantMemberWriteModelByID(ctx, member.ProjectID, member.UserID, member.GrantID, existingGrant.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-37fug", "Errors.NotFound")
	}

	if err := c.checkPermissionUpdateProjectGrantMember(ctx, existingGrant.GrantedOrgID, existingMember.GrantID); err != nil {
		return nil, err
	}
	if slices.Compare(existingMember.Roles, member.Roles) == 0 {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}

	pushedEvents, err := c.eventstore.Push(
		ctx,
		project.NewProjectGrantMemberChangedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &existingMember.WriteModel),
			member.UserID,
			member.GrantID,
			member.Roles...,
		))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&existingMember.WriteModel), nil
}

func (c *Commands) RemoveProjectGrantMember(ctx context.Context, projectID, userID, grantID string) (*domain.ObjectDetails, error) {
	if projectID == "" || userID == "" || grantID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-66mHd", "Errors.Project.Member.Invalid")
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, "", projectID, "")
	if err != nil {
		return nil, err
	}
	existingMember, err := c.projectGrantMemberWriteModelByID(ctx, projectID, userID, grantID, existingGrant.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}
	if err := c.checkPermissionDeleteProjectGrantMember(ctx, existingGrant.GrantedOrgID, existingMember.GrantID); err != nil {
		return nil, err
	}

	removeEvent := c.removeProjectGrantMember(ctx,
		ProjectAggregateFromWriteModelWithCTX(ctx, &existingMember.WriteModel),
		userID,
		grantID,
		false,
	)
	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMember.WriteModel), nil
}

func (c *Commands) removeProjectGrantMember(ctx context.Context, projectAgg *eventstore.Aggregate, userID, grantID string, cascade bool) eventstore.Command {
	if cascade {
		return project.NewProjectGrantMemberCascadeRemovedEvent(
			ctx,
			projectAgg,
			userID,
			grantID)
	} else {
		return project.NewProjectGrantMemberRemovedEvent(ctx, projectAgg, userID, grantID)
	}
}

func (c *Commands) projectGrantMemberWriteModelByID(ctx context.Context, projectID, userID, grantID, resourceOwner string) (member *ProjectGrantMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGrantMemberWriteModel(projectID, userID, grantID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
