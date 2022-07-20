package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"time"

	cryptoDB "github.com/zitadel/zitadel/internal/crypto/database"

	"github.com/zitadel/zitadel/internal/id"

	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/internal/config/options"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/webauthn"

	"gopkg.in/yaml.v3"

	"github.com/zitadel/zitadel/internal/domain"

	"github.com/zitadel/logging"
)

var (
	//go:embed defaults.yaml
	e2edefaults []byte
)

type userData struct {
	desc, role, pw string
}

func main() {
	masterkey := flag.String("materkey", "MasterkeyNeedsToHave32Characters", "the ZITADEL installations masterkey")
	debug := flag.Bool("debug", false, "print information that is helpful for debugging")

	err := options.InitViper()
	logging.OnError(err).Fatalf("unable to initialize zitadel config: %s", err)

	flag.Parse()

	err = viper.MergeConfig(bytes.NewBuffer(e2edefaults))
	logging.OnError(err).Fatalf("unable to initialize e2e config: %s", err)

	conf := MustNewConfig(viper.GetViper())

	if *debug {
		printConfig("config", conf)
	}

	logging.New().OnError(err).Fatal("validating e2e config failed")

	startE2ESetup(conf, *masterkey)
}

func startE2ESetup(conf *Config, masterkey string) {

	id.Configure(conf.Machine)

	ctx := context.Background()

	dbClient, err := database.Connect(conf.Database)
	logging.New().OnError(err).Fatalf("cannot start client for projection: %s", err)

	zitadelProjectResourceID, instanceID, err := ids(ctx, conf.E2E, dbClient)
	logging.New().OnError(err).Fatalf("cannot get instance and project IDs: %s", err)

	keyStorage, err := cryptoDB.NewKeyStorage(dbClient, masterkey)
	logging.New().OnError(err).Fatalf("cannot start key storage: %s", err)

	keys, err := ensureEncryptionKeys(conf.EncryptionKeys, keyStorage)
	logging.New().OnError(err).Fatalf("failed ensuring encryption keys: %s", err)
	eventstoreClient, err := eventstore.Start(dbClient)
	logging.New().OnError(err).Fatalf("cannot start eventstore for queries: %s", err)

	storage, err := conf.AssetStorage.NewStorage(dbClient)
	logging.New().OnError(err).Fatalf("cannot start asset storage client: %s", err)

	webAuthNConfig := &webauthn.Config{
		DisplayName:    conf.WebAuthNName,
		ExternalSecure: conf.ExternalSecure,
	}

	commands, err := command.StartCommands(
		eventstoreClient,
		conf.SystemDefaults,
		conf.InternalAuthZ.RolePermissionMappings,
		storage,
		webAuthNConfig,
		conf.ExternalDomain,
		conf.ExternalSecure,
		conf.ExternalPort,
		keys.IDPConfig,
		keys.OTP,
		keys.SMTP,
		keys.SMS,
		keys.User,
		keys.DomainVerification,
		keys.OIDC,
	)
	logging.New().OnError(err).Errorf("cannot start commands: %s", err)

	users := []userData{{
		desc: "org_owner",
		pw:   conf.E2E.OrgOwnerPassword,
		role: domain.RoleOrgOwner,
	}, {
		desc: "org_owner_viewer",
		pw:   conf.E2E.OrgOwnerViewerPassword,
		role: domain.RoleOrgOwner,
	}, {
		desc: "org_project_creator",
		pw:   conf.E2E.OrgProjectCreatorPassword,
		role: domain.RoleOrgProjectCreator,
	}, {
		desc: "login_policy_user",
		pw:   conf.E2E.LoginPolicyUserPassword,
	}, {
		desc: "password_complexity_user",
		pw:   conf.E2E.PasswordComplexityUserPassword,
	}}

	err = execute(ctx, commands, *conf.E2E, users, instanceID)
	logging.New().OnError(err).Fatalf("failed to execute commands steps")

	eventualConsistencyCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	err = awaitConsistency(
		eventualConsistencyCtx,
		*conf.E2E,
		users,
		zitadelProjectResourceID,
	)
	logging.New().OnError(err).Fatal("failed to await consistency")
}

func printConfig(desc string, cfg interface{}) {
	bytes, err := yaml.Marshal(cfg)
	logging.New().OnError(err).Fatal("cannot marshal config")

	logging.New().Info("got the following ", desc, " config")
	fmt.Println(string(bytes))
}
