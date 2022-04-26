package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	authz_repo "github.com/zitadel/zitadel/internal/authz/repository"

	"github.com/zitadel/zitadel/internal/domain"

	"github.com/caos/logging"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type user struct {
	desc, role, pw string
}

func main() {
	flag.Var(e2eSetupPaths, "setup-files", "paths to the setup files")
	debug := flag.Bool("debug", false, "print information that is helpful for debugging")
	flag.Parse()
	startE2ESetup(e2eSetupPaths.Values(), *debug)
}

type dummyAuthZRepo struct {
	authz_repo.Repository
}

func (d dummyAuthZRepo) CheckOrgFeatures(_ context.Context, _ string, _ ...string) error {
	return nil
}

func startE2ESetup(configPaths []string, debug bool) {

	conf := new(setupConfig)
	err := config.Read(conf, configPaths...)
	logging.Log("MAIN-EAWlt").OnError(err).Fatal("cannot read config")

	if debug {
		printConfig("e2e", conf.E2E)
		printConfig("system defaults", conf.SystemDefaults)
		printConfig("authz", conf.InternalAuthZ)
		printConfig("eventstore", conf.Eventstore)
		printConfig("log", conf.Log)
	}
	err = conf.E2E.validate()
	logging.Log("MAIN-NoZIV").OnError(err).Fatal("validating e2e config failed")

	ctx := context.Background()

	es, err := eventstore.Start(conf.Eventstore)
	logging.Log("MAIN-wjQ8G").OnError(err).Fatal("cannot start eventstore")

	commands, err := command.StartCommands(
		es,
		conf.SystemDefaults,
		conf.InternalAuthZ,
		nil,
		dummyAuthZRepo{},
	)
	logging.Log("MAIN-54MLq").OnError(err).Fatal("cannot start command side")

	users := []user{{
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

	err = execute(ctx, commands, conf.E2E, users)
	logging.Log("MAIN-cgZ3p").OnError(err).Errorf("failed to execute commands steps")

	eventualConsistencyCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	err = awaitConsistency(
		eventualConsistencyCtx,
		conf.E2E,
		users,
	)
	logging.Log("MAIN-cgZ3p").OnError(err).Fatal("failed to await consistency")
}

func printConfig(desc string, cfg interface{}) {
	bytes, err := yaml.Marshal(cfg)
	logging.Log("MAIN-JYmQq").OnError(err).Fatal("cannot marshal config")

	logging.Log("MAIN-7u4dZ").Info("got the following ", desc, " config")
	fmt.Println(string(bytes))
}
