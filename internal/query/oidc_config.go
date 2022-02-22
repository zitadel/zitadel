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
	oidcConfigsTable = table{
		name: projection.OIDCConfigProjectionTable,
	}
	OIDCConfigColumnAggregateID = Column{
		name:  projection.OIDCConfigColumnAggregateID,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnCreationDate = Column{
		name:  projection.OIDCConfigColumnCreationDate,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnChangeDate = Column{
		name:  projection.OIDCConfigColumnChangeDate,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnResourceOwner = Column{
		name:  projection.OIDCConfigColumnResourceOwner,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnSequence = Column{
		name:  projection.OIDCConfigColumnSequence,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnAccessTokenLifetime = Column{
		name:  projection.OIDCConfigColumnAccessTokenLifetime,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnIdTokenLifetime = Column{
		name:  projection.OIDCConfigColumnIdTokenLifetime,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnRefreshTokenIdleExpiration = Column{
		name:  projection.OIDCConfigColumnRefreshTokenIdleExpiration,
		table: oidcConfigsTable,
	}
	OIDCConfigColumnRefreshTokenExpiration = Column{
		name:  projection.OIDCConfigColumnRefreshTokenExpiration,
		table: oidcConfigsTable,
	}
)

type OIDCConfig struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
}

func (q *Queries) OIDCConfigByAggID(ctx context.Context, aggregateID string) (*OIDCConfig, error) {
	stmt, scan := prepareOIDCConfigQuery()
	query, args, err := stmt.Where(sq.Eq{
		OIDCConfigColumnAggregateID.identifier(): aggregateID,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-s9nle", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareOIDCConfigQuery() (sq.SelectBuilder, func(*sql.Row) (*OIDCConfig, error)) {
	return sq.Select(
			OIDCConfigColumnAggregateID.identifier(),
			OIDCConfigColumnCreationDate.identifier(),
			OIDCConfigColumnChangeDate.identifier(),
			OIDCConfigColumnResourceOwner.identifier(),
			OIDCConfigColumnSequence.identifier(),
			OIDCConfigColumnAccessTokenLifetime.identifier(),
			OIDCConfigColumnIdTokenLifetime.identifier(),
			OIDCConfigColumnRefreshTokenIdleExpiration.identifier(),
			OIDCConfigColumnRefreshTokenExpiration.identifier()).
			From(oidcConfigsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*OIDCConfig, error) {
			oidcConfig := new(OIDCConfig)
			err := row.Scan(
				&oidcConfig.AggregateID,
				&oidcConfig.CreationDate,
				&oidcConfig.ChangeDate,
				&oidcConfig.ResourceOwner,
				&oidcConfig.Sequence,
				&oidcConfig.AccessTokenLifetime,
				&oidcConfig.IdTokenLifetime,
				&oidcConfig.RefreshTokenIdleExpiration,
				&oidcConfig.RefreshTokenExpiration,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-s9nlw", "Errors.OIDCConfig.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-9bf8s", "Errors.Internal")
			}
			return oidcConfig, nil
		}
}
