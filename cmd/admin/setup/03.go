package setup

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	crypto_db "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type DefaultInstance struct {
	InstanceName    string
	CustomDomain    string
	DefaultLanguage language.Tag
	Org             command.OrgSetup

	instanceSetup     command.InstanceSetup
	userEncryptionKey *crypto.KeyConfig
	masterKey         string
	db                *sql.DB
	es                *eventstore.Eventstore
	defaults          systemdefaults.SystemDefaults
	zitadelRoles      []authz.RoleMapping
	externalDomain    string
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
		mig.externalDomain,
		mig.externalSecure,
		mig.externalPort,
		nil,
		nil,
		nil,
		nil,
		userAlg,
		nil,
		nil,
		nil)

	if err != nil {
		return err
	}

	mig.instanceSetup.InstanceName = mig.InstanceName
	mig.instanceSetup.CustomDomain = mig.CustomDomain
	mig.instanceSetup.DefaultLanguage = mig.DefaultLanguage
	mig.instanceSetup.Org = mig.Org
	mig.instanceSetup.Org.Human.Email.Address = strings.TrimSpace(mig.instanceSetup.Org.Human.Email.Address)
	if mig.instanceSetup.Org.Human.Email.Address == "" {
		mig.instanceSetup.Org.Human.Email.Address = "admin@" + mig.instanceSetup.CustomDomain
	}

	_, _, err = cmd.SetUpInstance(ctx, &mig.instanceSetup)
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
