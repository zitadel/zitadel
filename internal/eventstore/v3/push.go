package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var appNamePrefix = dialect.DBPurposeEventPusher.AppName() + "_"

var pushTxOpts = &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

func (es *Eventstore) Push(ctx context.Context, client database.QueryExecuter, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var tx database.Tx
	switch c := client.(type) {
	case database.Tx:
		tx = c
	case database.Client:
		// We cannot use READ COMMITTED on CockroachDB because we use cluster_logical_timestamp() which is not supported in this isolation level
		var opts *sql.TxOptions
		if es.client.Database.Type() == "postgres" {
			opts = pushTxOpts
		}
		tx, err = c.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = database.CloseTransaction(tx, err)
		}()
	default:
		// We cannot use READ COMMITTED on CockroachDB because we use cluster_logical_timestamp() which is not supported in this isolation level
		var opts *sql.TxOptions
		if es.client.Database.Type() == "postgres" {
			opts = pushTxOpts
		}
		tx, err = es.client.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = database.CloseTransaction(tx, err)
		}()
	}
	// tx is not closed because [crdb.ExecuteInTx] takes care of that
	var (
		sequences []*latestSequence
	)

	// needs to be set like this because psql complains about parameters in the SET statement
	_, err = tx.ExecContext(ctx, "SET application_name = '"+appNamePrefix+authz.GetInstance(ctx).InstanceID()+"'")
	if err != nil {
		logging.WithError(err).Warn("failed to set application name")
		return nil, err
	}

	sequences, err = latestSequences(ctx, tx, commands)
	if err != nil {
		return nil, err
	}

	events, err = insertEvents(ctx, tx, sequences, commands)
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

//go:embed push.sql
var pushStmt string

func insertEvents(ctx context.Context, tx database.Tx, sequences []*latestSequence, commands []eventstore.Command) ([]eventstore.Event, error) {
	events, placeholders, args, err := mapCommands(commands, sequences)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprintf(pushStmt, strings.Join(placeholders, ", ")), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&events[i].(*event).createdAt, &events[i].(*event).position)
		if err != nil {
			logging.WithError(err).Warn("failed to scan events")
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		pgErr := new(pgconn.PgError)
		if errors.As(err, &pgErr) {
			// Check if push tries to write an event just written
			// by another transaction
			if pgErr.Code == "40001" {
				// TODO: @livio-a should we return the parent or not?
				return nil, zerrors.ThrowInvalidArgument(err, "V3-p5xAn", "Errors.AlreadyExists")
			}
		}
		logging.WithError(rows.Err()).Warn("failed to push events")
		return nil, zerrors.ThrowInternal(err, "V3-VGnZY", "Errors.Internal")
	}

	return events, nil
}

const argsPerCommand = 10

func mapCommands(commands []eventstore.Command, sequences []*latestSequence) (events []eventstore.Event, placeholders []string, args []any, err error) {
	events = make([]eventstore.Event, len(commands))
	args = make([]any, 0, len(commands)*argsPerCommand)
	placeholders = make([]string, len(commands))

	for i, command := range commands {
		sequence := searchSequenceByCommand(sequences, command)
		if sequence == nil {
			logging.WithFields(
				"aggType", command.Aggregate().Type,
				"aggID", command.Aggregate().ID,
				"instance", command.Aggregate().InstanceID,
			).Panic("no sequence found")
			// added return for linting
			return nil, nil, nil, nil
		}
		sequence.sequence++

		events[i], err = commandToEvent(sequence, command)
		if err != nil {
			return nil, nil, nil, err
		}

		placeholders[i] = fmt.Sprintf(pushPlaceholderFmt,
			i*argsPerCommand+1,
			i*argsPerCommand+2,
			i*argsPerCommand+3,
			i*argsPerCommand+4,
			i*argsPerCommand+5,
			i*argsPerCommand+6,
			i*argsPerCommand+7,
			i*argsPerCommand+8,
			i*argsPerCommand+9,
			i*argsPerCommand+10,
		)

		revision, err := strconv.Atoi(strings.TrimPrefix(string(events[i].(*event).aggregate.Version), "v"))
		if err != nil {
			return nil, nil, nil, zerrors.ThrowInternal(err, "V3-JoZEp", "Errors.Internal")
		}
		args = append(args,
			events[i].(*event).aggregate.InstanceID,
			events[i].(*event).aggregate.ResourceOwner,
			events[i].(*event).aggregate.Type,
			events[i].(*event).aggregate.ID,
			revision,
			events[i].(*event).creator,
			events[i].(*event).typ,
			events[i].(*event).payload,
			events[i].(*event).sequence,
			i,
		)
	}

	return events, placeholders, args, nil
}

type transaction struct {
	database.Tx
}

var _ crdb.Tx = (*transaction)(nil)

func (t *transaction) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.Tx.ExecContext(ctx, query, args...)
	return err
}

func (t *transaction) Commit(ctx context.Context) error {
	return t.Tx.Commit()
}

func (t *transaction) Rollback(ctx context.Context) error {
	return t.Tx.Rollback()
}
