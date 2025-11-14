package domain

import (
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type SAMLApp struct {
	models.ObjectRoot

	AppID        string
	AppName      string
	EntityID     string
	Metadata     []byte
	MetadataURL  *string
	LoginVersion *LoginVersion
	LoginBaseURI *string

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
	if a.MetadataURL != nil {
		return *a.MetadataURL
	}
	return ""
}

func (a *SAMLApp) IsValid() bool {
	if (a.MetadataURL == nil || *a.MetadataURL == "") && a.Metadata == nil {
		return false
	}
	return true
}
