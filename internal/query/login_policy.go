package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	ForceMFALocalOnly          bool
	SecondFactors              database.NumberArray[domain.SecondFactorType]
	MultiFactors               database.NumberArray[domain.MultiFactorType]
	PasswordlessType           domain.PasswordlessType
	IsDefault                  bool
	HidePasswordReset          bool
	IgnoreUnknownUsernames     bool
	AllowDomainDiscovery       bool
	DisableLoginWithEmail      bool
	DisableLoginWithPhone      bool
	DefaultRedirectURI         string
	PasswordCheckLifetime      database.Duration
	ExternalLoginCheckLifetime database.Duration
	MFAInitSkipLifetime        database.Duration
	SecondFactorCheckLifetime  database.Duration
	MultiFactorCheckLifetime   database.Duration
	IDPLinks                   []*IDPLoginPolicyLink
}

type SecondFactors struct {
	SearchResponse
	Factors database.NumberArray[domain.SecondFactorType]
}

type MultiFactors struct {
	SearchResponse
	Factors database.NumberArray[domain.MultiFactorType]
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
	LoginPolicyColumnForceMFALocalOnly = Column{
		name:  projection.LoginPolicyForceMFALocalOnlyCol,
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
	LoginPolicyColumnMultiFactorCheckLifetime = Column{
		name:  projection.MultiFactorCheckLifetimeCol,
		table: loginPolicyTable,
	}
	LoginPolicyColumnOwnerRemoved = Column{
		name:  projection.LoginPolicyOwnerRemovedCol,
		table: loginPolicyTable,
	}
)

func (q *Queries) LoginPolicyByID(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (policy *LoginPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerLoginPolicyProjection")
		ctx, err = projection.LoginPolicyProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}
	eq := sq.Eq{LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[LoginPolicyColumnOwnerRemoved.identifier()] = false
	}

	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(
		sq.And{
			eq,
			sq.Or{
				sq.Eq{LoginPolicyColumnOrgID.identifier(): orgID},
				sq.Eq{LoginPolicyColumnOrgID.identifier(): authz.GetInstance(ctx).InstanceID()},
			},
		}).Limit(1).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-scVHo", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		policy, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SWgr3", "Errors.Internal")
	}
	return policy, q.addLinksToLoginPolicy(ctx, policy)
}

func (q *Queries) addLinksToLoginPolicy(ctx context.Context, policy *LoginPolicy) error {
	links, err := q.IDPLoginPolicyLinks(ctx, policy.OrgID, &IDPLoginPolicyLinksSearchQuery{}, false)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-aa4Ve", "Errors.Internal")
	}
	policy.IDPLinks = append(policy.IDPLinks, links.Links...)
	return nil
}

func (q *Queries) DefaultLoginPolicy(ctx context.Context) (policy *LoginPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareLoginPolicyQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.identifier():      authz.GetInstance(ctx).InstanceID(),
		LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-t4TBK", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		policy, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SArt2", "Errors.Internal")
	}
	return policy, q.addLinksToLoginPolicy(ctx, policy)
}

func (q *Queries) SecondFactorsByOrg(ctx context.Context, orgID string) (factors *SecondFactors, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

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
		OrderBy(LoginPolicyColumnIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-scVHo", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		factors, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	factors.State, err = q.latestState(ctx, loginPolicyTable)
	return factors, err
}

func (q *Queries) DefaultSecondFactors(ctx context.Context) (factors *SecondFactors, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareLoginPolicy2FAsQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.identifier():      authz.GetInstance(ctx).InstanceID(),
		LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-CZ2Nv", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		factors, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	factors.State, err = q.latestState(ctx, loginPolicyTable)
	return factors, err
}

func (q *Queries) MultiFactorsByOrg(ctx context.Context, orgID string) (factors *MultiFactors, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

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
		OrderBy(LoginPolicyColumnIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-B4o7h", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		factors, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	factors.State, err = q.latestState(ctx, loginPolicyTable)
	return factors, err
}

func (q *Queries) DefaultMultiFactors(ctx context.Context) (factors *MultiFactors, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareLoginPolicyMFAsQuery()
	stmt, args, err := query.Where(sq.Eq{
		LoginPolicyColumnOrgID.identifier():      authz.GetInstance(ctx).InstanceID(),
		LoginPolicyColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).OrderBy(LoginPolicyColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-WxYjr", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		factors, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	factors.State, err = q.latestState(ctx, loginPolicyTable)
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
			LoginPolicyColumnForceMFALocalOnly.identifier(),
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
			LoginPolicyColumnMultiFactorCheckLifetime.identifier(),
		).From(loginPolicyTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*LoginPolicy, error) {
			p := new(LoginPolicy)
			defaultRedirectURI := sql.NullString{}
			for rows.Next() {
				err := rows.Scan(
					&p.OrgID,
					&p.CreationDate,
					&p.ChangeDate,
					&p.Sequence,
					&p.AllowRegister,
					&p.AllowUsernamePassword,
					&p.AllowExternalIDPs,
					&p.ForceMFA,
					&p.ForceMFALocalOnly,
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
				)
				if err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-YcC53", "Errors.Internal")
				}
			}
			if p.OrgID == "" {
				return nil, zerrors.ThrowNotFound(nil, "QUERY-QsUBJ", "Errors.LoginPolicy.NotFound")
			}
			p.DefaultRedirectURI = defaultRedirectURI.String
			return p, nil
		}
}

func prepareLoginPolicy2FAsQuery() (sq.SelectBuilder, func(*sql.Row) (*SecondFactors, error)) {
	return sq.Select(
			LoginPolicyColumnSecondFactors.identifier(),
		).From(loginPolicyTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SecondFactors, error) {
			p := new(SecondFactors)
			err := row.Scan(
				&p.Factors,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-yPqIZ", "Errors.LoginPolicy.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Mr6H3", "Errors.Internal")
			}

			p.Count = uint64(len(p.Factors))
			return p, nil
		}
}

func prepareLoginPolicyMFAsQuery() (sq.SelectBuilder, func(*sql.Row) (*MultiFactors, error)) {
	return sq.Select(
			LoginPolicyColumnMultiFactors.identifier(),
		).From(loginPolicyTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*MultiFactors, error) {
			p := new(MultiFactors)
			err := row.Scan(
				&p.Factors,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-yPqIZ", "Errors.LoginPolicy.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Mr6H3", "Errors.Internal")
			}

			p.Count = uint64(len(p.Factors))
			return p, nil
		}
}
