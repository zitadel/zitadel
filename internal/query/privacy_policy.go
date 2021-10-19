package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
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

	IsDefault bool
}

func (q *Queries) MyPrivacyPolicy(ctx context.Context, orgID string) (*PrivacyPolicy, error) {
	stmt, scan := preparePrivacyPolicyQuery()
	query, args, err := stmt.Where(
		sq.Or{
			sq.Eq{
				PrivacyColID.identifier(): orgID,
			},
			sq.Eq{
				PrivacyColID.identifier(): q.iamID,
			},
		}).
		OrderBy(PrivacyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-UXuPI", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultPrivacyPolicy(ctx context.Context) (*PrivacyPolicy, error) {
	stmt, scan := preparePrivacyPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		PrivacyColID.identifier(): q.iamID,
	}).
		OrderBy(PrivacyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-LkFZ7", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

var (
	privacyTable = table{
		name: projection.PrivacyPolicyTable,
	}
	PrivacyColID = Column{
		name: projection.PrivacyPolicyIDCol,
	}
	PrivacyColSequence = Column{
		name: projection.PrivacyPolicySequenceCol,
	}
	PrivacyColCreationDate = Column{
		name: projection.PrivacyPolicyCreationDateCol,
	}
	PrivacyColChangeDate = Column{
		name: projection.PrivacyPolicyChangeDateCol,
	}
	PrivacyColResourceOwner = Column{
		name: projection.PrivacyPolicyResourceOwnerCol,
	}
	PrivacyColPrivacyLink = Column{
		name: projection.PrivacyPolicyPrivacyLinkCol,
	}
	PrivacyColTOSLink = Column{
		name: projection.PrivacyPolicyTOSLinkCol,
	}
	PrivacyColIsDefault = Column{
		name: projection.PrivacyPolicyIsDefaultCol,
	}
	PrivacyColState = Column{
		name: projection.PrivacyPolicyStateCol,
	}
)

func preparePrivacyPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*PrivacyPolicy, error)) {
	return sq.Select(
			PrivacyColID.identifier(),
			PrivacyColSequence.identifier(),
			PrivacyColCreationDate.identifier(),
			PrivacyColChangeDate.identifier(),
			PrivacyColResourceOwner.identifier(),
			PrivacyColPrivacyLink.identifier(),
			PrivacyColTOSLink.identifier(),
			PrivacyColIsDefault.identifier(),
			PrivacyColState.identifier(),
		).
			From(privacyTable.identifier()).PlaceholderFormat(sq.Dollar),
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
