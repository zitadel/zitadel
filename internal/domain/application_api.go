package domain

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type APIApp struct {
	models.ObjectRoot

	AppID              string
	AppName            string
	ClientID           string
	EncodedHash        string
	ClientSecretString string
	AuthMethodType     APIAuthMethodType

	State AppState
}

func (a *APIApp) GetApplicationName() string {
	return a.AppName
}

func (a *APIApp) GetState() AppState {
	return a.State
}

type APIAuthMethodType int32

const (
	APIAuthMethodTypeBasic APIAuthMethodType = iota
	APIAuthMethodTypePrivateKeyJWT
)

func (a *APIApp) IsValid() bool {
	return a.AppName != ""
}

func (a *APIApp) setClientID(clientID string) {
	a.ClientID = clientID
}

func (a *APIApp) setClientSecret(encodedHash string) {
	a.EncodedHash = encodedHash
}

func (a *APIApp) requiresClientSecret() bool {
	return a.AuthMethodType == APIAuthMethodTypeBasic
}

func (a *APIApp) GenerateClientSecretIfNeeded(generator *crypto.HashGenerator) (plain string, err error) {
	if a.AuthMethodType == APIAuthMethodTypePrivateKeyJWT {
		return "", nil
	}
	a.EncodedHash, plain, err = generator.NewCode()
	if err != nil {
		return "", err
	}
	return plain, nil
}
