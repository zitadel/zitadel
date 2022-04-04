package setup

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caos/zitadel/internal/api/authz"
	command "github.com/caos/zitadel/internal/command/v2"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	crypto_db "github.com/caos/zitadel/internal/crypto/database"
	"github.com/caos/zitadel/internal/eventstore"
)

type DefaultInstance struct {
	InstanceSetup command.InstanceSetup

	userEncryptionKey *crypto.KeyConfig
	masterKey         string
	db                *sql.DB
	es                *eventstore.Eventstore
	iamDomain         string
	defaults          systemdefaults.SystemDefaults
	zitadelRoles      []authz.RoleMapping
}

func (mig *DefaultInstance) Execute(ctx context.Context) error {
	keyStorage, err := crypto_db.NewKeyStorage(mig.db, mig.masterKey)
	if err != nil {
		return fmt.Errorf("cannot start key storage: %w", err)
	}
	userAlg, err := crypto.NewAESCrypto(mig.userEncryptionKey, keyStorage)
	if err != nil {
		return err
	}

	cmd := command.New(mig.es, mig.iamDomain, mig.defaults, userAlg, mig.zitadelRoles)

	_, err = cmd.SetUpInstance(ctx, &mig.InstanceSetup)
	return err
}

func (mig *DefaultInstance) String() string {
	return "02_default_instance"
}
