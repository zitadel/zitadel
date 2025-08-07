package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (q *Queries) DomainPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (policy *DomainPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerDomainPolicyProjection")
		ctx, err = projection.DomainPolicyProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
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

	stmt, scan := prepareDomainPolicyQuery()
	query, args, err := stmt.Where(eq).OrderBy(DomainPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-D3CqT", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		policy, err = scan(row)
		return err
	}, query, args...)
	return policy, err
}

func (q *Queries) DefaultDomainPolicy(ctx context.Context) (policy *DomainPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareDomainPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		DomainPolicyColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		DomainPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(DomainPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-pM7lP", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		policy, err = scan(row)
		return err
	}, query, args...)
	return policy, err
}

func prepareDomainPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*DomainPolicy, error)) {
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
			From(domainPolicyTable.identifier()).
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
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-K0Jr5", "Errors.DomainPolicy.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-rIy6j", "Errors.Internal")
			}
			return policy, nil
		}
}
