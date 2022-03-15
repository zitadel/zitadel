package command

import (
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
)

type Command struct {
	es              *eventstore.Eventstore
	userPasswordAlg crypto.HashAlgorithm
	iamDomain       string
}

func New(es *eventstore.Eventstore, iamDomain string, defaults sd.SystemDefaults) *Command {
	return &Command{
		es:              es,
		iamDomain:       iamDomain,
		userPasswordAlg: crypto.NewBCrypt(defaults.SecretGenerators.PasswordSaltCost),
	}
}
