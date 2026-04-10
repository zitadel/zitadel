package domain

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type ResetLinkSettingsCommand struct {
	Instance       bool   `json:"instance"`
	OrganizationId string `json:"organization_id"`
	result         *ResetLinkSettingsCommandResult
}

type ResetLinkSettingsCommandResult struct {
	ChangeTime time.Time
	Links      []Link
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

// Validate implements [Commander].
func (q *ResetLinkSettingsCommand) Validate(ctx context.Context, opts *InvokeOpts) error {
	return nil
}

func (q *ResetLinkSettingsCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Commander].
func (q *ResetLinkSettingsCommand) String() string { return "ResetLinkSettingsCommand" }

func (q *ResetLinkSettingsCommand) Result() *ResetLinkSettingsCommandResult {
	return q.result
}

var (
	_ Commander     = (*ResetLinkSettingsCommand)(nil)
	_ Transactional = (*ResetLinkSettingsCommand)(nil)
)
