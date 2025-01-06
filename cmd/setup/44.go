package setup

import (
	"context"
	"embed"
	_ "embed"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var (
	//go:embed 44/01_table_definition.sql
	createTransactionalInstance string

	//go:embed 44/cockroach/*.sql
	//go:embed 44/postgres/*.sql
	createReduceInstanceTrigger embed.FS
)

type CreateTransactionalInstance struct {
	dbClient   *database.DB
	eventstore *eventstore.Eventstore
	BulkLimit  uint64
}

func (mig *CreateTransactionalInstance) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	_, err = mig.dbClient.ExecContext(ctx, createTransactionalInstance)
	if err != nil {
		return err
	}
	statements, err := readStatements(createReduceInstanceTrigger, "44", mig.dbClient.Type())
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		_, err = mig.dbClient.ExecContext(ctx, stmt.query)
		if err != nil {
			return err
		}
	}

	reducer := new(instanceEvents)
	for {
		err = mig.eventstore.FilterToReducer(ctx,
			eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				AwaitOpenTransactions().
				Limit(mig.BulkLimit).
				Offset(reducer.offset).
				OrderAsc().
				AddQuery().
				AggregateTypes(instance.AggregateType).
				EventTypes(
					instance.InstanceAddedEventType,
					instance.InstanceChangedEventType,
					instance.InstanceRemovedEventType,
					instance.DefaultLanguageSetEventType,
					instance.ProjectSetEventType,
					instance.ConsoleSetEventType,
					instance.DefaultOrgSetEventType,
				).
				Builder(),
			reducer,
		)
		if err != nil || len(reducer.events) == 0 {
			return err
		}

		tx, err := mig.dbClient.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, event := range reducer.events {
			switch e := event.(type) {
			case *instance.InstanceAddedEvent:
				_, err = tx.ExecContext(ctx,
					"SELECT reduce_instance_added($1, $2, $3, $4)",
					e.Aggregate().ID,
					e.Name,
					e.CreatedAt(),
					e.Position(),
				)
			case *instance.InstanceChangedEvent:
				_, err = tx.ExecContext(ctx,
					"SELECT reduce_instance_updated($1, $2, $3, $4)",
					e.Aggregate().ID,
					e.Name,
					e.CreatedAt(),
					e.Position(),
				)
			case *instance.InstanceRemovedEvent:
				_, err = tx.ExecContext(ctx,
					"SELECT reduce_instance_removed($1)",
					e.Aggregate().ID,
				)
			case *instance.DefaultLanguageSetEvent:
				_, err = tx.ExecContext(ctx,
					"SELECT reduce_instance_removed($1, $2, $3, $4)",
					e.Aggregate().ID,
					e.Language,
					e.CreatedAt(),
					e.Position(),
				)
			case *instance.ProjectSetEvent:
				_, err = tx.ExecContext(ctx,
					"SELECT reduce_instance_project_set($1, $2, $3, $4)",
					e.Aggregate().ID,
					e.ProjectID,
					e.CreatedAt(),
					e.Position(),
				)
			case *instance.ConsoleSetEvent:
				_, err = tx.ExecContext(ctx,
					"SELECT reduce_instance_console_set($1, $2, $3, $4, $5)",
					e.Aggregate().ID,
					e.AppID,
					e.ClientID,
					e.CreatedAt(),
					e.Position(),
				)
			case *instance.DefaultOrgSetEvent:
				_, err = tx.ExecContext(ctx,
					"SELECT reduce_instance_default_org_set($1, $2, $3, $4)",
					e.Aggregate().ID,
					e.OrgID,
					e.CreatedAt(),
					e.Position(),
				)
			}
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if err = tx.Commit(); err != nil {
				return err
			}
		}

		reducer.events = nil
		reducer.offset += uint32(len(reducer.events))
	}
}

func (mig *CreateTransactionalInstance) String() string {
	return "44_create_transactional_instance"
}

type instanceEvents struct {
	offset uint32
	events []eventstore.Event
}

// AppendEvents implements eventstore.reducer.
func (i *instanceEvents) AppendEvents(events ...eventstore.Event) {
	i.events = append(i.events, events...)
}

// Reduce implements eventstore.reducer.
func (i *instanceEvents) Reduce() error {
	return nil
}
