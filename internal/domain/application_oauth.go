package domain

import (
	"fmt"
	"strings"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/id"
)

type oAuthApplication interface {
	setClientID(clientID string)
	setClientSecret(secret *crypto.CryptoValue)
	requiresClientSecret() bool
}

//ClientID random_number@projectname (eg. 495894098234@zitadel)
func SetNewClientID(a oAuthApplication, idGenerator id.Generator, project *Project) error {
	rndID, err := idGenerator.Next()
	if err != nil {
		return err
	}

	a.setClientID(fmt.Sprintf("%v@%v", rndID, strings.ReplaceAll(strings.ToLower(project.Name), " ", "_")))
	return nil
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
		return nil, "", errors.ThrowInternal(err, "MODEL-gH2Wl", "Errors.Project.CouldNotGenerateClientSecret")
	}
	return cryptoValue, stringSecret, nil
}
