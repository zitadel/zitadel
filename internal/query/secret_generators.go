package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/query/projection"

	"github.com/caos/zitadel/internal/errors"
)

var (
	secretGeneratorsTable = table{
		name: projection.SecretGeneratorProjectionTable,
	}
	SecretGeneratorColumnGeneratorType = Column{
		name:  projection.SecretGeneratorColumnGeneratorType,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnCreationDate = Column{
		name:  projection.SecretGeneratorColumnCreationDate,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnChangeDate = Column{
		name:  projection.SecretGeneratorColumnChangeDate,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnResourceOwner = Column{
		name:  projection.SecretGeneratorColumnResourceOwner,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnSequence = Column{
		name:  projection.SecretGeneratorColumnSequence,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnLength = Column{
		name:  projection.SecretGeneratorColumnLength,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnExpiry = Column{
		name:  projection.SecretGeneratorColumnExpiry,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnIncludeLowerLetters = Column{
		name:  projection.SecretGeneratorColumnIncludeLowerLetters,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnIncludeUpperLetters = Column{
		name:  projection.SecretGeneratorColumnIncludeUpperLetters,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnIncludeDigits = Column{
		name:  projection.SecretGeneratorColumnIncludeDigits,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnIncludeSymbols = Column{
		name:  projection.SecretGeneratorColumnIncludeSymbols,
		table: secretGeneratorsTable,
	}
)

type SecretGenerators struct {
	SearchResponse
	SecretGenerators []*SecretGenerator
}

type SecretGenerator struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	GeneratorType       string
	Length              uint
	Expiry              time.Duration
	IncludeLowerLetters bool
	IncludeUpperLetters bool
	IncludeDigits       bool
	IncludeSymbols      bool
}

type SecretGeneratorSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) SecretGeneratorByType(ctx context.Context, generatorType string) (*SecretGenerator, error) {
	stmt, scan := prepareSecretGeneratorQuery()
	query, args, err := stmt.Where(sq.Eq{
		SecretGeneratorColumnGeneratorType.identifier(): generatorType,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-3k99f", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) SearchSecretGenerators(ctx context.Context, queries *SecretGeneratorSearchQueries) (secretGenerators *SecretGenerators, err error) {
	query, scan := prepareSecretGeneratorsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-sn9lw", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-4miii", "Errors.Internal")
	}
	secretGenerators, err = scan(rows)
	if err != nil {
		return nil, err
	}
	secretGenerators.LatestSequence, err = q.latestSequence(ctx, secretGeneratorsTable)
	return secretGenerators, err
}

func (q *SecretGeneratorSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareSecretGeneratorQuery() (sq.SelectBuilder, func(*sql.Row) (*SecretGenerator, error)) {
	return sq.Select(
			SecretGeneratorColumnGeneratorType.identifier(),
			SecretGeneratorColumnCreationDate.identifier(),
			SecretGeneratorColumnChangeDate.identifier(),
			SecretGeneratorColumnResourceOwner.identifier(),
			SecretGeneratorColumnSequence.identifier(),
			SecretGeneratorColumnLength.identifier(),
			SecretGeneratorColumnExpiry.identifier(),
			SecretGeneratorColumnIncludeLowerLetters.identifier(),
			SecretGeneratorColumnIncludeUpperLetters.identifier(),
			SecretGeneratorColumnIncludeDigits.identifier(),
			SecretGeneratorColumnIncludeSymbols.identifier()).
			From(secretGeneratorsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SecretGenerator, error) {
			p := new(SecretGenerator)
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.Sequence,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-m9wff", "Errors.SecretGenerator.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-2k99d", "Errors.Internal")
			}
			return p, nil
		}
}

func prepareSecretGeneratorsQuery() (sq.SelectBuilder, func(*sql.Rows) (*SecretGenerators, error)) {
	return sq.Select(
			SecretGeneratorColumnGeneratorType.identifier(),
			SecretGeneratorColumnCreationDate.identifier(),
			SecretGeneratorColumnChangeDate.identifier(),
			SecretGeneratorColumnResourceOwner.identifier(),
			SecretGeneratorColumnSequence.identifier(),
			SecretGeneratorColumnLength.identifier(),
			SecretGeneratorColumnExpiry.identifier(),
			SecretGeneratorColumnIncludeLowerLetters.identifier(),
			SecretGeneratorColumnIncludeUpperLetters.identifier(),
			SecretGeneratorColumnIncludeDigits.identifier(),
			SecretGeneratorColumnIncludeSymbols.identifier(),
			countColumn.identifier()).
			From(secretGeneratorsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*SecretGenerators, error) {
			secretGenerators := make([]*SecretGenerator, 0)
			var count uint64
			for rows.Next() {
				secretGenerator := new(SecretGenerator)
				err := rows.Scan(
					&secretGenerator.ID,
					&secretGenerator.CreationDate,
					&secretGenerator.ChangeDate,
					&secretGenerator.ResourceOwner,
					&secretGenerator.Sequence,
					&secretGenerator.Length,
					&secretGenerator.Expiry,
					&secretGenerator.IncludeLowerLetters,
					&secretGenerator.IncludeUpperLetters,
					&secretGenerator.IncludeDigits,
					&secretGenerator.IncludeSymbols,
					&count,
				)
				if err != nil {
					return nil, err
				}
				secretGenerators = append(secretGenerators, secretGenerator)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-em9fs", "Errors.Query.CloseRows")
			}

			return &SecretGenerators{
				SecretGenerators: secretGenerators,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
