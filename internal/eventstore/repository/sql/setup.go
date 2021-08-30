package sql

import (
	"context"

	repo "github.com/caos/zitadel/internal/eventstore/repository"
)

func (db *CRDB) Step20(ctx context.Context, latestSequence uint64) error {
	currentSequence := uint64(1)
	maxSequence := uint64(1000)
	previousSequences := make(map[repo.AggregateType]uint64)
	for maxSequence < latestSequence {
		events, err := db.Filter(ctx, &repo.SearchQuery{
			Columns: repo.ColumnsEvent,
			Filters: [][]*repo.Filter{
				{
					&repo.Filter{
						Field:     repo.FieldSequence,
						Operation: repo.OperationGreater,
						Value:     currentSequence - 1,
					},
					&repo.Filter{
						Field:     repo.FieldSequence,
						Operation: repo.OperationLess,
						Value:     maxSequence,
					},
				},
			},
		})
		if err != nil {
			return err
		}

		tx, err := db.client.Begin()
		if err != nil {
			return err
		}

		for _, event := range events {
			if _, err := tx.Exec("SAVEPOINT event_update"); err != nil {
				return err
			}
			if _, err = tx.Exec("UPDATE eventstore.events SET previous_aggregate_type_sequence = $1 WHERE event_sequence = $2", previousSequences[event.AggregateType], event.Sequence); err != nil {
				return err
			}
			if _, err = tx.Exec("RELEASE SAVEPOINT event_update"); err != nil {
				return err
			}
			previousSequences[event.AggregateType] = event.Sequence
			currentSequence = event.Sequence
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		maxSequence = currentSequence + 1000
	}
	return nil
}
