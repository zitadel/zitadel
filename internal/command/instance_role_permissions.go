package command

import (
	"context"
	"database/sql"
	_ "embed"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/permission"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SynchronizeRolePermission checks the current state of role permissions in the eventstore for the aggregate.
// It pushes the commands required to reach the desired state passed in target.
// For system level permissions aggregateID must be set to `SYSTEM`, else it is the instance ID.
func (c *Commands) SynchronizeRolePermission(ctx context.Context, aggregateID string, target []authz.RoleMapping) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	cmds, err := synchronizeRolePermissionCommands(ctx, c.eventstore.Client(), aggregateID,
		rolePermissionMappingsToDatabaseMap(target, aggregateID == "SYSTEM"),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "COMMA-Iej2r", "Errors.Internal")
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		logging.WithError(err).Error("failed to push role permission commands")
		return nil, zerrors.ThrowInternal(err, "COMMA-AiV3u", "Errors.Internal")
	}
	return pushedEventsToObjectDetails(events), nil
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

var (
	//go:embed instance_role_permissions_sync.sql
	instanceRolePermissionsSyncQuery string
)

// synchronizeRolePermissionCommands checks the current state of role permissions in the eventstore for the aggregate.
// It returns the commands required to reach the desired state passed in target.
// For system level permissions aggregateID must be set to `SYSTEM`, else it is the instance ID.
func synchronizeRolePermissionCommands(ctx context.Context, db *database.DB, aggregateID string, target database.Map[[]string]) (cmds []eventstore.Command, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = db.QueryContext(ctx,
		rolePermissionScanner(ctx, permission.NewAggregate(aggregateID), &cmds),
		instanceRolePermissionsSyncQuery,
		aggregateID, target)
	if err != nil {
		return nil, err
	}
	return cmds, nil
}

func rolePermissionScanner(ctx context.Context, aggregate *eventstore.Aggregate, cmds *[]eventstore.Command) func(rows *sql.Rows) error {
	return func(rows *sql.Rows) error {
		for rows.Next() {
			var operation, role, perm string
			if err := rows.Scan(&operation, &role, &perm); err != nil {
				return err
			}
			logging.WithFields("aggregate_id", aggregate.ID, "operation", operation, "role", role, "permission", perm).Debug("sync role permission")
			switch operation {
			case "add":
				*cmds = append(*cmds, permission.NewAddedEvent(ctx, aggregate, role, perm))
			case "remove":
				*cmds = append(*cmds, permission.NewRemovedEvent(ctx, aggregate, role, perm))
			}
		}
		return rows.Close()

	}
}
