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

type AddProjectMember struct {
	ResourceOwner string
	ProjectID     string
	UserID        string
	Roles         []string
}

func (i *AddProjectMember) IsValid(zitadelRoles []authz.RoleMapping) error {
	if i.ProjectID == "" || i.UserID == "" || len(i.Roles) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-W8m4l", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(i.Roles, domain.ProjectRolePrefix, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-3m9ds", "Errors.Project.Member.Invalid")
	}
	return nil
}

func (c *Commands) AddProjectMember(ctx context.Context, member *AddProjectMember) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := member.IsValid(c.zitadelRoles); err != nil {
		return nil, err
	}
	_, err = c.checkUserExists(ctx, member.UserID, "")
	if err != nil {
		return nil, err
	}
	projectResourceOwner, err := c.checkProjectExists(ctx, member.ProjectID, member.ResourceOwner)
	if err != nil {
		return nil, err
	}
	// resourceowner of the member if not provided is the resourceowner of the project
	if member.ResourceOwner == "" {
		member.ResourceOwner = projectResourceOwner
	}
	addedMember, err := c.projectMemberWriteModelByID(ctx, member.ProjectID, member.UserID, member.ResourceOwner)
	if err != nil {
		return nil, err
	}
	// error if provided resourceowner is not equal to the resourceowner of the project
	if projectResourceOwner != addedMember.ResourceOwner {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-0l10S9OmZV", "Errors.Project.Member.Invalid")
	}
	if err := c.checkPermissionUpdateProjectMember(ctx, addedMember.ResourceOwner, addedMember.AggregateID); err != nil {
		return nil, err
	}
	if addedMember.State.Exists() {
		return nil, zerrors.ThrowAlreadyExists(nil, "PROJECT-PtXi1", "Errors.Project.Member.AlreadyExists")
	}

	pushedEvents, err := c.eventstore.Push(ctx,
		project.NewProjectMemberAddedEvent(ctx,
			ProjectAggregateFromWriteModelWithCTX(ctx, &addedMember.WriteModel),
			member.UserID,
			member.Roles...,
		),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&addedMember.WriteModel), nil
}

type ChangeProjectMember struct {
	ResourceOwner string
	ProjectID     string
	UserID        string
	Roles         []string
}

func (i *ChangeProjectMember) IsValid(zitadelRoles []authz.RoleMapping) error {
	if i.ProjectID == "" || i.UserID == "" || len(i.Roles) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-LiaZi", "Errors.Project.Member.Invalid")
	}
	if len(domain.CheckForInvalidRoles(i.Roles, domain.ProjectRolePrefix, zitadelRoles)) > 0 {
		return zerrors.ThrowInvalidArgument(nil, "PROJECT-3m9d", "Errors.Project.Member.Invalid")
	}
	return nil
}

// ChangeProjectMember updates an existing member
func (c *Commands) ChangeProjectMember(ctx context.Context, member *ChangeProjectMember) (*domain.ObjectDetails, error) {
	if err := member.IsValid(c.zitadelRoles); err != nil {
		return nil, err
	}

	existingMember, err := c.projectMemberWriteModelByID(ctx, member.ProjectID, member.UserID, member.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.NotFound")
	}
	if err := c.checkPermissionUpdateProjectMember(ctx, existingMember.ResourceOwner, existingMember.AggregateID); err != nil {
		return nil, err
	}
	if slices.Compare(existingMember.Roles, member.Roles) == 0 {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}
	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &existingMember.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectMemberChangedEvent(ctx, projectAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&existingMember.WriteModel), nil
}

func (c *Commands) RemoveProjectMember(ctx context.Context, projectID, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-66mHd", "Errors.Project.Member.Invalid")
	}
	existingMember, err := c.projectMemberWriteModelByID(ctx, projectID, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingMember.State.Exists() {
		return writeModelToObjectDetails(&existingMember.WriteModel), nil
	}
	if err := c.checkPermissionDeleteProjectMember(ctx, existingMember.ResourceOwner, existingMember.AggregateID); err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &existingMember.WriteModel)
	removeEvent := c.removeProjectMember(ctx, projectAgg, userID, false)
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

func (c *Commands) removeProjectMember(ctx context.Context, projectAgg *eventstore.Aggregate, userID string, cascade bool) eventstore.Command {
	if cascade {
		return project.NewProjectMemberCascadeRemovedEvent(
			ctx,
			projectAgg,
			userID)
	} else {
		return project.NewProjectMemberRemovedEvent(ctx, projectAgg, userID)
	}
}

func (c *Commands) projectMemberWriteModelByID(ctx context.Context, projectID, userID, resourceOwner string) (member *ProjectMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectMemberWriteModel(projectID, userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
