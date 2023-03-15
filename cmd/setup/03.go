package setup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	crypto_db "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type FirstInstance struct {
	InstanceName    string
	DefaultLanguage language.Tag
	Org             command.OrgSetup
	MachineKeyPath  string

	instanceSetup     command.InstanceSetup
	userEncryptionKey *crypto.KeyConfig
	smtpEncryptionKey *crypto.KeyConfig
	masterKey         string
	db                *sql.DB
	es                *eventstore.Eventstore
	defaults          systemdefaults.SystemDefaults
	zitadelRoles      []authz.RoleMapping
	externalDomain    string
	externalSecure    bool
	externalPort      uint16
	domain            string
}

func (mig *FirstInstance) Execute(ctx context.Context) error {
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

	if err = verifyKey(mig.smtpEncryptionKey, keyStorage); err != nil {
		return err
	}
	smtpEncryption, err := crypto.NewAESCrypto(mig.smtpEncryptionKey, keyStorage)
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
		smtpEncryption,
		nil,
		userAlg,
		nil,
		nil,
		nil,
		nil,
	)

	if err != nil {
		return err
	}

	mig.instanceSetup.InstanceName = mig.InstanceName
	mig.instanceSetup.CustomDomain = mig.externalDomain
	mig.instanceSetup.DefaultLanguage = mig.DefaultLanguage
	mig.instanceSetup.Org = mig.Org
	// check if username is email style or else append @<orgname>.<custom-domain>
	//this way we have the same value as before changing `UserLoginMustBeDomain` to false
	if !mig.instanceSetup.DomainPolicy.UserLoginMustBeDomain && !strings.Contains(mig.instanceSetup.Org.Human.Username, "@") {
		mig.instanceSetup.Org.Human.Username = mig.instanceSetup.Org.Human.Username + "@" + domain.NewIAMDomainName(mig.instanceSetup.Org.Name, mig.instanceSetup.CustomDomain)
	}
	mig.instanceSetup.Org.Human.Email.Address = mig.instanceSetup.Org.Human.Email.Address.Normalize()
	if mig.instanceSetup.Org.Human.Email.Address == "" {
		mig.instanceSetup.Org.Human.Email.Address = domain.EmailAddress(mig.instanceSetup.Org.Human.Username)
		if !strings.Contains(string(mig.instanceSetup.Org.Human.Email.Address), "@") {
			mig.instanceSetup.Org.Human.Email.Address = domain.EmailAddress(mig.instanceSetup.Org.Human.Username + "@" + domain.NewIAMDomainName(mig.instanceSetup.Org.Name, mig.instanceSetup.CustomDomain))
		}
	}

	_, _, key, _, err := cmd.SetUpInstance(ctx, &mig.instanceSetup)
	if key == nil {
		return err
	}

	f := os.Stdout
	if mig.MachineKeyPath != "" {
		f, err = os.OpenFile(mig.MachineKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	keyDetails, err := key.Detail()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, string(keyDetails))
	return err
}

func (mig *FirstInstance) String() string {
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
