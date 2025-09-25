package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DeleteOrgCommand struct {
	ID string `json:"id"`
}

// Events implements Commander.
func (d *DeleteOrgCommand) Events(ctx context.Context) []eventstore.Command {
	return nil
}

// Execute implements Commander.
func (d *DeleteOrgCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	return nil
}

// String implements Commander.
func (d *DeleteOrgCommand) String() string {
	return "DeleteOrgCommand"
}

// Validate implements Commander.
func (d *DeleteOrgCommand) Validate(ctx context.Context, opts *CommandOpts) (err error) {
	instance := authz.GetInstance(ctx)

	if d.ID == instance.DefaultOrganisationID() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-LCkE69", "Errors.Org.DefaultOrgNotDeletable")
	}

	closeFunc, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}

	defer func() { err = closeFunc(ctx, err) }()

	// TODO(IAM-Marco): Implement this when projects are available on relational
	// // Check if the ZITADEL project exists on the input organization
	// projectRepo := opts.projectRepo(pool)
	// _, getErr := projectRepo.Get(ctx, database.WithCondition(projectRepo.IDCondition(instance.ProjectID())))
	// if getErr == nil {
	// 	return zerrors.ThrowPreconditionFailed(nil, "DOM-X7YXxC", "Errors.Org.ZitadelOrgNotDeletable")
	// }
	// // "precondition failed" error means the project does not exist, return other errors in case it's not that
	// if !zerrors.IsPreconditionFailed(getErr) {
	// 	err = getErr
	// 	return err
	// }

	orgRepo := opts.organizationRepo(pool)
	org, err := orgRepo.Get(ctx, database.WithCondition(orgRepo.IDCondition(d.ID)))
	if err != nil {
		return err
	}
	wantedStates := map[OrgState]bool{
		OrgStateActive:   true,
		OrgStateInactive: true,
	}

	// TODO(IAM-Marco): I'm not sure this is needed on relational.
	if !wantedStates[org.State] {
		err = zerrors.ThrowNotFound(nil, "DOM-8KYOH3", "Errors.Org.NotFound")
		return err
	}

	return err
}

func NewDeleteOrgCommand(organizationID string) *DeleteOrgCommand {
	return &DeleteOrgCommand{ID: organizationID}
}

var _ Commander = (*DeleteOrgCommand)(nil)
