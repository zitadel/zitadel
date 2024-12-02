package setup

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	crypto_db "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type FirstInstance struct {
	InstanceName    string
	DefaultLanguage language.Tag
	Org             command.InstanceOrgSetup
	MachineKeyPath  string
	PatPath         string
	Features        *command.InstanceFeatures

	Skip bool

	instanceSetup     command.InstanceSetup
	userEncryptionKey *crypto.KeyConfig
	smtpEncryptionKey *crypto.KeyConfig
	oidcEncryptionKey *crypto.KeyConfig
	masterKey         string
	db                *database.DB
	es                *eventstore.Eventstore
	defaults          systemdefaults.SystemDefaults
	zitadelRoles      []authz.RoleMapping
	externalDomain    string
	externalSecure    bool
	externalPort      uint16
	domain            string
}

func (mig *FirstInstance) Execute(ctx context.Context, _ eventstore.Event) error {
	if mig.Skip {
		return nil
	}
	keyStorage, err := mig.verifyEncryptionKeys(ctx)
	if err != nil {
		return err
	}
	userAlg, err := crypto.NewAESCrypto(mig.userEncryptionKey, keyStorage)
	if err != nil {
		return err
	}
	smtpEncryption, err := crypto.NewAESCrypto(mig.smtpEncryptionKey, keyStorage)
	if err != nil {
		return err
	}
	oidcEncryption, err := crypto.NewAESCrypto(mig.oidcEncryptionKey, keyStorage)
	if err != nil {
		return err
	}

	cmd, err := command.StartCommands(ctx,
		mig.es,
		connector.Connectors{},
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
		oidcEncryption,
		nil,
		nil,
		nil,
		nil,
		nil,
		0,
		0,
		0,
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
		orgDomain, err := domain.NewIAMDomainName(mig.instanceSetup.Org.Name, mig.instanceSetup.CustomDomain)
		if err != nil {
			return err
		}
		mig.instanceSetup.Org.Human.Username = mig.instanceSetup.Org.Human.Username + "@" + orgDomain
	}
	mig.instanceSetup.Org.Human.Email.Address = mig.instanceSetup.Org.Human.Email.Address.Normalize()
	if mig.instanceSetup.Org.Human.Email.Address == "" {
		mig.instanceSetup.Org.Human.Email.Address = domain.EmailAddress(mig.instanceSetup.Org.Human.Username)
		if !strings.Contains(string(mig.instanceSetup.Org.Human.Email.Address), "@") {
			orgDomain, err := domain.NewIAMDomainName(mig.instanceSetup.Org.Name, mig.instanceSetup.CustomDomain)
			if err != nil {
				return err
			}
			mig.instanceSetup.Org.Human.Email.Address = domain.EmailAddress(mig.instanceSetup.Org.Human.Username + "@" + orgDomain)
		}
	}

	_, token, key, _, err := cmd.SetUpInstance(ctx, &mig.instanceSetup)
	if err != nil {
		return err
	}
	if mig.instanceSetup.Org.Machine != nil &&
		((mig.instanceSetup.Org.Machine.Pat != nil && token == "") ||
			(mig.instanceSetup.Org.Machine.MachineKey != nil && key == nil)) {
		return err
	}
	return mig.outputMachineAuthentication(key, token)
}

func (mig *FirstInstance) verifyEncryptionKeys(ctx context.Context) (*crypto_db.Database, error) {
	keyStorage, err := crypto_db.NewKeyStorage(mig.db, mig.masterKey)
	if err != nil {
		return nil, fmt.Errorf("cannot start key storage: %w", err)
	}
	if err = verifyKey(ctx, mig.userEncryptionKey, keyStorage); err != nil {
		return nil, err
	}
	if err = verifyKey(ctx, mig.smtpEncryptionKey, keyStorage); err != nil {
		return nil, err
	}
	if err = verifyKey(ctx, mig.oidcEncryptionKey, keyStorage); err != nil {
		return nil, err
	}
	return keyStorage, nil
}

func (mig *FirstInstance) outputMachineAuthentication(key *command.MachineKey, token string) error {
	if key != nil {
		keyDetails, err := key.Detail()
		if err != nil {
			return err
		}
		if err := outputStdoutOrPath(mig.MachineKeyPath, string(keyDetails)); err != nil {
			return err
		}
	}
	if token != "" {
		if err := outputStdoutOrPath(mig.PatPath, token); err != nil {
			return err
		}
	}
	return nil
}

func outputStdoutOrPath(path string, content string) (err error) {
	f := os.Stdout
	if path != "" {
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	_, err = fmt.Fprintln(f, content)
	return err
}

func (mig *FirstInstance) String() string {
	return "03_default_instance"
}

func verifyKey(ctx context.Context, key *crypto.KeyConfig, storage crypto.KeyStorage) (err error) {
	_, err = crypto.LoadKey(key.EncryptionKeyID, storage)
	if err == nil {
		return nil
	}
	k, err := crypto.NewKey(key.EncryptionKeyID)
	if err != nil {
		return err
	}
	return storage.CreateKeys(ctx, k)
}
