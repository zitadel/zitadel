package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UpdateOrgCommand struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	OldDomainName       *string `json:"old_name"`
	IsOldDomainVerified *bool   `json:"is_old_domain_primary"`
}

// RequiresTransaction implements [Transactional].
func (cmd *UpdateOrgCommand) RequiresTransaction() {}

// Events implements Commander.
func (cmd *UpdateOrgCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	toReturn := []eventstore.Command{}

	if cmd.OldDomainName != nil && *cmd.OldDomainName != cmd.Name {
		toReturn = append(toReturn, org.NewOrgChangedEvent(ctx, &org.NewAggregate(cmd.ID).Aggregate, *cmd.OldDomainName, cmd.Name))
	}

	return toReturn, nil
}

var (
	_ Commander     = (*UpdateOrgCommand)(nil)
	_ Transactional = (*UpdateOrgCommand)(nil)
)

func NewUpdateOrgCommand(id, name string) *UpdateOrgCommand {
	return &UpdateOrgCommand{
		ID:   id,
		Name: name,
	}
}

func (cmd *UpdateOrgCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	organizationRepo := opts.organizationRepo.LoadDomains()

	org, err := organizationRepo.Get(ctx, opts.DB(), database.WithCondition(
		organizationRepo.PrimaryKeyCondition(authz.GetInstance(ctx).InstanceID(), cmd.ID),
	))
	if err != nil {
		return err
	}

	err = cmd.setDomainInfos(ctx, org)
	if err != nil {
		return err
	}

	updateCount, err := organizationRepo.Update(
		ctx,
		opts.DB(),
		organizationRepo.PrimaryKeyCondition(org.InstanceID, org.ID),
		database.NewChange(organizationRepo.NameColumn(), cmd.Name),
	)
	if err != nil {
		return err
	}

	if updateCount == 0 {
		err = zerrors.ThrowNotFound(nil, "DOM-7PfSUn", "Errors.Org.NotFound")
		return err
	}
	if updateCount > 1 {
		err = zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(1, updateCount), "DOM-QzITrx", "unexpected number of rows updated")
		return err
	}

	return err
}

func (cmd *UpdateOrgCommand) String() string {
	return "UpdateOrgCommand"
}

func (cmd *UpdateOrgCommand) Validate(ctx context.Context, opts *InvokeOpts) error {
	if cmd.ID = strings.TrimSpace(cmd.ID); cmd.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-lEMhVC", "invalid organization ID")
	}
	if cmd.Name = strings.TrimSpace(cmd.Name); cmd.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-wfUntW", "invalid organization name")
	}

	organizationRepo := opts.organizationRepo

	org, err := organizationRepo.Get(ctx, opts.DB(),
		database.WithCondition(
			opts.organizationRepo.PrimaryKeyCondition(authz.GetInstance(ctx).InstanceID(), cmd.ID),
		))
	if err != nil {
		return err
	}

	if org.Name == cmd.Name {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-nDzwIu", "Errors.Org.NotChanged")
	}

	return nil
}

func (cmd *UpdateOrgCommand) setDomainInfos(ctx context.Context, org *Organization) error {
	iamDomain := http_utils.DomainContext(ctx).RequestedDomain()
	oldDomainName, err := domain.NewIAMDomainName(org.Name, iamDomain)
	if err != nil {
		return err
	}

	for _, d := range org.Domains {
		if d.Domain == oldDomainName && !d.IsPrimary {
			cmd.IsOldDomainVerified = &d.IsVerified
			cmd.OldDomainName = &d.Domain
			return nil
		}
	}

	return nil
}
