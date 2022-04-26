package sql

import (
	"context"

	"github.com/zitadel/logging"
	repo "github.com/zitadel/zitadel/internal/eventstore/repository"
)

func (db *CRDB) Step20(ctx context.Context, latestSequence uint64) error {
	currentSequence := uint64(0)
	limit := uint64(500)
	previousSequences := make(map[repo.AggregateType]Sequence)
	for currentSequence < latestSequence {
		events, err := db.Filter(ctx, &repo.SearchQuery{
			Columns: repo.ColumnsEvent,
			Limit:   limit,
			Filters: [][]*repo.Filter{
				{
					&repo.Filter{
						Field:     repo.FieldSequence,
						Operation: repo.OperationGreater,
						Value:     currentSequence,
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
			seq := Sequence(previousSequences[event.AggregateType])
			if _, err = tx.Exec("UPDATE eventstore.events SET previous_aggregate_type_sequence = $1 WHERE event_sequence = $2", &seq, event.Sequence); err != nil {
				return err
			}
			if _, err = tx.Exec("RELEASE SAVEPOINT event_update"); err != nil {
				return err
			}
			previousSequences[event.AggregateType] = Sequence(event.Sequence)
			currentSequence = event.Sequence
		}

		if err = tx.Commit(); err != nil {
			return err
		}
		logging.LogWithFields("SQL-bXVwS", "currentSeq", currentSequence, "events", len(events)).Info("events updated")
	}
	return nil
}
