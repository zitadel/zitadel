package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
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
	TargetColumnSigningKey = Column{
		name:  projection.TargetSigningKey,
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
	domain.ObjectDetails

	Name             string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
	signingKey       *crypto.CryptoValue
	SigningKey       string
}

func (t *Target) decryptSigningKey(alg crypto.EncryptionAlgorithm) error {
	if t.signingKey == nil {
		return nil
	}
	keyValue, err := crypto.DecryptString(t.signingKey, alg)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-bxevy3YXwy", "Errors.Internal")
	}
	t.SigningKey = keyValue
	return nil
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

func (q *Queries) SearchTargets(ctx context.Context, queries *TargetSearchQueries) (*Targets, error) {
	eq := sq.Eq{
		TargetColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareTargetsQuery(ctx, q.client)
	targets, err := genericRowsQueryWithState[*Targets](ctx, q.client, targetTable, combineToWhereStmt(query, queries.toQuery, eq), scan)
	if err != nil {
		return nil, err
	}
	for i := range targets.Targets {
		if err := targets.Targets[i].decryptSigningKey(q.targetEncryptionAlgorithm); err != nil {
			return nil, err
		}
	}
	return targets, nil
}

func (q *Queries) GetTargetByID(ctx context.Context, id string) (*Target, error) {
	eq := sq.Eq{
		TargetColumnID.identifier():         id,
		TargetColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareTargetQuery(ctx, q.client)
	target, err := genericRowQuery[*Target](ctx, q.client, query.Where(eq), scan)
	if err != nil {
		return nil, err
	}
	if err := target.decryptSigningKey(q.targetEncryptionAlgorithm); err != nil {
		return nil, err
	}
	return target, nil
}

func NewTargetNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(TargetColumnName, value, method)
}

func NewTargetInIDsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(TargetColumnID, values)
}

func prepareTargetsQuery(context.Context, prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*Targets, error)) {
	return sq.Select(
			TargetColumnID.identifier(),
			TargetColumnCreationDate.identifier(),
			TargetColumnChangeDate.identifier(),
			TargetColumnResourceOwner.identifier(),
			TargetColumnName.identifier(),
			TargetColumnTargetType.identifier(),
			TargetColumnTimeout.identifier(),
			TargetColumnURL.identifier(),
			TargetColumnInterruptOnError.identifier(),
			TargetColumnSigningKey.identifier(),
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
					&target.CreationDate,
					&target.EventDate,
					&target.ResourceOwner,
					&target.Name,
					&target.TargetType,
					&target.Timeout,
					&target.Endpoint,
					&target.InterruptOnError,
					&target.signingKey,
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

func prepareTargetQuery(context.Context, prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (*Target, error)) {
	return sq.Select(
			TargetColumnID.identifier(),
			TargetColumnCreationDate.identifier(),
			TargetColumnChangeDate.identifier(),
			TargetColumnResourceOwner.identifier(),
			TargetColumnName.identifier(),
			TargetColumnTargetType.identifier(),
			TargetColumnTimeout.identifier(),
			TargetColumnURL.identifier(),
			TargetColumnInterruptOnError.identifier(),
			TargetColumnSigningKey.identifier(),
		).From(targetTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Target, error) {
			target := new(Target)
			err := row.Scan(
				&target.ID,
				&target.CreationDate,
				&target.EventDate,
				&target.ResourceOwner,
				&target.Name,
				&target.TargetType,
				&target.Timeout,
				&target.Endpoint,
				&target.InterruptOnError,
				&target.signingKey,
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
