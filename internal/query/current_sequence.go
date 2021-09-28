package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/errors"
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
					return nil, errors.ThrowNotFound(err, "QUERY-gmd9o", "errors.current_sequence.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-aAZ1D", "errors.internal")
			}
			return seq, nil
		}
}

func (q *Queries) latestSequence(ctx context.Context, projection table) (*LatestSequence, error) {
	query, scan := prepareLatestSequence()
	stmt, args, err := query.Where(sq.Eq{
		CurrentSequenceColProjectionName.identifier(): projection.name,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5CfX9", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

var (
	currentSequencesTable = table{
		name: "zitadel.projections.current_sequences",
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
