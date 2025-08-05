package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type OrgState string

const (
	OrgStateActive   OrgState = "active"
	OrgStateInactive OrgState = "inactive"
)

type Organization struct {
	ID         string    `json:"id,omitempty" db:"id"`
	Name       string    `json:"name,omitempty" db:"name"`
	InstanceID string    `json:"instanceId,omitempty" db:"instance_id"`
	State      OrgState  `json:"state,omitempty" db:"state"`
	CreatedAt  time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt,omitzero" db:"updated_at"`

	Domains []*OrganizationDomain `json:"domains,omitempty" db:"-"` // domains need to be handled separately
}

var _ Commander = (*CreateOrganizationCommand)(nil)

type CreateOrganizationCommand struct {
	InstanceID string `json:"instanceId"`
	// ID is optional, if not set a new ID will be generated.
	// It can be set using the [WithOrganizationID] option in [NewCreateOrganizationCommand].
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`

	// CreatedAt MUST NOT be set by the caller.
	CreatedAt time.Time `json:"createdAt,omitzero"`

	// Admins represent the commands to create the administrators.
	// The Commanders MUST either be [AddOrgMemberCommand] or [CreateOrgMemberCommand].
	Admins []Commander `json:"admins,omitempty"`
}

type CreateOrganizationCommandOpts interface {
	applyOnCreateOrganizationCommand(cmd *CreateOrganizationCommand)
}

func NewCreateOrganizationCommand(instanceID, name string, opts ...CreateOrganizationCommandOpts) *CreateOrganizationCommand {
	cmd := &CreateOrganizationCommand{
		InstanceID: instanceID,
		Name:       name,
	}
	for _, opt := range opts {
		opt.applyOnCreateOrganizationCommand(cmd)
	}
	return cmd
}

// Execute implements [Commander].
//
// DISCUSS(adlerhurst): As we need to do validation to make sure a command contains all the data required
// we can consider the following options:
// 1. Validate the command before executing it, which is what we do here.
// 2. Create an invoker which checks if the struct has a `Validate() error` method and call it in the chain of invokers.
// While the the first one is more straightforward it bloats the execute method with validation logic.
// The second one would allow us to keep the execute method clean, but could be more error prone if the method gets missed during implementation.
func (cmd *CreateOrganizationCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	if cmd.ID == "" {
		cmd.ID, err = generateID()
		if err != nil {
			return err
		}
	}
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	err = orgRepo(opts.DB).Create(ctx, cmd)
	if err != nil {
		return err
	}

	for _, admin := range cmd.Admins {
		if err = opts.Invoke(ctx, admin); err != nil {
			return err
		}
	}

	return nil
}

// String implements [Commander].
func (CreateOrganizationCommand) String() string {
	return "CreateOrganizationCommand"
}

var (
	_ Commander = (*ActivateOrganizationCommand)(nil)
)

type ActivateOrganizationCommand struct {
	InstanceID string `json:"instanceId"`
	OrgID      string `json:"orgId"`

	// UpdatedAt MUST NOT be set by the caller.
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
}

func NewActivateOrganizationCommand(instanceID, orgID string) *ActivateOrganizationCommand {
	return &ActivateOrganizationCommand{
		InstanceID: instanceID,
		OrgID:      orgID,
	}
}

// Execute implements [Commander].
func (cmd *ActivateOrganizationCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	repo := orgRepo(opts.DB)
	_, err = repo.Update(ctx,
		database.And(
			repo.InstanceIDCondition(cmd.InstanceID),
			repo.IDCondition(cmd.OrgID),
		),
		repo.SetState(OrgStateActive),
	)
	return err
}

// String implements [Commander].
func (ActivateOrganizationCommand) String() string {
	return "ActivateOrganizationCommand"
}

var (
	_ Commander = (*DeactivateOrganizationCommand)(nil)
)

type DeactivateOrganizationCommand struct {
	InstanceID string `json:"instanceId"`
	OrgID      string `json:"orgId"`

	// UpdatedAt MUST NOT be set by the caller.
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
}

func NewDeactivateOrganizationCommand(instanceID, orgID string) *DeactivateOrganizationCommand {
	return &DeactivateOrganizationCommand{
		InstanceID: instanceID,
		OrgID:      orgID,
	}
}

// Execute implements [Commander].
func (cmd *DeactivateOrganizationCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	repo := orgRepo(opts.DB)
	_, err = repo.Update(ctx,
		database.And(
			repo.InstanceIDCondition(cmd.InstanceID),
			repo.IDCondition(cmd.OrgID),
		),
		repo.SetState(OrgStateInactive),
	)
	return err
}

// String implements [Commander].
func (DeactivateOrganizationCommand) String() string {
	return "DeactivateOrganizationCommand"
}

var (
	_ Commander = (*DeleteOrganizationCommand)(nil)
)

type DeleteOrganizationCommand struct {
	InstanceID string `json:"instanceId"`
	OrgID      string `json:"orgId"`
}

func NewDeleteOrganizationCommand(instanceID, orgID string) *DeleteOrganizationCommand {
	return &DeleteOrganizationCommand{
		InstanceID: instanceID,
		OrgID:      orgID,
	}
}

// Execute implements [Commander].
func (cmd *DeleteOrganizationCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	repo := orgRepo(opts.DB)
	_, err = repo.Delete(ctx,
		database.And(
			repo.InstanceIDCondition(cmd.InstanceID),
			repo.IDCondition(cmd.OrgID),
		),
	)
	return err
}

// String implements [Commander].
func (DeleteOrganizationCommand) String() string {
	return "DeleteOrganizationCommand"
}

var _ Commander = (*UpdateOrganizationCommand)(nil)

type UpdateOrganizationCommand struct {
	InstanceID string `json:"instanceId"`
	OrgID      string `json:"orgId"`

	repo    OrganizationRepository
	changes database.Changes
	opts    []UpdateOrganizationCommandOpts
}

func NewUpdateOrganizationCommand(instanceID, orgID string, opts ...UpdateOrganizationCommandOpts) *UpdateOrganizationCommand {
	return &UpdateOrganizationCommand{
		InstanceID: instanceID,
		OrgID:      orgID,
		opts:       opts,
	}
}

type UpdateOrganizationCommandOpts interface {
	applyOnUpdateOrganizationCommand(cmd *UpdateOrganizationCommand)
}

// Execute implements [Commander].
func (cmd *UpdateOrganizationCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	cmd.repo = orgRepo(opts.DB)
	for _, opt := range cmd.opts {
		opt.applyOnUpdateOrganizationCommand(cmd)
	}

	if len(cmd.changes) == 0 {
		return nil // No update needed if no changes are provided.
	}

	_, err = cmd.repo.Update(ctx,
		database.And(
			cmd.repo.InstanceIDCondition(cmd.InstanceID),
			cmd.repo.IDCondition(cmd.OrgID),
		),
		cmd.changes...,
	)
	return err
}

// String implements [Commander].
func (UpdateOrganizationCommand) String() string {
	return "UpdateOrganizationCommand"
}

type OrgsQueryOpts interface {
	applyOnOrgsQuery(query *OrgsQuery)
}

var _ Commander = (*OrgsQuery)(nil)

type OrgsQuery struct {
	InstanceID string

	opts       []OrgsQueryOpts
	repo       OrganizationRepository
	domainRepo OrganizationDomainRepository
	conditions []database.Condition
	pagination Pagination

	Result []*Organization
}

func NewOrgsQuery(instanceID string, opts ...OrgsQueryOpts) *OrgsQuery {
	return &OrgsQuery{
		InstanceID: instanceID,
		opts:       opts,
	}
}

// Execute implements [Commander].
func (q *OrgsQuery) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	q.repo = orgRepo(opts.DB)
	q.domainRepo = q.repo.Domains(true)
	q.conditions = append(q.conditions, q.repo.InstanceIDCondition(q.InstanceID))
	for _, opt := range q.opts {
		opt.applyOnOrgsQuery(q)
	}

	q.Result, err = q.repo.List(ctx,
		database.WithCondition(database.And(q.conditions...)),
		database.WithLimit(q.pagination.Limit),
		database.WithOffset(q.pagination.Offset),
		database.WithOrderBy(!q.pagination.Ascending, q.pagination.OrderColumns...),
	)
	return err
}

// String implements [Commander].
func (OrgsQuery) String() string {
	return "OrgsQuery"
}
