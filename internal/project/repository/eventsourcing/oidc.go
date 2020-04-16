package eventsourcing

import (
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/sethvargo/go-password/password"
	"github.com/sony/sonyflake"
	"strings"
)

////ClientID random_number@projectname (eg. 495894098234@zitadel)
func generateNewClientID(idGenerator *sonyflake.Sonyflake, project *model.Project) (string, error) {
	rndID, err := idGenerator.NextID()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v@%v", rndID, strings.ReplaceAll(strings.ToLower(project.Name), " ", "_")), nil
}

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
