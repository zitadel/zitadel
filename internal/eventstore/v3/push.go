package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (es *Eventstore) Push(ctx context.Context, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return nil, nil
}

var (
	//go:embed push.sql
	pushStmt string
	//go:embed push2.sql
	push2Stmt string
)

func (es *Eventstore) writeEvents(ctx context.Context, commands []eventstore.Command) (_ []eventstore.Event, err error) {
	tx, err := es.client.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
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
