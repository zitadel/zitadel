package domain

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type ResetLinkSettingsCommand struct {
	InstanceID     string `json:"instance_id"`
	OrganizationID string `json:"organization_id"`
	result         *ResetLinkSettingsCommandResult
}

type ResetLinkSettingsCommandResult struct {
	ChangeTime time.Time
	Links      []Link
}

func NewResetLinkSettingsCommand(instanceID string, organizationID string) *ResetLinkSettingsCommand {
	return &ResetLinkSettingsCommand{
		InstanceID:     instanceID,
		OrganizationID: organizationID,
	}
}

func (cmd *ResetLinkSettingsCommand) RequiresTransaction() {}

func (cmd *ResetLinkSettingsCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	var agg eventstore.Aggregate
	if cmd.OrganizationID == "" {
		agg = org.NewAggregate(cmd.OrganizationID).Aggregate
	} else {
		agg = instance.NewAggregate(authz.GetInstance(ctx).InstanceID()).Aggregate
	}

	return []eventstore.Command{
		NewLinkSettingsChangedEvent(
			eventstore.NewBaseEventForPush(ctx, &agg, LinkSettingsChangedEventType),
			cmd.InstanceID,
			cmd.OrganizationID,
			nil,
		),
	}, nil
}

// Validate implements [Commander].
func (cmd *ResetLinkSettingsCommand) Validate(ctx context.Context, opts *InvokeOpts) error {
	return nil
}

func (cmd *ResetLinkSettingsCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Commander].
func (cmd *ResetLinkSettingsCommand) String() string { return "ResetLinkSettingsCommand" }

func (cmd *ResetLinkSettingsCommand) Result() *ResetLinkSettingsCommandResult {
	return cmd.result
}

var (
	_ Commander     = (*ResetLinkSettingsCommand)(nil)
	_ Transactional = (*ResetLinkSettingsCommand)(nil)
)
