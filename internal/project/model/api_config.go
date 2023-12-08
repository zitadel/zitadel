package model

import (
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (c *APIConfig) IsValid() bool {
	return true
}

// ClientID random_number@projectname (eg. 495894098234@zitadel)
func (c *APIConfig) GenerateNewClientID(idGenerator id.Generator, project *Project) error {
	rndID, err := idGenerator.Next()
	if err != nil {
		return err
	}

	c.ClientID = fmt.Sprintf("%v@%v", rndID, strings.ReplaceAll(strings.ToLower(project.Name), " ", "_"))
	return nil
}

func (c *APIConfig) GenerateClientSecretIfNeeded(generator crypto.Generator) (string, error) {
	if c.AuthMethodType == APIAuthMethodTypeBasic {
		return c.GenerateNewClientSecret(generator)
	}
	return "", nil
}

func (c *APIConfig) GenerateNewClientSecret(generator crypto.Generator) (string, error) {
	cryptoValue, stringSecret, err := crypto.NewCode(generator)
	if err != nil {
		logging.Log("MODEL-ADvd2").OnError(err).Error("unable to create client secret")
		return "", zerrors.ThrowInternal(err, "MODEL-dsvr43", "Errors.Project.CouldNotGenerateClientSecret")
	}
	c.ClientSecret = cryptoValue
	return stringSecret, nil
}
