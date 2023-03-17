package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type DomainPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	UserLoginMustBeDomain                  bool
	ValidateOrgDomains                     bool
	SMTPSenderAddressMatchesInstanceDomain bool

	IsDefault bool
}

var (
	domainPolicyTable = table{
		name:          projection.DomainPolicyTable,
		instanceIDCol: projection.DomainPolicyInstanceIDCol,
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
	DomainPolicyColValidateOrgDomains = Column{
		name:  projection.DomainPolicyValidateOrgDomainsCol,
		table: domainPolicyTable,
	}
	DomainPolicyColSMTPSenderAddressMatchesInstanceDomain = Column{
		name:  projection.DomainPolicySMTPSenderAddressMatchesInstanceDomainCol,
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
	DomainPolicyColOwnerRemoved = Column{
		name:  projection.DomainPolicyOwnerRemovedCol,
		table: domainPolicyTable,
	}
)

func (q *Queries) DomainPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (_ *DomainPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.DomainPolicyProjection.Trigger(ctx)
	}
	eq := sq.And{
		sq.Eq{DomainPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()},
		sq.Or{
			sq.Eq{DomainPolicyColID.identifier(): orgID},
			sq.Eq{DomainPolicyColID.identifier(): authz.GetInstance(ctx).InstanceID()},
		},
	}
	if !withOwnerRemoved {
		eq = sq.And{
			sq.Eq{
				DomainPolicyColInstanceID.identifier():   authz.GetInstance(ctx).InstanceID(),
				DomainPolicyColOwnerRemoved.identifier(): false,
			},
			sq.Or{
				sq.Eq{DomainPolicyColID.identifier(): orgID},
				sq.Eq{DomainPolicyColID.identifier(): authz.GetInstance(ctx).InstanceID()},
			},
		}
	}

	stmt, scan := prepareDomainPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(eq).OrderBy(DomainPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-D3CqT", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultDomainPolicy(ctx context.Context) (_ *DomainPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareDomainPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		DomainPolicyColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		DomainPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(DomainPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-pM7lP", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareDomainPolicyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*DomainPolicy, error)) {
	return sq.Select(
			DomainPolicyColID.identifier(),
			DomainPolicyColSequence.identifier(),
			DomainPolicyColCreationDate.identifier(),
			DomainPolicyColChangeDate.identifier(),
			DomainPolicyColResourceOwner.identifier(),
			DomainPolicyColUserLoginMustBeDomain.identifier(),
			DomainPolicyColValidateOrgDomains.identifier(),
			DomainPolicyColSMTPSenderAddressMatchesInstanceDomain.identifier(),
			DomainPolicyColIsDefault.identifier(),
			DomainPolicyColState.identifier(),
		).
			From(domainPolicyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*DomainPolicy, error) {
			policy := new(DomainPolicy)
			err := row.Scan(
				&policy.ID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.ResourceOwner,
				&policy.UserLoginMustBeDomain,
				&policy.ValidateOrgDomains,
				&policy.SMTPSenderAddressMatchesInstanceDomain,
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
