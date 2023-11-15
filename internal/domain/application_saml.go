package domain

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type SAMLApp struct {
	models.ObjectRoot

	AppID       string
	AppName     string
	EntityID    string
	Metadata    []byte
	MetadataURL string

	State AppState
}

func (a *SAMLApp) GetApplicationName() string {
	return a.AppName
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
	return true
}
