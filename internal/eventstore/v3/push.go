package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"slices"
	"strconv"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var pushTxOpts = &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

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
	if len(commands) == 0 {
		return nil, nil
	}

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

	instanceIDs := make([]any, 0, len(commands))
	locks := make([]string, 0, len(commands))
	for _, cmd := range commands {
		if !slices.Contains(instanceIDs, any(cmd.Aggregate().InstanceID)) {
			instanceIDs = append(instanceIDs, cmd.Aggregate().InstanceID)
			locks = append(locks, "pg_advisory_xact_lock_shared('eventstore.events2'::REGCLASS::OID::INTEGER, hashtext($"+strconv.Itoa(len(instanceIDs))+"))")
		}
	}
	lockStmt := "SELECT " + strings.Join(locks, ", ")

	tx, close, err := es.pushTx(ctx, client)
	if err != nil {
		return nil, err
	}
	if close != nil {
		defer func() {
			err = close(err)
		}()
	}

	_, err = tx.ExecContext(ctx, lockStmt, instanceIDs...)
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
