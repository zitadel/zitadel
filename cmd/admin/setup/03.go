package setup

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
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
	domain            string
	defaults          systemdefaults.SystemDefaults
	zitadelRoles      []authz.RoleMapping
	externalSecure    bool
	externalPort      uint16
}

func (mig *DefaultInstance) Execute(ctx context.Context) error {
	keyStorage, err := crypto_db.NewKeyStorage(mig.db, mig.masterKey)
	if err != nil {
		return fmt.Errorf("cannot start key storage: %w", err)
	}
	if err = verifyKey(mig.userEncryptionKey, keyStorage); err != nil {
		return err
	}

	userAlg, err := crypto.NewAESCrypto(mig.userEncryptionKey, keyStorage)
	if err != nil {
		return err
	}

	cmd, err := command.StartCommands(mig.es,
		mig.defaults,
		mig.zitadelRoles,
		nil,
		nil,
		nil,
		mig.externalSecure,
		mig.externalPort,
		nil,
		nil,
		nil,
		nil,
		userAlg,
		nil,
		nil)

	if err != nil {
		return err
	}
	ctx = authz.WithRequestedDomain(ctx, mig.domain)

	_, _, err = cmd.SetUpInstance(ctx, &mig.InstanceSetup, mig.externalSecure)
	return err
}

func (mig *DefaultInstance) String() string {
	return "03_default_instance"
}

func verifyKey(key *crypto.KeyConfig, storage crypto.KeyStorage) (err error) {
	_, err = crypto.LoadKey(key.EncryptionKeyID, storage)
	if err == nil {
		return nil
	}
	k, err := crypto.NewKey(key.EncryptionKeyID)
	if err != nil {
		return err
	}
	return storage.CreateKeys(k)
}
