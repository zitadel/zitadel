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

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (es *Eventstore) Push(ctx context.Context, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	ctx, spanBeginTx := tracing.NewNamedSpan(ctx, "db.BeginTx")
	tx, err := es.client.BeginTx(ctx, nil)
	spanBeginTx.EndWithError(err)
	if err != nil {
		return nil, err
	}
	// tx is not closed because [crdb.ExecuteInTx] takes care of that
	var (
		sequences []*latestSequence
	)

	err = crdb.ExecuteInTx(ctx, &transaction{tx}, func() error {
		sequences, err = latestSequences(ctx, tx, commands)
		if err != nil {
			return err
		}

		events, err = insertEvents(ctx, tx, sequences, commands)
		if err != nil {
			return err
		}

		if err = handleUniqueConstraints(ctx, tx, commands); err != nil {
			return err
		}

		// CockroachDB by default does not allow multiple modifications of the same table using ON CONFLICT
		// Thats why we enable it manually
		if es.client.Type() == "cockroach" {
			_, err = tx.Exec("SET enable_multiple_modifications_of_table = on")
			if err != nil {
				return err
			}
		}

		return handleFieldCommands(ctx, tx, commands)
	})

	if err != nil {
		return nil, err
	}

	return events, nil
}

//go:embed push.sql
var pushStmt string

func insertEvents(ctx context.Context, tx *sql.Tx, sequences []*latestSequence, commands []eventstore.Command) ([]eventstore.Event, error) {
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
	*sql.Tx
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
