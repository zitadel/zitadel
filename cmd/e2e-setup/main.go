package main

import (
	"context"
	"flag"
	"io/ioutil"
	"time"

	internal_authz "github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/internal/eventstore/v1/models"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config"
	"github.com/caos/zitadel/internal/eventstore"
)

type E2EConfig struct {
	Org                            string
	MachineKeyPath                 string
	OrgOwnerPassword               string
	OrgOwnerViewerPassword         string
	OrgProjectCreatorPassword      string
	PasswordComplexityUserPassword string
	LoginPolicyUserPassword        string
}

type setupConfig struct {
	E2E E2EConfig

	Log logging.Config

	Eventstore     types.SQL
	SystemDefaults sd.SystemDefaults
	InternalAuthZ  internal_authz.Config
}

var (
	e2eSetupPaths = config.NewArrayFlags("authz.yaml", "system-defaults.yaml", "setup.yaml", "e2e.yaml")
)

func main() {
	flag.Var(e2eSetupPaths, "setup-files", "paths to the setup files")
	flag.Parse()
	startE2ESetup(e2eSetupPaths.Values())
}

func startE2ESetup(configPaths []string) {

	conf := new(setupConfig)
	err := config.Read(conf, configPaths...)
	logging.Log("MAIN-EAWlt").OnError(err).Fatal("cannot read config")

	ctx := context.Background()

	es, err := eventstore.Start(conf.Eventstore)
	logging.Log("MAIN-wjQ8G").OnError(err).Fatal("cannot start eventstore")

	commands, err := command.StartCommands(
		es,
		conf.SystemDefaults,
		conf.InternalAuthZ,
		nil,
		command.OrgFeatureCheckerFunc(func(_ context.Context, _ string, _ ...string) error { return nil }),
	)
	logging.Log("MAIN-54MLq").OnError(err).Fatal("cannot start command side")

	err = execute(ctx, commands, conf.E2E)
	logging.Log("MAIN-cgZ3p").OnError(err).Errorf("failed to execute commands steps")
}

func execute(ctx context.Context, commands *command.Commands, cfg E2EConfig) error {

	orgOwner := newHuman("org_owner", cfg.OrgOwnerPassword)

	org, err := commands.SetUpOrg(ctx, &domain.Org{
		Name:    cfg.Org,
		Domains: []*domain.OrgDomain{{Domain: "localhost"}},
	}, orgOwner, nil, false)
	if err != nil {
		return err
	}

	// Avoids the MFA nudge
	if _, err = commands.AddLoginPolicy(ctx, org.ResourceOwner, &domain.LoginPolicy{
		AllowUsernamePassword: true,
	}); err != nil {
		return err
	}

	// Avoids the change password screen
	if _, err = commands.ChangePassword(ctx, org.ResourceOwner, orgOwner.AggregateID, cfg.OrgOwnerPassword, cfg.OrgOwnerPassword, ""); err != nil {
		return err
	}

	sa, err := commands.AddMachine(ctx, org.ResourceOwner, &domain.Machine{
		Username:    "e2e",
		Name:        "e2e",
		Description: "User who calls the ZITADEL API for preparing end-to-end tests",
	})
	if err != nil {
		return err
	}

	if _, err = commands.AddOrgMember(ctx, domain.NewMember(org.ResourceOwner, sa.AggregateID, domain.RoleOrgOwner)); err != nil {
		return err
	}

	key, err := commands.AddUserMachineKey(ctx, &domain.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: sa.AggregateID,
		},
		ExpirationDate: time.Now().Add(30 * 24 * time.Hour),
		Type:           domain.AuthNKeyTypeJSON,
	}, org.ResourceOwner)
	if err != nil {
		return err
	}

	json, err := key.MarshalJSON()
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(cfg.MachineKeyPath, json, 0600); err != nil {
		return err
	}

	for _, user := range []struct{ desc, role, pw string }{{
		desc: "org_owner_viewer",
		pw:   cfg.OrgOwnerViewerPassword,
		role: domain.RoleOrgOwner,
	}, {
		desc: "org_project_creator",
		pw:   cfg.OrgProjectCreatorPassword,
		role: domain.RoleOrgProjectCreator,
	}, {
		desc: "login_policy_user",
		pw:   cfg.LoginPolicyUserPassword,
	}, {
		desc: "password_complexity_user",
		pw:   cfg.PasswordComplexityUserPassword,
	}} {

		newHuman, err := commands.AddHuman(ctx, org.ResourceOwner, newHuman(user.desc, user.pw))
		if err != nil {
			return err
		}

		// Avoids the change password screen
		if _, err = commands.ChangePassword(ctx, org.ResourceOwner, newHuman.AggregateID, user.pw, user.pw, ""); err != nil {
			return err
		}

		if user.role != "" {
			if _, err = commands.AddOrgMember(ctx, domain.NewMember(org.ResourceOwner, newHuman.AggregateID, user.role)); err != nil {
				return err
			}
		}
	}
	return nil
}

func newHuman(desc, pw string) *domain.Human {
	return &domain.Human{
		Username: desc + "_user_name",
		Profile: &domain.Profile{
			FirstName: desc + "_first_name",
			LastName:  desc + "_last_name",
		},
		Password: &domain.Password{
			SecretString:   pw,
			ChangeRequired: false,
		},
		Email: &domain.Email{
			EmailAddress:    desc + ".e2e@caos.ch",
			IsEmailVerified: true,
		},
	}
}
