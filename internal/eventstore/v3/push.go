package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func (es *Eventstore) Push(ctx context.Context, commands ...eventstore.Command) (_ []eventstore.Event, err error) {
	tx, err := es.client.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			logging.OnError(txErr).Debug("unable to rollback transaction")
			return
		}
		err = tx.Commit()
	}()
	sequences, err := latestSequences(ctx, tx, commands)
	if err != nil {
		return nil, err
	}

	events, err := insertEvents(ctx, tx, sequences, commands)
	if err != nil {
		return nil, err
	}

	if err = handleUniqueConstraints(ctx, tx, commands); err != nil {
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
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		args = append(args,
			events[i].(*event).aggregate.InstanceID,
			events[i].(*event).aggregate.ResourceOwner,
			events[i].(*event).aggregate.Type,
			events[i].(*event).aggregate.ID,
			events[i].(*event).aggregate.Version,
			events[i].(*event).creator,
			events[i].(*event).typ,
			events[i].(*event).payload,
			events[i].(*event).sequence,
			i,
		)
	}

	return events, placeholders, args, nil
}
