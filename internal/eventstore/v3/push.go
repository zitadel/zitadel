package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var appNamePrefix = dialect.DBPurposeEventPusher.AppName() + "_"

func (es *Eventstore) Push(ctx context.Context, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	events, err = es.writeCommands(ctx, commands)
	if isSetupNotExecutedError(err) {
		return es.pushWithoutFunc(ctx, commands...)
	}

	return events, err
}

func (es *Eventstore) writeCommands(ctx context.Context, commands []eventstore.Command) (_ []eventstore.Event, err error) {
	conn, err := es.client.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err = checkExecutionPlan(ctx, conn); err != nil {
		return nil, err
	}

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	// needs to be set like this because psql complains about parameters in the SET statement
	_, err = tx.ExecContext(ctx, "SET application_name = '"+appNamePrefix+authz.GetInstance(ctx).InstanceID()+"'")
	if err != nil {
		logging.WithError(err).Warn("failed to set application name")
		return nil, err
	}

	// needs to be set like this because psql complains about parameters in the SET statement
	_, err = tx.ExecContext(ctx, "SET application_name = '"+appNamePrefix+authz.GetInstance(ctx).InstanceID()+"'")
	if err != nil {
		logging.WithError(err).Warn("failed to set application name")
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			logging.WithError(err).Warn("failed to commit transaction")
		}
	}()

	events, err := writeEvents(ctx, tx, commands)
	if err != nil {
		return nil, err
	}

	if err = handleUniqueConstraints(ctx, tx, commands); err != nil {
		return nil, err
	}

	// CockroachDB by default does not allow multiple modifications of the same table using ON CONFLICT
	// Thats why we enable it manually
	if es.client.Type() == "cockroach" {
		_, err = tx.Exec("SET enable_multiple_modifications_of_table = on")
		if err != nil {
			return nil, err
		}
	}

	err = handleFieldCommands(ctx, tx, commands)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func writeEvents(ctx context.Context, tx *sql.Tx, commands []eventstore.Command) (_ []eventstore.Event, err error) {
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

var errTypesNotFound = errors.New("types not found")

func checkExecutionPlan(ctx context.Context, conn *sql.Conn) error {
	return conn.Raw(func(driverConn any) error {
		conn, ok := driverConn.(*stdlib.Conn)
		if !ok {
			return errTypesNotFound
		}

		var cmd *command
		if _, ok := conn.Conn().TypeMap().TypeForValue(cmd); ok {
			return nil
		}
		return registerEventstoreTypes(ctx, conn.Conn())
	})
}
