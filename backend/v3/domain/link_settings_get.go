package domain

import (
	"context"
	"errors"
)

// -------------------------------------------
// TODO(wim): remove these in favor of the models in the data layer
// -------------------------------------------

type LinkSettings struct {
	Links  []Link
	Source SettingsSource
}

type Link struct {
	Type           LinkType   `json:"type"`
	Url            string     `json:"url"`
	TranslationKey string     `json:"translation_key"`
	Target         LinkTarget `json:"target"`
}

type LinkType int

const (
	LinkTypeUnspecified = iota
	LinkTypeTermsOfService
	LinkTypePrivacyPolicy
	LinkTypeHelp
	LinkTypeSupport
	LinkTypeDocs
	LinkTypeCustom
)

type LinkTarget int

const (
	LinkTargetUnspecified = iota
	LinkTargetSelf
	LinkTargetBlank
)

type SettingsSource int

const (
	SettingsSourceUnspecified = iota
	SettingsSourceSystem
	SettingsSourceInstance
	SettingsSourceOrganization
	SettingsSourceProject
	SettingsSourceApplication
)

// -------------------------------------------
// QUERY
// -------------------------------------------

type GetLinkSettingsQuery struct {
	instanceID     string
	organizationID string
	result         *LinkSettings
}

func NewGetLinkSettingsQuery(instanceID string, organizationID string) *GetLinkSettingsQuery {
	return &GetLinkSettingsQuery{
		instanceID:     instanceID,
		organizationID: organizationID,
	}
}

func (q *GetLinkSettingsQuery) Result() *LinkSettings {
	return q.result
}

// Validate implements [Querier].
func (q *GetLinkSettingsQuery) Validate(ctx context.Context, opts *InvokeOpts) error { return nil }

func (q *GetLinkSettingsQuery) Execute(ctx context.Context, opts *InvokeOpts) error {
	// TODO(wim) implement this
	return errors.New("NOT YET IMPLEMENTED")
}

// String implements [Querier].
func (q *GetLinkSettingsQuery) String() string { return "GetLinkSettingsQuery" }

var _ Querier[*LinkSettings] = (*GetLinkSettingsQuery)(nil)
