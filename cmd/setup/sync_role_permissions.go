package setup

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/permission"
)

var (
	//go:embed sync_role_permissions.sql
	getRolePermissionOperationsQuery string
)

// SyncRolePermissions is a repeatable step which synchronizes the InternalAuthZ
// RolePermissionMappings from the configuration to the database.
// This is needed until role permissions are manageable over the API.
type SyncRolePermissions struct {
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

	target := rolePermissionMappingsToDatabaseMap(mig.rolePermissionMappings, true)
	cmds, err := mig.synchronizeCommands(ctx, "SYSTEM", target)
	if err != nil {
		return err
	}
	events, err := mig.eventstore.Push(ctx, cmds...)
	if err != nil {
		return err
	}

	logging.WithFields("migration", mig.String(), "pushed_events", len(events)).Info("pushed system role permission sync events")
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
	target := rolePermissionMappingsToDatabaseMap(mig.rolePermissionMappings, false)
	for i, instanceID := range instances {
		logging.WithFields("instance_id", instanceID, "migration", mig.String(), "progress", fmt.Sprintf("%d/%d", i+1, len(instances))).Info("prepare instance role permission sync events")
		cmds, err := mig.synchronizeCommands(ctx, instanceID, target)
		if err != nil {
			return err
		}
		events, err := mig.eventstore.Push(ctx, cmds...)
		if err != nil {
			return err
		}
		logging.WithFields("instance_id", instanceID, "migration", mig.String(), "pushed_events", len(events)).Info("pushed instance role permission sync events")
	}
	return nil
}

// synchronizeCommands checks the current state of role permissions in the eventstore for the aggregate.
// It returns the commands required to reach the desired state passed in target.
// For system level permissions aggregateID must be set to `SYSTEM`,
// else it is the instance ID.
func (mig *SyncRolePermissions) synchronizeCommands(ctx context.Context, aggregateID string, target database.Map[[]string]) (cmds []eventstore.Command, err error) {
	aggregate := permission.NewAggregate(aggregateID)
	err = mig.eventstore.Client().QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var operation, role, perm string
			if err := rows.Scan(&operation, &role, &perm); err != nil {
				return err
			}
			logging.WithFields("aggregate_id", aggregateID, "migration", mig.String(), "operation", operation, "role", role, "permission", perm).Debug("sync role permission")
			switch operation {
			case "add":
				cmds = append(cmds, permission.NewAddedEvent(ctx, aggregate, role, perm))
			case "remove":
				cmds = append(cmds, permission.NewRemovedEvent(ctx, aggregate, role, perm))
			}
		}
		return rows.Close()

	}, getRolePermissionOperationsQuery, aggregateID, target)
	if err != nil {
		return nil, err
	}
	return cmds, err
}

func (*SyncRolePermissions) String() string {
	return "repeatable_sync_role_permissions"
}

func (*SyncRolePermissions) Check(lastRun map[string]interface{}) bool {
	return true
}

func rolePermissionMappingsToDatabaseMap(mappings []authz.RoleMapping, system bool) database.Map[[]string] {
	out := make(database.Map[[]string], len(mappings))
	for _, m := range mappings {
		if system == strings.HasPrefix(m.Role, "SYSTEM") {
			out[m.Role] = m.Permissions
		}
	}
	return out
}
