package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type SAMLApp struct {
	models.ObjectRoot

	AppID       string
	AppName     string
	Metadata    string
	MetadataURL string

	State AppState
}

func (a *SAMLApp) GetApplicationName() string {
	return a.AppName
}

func (a *SAMLApp) GetState() AppState {
	return a.State
}

func (a *SAMLApp) GetMetadata() string {
	return a.Metadata
}

func (a *SAMLApp) GetMetadataURL() string {
	return a.MetadataURL
}


func (a *SAMLApp) IsValid() bool {
	if a.MetadataURL == "" && a.Metadata == ""{
		return false
	}
	return true
}