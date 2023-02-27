package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var (
	securityPolicyTable = table{
		name:          projection.SecurityPolicyProjectionTable,
		instanceIDCol: projection.SecurityPolicyColumnInstanceID,
	}
	SecurityPolicyColumnCreationDate = Column{
		name:  projection.SecurityPolicyColumnCreationDate,
		table: securityPolicyTable,
	}
	SecurityPolicyColumnChangeDate = Column{
		name:  projection.SecurityPolicyColumnChangeDate,
		table: securityPolicyTable,
	}
	SecurityPolicyColumnInstanceID = Column{
		name:  projection.SecurityPolicyColumnInstanceID,
		table: securityPolicyTable,
	}
	SecurityPolicyColumnSequence = Column{
		name:  projection.SecurityPolicyColumnSequence,
		table: securityPolicyTable,
	}
	SecurityPolicyColumnEnabled = Column{
		name:  projection.SecurityPolicyColumnEnabled,
		table: securityPolicyTable,
	}
	SecurityPolicyColumnAllowedOrigins = Column{
		name:  projection.SecurityPolicyColumnAllowedOrigins,
		table: securityPolicyTable,
	}
)

type SecurityPolicy struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	Enabled        bool
	AllowedOrigins database.StringArray
}

func (q *Queries) SecurityPolicy(ctx context.Context) (*SecurityPolicy, error) {
	stmt, scan := prepareSecurityPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		SecurityPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Sf6d1", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareSecurityPolicyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*SecurityPolicy, error)) {
	return sq.Select(
			SecurityPolicyColumnInstanceID.identifier(),
			SecurityPolicyColumnCreationDate.identifier(),
			SecurityPolicyColumnChangeDate.identifier(),
			SecurityPolicyColumnInstanceID.identifier(),
			SecurityPolicyColumnSequence.identifier(),
			SecurityPolicyColumnEnabled.identifier(),
			SecurityPolicyColumnAllowedOrigins.identifier()).
			From(securityPolicyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SecurityPolicy, error) {
			securityPolicy := new(SecurityPolicy)
			err := row.Scan(
				&securityPolicy.AggregateID,
				&securityPolicy.CreationDate,
				&securityPolicy.ChangeDate,
				&securityPolicy.ResourceOwner,
				&securityPolicy.Sequence,
				&securityPolicy.Enabled,
				&securityPolicy.AllowedOrigins,
			)
			if err != nil && !errs.Is(err, sql.ErrNoRows) { // ignore not found errors
				return nil, errors.ThrowInternal(err, "QUERY-Dfrt2", "Errors.Internal")
			}
			return securityPolicy, nil
		}
}
