package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DeleteOrgCommand struct {
	OrganizationName string `json: "organization_name"`
	ID               string `json:"id"`
}

func NewDeleteOrgCommand(organizationID string) *DeleteOrgCommand {
	return &DeleteOrgCommand{ID: organizationID}
}

// Events implements Commander.
func (d *DeleteOrgCommand) Events(ctx context.Context, opts *CommandOpts) ([]eventstore.Command, error) {
	return nil, nil
}

// Execute implements Commander.
func (d *DeleteOrgCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	closeFunc, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}

	defer func() { err = closeFunc(ctx, err) }()

	orgRepo := opts.organizationRepo(pool)

	orgToDelete, err := orgRepo.Get(ctx, database.WithCondition(orgRepo.IDCondition(d.ID)))
	if err != nil {
		return err
	}
	d.OrganizationName = orgToDelete.Name

	deletedRows, err := orgRepo.Delete(ctx, orgRepo.IDCondition(d.ID), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return err
	}

	if deletedRows > 1 {
		err = zerrors.ThrowInternalf(nil, "DOM-5cE9u6", "expecting 1 row deleted, got %d", deletedRows)
		return err
	}

	if deletedRows < 1 {
		err = zerrors.ThrowNotFoundf(nil, "DOM-ur6Qyv", "organization not found")
	}
	return err
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

var _ Commander = (*DeleteOrgCommand)(nil)
