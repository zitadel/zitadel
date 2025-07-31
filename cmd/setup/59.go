package setup

import (
	"context"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type SetupWebkeys struct {
	eventstore *eventstore.Eventstore
	commands   *command.Commands
}

func (mig *SetupWebkeys) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceAddedEventType).
			Builder().ExcludeAggregateIDs().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceRemovedEventType).
			Builder(),
	)
	if err != nil {
		return fmt.Errorf("%s get instance IDs: %w", mig, err)
	}
	conf := &crypto.WebKeyRSAConfig{
		Bits:   crypto.RSABits2048,
		Hasher: crypto.RSAHasherSHA256,
	}

	for _, instance := range instances {
		ctx := authz.WithInstanceID(ctx, instance)
		logging.Info("prepare initial webkeys for instance", "instance_id", instance, "migration", mig)
		if err := mig.commands.GenerateInitialWebKeys(ctx, conf); err != nil {
			return fmt.Errorf("%s generate initial webkeys: %w", mig, err)
		}
	}
	return nil
}

func (mig *SetupWebkeys) String() string {
	return "59_setup_webkeys"
}
