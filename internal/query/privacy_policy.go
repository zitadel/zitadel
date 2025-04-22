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

type PrivacyPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	TOSLink        string
	PrivacyLink    string
	HelpLink       string
	SupportEmail   domain.EmailAddress
	DocsLink       string
	CustomLink     string
	CustomLinkText string

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
	PrivacyColSupportEmail = Column{
		name:  projection.PrivacyPolicySupportEmailCol,
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
	PrivacyColDocsLink = Column{
		name:  projection.PrivacyPolicyDocsLinkCol,
		table: privacyTable,
	}
	PrivacyColCustomLink = Column{
		name:  projection.PrivacyPolicyCustomLinkCol,
		table: privacyTable,
	}
	PrivacyColCustomLinkText = Column{
		name:  projection.PrivacyPolicyCustomLinkTextCol,
		table: privacyTable,
	}
)

func (q *Queries) PrivacyPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (policy *PrivacyPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerPrivacyPolicyProjection")
		ctx, err = projection.PrivacyPolicyProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}
	eq := sq.Eq{PrivacyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[PrivacyColOwnerRemoved.identifier()] = false
	}
	stmt, scan := preparePrivacyPolicyQuery()
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
		return nil, zerrors.ThrowInternal(err, "QUERY-UXuPI", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		policy, err = scan(row)
		return err
	}, query, args...)
	return policy, err
}

func (q *Queries) DefaultPrivacyPolicy(ctx context.Context, shouldTriggerBulk bool) (policy *PrivacyPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerPrivacyPolicyProjection")
		ctx, err = projection.PrivacyPolicyProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := preparePrivacyPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		PrivacyColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		PrivacyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(PrivacyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-LkFZ7", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		policy, err = scan(row)
		return err
	}, query, args...)
	return policy, err
}

func preparePrivacyPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*PrivacyPolicy, error)) {
	return sq.Select(
			PrivacyColID.identifier(),
			PrivacyColSequence.identifier(),
			PrivacyColCreationDate.identifier(),
			PrivacyColChangeDate.identifier(),
			PrivacyColResourceOwner.identifier(),
			PrivacyColPrivacyLink.identifier(),
			PrivacyColTOSLink.identifier(),
			PrivacyColHelpLink.identifier(),
			PrivacyColSupportEmail.identifier(),
			PrivacyColDocsLink.identifier(),
			PrivacyColCustomLink.identifier(),
			PrivacyColCustomLinkText.identifier(),
			PrivacyColIsDefault.identifier(),
			PrivacyColState.identifier(),
		).
			From(privacyTable.identifier()).
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
				&policy.SupportEmail,
				&policy.DocsLink,
				&policy.CustomLink,
				&policy.CustomLinkText,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-vNMHL", "Errors.PrivacyPolicy.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-csrdo", "Errors.Internal")
			}
			return policy, nil
		}
}

func (p *PrivacyPolicy) ToDomain() *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		TOSLink:        p.TOSLink,
		PrivacyLink:    p.PrivacyLink,
		HelpLink:       p.HelpLink,
		SupportEmail:   p.SupportEmail,
		Default:        p.IsDefault,
		DocsLink:       p.DocsLink,
		CustomLink:     p.CustomLink,
		CustomLinkText: p.CustomLinkText,
	}
}
