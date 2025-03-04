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
	deviceAuthRequestTable = table{
		name:          projection.DeviceAuthRequestProjectionTable,
		instanceIDCol: projection.DeviceAuthRequestColumnInstanceID,
	}
	DeviceAuthRequestColumnClientID = Column{
		name:  projection.DeviceAuthRequestColumnClientID,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnDeviceCode = Column{
		name:  projection.DeviceAuthRequestColumnDeviceCode,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnUserCode = Column{
		name:  projection.DeviceAuthRequestColumnUserCode,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnScopes = Column{
		name:  projection.DeviceAuthRequestColumnScopes,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnAudience = Column{
		name:  projection.DeviceAuthRequestColumnAudience,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnCreationDate = Column{
		name:  projection.DeviceAuthRequestColumnCreationDate,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnChangeDate = Column{
		name:  projection.DeviceAuthRequestColumnChangeDate,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnSequence = Column{
		name:  projection.DeviceAuthRequestColumnSequence,
		table: deviceAuthRequestTable,
	}
	DeviceAuthRequestColumnInstanceID = Column{
		name:  projection.DeviceAuthRequestColumnInstanceID,
		table: deviceAuthRequestTable,
	}
)

// DeviceAuthRequestByUserCode finds a Device Authorization request by User-Code from the `device_auth_requests` projection.
func (q *Queries) DeviceAuthRequestByUserCode(ctx context.Context, userCode string) (authReq *domain.AuthRequestDevice, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareDeviceAuthQuery()
	eq := sq.Eq{
		DeviceAuthRequestColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		DeviceAuthRequestColumnUserCode.identifier():   userCode,
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Axu7l", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		authReq, err = scan(row)
		return err
	}, query, args...)
	return authReq, err
}

var deviceAuthSelectColumns = []string{
	DeviceAuthRequestColumnClientID.identifier(),
	DeviceAuthRequestColumnDeviceCode.identifier(),
	DeviceAuthRequestColumnUserCode.identifier(),
	DeviceAuthRequestColumnScopes.identifier(),
	DeviceAuthRequestColumnAudience.identifier(),
	AppColumnName.identifier(),
	ProjectColumnName.identifier(),
}

func prepareDeviceAuthQuery() (sq.SelectBuilder, func(*sql.Row) (*domain.AuthRequestDevice, error)) {
	return sq.Select(deviceAuthSelectColumns...).
			From(deviceAuthRequestTable.identifier()).
			LeftJoin(join(AppOIDCConfigColumnClientID, DeviceAuthRequestColumnClientID)).
			LeftJoin(join(AppColumnID, AppOIDCConfigColumnAppID)).
			LeftJoin(join(ProjectColumnID, AppColumnProjectID)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*domain.AuthRequestDevice, error) {
			dst := new(domain.AuthRequestDevice)
			var (
				scopes      database.TextArray[string]
				audience    database.TextArray[string]
				appName     sql.NullString
				projectName sql.NullString
			)

			err := row.Scan(
				&dst.ClientID,
				&dst.DeviceCode,
				&dst.UserCode,
				&scopes,
				&audience,
				&appName,
				&projectName,
			)
			if errors.Is(err, sql.ErrNoRows) {
				return nil, zerrors.ThrowNotFound(err, "QUERY-Sah9a", "Errors.DeviceAuth.NotFound")
			}
			if err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-Voo3o", "Errors.Internal")
			}
			dst.Scopes = scopes
			dst.Audience = audience
			dst.AppName = appName.String
			dst.ProjectName = projectName.String

			return dst, nil
		}
}
