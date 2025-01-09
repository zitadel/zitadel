package setup

import (
	"context"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/permission"
)

type AddRolePermissions struct {
	eventstore             *eventstore.Eventstore
	rolePermissionMappings []authz.RoleMapping
}

func (mig *AddRolePermissions) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceAddedEventType).
			Builder().
			ExcludeAggregateIDs().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceRemovedEventType).
			Builder().
			ExcludeAggregateIDs(). // make sure we don't try to re-push if the migration failed half way.
			AggregateTypes(permission.AggregateType).
			EventTypes(permission.AddedType).
			Builder(),
	)
	if err != nil {
		return err
	}
	for i, instanceID := range instances {
		logging.WithFields("instance_id", instanceID, "migration", mig.String(), "progress", fmt.Sprintf("%d/%d", i+1, len(instances))).Info("prepare role permission events")
		cmds := mig.instanceCommands(ctx, instanceID, mig.rolePermissionMappings)
		events, err := mig.eventstore.Push(ctx, cmds...)
		if err != nil {
			return err
		}
		logging.WithFields("instance_id", instanceID, "migration", mig.String(), "pushed_events", len(events)).Info("pushed role permission events")
	}
	return nil
}

func (*AddRolePermissions) instanceCommands(ctx context.Context, instanceID string, mapping []authz.RoleMapping) (cmds []eventstore.Command) {
	for _, m := range mapping {
		aggregate := permission.NewAggregate(instanceID)
		for _, p := range m.Permissions {
			cmds = append(cmds, permission.NewAddedEvent(ctx, aggregate, m.Role, p))
		}
	}
	return cmds
}

func (*AddRolePermissions) String() string {
	return "46_add_role_permissions"
}
