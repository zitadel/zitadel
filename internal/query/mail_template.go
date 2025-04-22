package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type MailTemplate struct {
	AggregateID  string
	Sequence     uint64
	CreationDate time.Time
	ChangeDate   time.Time
	State        domain.PolicyState

	Template  []byte
	IsDefault bool
}

var (
	mailTemplateTable = table{
		name:          projection.MailTemplateTable,
		instanceIDCol: projection.MailTemplateInstanceIDCol,
	}
	MailTemplateColAggregateID = Column{
		name:  projection.MailTemplateAggregateIDCol,
		table: mailTemplateTable,
	}
	MailTemplateColInstanceID = Column{
		name:  projection.MailTemplateInstanceIDCol,
		table: mailTemplateTable,
	}
	MailTemplateColSequence = Column{
		name:  projection.MailTemplateSequenceCol,
		table: mailTemplateTable,
	}
	MailTemplateColCreationDate = Column{
		name:  projection.MailTemplateCreationDateCol,
		table: mailTemplateTable,
	}
	MailTemplateColChangeDate = Column{
		name:  projection.MailTemplateChangeDateCol,
		table: mailTemplateTable,
	}
	MailTemplateColTemplate = Column{
		name:  projection.MailTemplateTemplateCol,
		table: mailTemplateTable,
	}
	MailTemplateColIsDefault = Column{
		name:  projection.MailTemplateIsDefaultCol,
		table: mailTemplateTable,
	}
	MailTemplateColState = Column{
		name:  projection.MailTemplateStateCol,
		table: mailTemplateTable,
	}
	MailTemplateColOwnerRemoved = Column{
		name:  projection.MailTemplateOwnerRemovedCol,
		table: mailTemplateTable,
	}
)

func (q *Queries) MailTemplateByOrg(ctx context.Context, orgID string, withOwnerRemoved bool) (template *MailTemplate, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareMailTemplateQuery()
	eq := sq.Eq{MailTemplateColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[MailTemplateColOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(
		sq.And{
			eq,
			sq.Or{
				sq.Eq{MailTemplateColAggregateID.identifier(): orgID},
				sq.Eq{MailTemplateColAggregateID.identifier(): authz.GetInstance(ctx).InstanceID()},
			},
		}).
		OrderBy(MailTemplateColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-m0sJg", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		template, err = scan(row)
		return err
	}, query, args...)
	return template, err
}

func (q *Queries) DefaultMailTemplate(ctx context.Context) (template *MailTemplate, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareMailTemplateQuery()
	query, args, err := stmt.Where(sq.Eq{
		MailTemplateColAggregateID.identifier(): authz.GetInstance(ctx).InstanceID(),
		MailTemplateColInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(MailTemplateColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-2m0fH", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		template, err = scan(row)
		return err
	}, query, args...)
	return template, err
}

func prepareMailTemplateQuery() (sq.SelectBuilder, func(*sql.Row) (*MailTemplate, error)) {
	return sq.Select(
			MailTemplateColAggregateID.identifier(),
			MailTemplateColSequence.identifier(),
			MailTemplateColCreationDate.identifier(),
			MailTemplateColChangeDate.identifier(),
			MailTemplateColTemplate.identifier(),
			MailTemplateColIsDefault.identifier(),
			MailTemplateColState.identifier(),
		).
			From(mailTemplateTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*MailTemplate, error) {
			policy := new(MailTemplate)
			err := row.Scan(
				&policy.AggregateID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.Template,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-2NO0g", "Errors.MailTemplate.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-4Nisf", "Errors.Internal")
			}
			return policy, nil
		}
}
