package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UpdateOrgCommand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var _ Commander = (*UpdateOrgCommand)(nil)

func NewUpdateOrgCommand(id, name string) *UpdateOrgCommand {
	return &UpdateOrgCommand{
		ID:   id,
		Name: name,
	}
}

func (u *UpdateOrgCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	organizationRepo := opts.orgRepo()

	org, err := organizationRepo.Get(ctx, database.WithCondition(organizationRepo.IDCondition(u.ID)))
	if err != nil {
		return err
	}

	if org.Name == u.Name {
		err = NewOrgNameNotChangedError("DOM-nDzwIu")
		return err
	}

	if org.State == OrgStateInactive {
		err = NewOrgNotFoundError("DOM-OcA1jq")
		return err
	}

	updateCount, err := organizationRepo.Update(
		ctx,
		organizationRepo.IDCondition(org.ID),
		org.InstanceID,
		database.NewChange(organizationRepo.NameColumn(), u.Name),
	)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		err = NewOrgNotFoundError("DOM-7PfSUn")
		return err
	}
	if updateCount > 1 {
		err = NewMultipleOrgsUpdatedError("DOM-QzITrx", 1, updateCount)
		return err
	}

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
