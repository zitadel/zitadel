package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (es *Eventstore) Push(ctx context.Context, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	tx, err := es.client.BeginTx(ctx, nil)
	if err != nil {
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

	// err = crdb.ExecuteInTx(ctx, &transaction{tx}, func() (err error) {
	events, err = insertEvents2(ctx, tx, commands)
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
	// })

	if err != nil {
		return nil, err
	}

	return events, nil
}

var (
	//go:embed push.sql
	pushStmt string
	//go:embed push2.sql
	push2Stmt string
)

func insertEvents2(ctx context.Context, tx *sql.Tx, commands []eventstore.Command) (_ []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	events, cmds, err := commandsToEvents2(ctx, commands)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, push2Stmt, cmds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&events[i].(*event).createdAt, &events[i].(*event).sequence, &events[i].(*event).position)
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

func commandsToEvents2(ctx context.Context, cmds []eventstore.Command) (_ []eventstore.Event, _ []*command, err error) {
	events := make([]eventstore.Event, len(cmds))
	commands := make([]*command, len(cmds))
	for i, cmd := range cmds {
		if cmd.Aggregate().InstanceID == "" {
			cmd.Aggregate().InstanceID = authz.GetInstance(ctx).InstanceID()
		}
		events[i], err = commandToEvent2(ctx, cmd)
		if err != nil {
			return nil, nil, err
		}
		commands[i] = events[i].(*event).command
	}
	return events, commands, nil
}

func commandToEvent2(ctx context.Context, cmd eventstore.Command) (_ eventstore.Event, err error) {
	var payload Payload
	if cmd.Payload() != nil {
		payload, err = json.Marshal(cmd.Payload())
		if err != nil {
			logging.WithError(err).Warn("marshal payload failed")
			return nil, zerrors.ThrowInternal(err, "V3-MInPK", "Errors.Internal")
		}
	}

	command := &command{
		InstanceID:    cmd.Aggregate().InstanceID,
		AggregateType: string(cmd.Aggregate().Type),
		AggregateID:   cmd.Aggregate().ID,
		CommandType:   string(cmd.Type()),
		Revision:      cmd.Revision(),
		Payload:       payload,
		Creator:       cmd.Creator(),
		Owner:         cmd.Aggregate().ResourceOwner,
	}

	return &event{
		aggregate: cmd.Aggregate(),
		command:   command,
	}, nil
}

// func insertEvents(ctx context.Context, tx *sql.Tx, sequences []*latestSequence, commands []eventstore.Command) ([]eventstore.Event, error) {
// 	events, placeholders, args, err := mapCommands(commands, sequences)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows, err := tx.QueryContext(ctx, fmt.Sprintf(pushStmt, strings.Join(placeholders, ", ")), args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for i := 0; rows.Next(); i++ {
// 		err = rows.Scan(&events[i].(*event).Created, &events[i].(*event).Pos)
// 		if err != nil {
// 			logging.WithError(err).Warn("failed to scan events")
// 			return nil, err
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		pgErr := new(pgconn.PgError)
// 		if errors.As(err, &pgErr) {
// 			// Check if push tries to write an event just written
// 			// by another transaction
// 			if pgErr.Code == "40001" {
// 				// TODO: @livio-a should we return the parent or not?
// 				return nil, zerrors.ThrowInvalidArgument(err, "V3-p5xAn", "Errors.AlreadyExists")
// 			}
// 		}
// 		logging.WithError(rows.Err()).Warn("failed to push events")
// 		return nil, zerrors.ThrowInternal(err, "V3-VGnZY", "Errors.Internal")
// 	}

// 	return events, nil
// }

const argsPerCommand = 10

// func mapCommands(commands []eventstore.Command, sequences []*latestSequence) (events []eventstore.Event, placeholders []string, args []any, err error) {
// 	events = make([]eventstore.Event, len(commands))
// 	args = make([]any, 0, len(commands)*argsPerCommand)
// 	placeholders = make([]string, len(commands))

// 	for i, command := range commands {
// 		sequence := searchSequenceByCommand(sequences, command)
// 		if sequence == nil {
// 			logging.WithFields(
// 				"aggType", command.Aggregate().Type,
// 				"aggID", command.Aggregate().ID,
// 				"instance", command.Aggregate().InstanceID,
// 			).Panic("no sequence found")
// 			// added return for linting
// 			return nil, nil, nil, nil
// 		}
// 		sequence.sequence++

// 		events[i], err = commandToEvent(sequence, command)
// 		if err != nil {
// 			return nil, nil, nil, err
// 		}

// 		placeholders[i] = fmt.Sprintf(pushPlaceholderFmt,
// 			i*argsPerCommand+1,
// 			i*argsPerCommand+2,
// 			i*argsPerCommand+3,
// 			i*argsPerCommand+4,
// 			i*argsPerCommand+5,
// 			i*argsPerCommand+6,
// 			i*argsPerCommand+7,
// 			i*argsPerCommand+8,
// 			i*argsPerCommand+9,
// 			i*argsPerCommand+10,
// 		)

// 		revision, err := strconv.Atoi(strings.TrimPrefix(string(events[i].(*event).Rev), "v"))
// 		if err != nil {
// 			return nil, nil, nil, zerrors.ThrowInternal(err, "V3-JoZEp", "Errors.Internal")
// 		}
// 		args = append(args,
// 			events[i].(*event).InstanceID,
// 			events[i].(*event).ResourceOwner,
// 			events[i].(*event).AggregateType,
// 			events[i].(*event).AggregateID,
// 			revision,
// 			events[i].(*event).CreatorUser,
// 			events[i].(*event).Typ,
// 			events[i].(*event).Data,
// 			events[i].(*event).Seq,
// 			i,
// 		)
// 	}

// 	return events, placeholders, args, nil
// }

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
