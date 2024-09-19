package model

import (
	"github.com/zitadel/zitadel/v2/internal/crypto"
	es_models "github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type APIConfig struct {
	es_models.ObjectRoot
	AppID              string
	ClientID           string
	ClientSecret       *crypto.CryptoValue
	ClientSecretString string
	AuthMethodType     APIAuthMethodType
}

type APIAuthMethodType int32

const (
	APIAuthMethodTypeBasic APIAuthMethodType = iota
	APIAuthMethodTypePrivateKeyJWT
)
