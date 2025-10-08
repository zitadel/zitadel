package domain

import (
	"context"
	"errors"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ Commander = (*DeactivateOrgCommand)(nil)

type DeactivateOrgCommand struct {
	ID string `json:"id"`
}

func NewDeactivateOrgCommand(organizationID string) *DeactivateOrgCommand {
	return &DeactivateOrgCommand{ID: organizationID}
}

// Events implements Commander.
func (d *DeactivateOrgCommand) Events(ctx context.Context, opts *CommandOpts) ([]eventstore.Command, error) {
	return []eventstore.Command{org.NewOrgDeactivatedEvent(ctx, &org.NewAggregate(d.ID).Aggregate)}, nil
}

// Execute implements Commander.
func (d *DeactivateOrgCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	organizationRepo := opts.organizationRepo

	updateCount, err := organizationRepo.Update(ctx, pool,
		database.And(
			organizationRepo.IDCondition(d.ID),
			organizationRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
		),
		database.NewChange(organizationRepo.StateColumn(), OrgStateInactive),
	)

	if err != nil {
		return err
	}

	if updateCount == 0 {
		err = zerrors.ThrowNotFound(nil, "DOM-vWPy7D", "Errors.Org.NotFound")
		return err
	}
	if updateCount > 1 {
		err = zerrors.ThrowInternal(NewMultipleObjecstUpdatedError(1, updateCount), "DOM-dXl1kJ", "unexpected number of rows updated")
		return err
	}

	return nil
}

// String implements Commander.
func (d *DeactivateOrgCommand) String() string {
	return "DeactivateOrgCommand"
}

// Validate implements Commander.
func (d *DeactivateOrgCommand) Validate(ctx context.Context, opts *CommandOpts) (err error) {
	if strings.TrimSpace(d.ID) == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-Qc3T1r", "invalid organization ID")
	}

	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	organizationRepo := opts.organizationRepo

	org, err := organizationRepo.Get(ctx, pool, database.WithCondition(
		database.And(
			organizationRepo.IDCondition(d.ID),
			organizationRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
		),
	))
	if err != nil {
		var notFoundError *database.NoRowFoundError
		if errors.As(err, &notFoundError) {
			err = zerrors.ThrowNotFound(err, "DOM-QEjfpz", "Errors.Org.NotFound")
		}
		return err
	}

	if org.State == OrgStateInactive {
		err = zerrors.ThrowPreconditionFailed(nil, "DOM-Z2dzsT", "Errors.Org.AlreadyDeactivated")
		return err
	}
	return nil
}
