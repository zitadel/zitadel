package main

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

func execute(ctx context.Context, cmd *command.Commands, cfg E2EConfig, users []userData, instanceID string) error {

	ctx = authz.WithInstanceID(ctx, instanceID)
	ctx = authz.WithRequestedDomain(ctx, "localhost")

	orgOwner := newHuman(users[0])

	baseUrl, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return err
	}

	orgOwnerID, org, err := cmd.SetUpOrg(ctx, &command.OrgSetup{
		Name:         cfg.Org,
		CustomDomain: baseUrl.Host,
		Human:        *orgOwner,
	})
	if err != nil {
		// TODO: Why is this error not typed?
		if strings.Contains(err.Error(), "Errors.Org.AlreadyExists") {
			logging.New().Info("Looks like setup is already done")
			err = nil
		}
		return err
	}

	// Avoids the MFA nudge
	if _, err = cmd.AddLoginPolicy(ctx, org.ResourceOwner, &domain.LoginPolicy{
		AllowUsernamePassword:      true,
		ExternalLoginCheckLifetime: 24 * 365 * time.Hour, // 1 year
		MFAInitSkipLifetime:        24 * 365 * time.Hour, // 1 year
		MultiFactorCheckLifetime:   24 * 365 * time.Hour, // 1 year
		PasswordCheckLifetime:      24 * 365 * time.Hour, // 1 year
		SecondFactorCheckLifetime:  24 * 365 * time.Hour, // 1 year
	}); err != nil {
		return err
	}

	if err = initHuman(ctx, cmd, orgOwnerID, users[0], org.ResourceOwner); err != nil {
		return err
	}

	sa, err := cmd.AddMachine(ctx, org.ResourceOwner, &domain.Machine{
		Username:    "e2e",
		Name:        "e2e",
		Description: "User who calls the ZITADEL API for preparing end-to-end tests",
	})
	if err != nil {
		return err
	}

	if _, err = cmd.AddOrgMember(ctx, org.ResourceOwner, sa.AggregateID, domain.RoleOrgOwner); err != nil {
		return err
	}

	key, err := cmd.AddUserMachineKey(ctx, &domain.MachineKey{
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

	if err = os.MkdirAll(filepath.Dir(cfg.MachineKeyPath), 0700); err != nil {
		return err
	}

	if err = ioutil.WriteFile(cfg.MachineKeyPath, json, 0600); err != nil {
		return err
	}

	for idx := range users[1:] {
		user := users[idx+1]

		createdHuman, err := cmd.AddHuman(ctx, org.ResourceOwner, newHuman(user))
		if err != nil {
			return err
		}

		if err = initHuman(ctx, cmd, createdHuman.ID, user, org.ResourceOwner); err != nil {
			return err
		}

		if user.role != "" {
			if _, err := cmd.AddOrgMember(ctx, org.ResourceOwner, createdHuman.ID, user.role); err != nil {
				return err
			}
		}
	}
	return nil
}

func newHuman(u userData) *command.AddHuman {
	return &command.AddHuman{
		Username:  u.desc + "_user_name",
		FirstName: u.desc + "_first_name",
		LastName:  u.desc + "_last_name",
		Password:  u.pw,
		Email: command.Email{
			Address:  u.desc + ".e2e@zitadel.com",
			Verified: true,
		},
		PasswordChangeRequired: false,
		Register:               false,
	}
}

// initHuman skips the MFA and change password screens
func initHuman(ctx context.Context, cmd *command.Commands, userID string, user userData, orgID string) error {
	// skip mfa
	if err := cmd.HumanSkipMFAInit(ctx, userID, orgID); err != nil {
		return err
	}

	// Avoids the change password screen
	_, err := cmd.ChangePassword(ctx, orgID, userID, user.pw, user.pw, "")
	return err
}
