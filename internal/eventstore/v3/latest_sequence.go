package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
)

type latestSequence struct {
	aggregate *Aggregate
	sequence  uint64
}

//go:embed sequences_query.sql
var latestSequencesStmt string

func latestSequences(ctx context.Context, tx *sql.Tx, commands []Command) ([]*latestSequence, error) {
	sequences := make([]*latestSequence, 0, len(commands))

	for _, command := range commands {
		if searchSequenceByCommand(sequences, command) != nil {
			continue
		}

		if command.Aggregate().InstanceID == "" {
			command.Aggregate().InstanceID = authz.GetInstance(ctx).InstanceID()
		}
		sequences = append(sequences, &latestSequence{
			aggregate: command.Aggregate(),
		})
	}

	const argsPerCondition = 3
	args := make([]interface{}, 0, len(sequences)*argsPerCondition)
	conditions := make([]string, len(sequences))

	for i, sequence := range sequences {
		conditions[i] = fmt.Sprintf("(instance_id = $%d AND aggregate_type = $%d AND aggregate_id = $%d)",
			i*argsPerCondition+1,
			i*argsPerCondition+2,
			i*argsPerCondition+3,
		)
		args = append(args, sequence.aggregate.InstanceID, sequence.aggregate.Type, sequence.aggregate.ID)
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprintf(latestSequencesStmt, strings.Join(conditions, " OR ")), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var aggregateType AggregateType
		var aggregateID, instanceID string
		var currentSequence uint64

		err := rows.Scan(&instanceID, &aggregateType, &aggregateID, &currentSequence)
		if err != nil {
			return nil, err
		}

		sequence := searchSequence(sequences, aggregateType, aggregateID, instanceID)
		sequence.sequence = currentSequence
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return sequences, nil
}

func searchSequenceByCommand(sequences []*latestSequence, command Command) *latestSequence {
	for _, sequence := range sequences {
		if sequence.aggregate.Type == command.Aggregate().Type &&
			sequence.aggregate.ID == command.Aggregate().ID &&
			sequence.aggregate.InstanceID == command.Aggregate().InstanceID {
			return sequence
		}
	}
	return nil
}

func searchSequence(sequences []*latestSequence, aggregateType AggregateType, aggregateID, instanceID string) *latestSequence {
	for _, sequence := range sequences {
		if sequence.aggregate.Type == aggregateType &&
			sequence.aggregate.ID == aggregateID &&
			sequence.aggregate.InstanceID == instanceID {
			return sequence
		}
	}
	return nil
}
