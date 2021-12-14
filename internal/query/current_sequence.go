package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type LatestSequence struct {
	Sequence  uint64
	Timestamp time.Time
}

func prepareLatestSequence() (sq.SelectBuilder, func(*sql.Row) (*LatestSequence, error)) {
	return sq.Select(
			CurrentSequenceColCurrentSequence.identifier(),
			CurrentSequenceColTimestamp.identifier()).
			From(currentSequencesTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LatestSequence, error) {
			seq := new(LatestSequence)
			err := row.Scan(
				&seq.Sequence,
				&seq.Timestamp,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-gmd9o", "Errors.CurrentSequence.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-aAZ1D", "Errors.Internal")
			}
			return seq, nil
		}
}

func (q *Queries) latestSequence(ctx context.Context, projections ...table) (*LatestSequence, error) {
	query, scan := prepareLatestSequence()
	or := make(sq.Or, len(projections))
	for i, projection := range projections {
		or[i] = sq.Eq{CurrentSequenceColProjectionName.identifier(): projection.name}
	}
	stmt, args, err := query.
		Where(or).
		OrderBy(CurrentSequenceColCurrentSequence.identifier()).
		ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5CfX9", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

var (
	currentSequencesTable = table{
		name: projection.CurrentSeqTable,
	}
	CurrentSequenceColAggregateType = Column{
		name:  "aggregate_type",
		table: currentSequencesTable,
	}
	CurrentSequenceColCurrentSequence = Column{
		name:  "current_sequence",
		table: currentSequencesTable,
	}
	CurrentSequenceColTimestamp = Column{
		name:  "timestamp",
		table: currentSequencesTable,
	}
	CurrentSequenceColProjectionName = Column{
		name:  "projection_name",
		table: currentSequencesTable,
	}
)
