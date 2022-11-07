package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
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
	SecondFactors              database.EnumArray[domain.SecondFactorType]
	MultiFactors               database.EnumArray[domain.MultiFactorType]
	PasswordlessType           domain.PasswordlessType
	IsDefault                  bool
	HidePasswordReset          bool
	IgnoreUnknownUsernames     bool
	AllowDomainDiscovery       bool
	DisableLoginWithEmail      bool
	DisableLoginWithPhone      bool
	DefaultRedirectURI         string
	PasswordCheckLifetime      time.Duration
	ExternalLoginCheckLifetime time.Duration
	MFAInitSkipLifetime        time.Duration
	SecondFactorCheckLifetime  time.Duration
	MultiFactorCheckLifetime   time.Duration
	IDPLinks                   []*IDPLoginPolicyLink
}

type SecondFactors struct {
	SearchResponse
	Factors database.EnumArray[domain.SecondFactorType]
}

type MultiFactors struct {
	SearchResponse
	Factors database.EnumArray[domain.MultiFactorType]
}

var (
	loginPolicyTable = table{
		name:          projection.LoginPolicyTable,
		instanceIDCol: projection.LoginPolicyInstanceIDCol,
	}
	LoginPolicyColumnOrgID = Column{
		name:  projection.LoginPolicyIDCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnInstanceID = Column{
		name:  projection.LoginPolicyInstanceIDCol,
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
	LoginPolicyColumnIgnoreUnknownUsernames = Column{
		name:  projection.IgnoreUnknownUsernames,
		table: loginPolicyTable,
	}
	LoginPolicyColumnAllowDomainDiscovery = Column{
		name:  projection.AllowDomainDiscovery,
		table: loginPolicyTable,
	}
	LoginPolicyColumnDisableLoginWithEmail = Column{
		name:  projection.DisableLoginWithEmail,
		table: loginPolicyTable,
	}
	LoginPolicyColumnDisableLoginWithPhone = Column{
		name:  projection.DisableLoginWithPhone,
		table: loginPolicyTable,
	}
	LoginPolicyColumnDefaultRedirectURI = Column{
		name:  projection.DefaultRedirectURI,
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

func (q *Queries) LoginPolicyByID(ctx context.Context, shouldTriggerBulk bool, orgID string) (*LoginPolicy, error) {
	if shouldTriggerBulk {
		projection.LoginPolicyProjection.Trigger(ctx)
	}

	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{
				LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			},
			sq.Or{
				sq.Eq{
					LoginPolicyColumnOrgID.identifier(): orgID,
				},
				sq.Eq{
					LoginPolicyColumnOrgID.identifier(): authz.GetInstance(ctx).InstanceID(),
				},
			},
		}).OrderBy(LoginPolicyColumnIsDefault.identifier() + " DESC").ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-scVHo", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SWgr3", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) DefaultLoginPolicy(ctx context.Context) (*LoginPolicy, error) {
	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.identifier():      authz.GetInstance(ctx).InstanceID(),
		LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-t4TBK", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SArt2", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) SecondFactorsByOrg(ctx context.Context, orgID string) (*SecondFactors, error) {
	query, scan := prepareLoginPolicy2FAsQuery()
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{
				LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			},
			sq.Or{
				sq.Eq{
					LoginPolicyColumnOrgID.identifier(): orgID,
				},
				sq.Eq{
					LoginPolicyColumnOrgID.identifier(): authz.GetInstance(ctx).InstanceID(),
				},
			},
		}).
		OrderBy(LoginPolicyColumnIsDefault.identifier() + " DESC").
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
		LoginPolicyColumnOrgID.identifier():      authz.GetInstance(ctx).InstanceID(),
		LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
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
		sq.And{
			sq.Eq{
				LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			},
			sq.Or{
				sq.Eq{
					LoginPolicyColumnOrgID.identifier(): orgID,
				},
				sq.Eq{
					LoginPolicyColumnOrgID.identifier(): authz.GetInstance(ctx).InstanceID(),
				},
			},
		}).
		OrderBy(LoginPolicyColumnIsDefault.identifier() + " DESC").
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
		LoginPolicyColumnOrgID.identifier():      authz.GetInstance(ctx).InstanceID(),
		LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
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

func prepareLoginPolicyQuery() (sq.SelectBuilder, func(*sql.Rows) (*LoginPolicy, error)) {
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
			LoginPolicyColumnIgnoreUnknownUsernames.identifier(),
			LoginPolicyColumnAllowDomainDiscovery.identifier(),
			LoginPolicyColumnDisableLoginWithEmail.identifier(),
			LoginPolicyColumnDisableLoginWithPhone.identifier(),
			LoginPolicyColumnDefaultRedirectURI.identifier(),
			LoginPolicyColumnPasswordCheckLifetime.identifier(),
			LoginPolicyColumnExternalLoginCheckLifetime.identifier(),
			LoginPolicyColumnMFAInitSkipLifetime.identifier(),
			LoginPolicyColumnSecondFactorCheckLifetime.identifier(),
			LoginPolicyColumnMultiFacotrCheckLifetime.identifier(),
			IDPLoginPolicyLinkIDPIDCol.identifier(),
			IDPNameCol.identifier(),
			IDPTypeCol.identifier(),
		).From(loginPolicyTable.identifier()).
			LeftJoin(join(IDPLoginPolicyLinkAggregateIDCol, LoginPolicyColumnOrgID)).
			LeftJoin(join(IDPIDCol, IDPLoginPolicyLinkIDPIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*LoginPolicy, error) {
			p := new(LoginPolicy)
			defaultRedirectURI := sql.NullString{}
			links := make([]*IDPLoginPolicyLink, 0)
			for rows.Next() {
				var (
					idpID   = sql.NullString{}
					idpName = sql.NullString{}
					idpType = sql.NullInt16{}
				)
				err := rows.Scan(
					&p.OrgID,
					&p.CreationDate,
					&p.ChangeDate,
					&p.Sequence,
					&p.AllowRegister,
					&p.AllowUsernamePassword,
					&p.AllowExternalIDPs,
					&p.ForceMFA,
					&p.SecondFactors,
					&p.MultiFactors,
					&p.PasswordlessType,
					&p.IsDefault,
					&p.HidePasswordReset,
					&p.IgnoreUnknownUsernames,
					&p.AllowDomainDiscovery,
					&p.DisableLoginWithEmail,
					&p.DisableLoginWithPhone,
					&defaultRedirectURI,
					&p.PasswordCheckLifetime,
					&p.ExternalLoginCheckLifetime,
					&p.MFAInitSkipLifetime,
					&p.SecondFactorCheckLifetime,
					&p.MultiFactorCheckLifetime,
					&idpID,
					&idpName,
					&idpType,
				)
				if err != nil {
					return nil, errors.ThrowInternal(err, "QUERY-YcC53", "Errors.Internal")
				}
				var link IDPLoginPolicyLink
				if idpID.Valid {
					link = IDPLoginPolicyLink{IDPID: idpID.String}

					link.IDPName = idpName.String
					//IDPType 0 is oidc so we have to set unspecified manually
					if idpType.Valid {
						link.IDPType = domain.IDPConfigType(idpType.Int16)
					} else {
						link.IDPType = domain.IDPConfigTypeUnspecified
					}
					links = append(links, &link)
				}
			}
			if p.OrgID == "" {
				return nil, errors.ThrowNotFound(nil, "QUERY-QsUBJ", "Errors.LoginPolicy.NotFound")
			}
			p.DefaultRedirectURI = defaultRedirectURI.String
			p.IDPLinks = links
			return p, nil
		}
}

func prepareLoginPolicy2FAsQuery() (sq.SelectBuilder, func(*sql.Row) (*SecondFactors, error)) {
	return sq.Select(
			LoginPolicyColumnSecondFactors.identifier(),
		).From(loginPolicyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SecondFactors, error) {
			p := new(SecondFactors)
			err := row.Scan(
				&p.Factors,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-yPqIZ", "Errors.LoginPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Mr6H3", "Errors.Internal")
			}

			p.Count = uint64(len(p.Factors))
			return p, nil
		}
}

func prepareLoginPolicyMFAsQuery() (sq.SelectBuilder, func(*sql.Row) (*MultiFactors, error)) {
	return sq.Select(
			LoginPolicyColumnMultiFactors.identifier(),
		).From(loginPolicyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*MultiFactors, error) {
			p := new(MultiFactors)
			err := row.Scan(
				&p.Factors,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-yPqIZ", "Errors.LoginPolicy.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-Mr6H3", "Errors.Internal")
			}

			p.Count = uint64(len(p.Factors))
			return p, nil
		}
}
