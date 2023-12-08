package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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

	AccessTokenLifetime        time.Duration `json:"access_token_lifetime,omitempty"`
	IdTokenLifetime            time.Duration `json:"id_token_lifetime,omitempty"`
	RefreshTokenIdleExpiration time.Duration `json:"refresh_token_idle_expiration,omitempty"`
	RefreshTokenExpiration     time.Duration `json:"refresh_token_expiration,omitempty"`
}

func (q *Queries) OIDCSettingsByAggID(ctx context.Context, aggregateID string) (settings *OIDCSettings, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareOIDCSettingsQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		OIDCSettingsColumnAggregateID.identifier(): aggregateID,
		OIDCSettingsColumnInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-s9nle", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		settings, err = scan(row)
		return err
	}, query, args...)
	return settings, err
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
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-s9nlw", "Errors.OIDCSettings.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-9bf8s", "Errors.Internal")
			}
			return oidcSettings, nil
		}
}
