package model

import (
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Application struct {
	es_models.ObjectRoot

	AppID      string
	State      AppState
	Name       string
	Type       AppType
	OIDCConfig *OIDCConfig
	APIConfig  *APIConfig
	SAMLConfig *SAMLConfig
}

type AppState int32

const (
	AppStateActive AppState = iota
	AppStateInactive
	AppStateRemoved
)

type AppType int32

const (
	AppTypeUnspecified AppType = iota
	AppTypeOIDC
	AppTypeSAML
	AppTypeAPI
)
