package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type DomainPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	UserLoginMustBeDomain bool

	IsDefault bool
}

var (
	domainPolicyTable = table{
		name: projection.DomainPolicyTable,
	}
	DomainPolicyColID = Column{
		name:  projection.DomainPolicyIDCol,
		table: domainPolicyTable,
	}
	DomainPolicyColSequence = Column{
		name:  projection.DomainPolicySequenceCol,
		table: domainPolicyTable,
	}
	DomainPolicyColCreationDate = Column{
		name:  projection.DomainPolicyCreationDateCol,
		table: domainPolicyTable,
	}
	DomainPolicyColChangeDate = Column{
		name:  projection.DomainPolicyChangeDateCol,
		table: domainPolicyTable,
	}
	DomainPolicyColResourceOwner = Column{
		name:  projection.DomainPolicyResourceOwnerCol,
		table: domainPolicyTable,
	}
	DomainPolicyColInstanceID = Column{
		name:  projection.DomainPolicyInstanceIDCol,
		table: domainPolicyTable,
	}
	DomainPolicyColUserLoginMustBeDomain = Column{
		name:  projection.DomainPolicyUserLoginMustBeDomainCol,
		table: domainPolicyTable,
	}
	DomainPolicyColIsDefault = Column{
		name:  projection.DomainPolicyIsDefaultCol,
		table: domainPolicyTable,
	}
	DomainPolicyColState = Column{
		name:  projection.DomainPolicyStateCol,
		table: domainPolicyTable,
	}
)

func (q *Queries) DomainPolicyByOrg(ctx context.Context, orgID string) (*DomainPolicy, error) {
	stmt, scan := prepareDomainPolicyQuery()
	query, args, err := stmt.Where(
		sq.And{
			sq.Eq{
				DomainPolicyColInstanceID.identifier(): authz.GetInstance(ctx).ID,
			},
			sq.Or{
				sq.Eq{
					DomainPolicyColID.identifier(): orgID,
				},
				sq.Eq{
					DomainPolicyColID.identifier(): domain.IAMID,
				},
			},
		}).
		OrderBy(DomainPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-D3CqT", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultDomainPolicy(ctx context.Context) (*DomainPolicy, error) {
	stmt, scan := prepareDomainPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		DomainPolicyColID.identifier():         domain.IAMID,
		DomainPolicyColInstanceID.identifier(): authz.GetInstance(ctx).ID,
	}).
		OrderBy(DomainPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-pM7lP", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareDomainPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*DomainPolicy, error)) {
	return sq.Select(
			DomainPolicyColID.identifier(),
			DomainPolicyColSequence.identifier(),
			DomainPolicyColCreationDate.identifier(),
			DomainPolicyColChangeDate.identifier(),
			DomainPolicyColResourceOwner.identifier(),
			DomainPolicyColUserLoginMustBeDomain.identifier(),
			DomainPolicyColIsDefault.identifier(),
			DomainPolicyColState.identifier(),
		).
			From(domainPolicyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*DomainPolicy, error) {
			policy := new(DomainPolicy)
			err := row.Scan(
				&policy.ID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.ResourceOwner,
				&policy.UserLoginMustBeDomain,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-K0Jr5", "Errors.DomainPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-rIy6j", "Errors.Internal")
			}
			return policy, nil
		}
}
