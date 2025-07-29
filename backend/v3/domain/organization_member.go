package domain

import "context"

var _ Commander = (*AddOrgMemberCommand)(nil)

// AddOrgMemberCommand adds an existing user as an organization member.
type AddOrgMemberCommand struct {
	InstanceID string   `json:"instanceId"`
	OrgID      string   `json:"orgId"`
	UserID     string   `json:"userId"`
	Roles      []string `json:"roles"`
}

// Execute implements [Commander].
func (a *AddOrgMemberCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	panic("unimplemented")
}

// String implements [Commander].
func (a *AddOrgMemberCommand) String() string {
	return "AddOrgMemberCommand"
}

var _ Commander = (*CreateOrgMemberCommand)(nil)

// CreateOrgMemberCommand creates a new user and adds them as an organization member.
type CreateOrgMemberCommand struct{}

// Execute implements [Commander].
func (c *CreateOrgMemberCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	panic("unimplemented")
}

// String implements [Commander].
func (c *CreateOrgMemberCommand) String() string {
	return "CreateOrgMemberCommand"
}

// MemberRepository is a sub repository of the org repository and maybe the instance repository.
type MemberRepository interface {
	AddMember(ctx context.Context, orgID, userID string, roles []string) error
	SetMemberRoles(ctx context.Context, orgID, userID string, roles []string) error
	RemoveMember(ctx context.Context, orgID, userID string) error
}
