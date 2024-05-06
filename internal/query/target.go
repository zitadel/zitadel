package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	targetTable = table{
		name:          projection.TargetTable,
		instanceIDCol: projection.TargetInstanceIDCol,
	}
	TargetColumnID = Column{
		name:  projection.TargetIDCol,
		table: targetTable,
	}
	TargetColumnCreationDate = Column{
		name:  projection.TargetCreationDateCol,
		table: targetTable,
	}
	TargetColumnChangeDate = Column{
		name:  projection.TargetChangeDateCol,
		table: targetTable,
	}
	TargetColumnResourceOwner = Column{
		name:  projection.TargetResourceOwnerCol,
		table: targetTable,
	}
	TargetColumnInstanceID = Column{
		name:  projection.TargetInstanceIDCol,
		table: targetTable,
	}
	TargetColumnSequence = Column{
		name:  projection.TargetSequenceCol,
		table: targetTable,
	}
	TargetColumnName = Column{
		name:  projection.TargetNameCol,
		table: targetTable,
	}
	TargetColumnTargetType = Column{
		name:  projection.TargetTargetType,
		table: targetTable,
	}
	TargetColumnURL = Column{
		name:  projection.TargetEndpointCol,
		table: targetTable,
	}
	TargetColumnTimeout = Column{
		name:  projection.TargetTimeoutCol,
		table: targetTable,
	}
	TargetColumnInterruptOnError = Column{
		name:  projection.TargetInterruptOnErrorCol,
		table: targetTable,
	}
)

type Targets struct {
	SearchResponse
	Targets []*Target
}

func (t *Targets) SetState(s *State) {
	t.State = s
}

type Target struct {
	ID string
	domain.ObjectDetails

	Name             string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
}

type TargetSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *TargetSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchTargets(ctx context.Context, queries *TargetSearchQueries) (targets *Targets, err error) {
	eq := sq.Eq{
		TargetColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareTargetsQuery(ctx, q.client)
	return genericRowsQueryWithState[*Targets](ctx, q.client, targetTable, combineToWhereStmt(query, queries.toQuery, eq), scan)
}

func (q *Queries) GetTargetByID(ctx context.Context, id string) (target *Target, err error) {
	eq := sq.Eq{
		TargetColumnID.identifier():         id,
		TargetColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareTargetQuery(ctx, q.client)
	return genericRowQuery[*Target](ctx, q.client, query.Where(eq), scan)
}

func NewTargetNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(TargetColumnName, value, method)
}

func NewTargetInIDsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(TargetColumnID, values)
}

func prepareTargetsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*Targets, error)) {
	return sq.Select(
			TargetColumnID.identifier(),
			TargetColumnChangeDate.identifier(),
			TargetColumnResourceOwner.identifier(),
			TargetColumnSequence.identifier(),
			TargetColumnName.identifier(),
			TargetColumnTargetType.identifier(),
			TargetColumnTimeout.identifier(),
			TargetColumnURL.identifier(),
			TargetColumnInterruptOnError.identifier(),
			countColumn.identifier(),
		).From(targetTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Targets, error) {
			targets := make([]*Target, 0)
			var count uint64
			for rows.Next() {
				target := new(Target)
				err := rows.Scan(
					&target.ID,
					&target.EventDate,
					&target.ResourceOwner,
					&target.Sequence,
					&target.Name,
					&target.TargetType,
					&target.Timeout,
					&target.Endpoint,
					&target.InterruptOnError,
					&count,
				)
				if err != nil {
					return nil, err
				}
				targets = append(targets, target)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-fzwi6cgxos", "Errors.Query.CloseRows")
			}

			return &Targets{
				Targets: targets,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareTargetQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (*Target, error)) {
	return sq.Select(
			TargetColumnID.identifier(),
			TargetColumnChangeDate.identifier(),
			TargetColumnResourceOwner.identifier(),
			TargetColumnSequence.identifier(),
			TargetColumnName.identifier(),
			TargetColumnTargetType.identifier(),
			TargetColumnTimeout.identifier(),
			TargetColumnURL.identifier(),
			TargetColumnInterruptOnError.identifier(),
		).From(targetTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Target, error) {
			target := new(Target)
			err := row.Scan(
				&target.ID,
				&target.EventDate,
				&target.ResourceOwner,
				&target.Sequence,
				&target.Name,
				&target.TargetType,
				&target.Timeout,
				&target.Endpoint,
				&target.InterruptOnError,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-hj5oaniyrz", "Errors.Target.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-5qhc19sc49", "Errors.Internal")
			}
			return target, nil
		}
}
