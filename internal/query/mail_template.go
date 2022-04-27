package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
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
		name: projection.MailTemplateTable,
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
)

func (q *Queries) MailTemplateByOrg(ctx context.Context, orgID string) (*MailTemplate, error) {
	stmt, scan := prepareMailTemplateQuery()
	query, args, err := stmt.Where(
		sq.And{
			sq.Eq{
				MailTemplateColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			},
			sq.Or{
				sq.Eq{
					MailTemplateColAggregateID.identifier(): orgID,
				},
				sq.Eq{
					MailTemplateColAggregateID.identifier(): authz.GetInstance(ctx).InstanceID(),
				},
			},
		}).
		OrderBy(MailTemplateColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-m0sJg", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultMailTemplate(ctx context.Context) (*MailTemplate, error) {
	stmt, scan := prepareMailTemplateQuery()
	query, args, err := stmt.Where(sq.Eq{
		MailTemplateColAggregateID.identifier(): authz.GetInstance(ctx).InstanceID(),
		MailTemplateColInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(MailTemplateColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-2m0fH", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
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
			From(mailTemplateTable.identifier()).PlaceholderFormat(sq.Dollar),
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
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-2NO0g", "Errors.MailTemplate.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-4Nisf", "Errors.Internal")
			}
			return policy, nil
		}
}
