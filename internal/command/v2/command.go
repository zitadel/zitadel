package command

import (
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/action"
	iam_repo "github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/keypair"
	"github.com/caos/zitadel/internal/repository/org"
	proj_repo "github.com/caos/zitadel/internal/repository/project"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	usr_grant_repo "github.com/caos/zitadel/internal/repository/usergrant"
)

type Command struct {
	es              *eventstore.Eventstore
	userPasswordAlg crypto.HashAlgorithm
	iamDomain       string
	phoneAlg        crypto.EncryptionAlgorithm
	initCodeAlg     crypto.EncryptionAlgorithm
}

func New(es *eventstore.Eventstore, iamDomain string, defaults sd.SystemDefaults) *Command {
	iam_repo.RegisterEventMappers(es)
	org.RegisterEventMappers(es)
	usr_repo.RegisterEventMappers(es)
	usr_grant_repo.RegisterEventMappers(es)
	proj_repo.RegisterEventMappers(es)
	keypair.RegisterEventMappers(es)
	action.RegisterEventMappers(es)

	return &Command{
		es:              es,
		iamDomain:       iamDomain,
		userPasswordAlg: crypto.NewBCrypt(defaults.SecretGenerators.PasswordSaltCost),
	}
}
