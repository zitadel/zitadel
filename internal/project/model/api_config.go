package model

import (
	"fmt"
	"strings"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
)

type APIConfig struct {
	es_models.ObjectRoot
	AppID              string
	ClientID           string
	ClientSecret       *crypto.CryptoValue
	ClientSecretString string
	AuthMethodType     APIAuthMethodType
	ClientKeys         []*ClientKey
}

type APIAuthMethodType int32

const (
	APIAuthMethodTypeBasic APIAuthMethodType = iota
	APIAuthMethodTypePrivateKeyJWT
)

func (c *APIConfig) IsValid() bool {
	return true
}

//ClientID random_number@projectname (eg. 495894098234@zitadel)
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
		return "", errors.ThrowInternal(err, "MODEL-dsvr43", "Errors.Project.CouldNotGenerateClientSecret")
	}
	c.ClientSecret = cryptoValue
	return stringSecret, nil
}
