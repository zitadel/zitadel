package setup

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type DefaultNotificationPolicies struct {
	PasswordChange bool

	es             *eventstore.Eventstore
	defaults       systemdefaults.SystemDefaults
	zitadelRoles   []authz.RoleMapping
	externalDomain string
	externalSecure bool
	externalPort   uint16
}

func (mig *DefaultNotificationPolicies) Execute(ctx context.Context) error {

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
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	instances, err := command.ListInstances(ctx, mig.es.Filter)
	if err != nil {
		return err
	}

	for _, instanceID := range instances {
		alreadyExists, err := command.ExistsDefaultNotificationPolicy(ctx, mig.es.Filter, instanceID)
		if err != nil {
			return err
		}
		if !alreadyExists {
			_, err := cmd.AddDefaultNotificationPolicy(authz.WithInstanceID(ctx, instanceID), instanceID, mig.PasswordChange)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (mig *DefaultNotificationPolicies) String() string {
	return "07_default_notification_policy"
}
