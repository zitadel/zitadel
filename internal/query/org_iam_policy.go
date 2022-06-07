package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
)

type OrgIAMPolicy struct {
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
	orgIAMTable = table{
		name: projection.OrgIAMPolicyTable,
	}
	OrgIAMColID = Column{
		name:  projection.OrgIAMPolicyIDCol,
		table: orgIAMTable,
	}
	OrgIAMColSequence = Column{
		name:  projection.OrgIAMPolicySequenceCol,
		table: orgIAMTable,
	}
	OrgIAMColCreationDate = Column{
		name:  projection.OrgIAMPolicyCreationDateCol,
		table: orgIAMTable,
	}
	OrgIAMColChangeDate = Column{
		name:  projection.OrgIAMPolicyChangeDateCol,
		table: orgIAMTable,
	}
	OrgIAMColResourceOwner = Column{
		name:  projection.OrgIAMPolicyResourceOwnerCol,
		table: orgIAMTable,
	}
	OrgIAMColUserLoginMustBeDomain = Column{
		name:  projection.OrgIAMPolicyUserLoginMustBeDomainCol,
		table: orgIAMTable,
	}
	OrgIAMColIsDefault = Column{
		name:  projection.OrgIAMPolicyIsDefaultCol,
		table: orgIAMTable,
	}
	OrgIAMColState = Column{
		name:  projection.OrgIAMPolicyStateCol,
		table: orgIAMTable,
	}
)

func (q *Queries) OrgIAMPolicyByOrg(ctx context.Context, shouldRealTime bool, orgID string) (*OrgIAMPolicy, error) {
	if shouldRealTime {
		projection.OrgIAMPolicyProjection.TriggerBulk(ctx)
	}
	stmt, scan := prepareOrgIAMPolicyQuery()
	query, args, err := stmt.Where(
		sq.Or{
			sq.Eq{
				OrgIAMColID.identifier(): orgID,
			},
			sq.Eq{
				OrgIAMColID.identifier(): q.iamID,
			},
		}).
		OrderBy(OrgIAMColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-D3CqT", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultOrgIAMPolicy(ctx context.Context) (*OrgIAMPolicy, error) {
	stmt, scan := prepareOrgIAMPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		OrgIAMColID.identifier(): q.iamID,
	}).
		OrderBy(OrgIAMColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-pM7lP", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareOrgIAMPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*OrgIAMPolicy, error)) {
	return sq.Select(
			OrgIAMColID.identifier(),
			OrgIAMColSequence.identifier(),
			OrgIAMColCreationDate.identifier(),
			OrgIAMColChangeDate.identifier(),
			OrgIAMColResourceOwner.identifier(),
			OrgIAMColUserLoginMustBeDomain.identifier(),
			OrgIAMColIsDefault.identifier(),
			OrgIAMColState.identifier(),
		).
			From(orgIAMTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*OrgIAMPolicy, error) {
			policy := new(OrgIAMPolicy)
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
					return nil, errors.ThrowNotFound(err, "QUERY-K0Jr5", "Errors.OrgIAMPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-rIy6j", "Errors.Internal")
			}
			return policy, nil
		}
}
