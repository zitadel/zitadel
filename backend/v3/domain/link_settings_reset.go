package domain

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// -------------------------------------------
// COMMAND
// -------------------------------------------

type ResetLinkSettingsCommand struct {
	Instance       bool   `json:"instance"`
	OrganizationId string `json:"organization_id"`
}

func NewResetLinkSettingsCommand(instance bool, organizationId string) *ResetLinkSettingsCommand {
	return &ResetLinkSettingsCommand{
		Instance:       instance,
		OrganizationId: organizationId,
	}
}

func (cmd *ResetLinkSettingsCommand) RequiresTransaction() {}

func (cmd *ResetLinkSettingsCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	// TODO(wim) implement this
	return nil, errors.New("NOT YET IMPLEMENTED")
}

// Validate implements [Querier].
func (q *ResetLinkSettingsCommand) Validate(ctx context.Context, opts *InvokeOpts) error { return nil }

func (q *ResetLinkSettingsCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Querier].
func (q *ResetLinkSettingsCommand) String() string { return "ResetLinkSettingsCommand" }

var (
	_ Commander     = (*ResetLinkSettingsCommand)(nil)
	_ Transactional = (*ResetLinkSettingsCommand)(nil)
)
