package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/queue"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (es *Eventstore) Push(ctx context.Context, client database.ContextQueryExecuter, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	events, err = es.writeCommands(ctx, client, commands)
	if isSetupNotExecutedError(err) {
		return es.pushWithoutFunc(ctx, client, commands...)
	}

	return events, err
}

func (es *Eventstore) writeCommands(ctx context.Context, client database.ContextQueryExecuter, commands []eventstore.Command) (_ []eventstore.Event, err error) {
	var conn *sql.Conn
	switch c := client.(type) {
	case database.Client:
		conn, err = c.Conn(ctx)
	case nil:
		conn, err = es.client.Conn(ctx)
		client = conn
	}
	if err != nil {
		return nil, err
	}
	if conn != nil {
		defer conn.Close()
	}

	tx, close, err := es.pushTx(ctx, client)
	if err != nil {
		return nil, err
	}
	if close != nil {
		defer func() {
			err = close(err)
		}()
	}

	// lock the instance for reading events if await events is set for the duration of the transaction.
	_, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock_shared('eventstore.events2'::REGCLASS::OID::INTEGER, hashtext($1))", authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}

	events, err := writeEvents(ctx, tx, commands)
	if err != nil {
		return nil, err
	}
	if err = handleUniqueConstraints(ctx, tx, commands); err != nil {
		return nil, err
	}

	err = es.handleFieldCommands(ctx, tx, commands)
	if err != nil {
		return nil, err
	}

	err = es.queueExecutions(ctx, tx, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func writeEvents(ctx context.Context, tx database.Tx, commands []eventstore.Command) (_ []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	events, cmds, err := commandsToEvents(ctx, commands)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `select owner, created_at, "sequence", position from eventstore.push($1::eventstore.command[])`, cmds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&events[i].(*event).command.Owner, &events[i].(*event).createdAt, &events[i].(*event).sequence, &events[i].(*event).position)
		if err != nil {
			logging.WithError(err).Warn("failed to scan events")
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (es *Eventstore) queueExecutions(ctx context.Context, tx database.Tx, events []eventstore.Event) error {
	if es.queue == nil {
		return nil
	}

	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		types := make([]string, len(events))
		for i, event := range events {
			types[i] = string(event.Type())
		}
		logging.WithFields("event_types", types).Warningf("event executions skipped: wrong type of transaction %T", tx)
		return nil
	}
	jobArgs, err := eventsToJobArgs(ctx, events)
	if err != nil {
		return err
	}
	if len(jobArgs) == 0 {
		return nil
	}
	return es.queue.InsertManyFastTx(
		ctx, sqlTx, jobArgs,
		queue.WithQueueName(exec_repo.QueueName),
	)
}

func eventsToJobArgs(ctx context.Context, events []eventstore.Event) ([]river.JobArgs, error) {
	if len(events) == 0 {
		return nil, nil
	}
	router := authz.GetInstance(ctx).ExecutionRouter()
	if router.IsZero() {
		return nil, nil
	}

	jobArgs := make([]river.JobArgs, 0, len(events))
	for _, event := range events {
		targets, ok := router.GetEventBestMatch(fmt.Sprintf("event/%s", event.Type()))
		if !ok {
			continue
		}
		req, err := exec_repo.NewRequest(event, targets)
		if err != nil {
			return nil, err
		}
		jobArgs = append(jobArgs, req)
	}
	return jobArgs, nil
}
