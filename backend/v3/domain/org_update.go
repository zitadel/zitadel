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

// Events implements Commander.
func (u *UpdateOrgCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	toReturn := []eventstore.Command{}

	if u.OldDomainName != nil && *u.OldDomainName != u.Name {
		toReturn = append(toReturn, org.NewOrgChangedEvent(ctx, &org.NewAggregate(u.ID).Aggregate, *u.OldDomainName, u.Name))
	}

	return toReturn, nil
}

var _ Commander = (*UpdateOrgCommand)(nil)

func NewUpdateOrgCommand(id, name string) *UpdateOrgCommand {
	return &UpdateOrgCommand{
		ID:   id,
		Name: name,
	}
}

func (u *UpdateOrgCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	organizationRepo := opts.organizationRepo.LoadDomains()

	org, err := organizationRepo.Get(ctx, pool, database.WithCondition(
		database.And(
			organizationRepo.IDCondition(u.ID),
			organizationRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
		),
	))
	if err != nil {
		return err
	}

	err = u.setDomainInfos(ctx, org)
	if err != nil {
		return err
	}

	updateCount, err := organizationRepo.Update(
		ctx,
		pool,
		database.And(
			organizationRepo.IDCondition(org.ID),
			organizationRepo.InstanceIDCondition(org.InstanceID),
		),
		database.NewChange(organizationRepo.NameColumn(), u.Name),
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

func (u *UpdateOrgCommand) String() string {
	return "UpdateOrgCommand"
}

func (u *UpdateOrgCommand) Validate(ctx context.Context, opts *InvokeOpts) error {
	if u.ID = strings.TrimSpace(u.ID); u.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-lEMhVC", "invalid organization ID")
	}
	if u.Name = strings.TrimSpace(u.Name); u.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-wfUntW", "invalid organization name")
	}

	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	organizationRepo := opts.organizationRepo

	org, err := organizationRepo.Get(ctx, pool, database.WithCondition(
		database.And(
			organizationRepo.IDCondition(u.ID),
			organizationRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
		),
	))
	if err != nil {
		return err
	}

	if org.Name == u.Name {
		err = zerrors.ThrowPreconditionFailed(nil, "DOM-nDzwIu", "Errors.Org.NotChanged")
		return err
	}

	return nil
}

func (u *UpdateOrgCommand) setDomainInfos(ctx context.Context, org *Organization) error {
	iamDomain := http_utils.DomainContext(ctx).RequestedDomain()
	oldDomainName, err := domain.NewIAMDomainName(org.Name, iamDomain)
	if err != nil {
		return err
	}

	for _, d := range org.Domains {
		if d.Domain == oldDomainName && !d.IsPrimary {
			u.IsOldDomainVerified = &d.IsVerified
			u.OldDomainName = &d.Domain
			return nil
		}
	}

	return nil
}
