package eventstore

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type latestSequence struct {
	aggregate *eventstore.Aggregate
	sequence  uint64
}

//go:embed sequences_query.sql
var latestSequencesStmt string

func latestSequences(ctx context.Context, tx database.Tx, commands []eventstore.Command) ([]*latestSequence, error) {
	sequences := commandsToSequences(ctx, commands)

	conditions, args := sequencesToSql(sequences)
	rows, err := tx.QueryContext(ctx, fmt.Sprintf(latestSequencesStmt, strings.Join(conditions, " UNION ALL ")), args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "V3-5jU5z", "Errors.Internal")
	}
	defer rows.Close()

	for rows.Next() {
		if err := scanToSequence(rows, sequences); err != nil {
			return nil, zerrors.ThrowInternal(err, "V3-Ydiwv", "Errors.Internal")
		}
	}

	if rows.Err() != nil {
		return nil, zerrors.ThrowInternal(rows.Err(), "V3-XApDk", "Errors.Internal")
	}
	return sequences, nil
}

func searchSequenceByCommand(sequences []*latestSequence, command eventstore.Command) *latestSequence {
	for _, sequence := range sequences {
		if sequence.aggregate.Type == command.Aggregate().Type &&
			sequence.aggregate.ID == command.Aggregate().ID &&
			sequence.aggregate.InstanceID == command.Aggregate().InstanceID {
			return sequence
		}
	}
	return nil
}

func searchSequence(sequences []*latestSequence, aggregateType eventstore.AggregateType, aggregateID, instanceID string) *latestSequence {
	for _, sequence := range sequences {
		if sequence.aggregate.Type == aggregateType &&
			sequence.aggregate.ID == aggregateID &&
			sequence.aggregate.InstanceID == instanceID {
			return sequence
		}
	}
	return nil
}

func commandsToSequences(ctx context.Context, commands []eventstore.Command) []*latestSequence {
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

	return sequences
}

const argsPerCondition = 3

func sequencesToSql(sequences []*latestSequence) (conditions []string, args []any) {
	args = make([]interface{}, 0, len(sequences)*argsPerCondition)
	conditions = make([]string, len(sequences))

	for i, sequence := range sequences {
		conditions[i] = fmt.Sprintf(`(SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $%d AND aggregate_type = $%d AND aggregate_id = $%d ORDER BY "sequence" DESC LIMIT 1)`,
			i*argsPerCondition+1,
			i*argsPerCondition+2,
			i*argsPerCondition+3,
		)
		args = append(args, sequence.aggregate.InstanceID, sequence.aggregate.Type, sequence.aggregate.ID)
	}

	return conditions, args
}

func scanToSequence(rows *sql.Rows, sequences []*latestSequence) error {
	var aggregateType eventstore.AggregateType
	var aggregateID, instanceID string
	var currentSequence uint64
	var resourceOwner string

	if err := rows.Scan(&instanceID, &resourceOwner, &aggregateType, &aggregateID, &currentSequence); err != nil {
		return zerrors.ThrowInternal(err, "V3-OIWqj", "Errors.Internal")
	}

	sequence := searchSequence(sequences, aggregateType, aggregateID, instanceID)
	if sequence == nil {
		logging.WithFields(
			"aggType", aggregateType,
			"aggID", aggregateID,
			"instance", instanceID,
		).Panic("no sequence found")
		// added return for linting
		return nil
	}
	sequence.sequence = currentSequence
	if sequence.aggregate.ResourceOwner == "" {
		sequence.aggregate.ResourceOwner = resourceOwner
	}

	return nil
}
