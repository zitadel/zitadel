package eventstore

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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

// checks whether the error is caused because setup step 39 was not executed
func isSetupNotExecutedError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return (pgErr.Code == "42704" && strings.Contains(pgErr.Message, "eventstore.command")) ||
			(pgErr.Code == "42883" && strings.Contains(pgErr.Message, "eventstore.push"))
	}
	return errors.Is(err, errTypesNotFound)
}

var (
	//go:embed push.sql
	pushStmt string
)

// pushWithoutFunc implements pushing events before setup step 39 was introduced.
// TODO: remove with v3
func (es *Eventstore) pushWithoutFunc(ctx context.Context, client database.ContextQueryExecuter, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	tx, closeTx, err := es.pushTx(ctx, client)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = closeTx(err)
	}()

	// tx is not closed because [crdb.ExecuteInTx] takes care of that
	var (
		sequences []*latestSequence
	)
	sequences, err = latestSequences(ctx, tx, commands)
	if err != nil {
		return nil, err
	}

	events, err = es.writeEventsOld(ctx, tx, sequences, commands)
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

func (es *Eventstore) writeEventsOld(ctx context.Context, tx database.Tx, sequences []*latestSequence, commands []eventstore.Command) ([]eventstore.Event, error) {
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

		events[i], err = commandToEventOld(sequence, command)
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

		args = append(args,
			events[i].(*event).command.InstanceID,
			events[i].(*event).command.Owner,
			events[i].(*event).command.AggregateType,
			events[i].(*event).command.AggregateID,
			events[i].(*event).command.Revision,
			events[i].(*event).command.Creator,
			events[i].(*event).command.CommandType,
			events[i].(*event).command.Payload,
			events[i].(*event).sequence,
			i,
		)
	}

	return events, placeholders, args, nil
}
