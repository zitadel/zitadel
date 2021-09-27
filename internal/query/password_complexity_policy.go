package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type PasswordComplexityPolicy struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string

	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool

	IsDefault bool
}

func (q *Queries) MyPasswordComplexityPolicy(ctx context.Context, orgID string) (*PasswordComplexityPolicy, error) {
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
		return nil, errors.ThrowInternal(err, "QUERY-lDnrk", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultPasswordComplexityPolicy(ctx context.Context) (*PasswordComplexityPolicy, error) {
	stmt, scan := preparePasswordComplexityPolicyQuery()
	query, args, err := stmt.Where(sq.Eq{
		PasswordComplexityColID.identifier(): q.iamID,
	}).
		OrderBy(PasswordComplexityColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-h4Uyr", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

var (
	passwordComplexityTable = table{
		name: projection.PasswordComplexityTable,
	}
	PasswordComplexityColID = Column{
		name: projection.ComplexityPolicyIDCol,
	}
	PasswordComplexityColSequence = Column{
		name: projection.ComplexityPolicySequenceCol,
	}
	PasswordComplexityColCreationDate = Column{
		name: projection.ComplexityPolicyCreationDateCol,
	}
	PasswordComplexityColChangeDate = Column{
		name: projection.ComplexityPolicyChangeDateCol,
	}
	PasswordComplexityColResourceOwner = Column{
		name: projection.ComplexityPolicyResourceOwnerCol,
	}
	PasswordComplexityColMinLength = Column{
		name: projection.ComplexityPolicyMinLengthCol,
	}
	PasswordComplexityColHasLowercase = Column{
		name: projection.ComplexityPolicyHasLowercaseCol,
	}
	PasswordComplexityColHasUpperCase = Column{
		name: projection.ComplexityPolicyHasUppercaseCol,
	}
	PasswordComplexityColHasNumber = Column{
		name: projection.ComplexityPolicyHasNumberCol,
	}
	PasswordComplexityColHasSymbol = Column{
		name: projection.ComplexityPolicyHasSymbolCol,
	}
	PasswordComplexityColIsDefault = Column{
		name: projection.ComplexityPolicyIsDefaultCol,
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
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-63mtI", "errors.policy.password.complexity.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-uulCZ", "errors.internal")
			}
			return policy, nil
		}
}
