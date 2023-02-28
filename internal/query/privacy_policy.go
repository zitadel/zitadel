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

type PrivacyPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	TOSLink     string
	PrivacyLink string
	HelpLink    string

	IsDefault bool
}

var (
	privacyTable = table{
		name:          projection.PrivacyPolicyTable,
		instanceIDCol: projection.PrivacyPolicyInstanceIDCol,
	}
	PrivacyColID = Column{
		name:  projection.PrivacyPolicyIDCol,
		table: privacyTable,
	}
	PrivacyColSequence = Column{
		name:  projection.PrivacyPolicySequenceCol,
		table: privacyTable,
	}
	PrivacyColCreationDate = Column{
		name:  projection.PrivacyPolicyCreationDateCol,
		table: privacyTable,
	}
	PrivacyColChangeDate = Column{
		name:  projection.PrivacyPolicyChangeDateCol,
		table: privacyTable,
	}
	PrivacyColResourceOwner = Column{
		name:  projection.PrivacyPolicyResourceOwnerCol,
		table: privacyTable,
	}
	PrivacyColInstanceID = Column{
		name:  projection.PrivacyPolicyInstanceIDCol,
		table: privacyTable,
	}
	PrivacyColPrivacyLink = Column{
		name:  projection.PrivacyPolicyPrivacyLinkCol,
		table: privacyTable,
	}
	PrivacyColTOSLink = Column{
		name:  projection.PrivacyPolicyTOSLinkCol,
		table: privacyTable,
	}
	PrivacyColHelpLink = Column{
		name:  projection.PrivacyPolicyHelpLinkCol,
		table: privacyTable,
	}
	PrivacyColIsDefault = Column{
		name:  projection.PrivacyPolicyIsDefaultCol,
		table: privacyTable,
	}
	PrivacyColState = Column{
		name:  projection.PrivacyPolicyStateCol,
		table: privacyTable,
	}
	PrivacyColOwnerRemoved = Column{
		name:  projection.PrivacyPolicyOwnerRemovedCol,
		table: privacyTable,
	}
)

func (q *Queries) PrivacyPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (_ *PrivacyPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.PrivacyPolicyProjection.Trigger(ctx)
	}
	eq := sq.Eq{PrivacyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[PrivacyColOwnerRemoved.identifier()] = false
	}
	stmt, scan := preparePrivacyPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(
		sq.And{
			eq,
			sq.Or{
				sq.Eq{PrivacyColID.identifier(): orgID},
				sq.Eq{PrivacyColID.identifier(): authz.GetInstance(ctx).InstanceID()},
			},
		}).
		OrderBy(PrivacyColIsDefault.identifier()).Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-UXuPI", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultPrivacyPolicy(ctx context.Context, shouldTriggerBulk bool) (_ *PrivacyPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.PrivacyPolicyProjection.Trigger(ctx)
	}

	stmt, scan := preparePrivacyPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		PrivacyColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		PrivacyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(PrivacyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-LkFZ7", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func preparePrivacyPolicyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*PrivacyPolicy, error)) {
	return sq.Select(
			PrivacyColID.identifier(),
			PrivacyColSequence.identifier(),
			PrivacyColCreationDate.identifier(),
			PrivacyColChangeDate.identifier(),
			PrivacyColResourceOwner.identifier(),
			PrivacyColPrivacyLink.identifier(),
			PrivacyColTOSLink.identifier(),
			PrivacyColHelpLink.identifier(),
			PrivacyColIsDefault.identifier(),
			PrivacyColState.identifier(),
		).
			From(privacyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*PrivacyPolicy, error) {
			policy := new(PrivacyPolicy)
			err := row.Scan(
				&policy.ID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.ResourceOwner,
				&policy.PrivacyLink,
				&policy.TOSLink,
				&policy.HelpLink,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-vNMHL", "Errors.PrivacyPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-csrdo", "Errors.Internal")
			}
			return policy, nil
		}
}

func (p *PrivacyPolicy) ToDomain() *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		TOSLink:     p.TOSLink,
		PrivacyLink: p.PrivacyLink,
		HelpLink:    p.HelpLink,
		Default:     p.IsDefault,
	}
}
