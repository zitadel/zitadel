package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type IDP struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string
	ID            string
	State         domain.IDPConfigState
	Name          string
	StylingType   domain.IDPConfigStylingType
	OwnerType     domain.IdentityProviderType
	AutoRegister  bool
	*OIDCIDP
	*JWTIDP
}

type IDPs struct {
	SearchResponse
	IDPs []*IDP
}

type OIDCIDP struct {
	IDPID                 string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	Issuer                string
	Scopes                database.StringArray
	DisplayNameMapping    domain.OIDCMappingField
	UsernameMapping       domain.OIDCMappingField
	AuthorizationEndpoint string
	TokenEndpoint         string
}

type JWTIDP struct {
	IDPID        string
	Issuer       string
	KeysEndpoint string
	HeaderName   string
	Endpoint     string
}

var (
	idpTable = table{
		name:          projection.IDPTable,
		instanceIDCol: projection.IDPInstanceIDCol,
	}
	IDPIDCol = Column{
		name:  projection.IDPIDCol,
		table: idpTable,
	}
	IDPCreationDateCol = Column{
		name:  projection.IDPCreationDateCol,
		table: idpTable,
	}
	IDPChangeDateCol = Column{
		name:  projection.IDPChangeDateCol,
		table: idpTable,
	}
	IDPSequenceCol = Column{
		name:  projection.IDPSequenceCol,
		table: idpTable,
	}
	IDPResourceOwnerCol = Column{
		name:  projection.IDPResourceOwnerCol,
		table: idpTable,
	}
	IDPInstanceIDCol = Column{
		name:  projection.IDPInstanceIDCol,
		table: idpTable,
	}
	IDPStateCol = Column{
		name:  projection.IDPStateCol,
		table: idpTable,
	}
	IDPNameCol = Column{
		name:  projection.IDPNameCol,
		table: idpTable,
	}
	IDPStylingTypeCol = Column{
		name:  projection.IDPStylingTypeCol,
		table: idpTable,
	}
	IDPOwnerTypeCol = Column{
		name:  projection.IDPOwnerTypeCol,
		table: idpTable,
	}
	IDPAutoRegisterCol = Column{
		name:  projection.IDPAutoRegisterCol,
		table: idpTable,
	}
	IDPTypeCol = Column{
		name:  projection.IDPTypeCol,
		table: idpTable,
	}
	IDPOwnerRemovedCol = Column{
		name:  projection.IDPOwnerRemovedCol,
		table: idpTable,
	}
)

var (
	oidcIDPTable = table{
		name:          projection.IDPOIDCTable,
		instanceIDCol: projection.OIDCConfigInstanceIDCol,
	}
	OIDCIDPColIDPID = Column{
		name:  projection.OIDCConfigIDPIDCol,
		table: oidcIDPTable,
	}
	OIDCIDPColClientID = Column{
		name:  projection.OIDCConfigClientIDCol,
		table: oidcIDPTable,
	}
	OIDCIDPColClientSecret = Column{
		name:  projection.OIDCConfigClientSecretCol,
		table: oidcIDPTable,
	}
	OIDCIDPColIssuer = Column{
		name:  projection.OIDCConfigIssuerCol,
		table: oidcIDPTable,
	}
	OIDCIDPColScopes = Column{
		name:  projection.OIDCConfigScopesCol,
		table: oidcIDPTable,
	}
	OIDCIDPColDisplayNameMapping = Column{
		name:  projection.OIDCConfigDisplayNameMappingCol,
		table: oidcIDPTable,
	}
	OIDCIDPColUsernameMapping = Column{
		name:  projection.OIDCConfigUsernameMappingCol,
		table: oidcIDPTable,
	}
	OIDCIDPColAuthorizationEndpoint = Column{
		name:  projection.OIDCConfigAuthorizationEndpointCol,
		table: oidcIDPTable,
	}
	OIDCIDPColTokenEndpoint = Column{
		name:  projection.OIDCConfigTokenEndpointCol,
		table: oidcIDPTable,
	}
)

var (
	jwtIDPTable = table{
		name:          projection.IDPJWTTable,
		instanceIDCol: projection.JWTConfigInstanceIDCol,
	}
	JWTIDPColIDPID = Column{
		name:  projection.JWTConfigIDPIDCol,
		table: jwtIDPTable,
	}
	JWTIDPColIssuer = Column{
		name:  projection.JWTConfigIssuerCol,
		table: jwtIDPTable,
	}
	JWTIDPColKeysEndpoint = Column{
		name:  projection.JWTConfigKeysEndpointCol,
		table: jwtIDPTable,
	}
	JWTIDPColHeaderName = Column{
		name:  projection.JWTConfigHeaderNameCol,
		table: jwtIDPTable,
	}
	JWTIDPColEndpoint = Column{
		name:  projection.JWTConfigEndpointCol,
		table: jwtIDPTable,
	}
)

// IDPByIDAndResourceOwner searches for the requested id in the context of the resource owner and IAM
func (q *Queries) IDPByIDAndResourceOwner(ctx context.Context, shouldTriggerBulk bool, id, resourceOwner string, withOwnerRemoved bool) (_ *IDP, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.IDPProjection.Trigger(ctx)
	}

	eq := sq.Eq{
		IDPIDCol.identifier():         id,
		IDPInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[IDPOwnerRemovedCol.identifier()] = false
	}
	where := sq.And{
		eq,
		sq.Or{
			sq.Eq{IDPResourceOwnerCol.identifier(): resourceOwner},
			sq.Eq{IDPResourceOwnerCol.identifier(): authz.GetInstance(ctx).InstanceID()},
		},
	}
	stmt, scan := prepareIDPByIDQuery(ctx, q.client)
	query, args, err := stmt.Where(where).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-0gocI", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

// IDPs searches idps matching the query
func (q *Queries) IDPs(ctx context.Context, queries *IDPSearchQueries, withOwnerRemoved bool) (idps *IDPs, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareIDPsQuery(ctx, q.client)
	eq := sq.Eq{
		IDPInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[IDPOwnerRemovedCol.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-X6X7y", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-xPlVH", "Errors.Internal")
	}
	idps, err = scan(rows)
	if err != nil {
		return nil, err
	}
	idps.LatestSequence, err = q.latestSequence(ctx, idpTable)
	return idps, err
}

type IDPSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewIDPIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(IDPIDCol, id, TextEquals)
}

func NewIDPOwnerTypeSearchQuery(ownerType domain.IdentityProviderType) (SearchQuery, error) {
	return NewNumberQuery(IDPOwnerTypeCol, ownerType, NumberEquals)
}

func NewIDPNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(IDPNameCol, value, method)
}

func NewIDPResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(IDPResourceOwnerCol, value, TextEquals)
}

func NewIDPResourceOwnerListSearchQuery(ids ...string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(IDPResourceOwnerCol, list, ListIn)
}

func (q *IDPSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareIDPByIDQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*IDP, error)) {
	return sq.Select(
			IDPIDCol.identifier(),
			IDPResourceOwnerCol.identifier(),
			IDPCreationDateCol.identifier(),
			IDPChangeDateCol.identifier(),
			IDPSequenceCol.identifier(),
			IDPStateCol.identifier(),
			IDPNameCol.identifier(),
			IDPStylingTypeCol.identifier(),
			IDPOwnerTypeCol.identifier(),
			IDPAutoRegisterCol.identifier(),
			OIDCIDPColIDPID.identifier(),
			OIDCIDPColClientID.identifier(),
			OIDCIDPColClientSecret.identifier(),
			OIDCIDPColIssuer.identifier(),
			OIDCIDPColScopes.identifier(),
			OIDCIDPColDisplayNameMapping.identifier(),
			OIDCIDPColUsernameMapping.identifier(),
			OIDCIDPColAuthorizationEndpoint.identifier(),
			OIDCIDPColTokenEndpoint.identifier(),
			JWTIDPColIDPID.identifier(),
			JWTIDPColIssuer.identifier(),
			JWTIDPColKeysEndpoint.identifier(),
			JWTIDPColHeaderName.identifier(),
			JWTIDPColEndpoint.identifier(),
		).From(idpTable.identifier()).
			LeftJoin(join(OIDCIDPColIDPID, IDPIDCol)).
			LeftJoin(join(JWTIDPColIDPID, IDPIDCol) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IDP, error) {
			idp := new(IDP)

			oidcIDPID := sql.NullString{}
			oidcClientID := sql.NullString{}
			oidcClientSecret := new(crypto.CryptoValue)
			oidcIssuer := sql.NullString{}
			oidcScopes := database.StringArray{}
			oidcDisplayNameMapping := sql.NullInt32{}
			oidcUsernameMapping := sql.NullInt32{}
			oidcAuthorizationEndpoint := sql.NullString{}
			oidcTokenEndpoint := sql.NullString{}

			jwtIDPID := sql.NullString{}
			jwtIssuer := sql.NullString{}
			jwtKeysEndpoint := sql.NullString{}
			jwtHeaderName := sql.NullString{}
			jwtEndpoint := sql.NullString{}

			err := row.Scan(
				&idp.ID,
				&idp.ResourceOwner,
				&idp.CreationDate,
				&idp.ChangeDate,
				&idp.Sequence,
				&idp.State,
				&idp.Name,
				&idp.StylingType,
				&idp.OwnerType,
				&idp.AutoRegister,
				&oidcIDPID,
				&oidcClientID,
				oidcClientSecret,
				&oidcIssuer,
				&oidcScopes,
				&oidcDisplayNameMapping,
				&oidcUsernameMapping,
				&oidcAuthorizationEndpoint,
				&oidcTokenEndpoint,
				&jwtIDPID,
				&jwtIssuer,
				&jwtKeysEndpoint,
				&jwtHeaderName,
				&jwtEndpoint,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-rhR2o", "Errors.IDPConfig.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-zE3Ro", "Errors.Internal")
			}

			if oidcIDPID.Valid {
				idp.OIDCIDP = &OIDCIDP{
					IDPID:                 oidcIDPID.String,
					ClientID:              oidcClientID.String,
					ClientSecret:          oidcClientSecret,
					Issuer:                oidcIssuer.String,
					Scopes:                oidcScopes,
					DisplayNameMapping:    domain.OIDCMappingField(oidcDisplayNameMapping.Int32),
					UsernameMapping:       domain.OIDCMappingField(oidcUsernameMapping.Int32),
					AuthorizationEndpoint: oidcAuthorizationEndpoint.String,
					TokenEndpoint:         oidcTokenEndpoint.String,
				}
			} else if jwtIDPID.Valid {
				idp.JWTIDP = &JWTIDP{
					IDPID:        jwtIDPID.String,
					Issuer:       jwtIssuer.String,
					KeysEndpoint: jwtKeysEndpoint.String,
					HeaderName:   jwtHeaderName.String,
					Endpoint:     jwtEndpoint.String,
				}
			}

			return idp, nil
		}
}

func prepareIDPsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*IDPs, error)) {
	return sq.Select(
			IDPIDCol.identifier(),
			IDPResourceOwnerCol.identifier(),
			IDPCreationDateCol.identifier(),
			IDPChangeDateCol.identifier(),
			IDPSequenceCol.identifier(),
			IDPStateCol.identifier(),
			IDPNameCol.identifier(),
			IDPStylingTypeCol.identifier(),
			IDPOwnerTypeCol.identifier(),
			IDPAutoRegisterCol.identifier(),
			OIDCIDPColIDPID.identifier(),
			OIDCIDPColClientID.identifier(),
			OIDCIDPColClientSecret.identifier(),
			OIDCIDPColIssuer.identifier(),
			OIDCIDPColScopes.identifier(),
			OIDCIDPColDisplayNameMapping.identifier(),
			OIDCIDPColUsernameMapping.identifier(),
			OIDCIDPColAuthorizationEndpoint.identifier(),
			OIDCIDPColTokenEndpoint.identifier(),
			JWTIDPColIDPID.identifier(),
			JWTIDPColIssuer.identifier(),
			JWTIDPColKeysEndpoint.identifier(),
			JWTIDPColHeaderName.identifier(),
			JWTIDPColEndpoint.identifier(),
			countColumn.identifier(),
		).From(idpTable.identifier()).
			LeftJoin(join(OIDCIDPColIDPID, IDPIDCol)).
			LeftJoin(join(JWTIDPColIDPID, IDPIDCol) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPs, error) {
			idps := make([]*IDP, 0)
			var count uint64
			for rows.Next() {
				idp := new(IDP)

				oidcIDPID := sql.NullString{}
				oidcClientID := sql.NullString{}
				oidcClientSecret := new(crypto.CryptoValue)
				oidcIssuer := sql.NullString{}
				oidcScopes := database.StringArray{}
				oidcDisplayNameMapping := sql.NullInt32{}
				oidcUsernameMapping := sql.NullInt32{}
				oidcAuthorizationEndpoint := sql.NullString{}
				oidcTokenEndpoint := sql.NullString{}

				jwtIDPID := sql.NullString{}
				jwtIssuer := sql.NullString{}
				jwtKeysEndpoint := sql.NullString{}
				jwtHeaderName := sql.NullString{}
				jwtEndpoint := sql.NullString{}

				err := rows.Scan(
					&idp.ID,
					&idp.ResourceOwner,
					&idp.CreationDate,
					&idp.ChangeDate,
					&idp.Sequence,
					&idp.State,
					&idp.Name,
					&idp.StylingType,
					&idp.OwnerType,
					&idp.AutoRegister,
					// oidc config
					&oidcIDPID,
					&oidcClientID,
					oidcClientSecret,
					&oidcIssuer,
					&oidcScopes,
					&oidcDisplayNameMapping,
					&oidcUsernameMapping,
					&oidcAuthorizationEndpoint,
					&oidcTokenEndpoint,
					// jwt config
					&jwtIDPID,
					&jwtIssuer,
					&jwtKeysEndpoint,
					&jwtHeaderName,
					&jwtEndpoint,
					&count,
				)

				if err != nil {
					return nil, err
				}

				if oidcIDPID.Valid {
					idp.OIDCIDP = &OIDCIDP{
						IDPID:                 oidcIDPID.String,
						ClientID:              oidcClientID.String,
						ClientSecret:          oidcClientSecret,
						Issuer:                oidcIssuer.String,
						Scopes:                oidcScopes,
						DisplayNameMapping:    domain.OIDCMappingField(oidcDisplayNameMapping.Int32),
						UsernameMapping:       domain.OIDCMappingField(oidcUsernameMapping.Int32),
						AuthorizationEndpoint: oidcAuthorizationEndpoint.String,
						TokenEndpoint:         oidcTokenEndpoint.String,
					}
				} else if jwtIDPID.Valid {
					idp.JWTIDP = &JWTIDP{
						IDPID:        jwtIDPID.String,
						Issuer:       jwtIssuer.String,
						KeysEndpoint: jwtKeysEndpoint.String,
						HeaderName:   jwtHeaderName.String,
						Endpoint:     jwtEndpoint.String,
					}
				}

				idps = append(idps, idp)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-iiBgK", "Errors.Query.CloseRows")
			}

			return &IDPs{
				IDPs: idps,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *Queries) GetOIDCIDPClientSecret(ctx context.Context, shouldRealTime bool, resourceowner, idpID string, withOwnerRemoved bool) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	idp, err := q.IDPByIDAndResourceOwner(ctx, shouldRealTime, idpID, resourceowner, withOwnerRemoved)
	if err != nil {
		return "", err
	}

	if idp.ClientSecret != nil && idp.ClientSecret.Crypted != nil {
		return crypto.DecryptString(idp.ClientSecret, q.idpConfigEncryption)
	}
	return "", errors.ThrowNotFound(nil, "QUERY-bsm2o", "Errors.Query.NotFound")
}
