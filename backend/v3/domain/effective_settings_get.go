package domain

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

type EffectiveSettings struct {
	Links []Link
}

type GetEffectiveSettingsQuery struct {
	OrganizationID string
	ProjectID      string
	ApplicationID  string
	Types          []settings.SettingsType

	result EffectiveSettings
}

func NewGetEffectiveSettingsQuery(organizationID string, projectID string, applicationID string, types []settings.SettingsType) *GetEffectiveSettingsQuery {
	return &GetEffectiveSettingsQuery{
		OrganizationID: organizationID,
		ProjectID:      projectID,
		ApplicationID:  applicationID,
		Types:          types,
	}
}

func (q *GetEffectiveSettingsQuery) Result() EffectiveSettings {
	return q.result
}

// Validate implements [Querier].
func (q *GetEffectiveSettingsQuery) Validate(ctx context.Context, opts *InvokeOpts) error { return nil }

func (q *GetEffectiveSettingsQuery) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Querier].
func (q *GetEffectiveSettingsQuery) String() string { return "GetLinkSettingsQuery" }

var _ Querier[EffectiveSettings] = (*GetEffectiveSettingsQuery)(nil)
