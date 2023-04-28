package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	secretGeneratorsTable = table{
		name:          projection.SecretGeneratorProjectionTable,
		instanceIDCol: projection.SecretGeneratorColumnInstanceID,
	}
	SecretGeneratorColumnAggregateID = Column{
		name:  projection.SecretGeneratorColumnAggregateID,
		table: secretGeneratorsTable,
	}
	SecretGeneratorColumnInstanceID = Column{
		name:  projection.SecretGeneratorColumnInstanceID,
		table: secretGeneratorsTable,
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
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	GeneratorType       domain.SecretGeneratorType
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

func (q *Queries) InitEncryptionGenerator(ctx context.Context, generatorType domain.SecretGeneratorType, algorithm crypto.EncryptionAlgorithm) (crypto.Generator, error) {
	generatorConfig, err := q.SecretGeneratorByType(ctx, generatorType)
	if err != nil {
		return nil, err
	}
	cryptoConfig := crypto.GeneratorConfig{
		Length:              generatorConfig.Length,
		Expiry:              generatorConfig.Expiry,
		IncludeLowerLetters: generatorConfig.IncludeLowerLetters,
		IncludeUpperLetters: generatorConfig.IncludeUpperLetters,
		IncludeDigits:       generatorConfig.IncludeDigits,
		IncludeSymbols:      generatorConfig.IncludeSymbols,
	}
	return crypto.NewEncryptionGenerator(cryptoConfig, algorithm), nil
}

func (q *Queries) InitHashGenerator(ctx context.Context, generatorType domain.SecretGeneratorType, algorithm crypto.HashAlgorithm) (crypto.Generator, error) {
	generatorConfig, err := q.SecretGeneratorByType(ctx, generatorType)
	if err != nil {
		return nil, err
	}
	cryptoConfig := crypto.GeneratorConfig{
		Length:              generatorConfig.Length,
		Expiry:              generatorConfig.Expiry,
		IncludeLowerLetters: generatorConfig.IncludeLowerLetters,
		IncludeUpperLetters: generatorConfig.IncludeUpperLetters,
		IncludeDigits:       generatorConfig.IncludeDigits,
		IncludeSymbols:      generatorConfig.IncludeSymbols,
	}
	return crypto.NewHashGenerator(cryptoConfig, algorithm), nil
}

func (q *Queries) SecretGeneratorByType(ctx context.Context, generatorType domain.SecretGeneratorType) (_ *SecretGenerator, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareSecretGeneratorQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		SecretGeneratorColumnGeneratorType.identifier(): generatorType,
		SecretGeneratorColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-3k99f", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) SearchSecretGenerators(ctx context.Context, queries *SecretGeneratorSearchQueries) (secretGenerators *SecretGenerators, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSecretGeneratorsQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			SecretGeneratorColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
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

func NewSecretGeneratorTypeSearchQuery(value int32) (SearchQuery, error) {
	return NewNumberQuery(SecretGeneratorColumnGeneratorType, value, NumberEquals)
}

func prepareSecretGeneratorQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*SecretGenerator, error)) {
	return sq.Select(
			SecretGeneratorColumnAggregateID.identifier(),
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
			From(secretGeneratorsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SecretGenerator, error) {
			secretGenerator := new(SecretGenerator)
			err := row.Scan(
				&secretGenerator.AggregateID,
				&secretGenerator.GeneratorType,
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
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-m9wff", "Errors.SecretGenerator.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-2k99d", "Errors.Internal")
			}
			return secretGenerator, nil
		}
}

func prepareSecretGeneratorsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*SecretGenerators, error)) {
	return sq.Select(
			SecretGeneratorColumnAggregateID.identifier(),
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
			From(secretGeneratorsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*SecretGenerators, error) {
			secretGenerators := make([]*SecretGenerator, 0)
			var count uint64
			for rows.Next() {
				secretGenerator := new(SecretGenerator)
				err := rows.Scan(
					&secretGenerator.AggregateID,
					&secretGenerator.GeneratorType,
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
