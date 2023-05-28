package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
)

func (es *Eventstore) Push(ctx context.Context, commands ...Command) (_ []Event, err error) {
	tx, err := es.client.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
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

	err = handleUniqueConstraints(ctx, tx, commands)
	if err != nil {
		return nil, err
	}

	return events, nil
}

//go:embed push.sql
var pushStmt string

func insertEvents(ctx context.Context, tx *sql.Tx, sequences []*latestSequence, commands []Command) (events []Event, err error) {
	const argsPerCommand = 9

	events = make([]Event, len(commands))
	args := make([]any, 0, len(commands)*argsPerCommand)
	placeHolders := make([]string, len(commands))

	for i, command := range commands {
		sequence := searchSequenceByCommand(sequences, command)
		if sequence == nil {
			panic("asdf")
		}
		sequence.sequence++

		events[i], err = commandToEvent(sequence, command)
		if err != nil {
			return nil, err
		}

		placeHolders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*argsPerCommand+1,
			i*argsPerCommand+2,
			i*argsPerCommand+3,
			i*argsPerCommand+4,
			i*argsPerCommand+5,
			i*argsPerCommand+6,
			i*argsPerCommand+7,
			i*argsPerCommand+8,
			i*argsPerCommand+9,
		)
		args = append(args,
			events[i].(*event).aggregate.InstanceID,
			events[i].(*event).aggregate.ResourceOwner,
			events[i].(*event).aggregate.Type,
			events[i].(*event).aggregate.ID,
			events[i].(*event).aggregate.Version,
			events[i].(*event).creator,
			events[i].(*event).typ,
			sql.NullString{String: string(events[i].(*event).payload), Valid: len(events[i].(*event).payload) > 0},
			events[i].(*event).sequence,
		)
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprintf(pushStmt, strings.Join(placeHolders, ", ")), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&events[i].(*event).createdAt)
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}
