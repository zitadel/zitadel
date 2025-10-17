package domain

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type SetPrimaryOrgDomainCommand struct {
	OrgID string `json:"org_id"`
	Name  string `json:"name"`
}

// Events implements Commander.
func (a *SetPrimaryOrgDomainCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	// TODO(IAM-Marco) Finish implementation in https://github.com/zitadel/zitadel/issues/10447
	return []eventstore.Command{
		org.NewDomainPrimarySetEvent(ctx, &org.NewAggregate(a.OrgID).Aggregate, a.Name),
	}, nil
}

func NewSetPrimaryOrgDomainCommand(orgID, domainName string) *SetPrimaryOrgDomainCommand {
	return &SetPrimaryOrgDomainCommand{
		OrgID: orgID,
		Name:  domainName,
	}
}

// Execute implements Commander.
func (a *SetPrimaryOrgDomainCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

// String implements Commander.
func (a *SetPrimaryOrgDomainCommand) String() string {
	return "SetPrimaryOrgDomainCommand"
}

// Validate implements Commander.
func (a *SetPrimaryOrgDomainCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

var _ Commander = (*SetPrimaryOrgDomainCommand)(nil)
