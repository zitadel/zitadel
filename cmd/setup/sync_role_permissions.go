package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

// SyncRolePermissions is a repeatable step which synchronizes the InternalAuthZ
// RolePermissionMappings from the configuration to the database.
// This is needed until role permissions are manageable over the API.
type SyncRolePermissions struct {
	commands               *command.Commands
	eventstore             *eventstore.Eventstore
	rolePermissionMappings []authz.RoleMapping
}

func (mig *SyncRolePermissions) Execute(ctx context.Context, _ eventstore.Event) error {
	if err := mig.executeSystem(ctx); err != nil {
		return err
	}
	return mig.executeInstances(ctx)
}

func (mig *SyncRolePermissions) executeSystem(ctx context.Context) error {
	logging.WithFields("migration", mig.String()).Info("prepare system role permission sync events")
	details, err := mig.commands.SynchronizeRolePermission(ctx, "SYSTEM", mig.rolePermissionMappings)
	if err != nil {
		return err
	}
	logging.WithFields("migration", mig.String(), "sequence", details.Sequence).Info("pushed system role permission sync events")
	return nil
}

func (mig *SyncRolePermissions) executeInstances(ctx context.Context) error {
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
			Builder(),
	)
	if err != nil {
		return err
	}
	for i, instanceID := range instances {
		logging.WithFields("instance_id", instanceID, "migration", mig.String(), "progress", fmt.Sprintf("%d/%d", i+1, len(instances))).Info("prepare instance role permission sync events")
		details, err := mig.commands.SynchronizeRolePermission(ctx, instanceID, mig.rolePermissionMappings)
		if err != nil {
			return err
		}
		logging.WithFields("instance_id", instanceID, "migration", mig.String(), "sequence", details.Sequence).Info("pushed instance role permission sync events")
	}
	return nil
}

func (*SyncRolePermissions) String() string {
	return "repeatable_sync_role_permissions"
}

func (*SyncRolePermissions) Check(lastRun map[string]interface{}) bool {
	return true
}
