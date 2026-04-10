package domain

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// -------------------------------------------
// COMMAND
// -------------------------------------------

type SetLinkSettingsCommand struct {
	Instance       bool         `json:"instance"`
	OrganizationId string       `json:"organization_id"`
	Settings       LinkSettings `json:"settings"`
}

func NewSetLinkSettingsCommand(instance bool, organizationId string, settings LinkSettings) *SetLinkSettingsCommand {
	return &SetLinkSettingsCommand{
		Instance:       instance,
		OrganizationId: organizationId,
		Settings:       settings,
	}
}

func (cmd *SetLinkSettingsCommand) RequiresTransaction() {}

func (cmd *SetLinkSettingsCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	// TODO(wim) implement this
	return nil, errors.New("NOT YET IMPLEMENTED")
}

// Validate implements [Querier].
func (q *SetLinkSettingsCommand) Validate(ctx context.Context, opts *InvokeOpts) error { return nil }

func (q *SetLinkSettingsCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Querier].
func (q *SetLinkSettingsCommand) String() string { return "SetLinkSettingsCommand" }

var (
	_ Commander     = (*SetLinkSettingsCommand)(nil)
	_ Transactional = (*SetLinkSettingsCommand)(nil)
)
