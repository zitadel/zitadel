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

var _ Commander = (*ActivateOrgCommand)(nil)

type ActivateOrgCommand struct {
	ID string `json:"id"`
}

func NewActivateOrgCommand(organizationID string) *ActivateOrgCommand {
	return &ActivateOrgCommand{ID: organizationID}
}

// Events implements Commander.
func (d *ActivateOrgCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	return []eventstore.Command{org.NewOrgReactivatedEvent(ctx, &org.NewAggregate(d.ID).Aggregate)}, nil
}

// Execute implements Commander.
func (d *ActivateOrgCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
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
		database.NewChange(organizationRepo.StateColumn(), OrgStateActive),
	)

	if err != nil {
		return err
	}

	if updateCount == 0 {
		err = zerrors.ThrowNotFound(nil, "DOM-CGumXG", "Errors.Org.NotFound")
		return err
	}
	if updateCount > 1 {
		err = zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-SEWCLp", "unexpected number of rows updated")
		return err
	}

	return nil
}

// String implements Commander.
func (d *ActivateOrgCommand) String() string {
	return "ActivateOrgCommand"
}

// Validate implements Commander.
func (d *ActivateOrgCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if strings.TrimSpace(d.ID) == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-hJuuAv", "invalid organization ID")
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
			err = zerrors.ThrowNotFound(err, "DOM-86HVfs", "Errors.Org.NotFound")
		}
		return err
	}

	if org.State == OrgStateActive {
		err = zerrors.ThrowPreconditionFailed(nil, "DOM-Ixfbxh", "Errors.Org.AlreadyActive")
		return err
	}

	return nil
}
