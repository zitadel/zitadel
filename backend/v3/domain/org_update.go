package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
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
func (u *UpdateOrgCommand) Events(ctx context.Context) []eventstore.Command {
	toReturn := []eventstore.Command{}

	if u.OldDomainName != nil && *u.OldDomainName != u.Name {
		toReturn = append(toReturn, org.NewOrgChangedEvent(ctx, &org.NewAggregate(u.ID).Aggregate, *u.OldDomainName, u.Name))
	}

	return toReturn
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

	organizationRepo := opts.organizationRepo(pool)
	organizationRepo.Domains(true)

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

	err = u.setDomainInfos(ctx, org)
	if err != nil {
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

func (u *UpdateOrgCommand) Validate(_ context.Context, _ *CommandOpts) error {
	if u.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-lEMhVC", "invalid organization ID")
	}
	if u.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-wfUntW", "invalid organization name")
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
