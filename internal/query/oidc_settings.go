package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	oidcSettingsTable = table{
		name:          projection.OIDCSettingsProjectionTable,
		instanceIDCol: projection.OIDCSettingsColumnInstanceID,
	}
	OIDCSettingsColumnAggregateID = Column{
		name:  projection.OIDCSettingsColumnAggregateID,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnCreationDate = Column{
		name:  projection.OIDCSettingsColumnCreationDate,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnChangeDate = Column{
		name:  projection.OIDCSettingsColumnChangeDate,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnResourceOwner = Column{
		name:  projection.OIDCSettingsColumnResourceOwner,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnInstanceID = Column{
		name:  projection.OIDCSettingsColumnInstanceID,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnSequence = Column{
		name:  projection.OIDCSettingsColumnSequence,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnAccessTokenLifetime = Column{
		name:  projection.OIDCSettingsColumnAccessTokenLifetime,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnIdTokenLifetime = Column{
		name:  projection.OIDCSettingsColumnIdTokenLifetime,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnRefreshTokenIdleExpiration = Column{
		name:  projection.OIDCSettingsColumnRefreshTokenIdleExpiration,
		table: oidcSettingsTable,
	}
	OIDCSettingsColumnRefreshTokenExpiration = Column{
		name:  projection.OIDCSettingsColumnRefreshTokenExpiration,
		table: oidcSettingsTable,
	}
)

type OIDCSettings struct {
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

func (q *Queries) OIDCSettingsByAggID(ctx context.Context, aggregateID string) (_ *OIDCSettings, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareOIDCSettingsQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		OIDCSettingsColumnAggregateID.identifier(): aggregateID,
		OIDCSettingsColumnInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-s9nle", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareOIDCSettingsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*OIDCSettings, error)) {
	return sq.Select(
			OIDCSettingsColumnAggregateID.identifier(),
			OIDCSettingsColumnCreationDate.identifier(),
			OIDCSettingsColumnChangeDate.identifier(),
			OIDCSettingsColumnResourceOwner.identifier(),
			OIDCSettingsColumnSequence.identifier(),
			OIDCSettingsColumnAccessTokenLifetime.identifier(),
			OIDCSettingsColumnIdTokenLifetime.identifier(),
			OIDCSettingsColumnRefreshTokenIdleExpiration.identifier(),
			OIDCSettingsColumnRefreshTokenExpiration.identifier()).
			From(oidcSettingsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*OIDCSettings, error) {
			oidcSettings := new(OIDCSettings)
			err := row.Scan(
				&oidcSettings.AggregateID,
				&oidcSettings.CreationDate,
				&oidcSettings.ChangeDate,
				&oidcSettings.ResourceOwner,
				&oidcSettings.Sequence,
				&oidcSettings.AccessTokenLifetime,
				&oidcSettings.IdTokenLifetime,
				&oidcSettings.RefreshTokenIdleExpiration,
				&oidcSettings.RefreshTokenExpiration,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-s9nlw", "Errors.OIDCSettings.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-9bf8s", "Errors.Internal")
			}
			return oidcSettings, nil
		}
}
