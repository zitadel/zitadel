package query

import (
	"context"
	"database/sql"
	errs "errors"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/constants"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/store"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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

func (q *Queries) MailTemplateByOrg(ctx context.Context, orgID string, withOwnerRemoved bool) (_ *MailTemplate, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareMailTemplateQuery(ctx, q.client, orgID)
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
		return nil, errors.ThrowInternal(err, "QUERY-m0sJg", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultMailTemplate(ctx context.Context) (_ *MailTemplate, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareMailTemplateQuery(ctx, q.client, "")
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

func prepareMailTemplateQuery(ctx context.Context, db prepareDatabase, orgID string) (sq.SelectBuilder, func(*sql.Row) (*MailTemplate, error)) {
	constants.InitConstants()
	return sq.Select(
			MailTemplateColAggregateID.identifier(),
			MailTemplateColSequence.identifier(),
			MailTemplateColCreationDate.identifier(),
			MailTemplateColChangeDate.identifier(),
			MailTemplateColTemplate.identifier(),
			MailTemplateColIsDefault.identifier(),
			MailTemplateColState.identifier(),
		).
			From(mailTemplateTable.identifier() + db.Timetravel(call.Took(ctx))).
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
			template, success := downloadTemplate(orgID)
			if success {
				policy.Template = template
			}
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-2NO0g", "Errors.MailTemplate.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-4Nisf", "Errors.Internal")
			}
			return policy, nil
		}
}

func downloadTemplate(orgId string) ([]byte, bool) {
	if len(orgId) <= 0 {
		return []byte{}, false
	}
	client, err := store.Setup(orgId)
	defer client.Close()

	if err != nil {
		return []byte{}, false
	}

	f, err := os.CreateTemp(".", "email")
	if err != nil {

		return []byte{}, false
	}

	tmpFileName := f.Name()

	defer os.Remove(tmpFileName)

	err = store.DownloadCaller(f, "init.html")
	if err != nil {
		_ = f.Close()
		return []byte{}, false
	}
	defer f.Close()

	buf, err := os.ReadFile(tmpFileName)

	if err != nil {
		_ = f.Close()
		return []byte{}, false
	}
	return buf, true
}
