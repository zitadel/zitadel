package domain

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetLinkSettingsCommand struct {
	Instance       bool   `json:"instance"`
	OrganizationId string `json:"organization_id"`
	Links          []Link `json:"links"`
	changeTime     time.Time
}

func NewSetLinkSettingsCommand(instance bool, organizationId string, links []Link) *SetLinkSettingsCommand {
	return &SetLinkSettingsCommand{
		Instance:       instance,
		OrganizationId: organizationId,
		Links:          links,
	}
}

func (cmd *SetLinkSettingsCommand) RequiresTransaction() {}

func (cmd *SetLinkSettingsCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	// TODO(wim) implement this
	return nil, errors.New("NOT YET IMPLEMENTED")
}

// Validate implements [Commander].
func (q *SetLinkSettingsCommand) Validate(ctx context.Context, opts *InvokeOpts) error {
	foundLinkTypes := make(map[LinkType]bool)

	for _, l := range q.Links {
		switch {
		case LinkTypeUnspecified:
			return zerrors.ThrowInvalidArgument(nil, "ccHyiN", "unspecified links are not allowed")
		case LinkTypeCustom:
			if l.TranslationKey == "" {
				return zerrors.ThrowInvalidArgument(nil, "RZp5mC", "a custom link should have translation key")
			}
		default:
			if _, ok := foundLinkTypes[l.Type]; ok {
				return zerrors.ThrowInvalidArgument(nil, "tSuHhT", "each link type is only allowed once")
			}
		}
	}

	return nil
}

func (q *SetLinkSettingsCommand) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Commander].
func (q *SetLinkSettingsCommand) String() string { return "SetLinkSettingsCommand" }

func (q *SetLinkSettingsCommand) Result() time.Time {
	return q.changeTime
}

var (
	_ Commander     = (*SetLinkSettingsCommand)(nil)
	_ Transactional = (*SetLinkSettingsCommand)(nil)
)
