package query

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	deviceAuthTable = table{
		name:          projection.DeviceAuthProjectionTable,
		instanceIDCol: projection.DeviceAuthColumnInstanceID,
	}
	DeviceAuthColumnID = Column{
		name:  projection.DeviceAuthColumnID,
		table: deviceAuthTable,
	}
	DeviceAuthColumnClientID = Column{
		name:  projection.DeviceAuthColumnClientID,
		table: deviceAuthTable,
	}
	DeviceAuthColumnDeviceCode = Column{
		name:  projection.DeviceAuthColumnDeviceCode,
		table: deviceAuthTable,
	}
	DeviceAuthColumnUserCode = Column{
		name:  projection.DeviceAuthColumnUserCode,
		table: deviceAuthTable,
	}
	DeviceAuthColumnExpires = Column{
		name:  projection.DeviceAuthColumnExpires,
		table: deviceAuthTable,
	}
	DeviceAuthColumnScopes = Column{
		name:  projection.DeviceAuthColumnScopes,
		table: deviceAuthTable,
	}
	DeviceAuthColumnState = Column{
		name:  projection.DeviceAuthColumnState,
		table: deviceAuthTable,
	}
	DeviceAuthColumnSubject = Column{
		name:  projection.DeviceAuthColumnSubject,
		table: deviceAuthTable,
	}
	DeviceAuthColumnCreationDate = Column{
		name:  projection.DeviceAuthColumnCreationDate,
		table: deviceAuthTable,
	}
	DeviceAuthColumnChangeDate = Column{
		name:  projection.DeviceAuthColumnChangeDate,
		table: deviceAuthTable,
	}
	DeviceAuthColumnSequence = Column{
		name:  projection.DeviceAuthColumnSequence,
		table: deviceAuthTable,
	}
	DeviceAuthColumnInstanceID = Column{
		name:  projection.DeviceAuthColumnInstanceID,
		table: deviceAuthTable,
	}
)

func (q *Queries) DeviceAuthByDeviceCode(ctx context.Context, clientID, deviceCode string) (deviceAuth *domain.DeviceAuth, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareDeviceAuthQuery(ctx, q.client)
	eq := sq.Eq{
		DeviceAuthColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		DeviceAuthColumnClientID.identifier():   clientID,
		DeviceAuthColumnDeviceCode.identifier(): deviceCode,
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-uk1Oh", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		deviceAuth, err = scan(row)
		return err
	}, query, args...)
	return deviceAuth, err
}

func (q *Queries) DeviceAuthByUserCode(ctx context.Context, userCode string) (deviceAuth *domain.DeviceAuth, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareDeviceAuthQuery(ctx, q.client)
	eq := sq.Eq{
		DeviceAuthColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		DeviceAuthColumnUserCode.identifier():   userCode,
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Axu7l", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		deviceAuth, err = scan(row)
		return err
	}, query, args...)
	return deviceAuth, err
}

var deviceAuthSelectColumns = []string{
	DeviceAuthColumnID.identifier(),
	DeviceAuthColumnClientID.identifier(),
	DeviceAuthColumnScopes.identifier(),
	DeviceAuthColumnExpires.identifier(),
	DeviceAuthColumnState.identifier(),
	DeviceAuthColumnSubject.identifier(),
}

func prepareDeviceAuthQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*domain.DeviceAuth, error)) {
	return sq.Select(deviceAuthSelectColumns...).From(deviceAuthTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*domain.DeviceAuth, error) {
			dst := new(domain.DeviceAuth)
			var scopes database.TextArray[string]

			err := row.Scan(
				&dst.AggregateID,
				&dst.ClientID,
				&scopes,
				&dst.Expires,
				&dst.State,
				&dst.Subject,
			)
			if errors.Is(err, sql.ErrNoRows) {
				return nil, zerrors.ThrowNotFound(err, "QUERY-Sah9a", "Errors.DeviceAuth.NotExisting")
			}
			if err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-Voo3o", "Errors.Internal")
			}

			dst.Scopes = scopes
			return dst, nil
		}
}
