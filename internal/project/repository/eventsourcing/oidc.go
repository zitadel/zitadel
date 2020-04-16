package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/sethvargo/go-password/password"
)

////ClientID random_number@project.domain (eg. 495894098234@citadel.caos.ch)
//func generateNewClientID(projectID string) (string, error) {
//	rndID, err := a.repository.GetRandomID()
//	if err != nil {
//		return "", err
//	}
//	//TODO: Get domain from ctx
//	domainID := ""
//	domain, err := a.repository.DomainByID(ctx, domainID)
//	if err != nil {
//		return "", err
//	}
//
//	project, err := a.repository.ProjectByID(ctx, projectID)
//	if err != nil {
//		return "", err
//	}
//
//	return fmt.Sprintf("%v@%v.%v", rndID, strings.ReplaceAll(project.Name, " ", "_"), domain.Name), nil
//}

func generateNewClientSecret(pwGenerator password.PasswordGenerator, alg crypto.HashAlgorithm) (string, *crypto.CryptoValue, error) {
	stringSecret, err := pwGenerator.Generate(64, 10, 10, false, false)
	if err != nil {
		logging.Log("APP-UpnTI").OnError(err).Error("unable to create client secret")
		return "", nil, errors.ThrowInternal(err, "APP-gH2Wl", "unable to create password")
	}
	secret, err := crypto.Hash([]byte(stringSecret), alg)
	if err != nil {
		return "", nil, errors.ThrowInternal(err, "APP-gH2Wl", "unable to hash password")
	}
	return stringSecret, secret, nil
}
