package domain

import (
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/id"
)

type oAuthApplication interface {
	setClientID(clientID string)
	setClientSecret(encodedHash string)
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

func SetNewClientSecretIfNeeded(a oAuthApplication, generator *crypto.HashGenerator) (string, error) {
	if !a.requiresClientSecret() {
		return "", nil
	}
	encodedHash, plain, err := generator.NewCode()
	if err != nil {
		return "", err
	}
	a.setClientSecret(encodedHash)
	return plain, nil
}
