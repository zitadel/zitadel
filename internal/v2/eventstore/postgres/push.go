package postgres

import (
	"context"
	"database/sql"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Push implements eventstore.Pusher.
func (s *Storage) Push(ctx context.Context, pushIntents ...eventstore.PushIntent) (err error) {
	tx, err := s.client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return err
	}
	defer func() {
		err = closeTx(tx, err)
	}()

	intents, err := lockAggregates(ctx, tx, pushIntents)
	if err != nil {
		return err
	}

	if !checkSequences(intents) {
		return zerrors.ThrowInvalidArgument(nil, "POSTG-KOM6E", "Errors.Internal.Eventstore.SequenceNotMatched")
	}

	var commands []*command
	for _, intent := range intents {
		additionalCommands, err := intentToCommands(intent)
		if err != nil {
			return err
		}
		commands = append(commands, additionalCommands...)
	}

	err = uniqueConstraints(ctx, tx, commands)
	if err != nil {
		return err
	}

	return push(ctx, tx, commands)
}

func lockAggregates(ctx context.Context, tx *sql.Tx, pushIntents []eventstore.PushIntent) ([]*intent, error) {
	var stmt statement

	stmt.builder.WriteString("WITH existing AS (")
	for i, intent := range pushIntents {
		stmt.builder.WriteString(`(SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = `)
		stmt.writeArgs(intent.Aggregate().Instance)
		stmt.builder.WriteString(` AND aggregate_type = `)
		stmt.writeArgs(intent.Aggregate().Type)
		stmt.builder.WriteString(` AND aggregate_id = `)
		stmt.writeArgs(intent.Aggregate().ID)
		stmt.builder.WriteString(` ORDER BY "sequence" DESC LIMIT 1)`)

		if i < len(pushIntents)-1 {
			stmt.builder.WriteString(" UNION ALL ")
		}
	}
	stmt.builder.WriteString(") SELECT e.instance_id, e.owner, e.aggregate_type, e.aggregate_id, e.sequence FROM eventstore.events2 e JOIN existing ON e.instance_id = existing.instance_id AND e.aggregate_type = existing.aggregate_type AND e.aggregate_id = existing.aggregate_id AND e.sequence = existing.sequence FOR UPDATE")

	rows, err := tx.QueryContext(ctx, stmt.builder.String(), stmt.args...)
	if err != nil {
		return nil, err
	}

	res := makeIntents(pushIntents)

	err = mapRowsToObject(rows, func(scan func(dest ...any) error) error {
		var sequence sql.Null[uint32]
		agg := new(eventstore.Aggregate)

		err := scan(
			&agg.Instance,
			&agg.Owner,
			&agg.Type,
			&agg.ID,
			&sequence,
		)
		if err != nil {
			return err
		}

		intentByAggregate(res, agg).sequence = sequence.V

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func push(ctx context.Context, tx *sql.Tx, commands []*command) error {
	var stmt statement

	stmt.builder.WriteString(`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES `)
	for i, cmd := range commands {
		stmt.builder.WriteString(`(`)
		stmt.writeArgs(cmd.aggregate.Instance, cmd.aggregate.Owner, cmd.aggregate.Type, cmd.aggregate.ID, cmd.revision, cmd.creator, cmd.typ, cmd.payload, cmd.sequence, i)
		stmt.builder.WriteString(", statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())")
		stmt.builder.WriteString(`)`)
		if i < len(commands)-1 {
			stmt.builder.WriteString(", ")
		}
	}
	stmt.builder.WriteString(` RETURNING created_at, "position"`)

	rows, err := tx.QueryContext(ctx, stmt.builder.String(), stmt.args...)
	if err != nil {
		return err
	}

	var i int
	return mapRowsToObject(rows, func(scan func(dest ...any) error) error {
		err := scan(
			&commands[i].createdAt,
			&commands[i].position,
		)
		i++
		return err
	})
}

func uniqueConstraints(ctx context.Context, tx *sql.Tx, commands []*command) error {
	var stmt statement

	for _, cmd := range commands {
		if len(cmd.uniqueConstraints) == 0 {
			continue
		}
		for _, constraint := range cmd.uniqueConstraints {
			stmt.reset()
			switch constraint.Action {
			case eventstore.UniqueConstraintAdd:
				stmt.builder.WriteString(`INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES (`)
				stmt.writeArgs(cmd.aggregate.Instance, constraint.UniqueType, constraint.UniqueField)
				stmt.builder.WriteString(`)`)
			case eventstore.UniqueConstraintInstanceRemove:
				stmt.builder.WriteString(`DELETE FROM eventstore.unique_constraints WHERE `)
				stmt.builder.WriteString(deleteUniqueConstraintClause)
				stmt.appendArgs(
					sql.Named("instanceId", cmd.aggregate.Instance),
					sql.Named("uniqueType", constraint.UniqueType),
					sql.Named("uniqueField", constraint.UniqueField),
				)
			case eventstore.UniqueConstraintRemove:
				stmt.builder.WriteString(`DELETE FROM eventstore.unique_constraints WHERE instance_id = `)
				stmt.writeArgs(cmd.aggregate.Instance)
			}
			_, err := tx.ExecContext(ctx, stmt.builder.String(), stmt.args...)
			if err != nil {
				logging.WithFields("action", constraint.Action).Warn("handling of unique constraint failed")
				errMessage := constraint.ErrorMessage
				if errMessage == "" {
					errMessage = "Errors.Internal"
				}
				return zerrors.ThrowAlreadyExists(err, "POSTG-QzjyP", errMessage)
			}
		}
	}

	return nil
}

// the query is so complex because we accidentally stored unique constraint case sensitive
// the query checks first if there is a case sensitive match and afterwards if there is a case insensitive match
var deleteUniqueConstraintClause = `
(instance_id = @instanceId AND unique_type = @uniqueType AND unique_field = (
    SELECT unique_field from (
        SELECT instance_id, unique_type, unique_field
        FROM eventstore.unique_constraints
        WHERE instance_id = @instanceId AND unique_type = @uniqueType AND unique_field = @uniqueField
    UNION ALL
        SELECT instance_id, unique_type, unique_field
        FROM eventstore.unique_constraints
        WHERE instance_id = @instanceId AND unique_type = @uniqueType AND unique_field = LOWER(@uniqueField)
    ) AS case_insensitive_constraints LIMIT 1)
)`
