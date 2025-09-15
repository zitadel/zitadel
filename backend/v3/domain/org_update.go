package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UpdateOrgCommand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var _ Commander = (*UpdateOrgCommand)(nil)

func (u *UpdateOrgCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	organizationRepo := opts.orgRepo()

	updateCount, err := organizationRepo.Update(
		ctx,
		organizationRepo.IDCondition(u.ID),
		authz.GetInstance(ctx).InstanceID(),
		database.NewChange(organizationRepo.NameColumn(), u.Name),
	)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		err = zerrors.ThrowNotFound(nil, "DOM-7PfSUn", "organization not found")
		return err
	}
	if updateCount > 1 {
		err = zerrors.ThrowInternalf(nil, "DOM-QzITrx", "expecting 1 row updated, got %d", updateCount)
		return err
	}

	orgCache.Set(ctx, &Organization{
		ID:   u.ID,
		Name: u.Name,
	})

	return err
}

func (u *UpdateOrgCommand) String() string {
	return "UpdateOrgCommand"
}

func (u *UpdateOrgCommand) Validate() error {
	if u.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-lEMhVC", "invalid organization ID")
	}
	if u.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-wfUntW", "invalid organization name")
	}
	return nil
}
