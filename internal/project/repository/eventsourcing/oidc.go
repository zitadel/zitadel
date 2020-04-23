package eventsourcing

import (
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/sony/sonyflake"
	"strings"
)

//ClientID random_number@projectname (eg. 495894098234@zitadel)
func generateNewClientID(idGenerator *sonyflake.Sonyflake, project *model.Project) (string, error) {
	rndID, err := idGenerator.NextID()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v@%v", rndID, strings.ReplaceAll(strings.ToLower(project.Name), " ", "_")), nil
}

func generateNewClientSecret(pwGenerator crypto.Generator) (string, *crypto.CryptoValue, error) {
	cryptoValue, stringSecret, err := crypto.NewCode(pwGenerator)
	if err != nil {
		logging.Log("APP-UpnTI").OnError(err).Error("unable to create client secret")
		return "", nil, errors.ThrowInternal(err, "APP-gH2Wl", "unable to create password")
	}
	return stringSecret, cryptoValue, nil
}
