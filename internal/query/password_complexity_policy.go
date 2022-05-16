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

type PasswordComplexityPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.PolicyState

	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool

	IsDefault bool
}

func (q *Queries) PasswordComplexityPolicyByOrg(ctx context.Context, shouldRealTime bool, orgID string) (*PasswordComplexityPolicy, error) {
	if shouldRealTime {
		projection.PasswordComplexityProjection.TriggerBulk(ctx)
	}
	stmt, scan := preparePasswordComplexityPolicyQuery()
	query, args, err := stmt.Where(
		sq.Or{
			sq.Eq{
				PasswordComplexityColID.identifier(): orgID,
			},
			sq.Eq{
				PasswordComplexityColID.identifier(): q.iamID,
			},
		}).
		OrderBy(PasswordComplexityColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-lDnrk", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultPasswordComplexityPolicy(ctx context.Context, shouldRealTime bool) (*PasswordComplexityPolicy, error) {
	if shouldRealTime {
		projection.PasswordComplexityProjection.TriggerBulk(ctx)
	}
	stmt, scan := preparePasswordComplexityPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		PasswordComplexityColID.identifier(): q.iamID,
	}).
		OrderBy(PasswordComplexityColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-h4Uyr", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

var (
	passwordComplexityTable = table{
		name: projection.PasswordComplexityTable,
	}
	PasswordComplexityColID = Column{
		name:  projection.ComplexityPolicyIDCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColSequence = Column{
		name:  projection.ComplexityPolicySequenceCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColCreationDate = Column{
		name:  projection.ComplexityPolicyCreationDateCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColChangeDate = Column{
		name:  projection.ComplexityPolicyChangeDateCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColResourceOwner = Column{
		name:  projection.ComplexityPolicyResourceOwnerCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColMinLength = Column{
		name:  projection.ComplexityPolicyMinLengthCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColHasLowercase = Column{
		name:  projection.ComplexityPolicyHasLowercaseCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColHasUpperCase = Column{
		name:  projection.ComplexityPolicyHasUppercaseCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColHasNumber = Column{
		name:  projection.ComplexityPolicyHasNumberCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColHasSymbol = Column{
		name:  projection.ComplexityPolicyHasSymbolCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColIsDefault = Column{
		name:  projection.ComplexityPolicyIsDefaultCol,
		table: passwordComplexityTable,
	}
	PasswordComplexityColState = Column{
		name:  projection.ComplexityPolicyStateCol,
		table: passwordComplexityTable,
	}
)

func preparePasswordComplexityPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*PasswordComplexityPolicy, error)) {
	return sq.Select(
			PasswordComplexityColID.identifier(),
			PasswordComplexityColSequence.identifier(),
			PasswordComplexityColCreationDate.identifier(),
			PasswordComplexityColChangeDate.identifier(),
			PasswordComplexityColResourceOwner.identifier(),
			PasswordComplexityColMinLength.identifier(),
			PasswordComplexityColHasLowercase.identifier(),
			PasswordComplexityColHasUpperCase.identifier(),
			PasswordComplexityColHasNumber.identifier(),
			PasswordComplexityColHasSymbol.identifier(),
			PasswordComplexityColIsDefault.identifier(),
			PasswordComplexityColState.identifier(),
		).
			From(passwordComplexityTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*PasswordComplexityPolicy, error) {
			policy := new(PasswordComplexityPolicy)
			err := row.Scan(
				&policy.ID,
				&policy.Sequence,
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.ResourceOwner,
				&policy.MinLength,
				&policy.HasLowercase,
				&policy.HasUppercase,
				&policy.HasNumber,
				&policy.HasSymbol,
				&policy.IsDefault,
				&policy.State,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-63mtI", "Errors.PasswordComplexity.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-uulCZ", "Errors.Internal")
			}
			return policy, nil
		}
}
