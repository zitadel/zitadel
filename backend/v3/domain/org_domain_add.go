package domain

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type AddOrgDomainCommand struct {
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	// ClaimedUserIDs []string `json:"claimed_user_ids"`
}

// Events implements Commander.
func (a *AddOrgDomainCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	// TODO(IAM-Marco) Finish implementation in https://github.com/zitadel/zitadel/issues/10447
	return []eventstore.Command{
		org.NewDomainAddedEvent(ctx, &org.NewAggregate(a.OrganizationID).Aggregate, a.Name),
		org.NewDomainVerifiedEvent(ctx, &org.NewAggregate(a.OrganizationID).Aggregate, a.Name),
	}, nil
}

func NewAddOrgDomainCommand(orgID, domainName string) *AddOrgDomainCommand {
	return &AddOrgDomainCommand{
		OrganizationID: orgID,
		Name:           domainName,
	}
}

// Execute implements Commander.
func (a *AddOrgDomainCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

// String implements Commander.
func (a *AddOrgDomainCommand) String() string {
	return "AddOrgDomainCommand"
}

// Validate implements Commander.
func (a *AddOrgDomainCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

var _ Commander = (*AddOrgDomainCommand)(nil)
