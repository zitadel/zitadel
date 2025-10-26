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

var (
	_ Commander     = (*ActivateOrgCommand)(nil)
	_ Transactional = (*ActivateOrgCommand)(nil)
)

type ActivateOrgCommand struct {
	ID string `json:"id"`
}

func NewActivateOrgCommand(organizationID string) *ActivateOrgCommand {
	return &ActivateOrgCommand{ID: organizationID}
}

// RequiresTransaction implements [Transactional].
func (cmd *ActivateOrgCommand) RequiresTransaction() {}

// Events implements [Commander].
func (cmd *ActivateOrgCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	return []eventstore.Command{org.NewOrgReactivatedEvent(ctx, &org.NewAggregate(cmd.ID).Aggregate)}, nil
}

// Execute implements [Commander].
func (cmd *ActivateOrgCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	organizationRepo := opts.organizationRepo

	updateCount, err := organizationRepo.Update(ctx, opts.DB(),
		database.And(
			organizationRepo.IDCondition(cmd.ID),
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

// String implements [Commander].
func (cmd *ActivateOrgCommand) String() string {
	return "ActivateOrgCommand"
}

// Validate implements [Commander].
func (cmd *ActivateOrgCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if cmd.ID = strings.TrimSpace(cmd.ID); cmd.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-hJuuAv", "invalid organization ID")
	}

	organizationRepo := opts.organizationRepo

	// TODO: lock entry as soon as https://github.com/zitadel/zitadel/issues/10930 is done
	org, err := organizationRepo.Get(ctx, opts.DB(), database.WithCondition(
		database.And(
			organizationRepo.IDCondition(cmd.ID),
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
