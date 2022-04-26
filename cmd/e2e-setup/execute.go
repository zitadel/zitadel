package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

func execute(ctx context.Context, commands *command.Commands, cfg E2EConfig, users []user) error {

	orgOwner := newHuman(users[0])

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

	if err = os.MkdirAll(filepath.Dir(cfg.MachineKeyPath), 0700); err != nil {
		return err
	}

	if err = ioutil.WriteFile(cfg.MachineKeyPath, json, 0600); err != nil {
		return err
	}

	for _, user := range users[1:] {

		newHuman, err := commands.AddHuman(ctx, org.ResourceOwner, newHuman(user))
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

func newHuman(u user) *domain.Human {
	return &domain.Human{
		Username: u.desc + "_user_name",
		Profile: &domain.Profile{
			FirstName: u.desc + "_first_name",
			LastName:  u.desc + "_last_name",
		},
		Password: &domain.Password{
			SecretString:   u.pw,
			ChangeRequired: false,
		},
		Email: &domain.Email{
			EmailAddress:    u.desc + ".e2e@caos.ch",
			IsEmailVerified: true,
		},
	}
}
