package domain

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type RemoveOrgDomainCommand struct {
	OrgID      string  `json:"org_id"`
	Name       *string `json:"name"`
	IsVerified *bool   `json:"is_verified"`
}

// Events implements Commander.
func (r *RemoveOrgDomainCommand) Events(ctx context.Context) []eventstore.Command {
	// TODO(IAM-Marco) Finish implementation in https://github.com/zitadel/zitadel/issues/10447
	oldDomainName := ""
	isVerified := false
	if r.Name != nil {
		oldDomainName = *r.Name
	}
	if r.IsVerified != nil {
		isVerified = *r.IsVerified
	}

	return []eventstore.Command{
		org.NewDomainRemovedEvent(ctx, &org.NewAggregate(r.OrgID).Aggregate, oldDomainName, isVerified),
	}
}

func NewRemoveOrgDomainCommand(orgID string, domainName *string, isDomainVerified *bool) *RemoveOrgDomainCommand {
	return &RemoveOrgDomainCommand{
		OrgID:      orgID,
		Name:       domainName,
		IsVerified: isDomainVerified,
	}
}

// Execute implements Commander.
func (r *RemoveOrgDomainCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	return nil
}

// String implements Commander.
func (r *RemoveOrgDomainCommand) String() string {
	return "RemoveOrgDomainCommand"
}

// Validate implements Commander.
func (r *RemoveOrgDomainCommand) Validate() (err error) {
	return nil
}

var _ Commander = (*RemoveOrgDomainCommand)(nil)
