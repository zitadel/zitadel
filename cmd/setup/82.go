package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

//go:embed 82/*.sql
var projectGrantMemberPermissions embed.FS

// ProjectGrantMemberPermissions lets the permission check in SQL grant permissions thru project grant memberships.
// It fills the roles of the existing project grant members into the fields table first,
// then adds the project_grant_members view and returns the permitted project grants from eventstore.permitted_projects.
type ProjectGrantMemberPermissions struct {
	dbClient   *database.DB
	eventstore *eventstore.Eventstore
}

func (mig *ProjectGrantMemberPermissions) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes("instance").
			EventTypes(instance.InstanceAddedEventType).
			Builder().ExcludeAggregateIDs().
			AggregateTypes("instance").
			EventTypes(instance.InstanceRemovedEventType).
			Builder(),
	)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		ctx := authz.WithInstanceID(ctx, instance)
		if err := projection.GrantMembershipFields.Trigger(ctx); err != nil {
			return fmt.Errorf("%s: fill fields of instance %s: %w", mig.String(), instance, err)
		}
	}

	statements, err := readStatements(projectGrantMemberPermissions, "82")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.Info(ctx, "execute statement", "file", stmt.file, "migration", mig.String())
		if _, err := mig.dbClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (*ProjectGrantMemberPermissions) String() string {
	return "82_project_grant_member_permissions"
}
