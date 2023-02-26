package domain

import (
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type SAMLApp struct {
	models.ObjectRoot

	AppID       string
	AppName     string
	ExternalURL string
	IsVisibleToEndUser bool
	EntityID    string
	Metadata    []byte
	MetadataURL string

	State AppState
}

func (a *SAMLApp) GetApplicationName() string {
	return a.AppName
}

func (a *SAMLApp) GetApplicationExternalURL() string {
	return a.ExternalURL
}

func (a *SAMLApp) GetApplicationIsVisibleToEndUser() bool {
	return a.IsVisibleToEndUser
}

func (a *SAMLApp) GetState() AppState {
	return a.State
}

func (a *SAMLApp) GetMetadata() []byte {
	return a.Metadata
}

func (a *SAMLApp) GetMetadataURL() string {
	return a.MetadataURL
}

func (a *SAMLApp) IsValid() bool {
	if a.MetadataURL == "" && a.Metadata == nil {
		return false
	}
	if a.ExternalURL != "" && !IsValidURL(a.ExternalURL) {
		return false
	}
	return true
}
