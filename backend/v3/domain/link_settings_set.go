package domain

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	linkSettingsEventPrefix      = "link_settings."
	LinkSettingsChangedEventType = linkSettingsEventPrefix + "changed"
)

type SetLinkSettingsCommand struct {
	Instance
	InstanceID     string `json:"instance_id"`
	OrganizationID string `json:"organization_id"`
	Links          []Link `json:"links"`
	changeTime     time.Time
}

type LinkSettingsSetEvent struct {
	eventstore.BaseEvent `json:"-"`
	InstanceID           string `json:"instance_id"`
	OrganizationID       string `json:"organization_id"`
	Links                []Link `json:"links"`
}

func NewLinkSettingsChangedEvent(
	base *eventstore.BaseEvent,
	instanceID string,
	organizationID string,
	links []Link,
) *LinkSettingsSetEvent {
	return &LinkSettingsSetEvent{
		BaseEvent:      *base,
		InstanceID:     instanceID,
		OrganizationID: organizationID,
		Links:          links,
	}
}

func (e *LinkSettingsSetEvent) Payload() interface{} {
	return e
}

func (e *LinkSettingsSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint { return nil }

func NewSetLinkSettingsCommand(instanceID string, organizationID string, links []Link) *SetLinkSettingsCommand {
	return &SetLinkSettingsCommand{
		InstanceID:     instanceID,
		OrganizationID: organizationID,
		Links:          links,
	}
}

func (cmd *SetLinkSettingsCommand) RequiresTransaction() {}

func (cmd *SetLinkSettingsCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
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
			cmd.Links,
		),
	}, nil
}

// Validate implements [Commander].
func (cmd *SetLinkSettingsCommand) Validate(ctx context.Context, opts *InvokeOpts) error {
	foundLinkTypes := make(map[LinkType]bool)

	for _, l := range cmd.Links {
		switch l.Type {
		case LinkTypeUnspecified:
			return zerrors.ThrowInvalidArgument(nil, "ccHyiN", "unspecified links are not allowed")
		case LinkTypeCustom:
			if l.TranslationKey == "" {
				return zerrors.ThrowInvalidArgument(nil, "RZp5mC", "a custom link should have translation key")
			}
			// don't check for duplicates since multiple custom links are allowed
		default:
			if _, ok := foundLinkTypes[l.Type]; ok {
				return zerrors.ThrowInvalidArgument(nil, "tSuHhT", "each link type is only allowed once")
			}
			foundLinkTypes[l.Type] = true
		}
	}

	return nil
}

func (cmd *SetLinkSettingsCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Commander].
func (cmd *SetLinkSettingsCommand) String() string { return "SetLinkSettingsCommand" }

func (cmd *SetLinkSettingsCommand) Result() time.Time {
	return cmd.changeTime
}

var (
	_ Commander     = (*SetLinkSettingsCommand)(nil)
	_ Transactional = (*SetLinkSettingsCommand)(nil)
)
