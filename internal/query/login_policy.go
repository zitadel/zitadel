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
	"github.com/lib/pq"
)

type LoginPolicy struct {
	OrgID                      string
	CreationDate               time.Time
	ChangeDate                 time.Time
	Sequence                   uint64
	AllowRegister              bool
	AllowUsernamePassword      bool
	AllowExternalIDPs          bool
	ForceMFA                   bool
	SecondFactors              []domain.SecondFactorType
	MultiFactors               []domain.MultiFactorType
	PasswordlessType           domain.PasswordlessType
	IsDefault                  bool
	HidePasswordReset          bool
	PasswordCheckLifetime      time.Duration
	ExternalLoginCheckLifetime time.Duration
	MFAInitSkipLifetime        time.Duration
	SecondFactorCheckLifetime  time.Duration
	MultiFactorCheckLifetime   time.Duration
}

type SecondFactors struct {
	SearchResponse
	Factors []domain.SecondFactorType
}

type MultiFactors struct {
	SearchResponse
	Factors []domain.MultiFactorType
}

var (
	loginPolicyTable = table{
		name: projection.LoginPolicyTable,
	}
	LoginPolicyColumnOrgID = Column{
		name:  projection.LoginPolicyIDCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnCreationDate = Column{
		name:  projection.LoginPolicyCreationDateCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnChangeDate = Column{
		name:  projection.LoginPolicyChangeDateCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnSequence = Column{
		name:  projection.LoginPolicySequenceCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnAllowRegister = Column{
		name:  projection.LoginPolicyAllowRegisterCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnAllowUsernamePassword = Column{
		name:  projection.LoginPolicyAllowUsernamePasswordCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnAllowExternalIDPs = Column{
		name:  projection.LoginPolicyAllowExternalIDPsCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnForceMFA = Column{
		name:  projection.LoginPolicyForceMFACol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnSecondFactors = Column{
		name:  projection.LoginPolicy2FAsCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnMultiFactors = Column{
		name:  projection.LoginPolicyMFAsCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnPasswordlessType = Column{
		name:  projection.LoginPolicyPasswordlessTypeCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnIsDefault = Column{
		name:  projection.LoginPolicyIsDefaultCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnHidePasswordReset = Column{
		name:  projection.LoginPolicyHidePWResetCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnPasswordCheckLifetime = Column{
		name:  projection.PasswordCheckLifetimeCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnExternalLoginCheckLifetime = Column{
		name:  projection.ExternalLoginCheckLifetimeCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnMFAInitSkipLifetime = Column{
		name:  projection.MFAInitSkipLifetimeCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnSecondFactorCheckLifetime = Column{
		name:  projection.SecondFactorCheckLifetimeCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnMultiFacotrCheckLifetime = Column{
		name:  projection.MultiFactorCheckLifetimeCol,
		table: loginPolicyTable,
	}
)

func (q *Queries) LoginPolicyByID(ctx context.Context, orgID string) (*LoginPolicy, error) {
	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(
		sq.Or{
			sq.Eq{
				LoginPolicyColumnOrgID.identifier(): orgID,
			},
			sq.Eq{
				LoginPolicyColumnOrgID.identifier(): domain.IAMID,
			},
		}).
		OrderBy(LoginPolicyColumnIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-scVHo", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) DefaultLoginPolicy(ctx context.Context) (*LoginPolicy, error) {
	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.identifier(): domain.IAMID,
	}).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-t4TBK", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) SecondFactorsByOrg(ctx context.Context, orgID string) (*SecondFactors, error) {
	query, scan := prepareLoginPolicy2FAsQuery()
	stmt, args, err := query.Where(
		sq.Or{
			sq.Eq{
				LoginPolicyColumnOrgID.identifier(): orgID,
			},
			sq.Eq{
				LoginPolicyColumnOrgID.identifier(): domain.IAMID,
			},
		}).
		OrderBy(LoginPolicyColumnIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-scVHo", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	factors, err := scan(row)
	if err != nil {
		return nil, err
	}
	factors.LatestSequence, err = q.latestSequence(ctx, loginPolicyTable)
	return factors, err
}

func (q *Queries) DefaultSecondFactors(ctx context.Context) (*SecondFactors, error) {
	query, scan := prepareLoginPolicy2FAsQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.identifier(): domain.IAMID,
	}).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-CZ2Nv", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	factors, err := scan(row)
	if err != nil {
		return nil, err
	}
	factors.LatestSequence, err = q.latestSequence(ctx, loginPolicyTable)
	return factors, err
}

func (q *Queries) MultiFactorsByOrg(ctx context.Context, orgID string) (*MultiFactors, error) {
	query, scan := prepareLoginPolicyMFAsQuery()
	stmt, args, err := query.Where(
		sq.Or{
			sq.Eq{
				LoginPolicyColumnOrgID.identifier(): orgID,
			},
			sq.Eq{
				LoginPolicyColumnOrgID.identifier(): domain.IAMID,
			},
		}).
		OrderBy(LoginPolicyColumnIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-B4o7h", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	factors, err := scan(row)
	if err != nil {
		return nil, err
	}
	factors.LatestSequence, err = q.latestSequence(ctx, loginPolicyTable)
	return factors, err
}

func (q *Queries) DefaultMultiFactors(ctx context.Context) (*MultiFactors, error) {
	query, scan := prepareLoginPolicyMFAsQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.identifier(): domain.IAMID,
	}).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-WxYjr", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	factors, err := scan(row)
	if err != nil {
		return nil, err
	}
	factors.LatestSequence, err = q.latestSequence(ctx, loginPolicyTable)
	return factors, err
}

func prepareLoginPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*LoginPolicy, error)) {
	return sq.Select(
			LoginPolicyColumnOrgID.identifier(),
			LoginPolicyColumnCreationDate.identifier(),
			LoginPolicyColumnChangeDate.identifier(),
			LoginPolicyColumnSequence.identifier(),
			LoginPolicyColumnAllowRegister.identifier(),
			LoginPolicyColumnAllowUsernamePassword.identifier(),
			LoginPolicyColumnAllowExternalIDPs.identifier(),
			LoginPolicyColumnForceMFA.identifier(),
			LoginPolicyColumnSecondFactors.identifier(),
			LoginPolicyColumnMultiFactors.identifier(),
			LoginPolicyColumnPasswordlessType.identifier(),
			LoginPolicyColumnIsDefault.identifier(),
			LoginPolicyColumnHidePasswordReset.identifier(),
			LoginPolicyColumnPasswordCheckLifetime.identifier(),
			LoginPolicyColumnExternalLoginCheckLifetime.identifier(),
			LoginPolicyColumnMFAInitSkipLifetime.identifier(),
			LoginPolicyColumnSecondFactorCheckLifetime.identifier(),
			LoginPolicyColumnMultiFacotrCheckLifetime.identifier(),
		).From(loginPolicyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LoginPolicy, error) {
			p := new(LoginPolicy)
			secondFactors := pq.Int32Array{}
			multiFactors := pq.Int32Array{}
			err := row.Scan(
				&p.OrgID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.Sequence,
				&p.AllowRegister,
				&p.AllowUsernamePassword,
				&p.AllowExternalIDPs,
				&p.ForceMFA,
				&secondFactors,
				&multiFactors,
				&p.PasswordlessType,
				&p.IsDefault,
				&p.HidePasswordReset,
				&p.PasswordCheckLifetime,
				&p.ExternalLoginCheckLifetime,
				&p.MFAInitSkipLifetime,
				&p.SecondFactorCheckLifetime,
				&p.MultiFactorCheckLifetime,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-QsUBJ", "Errors.LoginPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-YcC53", "Errors.Internal")
			}

			p.MultiFactors = make([]domain.MultiFactorType, len(multiFactors))
			for i, mfa := range multiFactors {
				p.MultiFactors[i] = domain.MultiFactorType(mfa)
			}
			p.SecondFactors = make([]domain.SecondFactorType, len(secondFactors))
			for i, mfa := range secondFactors {
				p.SecondFactors[i] = domain.SecondFactorType(mfa)
			}
			return p, nil
		}
}

func prepareLoginPolicy2FAsQuery() (sq.SelectBuilder, func(*sql.Row) (*SecondFactors, error)) {
	return sq.Select(
			LoginPolicyColumnSecondFactors.identifier(),
		).From(loginPolicyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SecondFactors, error) {
			p := new(SecondFactors)
			secondFactors := pq.Int32Array{}
			err := row.Scan(
				&secondFactors,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-yPqIZ", "Errors.LoginPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Mr6H3", "Errors.Internal")
			}

			p.Factors = make([]domain.SecondFactorType, len(secondFactors))
			p.Count = uint64(len(secondFactors))
			for i, mfa := range secondFactors {
				p.Factors[i] = domain.SecondFactorType(mfa)
			}
			return p, nil
		}
}

func prepareLoginPolicyMFAsQuery() (sq.SelectBuilder, func(*sql.Row) (*MultiFactors, error)) {
	return sq.Select(
			LoginPolicyColumnMultiFactors.identifier(),
		).From(loginPolicyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*MultiFactors, error) {
			p := new(MultiFactors)
			multiFactors := pq.Int32Array{}
			err := row.Scan(
				&multiFactors,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-yPqIZ", "Errors.LoginPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Mr6H3", "Errors.Internal")
			}

			p.Factors = make([]domain.MultiFactorType, len(multiFactors))
			p.Count = uint64(len(multiFactors))
			for i, mfa := range multiFactors {
				p.Factors[i] = domain.MultiFactorType(mfa)
			}
			return p, nil
		}
}
