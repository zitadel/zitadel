package domain

import (
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type oAuthApplication interface {
	setClientID(clientID string)
	setClientSecret(secret *crypto.CryptoValue)
	requiresClientSecret() bool
}

// ClientID random_number@projectname (eg. 495894098234@zitadel)
func SetNewClientID(a oAuthApplication, idGenerator id.Generator, project *Project) error {
	clientID, err := NewClientID(idGenerator, project.Name)
	if err != nil {
		return err
	}

	a.setClientID(clientID)
	return nil
}

func NewClientID(idGenerator id.Generator, projectName string) (string, error) {
	rndID, err := idGenerator.Next()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s@%s", rndID, strings.ReplaceAll(strings.ToLower(projectName), " ", "_")), nil
}

func SetNewClientSecretIfNeeded(a oAuthApplication, generator crypto.Generator) (string, error) {
	if !a.requiresClientSecret() {
		return "", nil
	}
	clientSecret, secretString, err := NewClientSecret(generator)
	if err != nil {
		return "", err
	}
	a.setClientSecret(clientSecret)
	return secretString, nil
}

func NewClientSecret(generator crypto.Generator) (*crypto.CryptoValue, string, error) {
	cryptoValue, stringSecret, err := crypto.NewCode(generator)
	if err != nil {
		logging.Log("MODEL-UpnTI").OnError(err).Error("unable to create client secret")
		return nil, "", zerrors.ThrowInternal(err, "MODEL-gH2Wl", "Errors.Project.CouldNotGenerateClientSecret")
	}
	return cryptoValue, stringSecret, nil
}
