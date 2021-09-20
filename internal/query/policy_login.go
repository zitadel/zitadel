package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/lib/pq"
)

const (
	loginPoliciesTable = "zitadel.projections.login_policies"
)

type LoginPolicyColumn int8

type LoginPolicy struct {
	OrgID                 string
	CreationDate          time.Time
	ChangeDate            time.Time
	Sequence              uint64
	AllowRegister         bool
	AllowUsernamePassword bool
	AllowExternalIDPs     bool
	ForceMFA              bool
	SecondFactors         []domain.SecondFactorType
	MultiFactors          []domain.MultiFactorType
	PasswordlessType      domain.PasswordlessType
	IsDefault             bool
	HidePasswordReset     bool
	UserLoginMustBeDomain bool
}

const (
	LoginPolicyColumnOrgID LoginPolicyColumn = iota + 1
	LoginPolicyColumnCreationDate
	LoginPolicyColumnChangeDate
	LoginPolicyColumnSequence
	LoginPolicyColumnAllowRegister
	LoginPolicyColumnAllowUsernamePassword
	LoginPolicyColumnAllowExternalIDPs
	LoginPolicyColumnForceMFA
	LoginPolicyColumnSecondFactors
	LoginPolicyColumnMultiFactors
	LoginPolicyColumnPasswordlessType
	LoginPolicyColumnIsDefault
	LoginPolicyColumnHidePasswordReset
	LoginPolicyColumnUserLoginMustBeDomain
)

func (q *Queries) LoginPolicyByID(ctx context.Context, orgID string) (*LoginPolicy, error) {
	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(
		sq.Or{
			sq.Eq{
				LoginPolicyColumnOrgID.toColumnName(): orgID,
			},
			sq.Eq{
				LoginPolicyColumnOrgID.toColumnName(): domain.IAMID,
			},
		}).
		OrderBy(LoginPolicyColumnIsDefault.toColumnName()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-scVHo", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) DefaultLoginPolicy(ctx context.Context) (*LoginPolicy, error) {
	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.toColumnName(): domain.IAMID,
	}).OrderBy(LoginPolicyColumnIsDefault.toColumnName()).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-t4TBK", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func prepareLoginPolicyQuery() (sq.SelectBuilder, func(*sql.Row) (*LoginPolicy, error)) {
	return sq.Select(
			LoginPolicyColumnOrgID.toColumnName(),
			LoginPolicyColumnCreationDate.toColumnName(),
			LoginPolicyColumnChangeDate.toColumnName(),
			LoginPolicyColumnSequence.toColumnName(),
			LoginPolicyColumnAllowRegister.toColumnName(),
			LoginPolicyColumnAllowUsernamePassword.toColumnName(),
			LoginPolicyColumnAllowExternalIDPs.toColumnName(),
			LoginPolicyColumnForceMFA.toColumnName(),
			LoginPolicyColumnSecondFactors.toColumnName(),
			LoginPolicyColumnMultiFactors.toColumnName(),
			LoginPolicyColumnPasswordlessType.toColumnName(),
			LoginPolicyColumnIsDefault.toColumnName(),
			LoginPolicyColumnHidePasswordReset.toColumnName(),
			LoginPolicyColumnUserLoginMustBeDomain.toColumnName(),
		).From(loginPoliciesTable).PlaceholderFormat(sq.Dollar),
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
				&p.UserLoginMustBeDomain,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-QsUBJ", "errors.orgs.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-YcC53", "errors.internal")
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

func (c LoginPolicyColumn) toColumnName() string {
	switch c {
	case LoginPolicyColumnOrgID:
		return "aggregate_id"
	case LoginPolicyColumnCreationDate:
		return "creation_date"
	case LoginPolicyColumnChangeDate:
		return "change_date"
	case LoginPolicyColumnSequence:
		return "sequence"
	case LoginPolicyColumnAllowRegister:
		return "allow_register"
	case LoginPolicyColumnAllowUsernamePassword:
		return "allow_username_password"
	case LoginPolicyColumnAllowExternalIDPs:
		return "allow_external_idps"
	case LoginPolicyColumnForceMFA:
		return "force_mfa"
	case LoginPolicyColumnSecondFactors:
		return "second_factors"
	case LoginPolicyColumnMultiFactors:
		return "multi_factors"
	case LoginPolicyColumnPasswordlessType:
		return "passwordless_type"
	case LoginPolicyColumnIsDefault:
		return "is_default"
	case LoginPolicyColumnHidePasswordReset:
		return "hide_password_reset"
	case LoginPolicyColumnUserLoginMustBeDomain:
		return "user_login_must_be_domain"
	default:
		return ""
	}
}
