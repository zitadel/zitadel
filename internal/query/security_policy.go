package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	SecurityPolicyColumnEnableIframeEmbedding = Column{
		name:  projection.SecurityPolicyColumnEnableIframeEmbedding,
		table: securityPolicyTable,
	}
	SecurityPolicyColumnAllowedOrigins = Column{
		name:  projection.SecurityPolicyColumnAllowedOrigins,
		table: securityPolicyTable,
	}
	SecurityPolicyColumnEnableImpersonation = Column{
		name:  projection.SecurityPolicyColumnEnableImpersonation,
		table: securityPolicyTable,
	}
)

type SecurityPolicy struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	EnableIframeEmbedding bool
	AllowedOrigins        database.TextArray[string]
	EnableImpersonation   bool
}

func (q *Queries) SecurityPolicy(ctx context.Context) (policy *SecurityPolicy, err error) {
	stmt, scan := prepareSecurityPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		SecurityPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Sf6d1", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		policy, err = scan(row)
		return err
	}, query, args...)
	return policy, err
}

func prepareSecurityPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*SecurityPolicy, error)) {
	return sq.Select(
			SecurityPolicyColumnInstanceID.identifier(),
			SecurityPolicyColumnCreationDate.identifier(),
			SecurityPolicyColumnChangeDate.identifier(),
			SecurityPolicyColumnInstanceID.identifier(),
			SecurityPolicyColumnSequence.identifier(),
			SecurityPolicyColumnEnableIframeEmbedding.identifier(),
			SecurityPolicyColumnAllowedOrigins.identifier(),
			SecurityPolicyColumnEnableImpersonation.identifier()).
			From(securityPolicyTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SecurityPolicy, error) {
			securityPolicy := new(SecurityPolicy)
			err := row.Scan(
				&securityPolicy.AggregateID,
				&securityPolicy.CreationDate,
				&securityPolicy.ChangeDate,
				&securityPolicy.ResourceOwner,
				&securityPolicy.Sequence,
				&securityPolicy.EnableIframeEmbedding,
				&securityPolicy.AllowedOrigins,
				&securityPolicy.EnableImpersonation,
			)
			if err != nil && !errors.Is(err, sql.ErrNoRows) { // ignore not found errors
				return nil, zerrors.ThrowInternal(err, "QUERY-Dfrt2", "Errors.Internal")
			}
			return securityPolicy, nil
		}
}
