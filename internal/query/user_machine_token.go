package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/query/projection"

	"github.com/caos/zitadel/internal/errors"
)

var (
	machineTokensTable = table{
		name: projection.MachineTokenProjectionTable,
	}
	MachineTokenColumnID = Column{
		name:  projection.MachineTokenColumnID,
		table: machineTokensTable,
	}
	MachineTokenColumnUserID = Column{
		name:  projection.MachineTokenColumnUserID,
		table: machineTokensTable,
	}
	MachineTokenColumnExpiration = Column{
		name:  projection.MachineTokenColumnExpiration,
		table: machineTokensTable,
	}
	MachineTokenColumnScopes = Column{
		name:  projection.MachineTokenColumnScopes,
		table: machineTokensTable,
	}
	MachineTokenColumnCreationDate = Column{
		name:  projection.MachineTokenColumnCreationDate,
		table: machineTokensTable,
	}
	MachineTokenColumnChangeDate = Column{
		name:  projection.MachineTokenColumnChangeDate,
		table: machineTokensTable,
	}
	MachineTokenColumnResourceOwner = Column{
		name:  projection.MachineTokenColumnResourceOwner,
		table: machineTokensTable,
	}
	MachineTokenColumnSequence = Column{
		name:  projection.MachineTokenColumnSequence,
		table: machineTokensTable,
	}
)

type MachineTokens struct {
	SearchResponse
	MachineTokens []*MachineToken
}

type MachineToken struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	UserID     string
	Expiration time.Time
	Scopes     []string
}

type MachineTokenSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) MachineTokenByID(ctx context.Context, id string, queries ...SearchQuery) (*MachineToken, error) {
	query, scan := prepareMachineTokenQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	stmt, args, err := query.Where(sq.Eq{
		MachineTokenColumnID.identifier(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgfb4", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) SearchMachineTokens(ctx context.Context, queries *MachineTokenSearchQueries) (machineTokens *MachineTokens, err error) {
	query, scan := prepareMachineTokensQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-Hjw2w", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Bmz63", "Errors.Internal")
	}
	machineTokens, err = scan(rows)
	if err != nil {
		return nil, err
	}
	machineTokens.LatestSequence, err = q.latestSequence(ctx, machineTokensTable)
	return machineTokens, err
}

func NewMachineTokenResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(MachineTokenColumnResourceOwner, value, TextEquals)
}

func NewMachineTokenUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(MachineTokenColumnUserID, value, TextEquals)
}

func (r *MachineTokenSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewMachineTokenResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (q *MachineTokenSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareMachineTokenQuery() (sq.SelectBuilder, func(*sql.Row) (*MachineToken, error)) {
	return sq.Select(
			MachineTokenColumnID.identifier(),
			MachineTokenColumnCreationDate.identifier(),
			MachineTokenColumnChangeDate.identifier(),
			MachineTokenColumnResourceOwner.identifier(),
			MachineTokenColumnSequence.identifier(),
			MachineTokenColumnUserID.identifier(),
			MachineTokenColumnExpiration.identifier(),
			MachineTokenColumnScopes.identifier()).
			From(machineTokensTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*MachineToken, error) {
			p := new(MachineToken)
			scopes := pq.StringArray{}
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.Sequence,
				&p.UserID,
				&p.Expiration,
				&scopes,
			)
			p.Scopes = scopes
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-fk2fs", "Errors.MachineToken.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-dj2FF", "Errors.Internal")
			}
			return p, nil
		}
}

func prepareMachineTokensQuery() (sq.SelectBuilder, func(*sql.Rows) (*MachineTokens, error)) {
	return sq.Select(
			MachineTokenColumnID.identifier(),
			MachineTokenColumnCreationDate.identifier(),
			MachineTokenColumnChangeDate.identifier(),
			MachineTokenColumnResourceOwner.identifier(),
			MachineTokenColumnSequence.identifier(),
			MachineTokenColumnUserID.identifier(),
			MachineTokenColumnExpiration.identifier(),
			MachineTokenColumnScopes.identifier(),
			countColumn.identifier()).
			From(machineTokensTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*MachineTokens, error) {
			machineTokens := make([]*MachineToken, 0)
			var count uint64
			for rows.Next() {
				token := new(MachineToken)
				scopes := pq.StringArray{}
				err := rows.Scan(
					&token.ID,
					&token.CreationDate,
					&token.ChangeDate,
					&token.ResourceOwner,
					&token.Sequence,
					&token.UserID,
					&token.Expiration,
					&scopes,
					&count,
				)
				if err != nil {
					return nil, err
				}
				token.Scopes = scopes
				machineTokens = append(machineTokens, token)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-QMXJv", "Errors.Query.CloseRows")
			}

			return &MachineTokens{
				MachineTokens: machineTokens,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
