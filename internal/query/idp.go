package query

import (
	"context"
	"database/sql"
	errs "errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

const (
	idpTable           = "zitadel.projections.idps"
	idpOIDCConfigTable = "zitadel.projections.idps_oidc_config"
	idpJWTConfigTable  = "zitadel.projections.idps_jwt_config"
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
	Scopes                []string
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

func (q *Queries) IDPByID(ctx context.Context, id string) (*IDP, error) {
	stmt, scan := prepareIDPByIDQuery()
	query, args, err := stmt.Where(sq.Eq{
		IDPIDCol.toColumnName(): id,
	}, id).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-0gocI", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) SearchIDPs(ctx context.Context, queries *IDPSearchQueries) (idps *IDPs, err error) {
	query, scan := prepareIDPsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-zC6gk", "Errors.idps.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		log.Println(err)
		return nil, errors.ThrowInternal(err, "QUERY-YTug9", "Errors.idps.internal")
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
	return nil, errors.ThrowUnimplemented(nil, "QUERY-8yZAI", "number comparison not implemented")
}

func NewIDPNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(IDPNameCol, value, method)
}

func (q *IDPSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func prepareIDPByIDQuery() (sq.SelectBuilder, func(*sql.Row) (*IDP, error)) {
	return sq.Select(
			IDPIDCol.toColumnName(),
			IDPStateCol.toColumnName(),
			IDPNameCol.toColumnName(),
			IDPStylingTypeCol.toColumnName(),
			IDPOwnerCol.toColumnName(),
			IDPAutoRegisterCol.toColumnName(),
			oidcIDPIDPIDCol.toColumnName(),
			oidcIDPClientIDCol.toColumnName(),
			oidcIDPClientSecretCol.toColumnName(),
			oidcIDPIssuerCol.toColumnName(),
			oidcIDPScopesCol.toColumnName(),
			oidcIDPDisplayNameMappingCol.toColumnName(),
			oidcIDPUsernameMappingCol.toColumnName(),
			oidcIDPAuthorizationEndpointCol.toColumnName(),
			oidcIDPTokenEndpointCol.toColumnName(),
			jwtIDPIDPIDCol.toColumnName(),
			jwtIDPIssuerCol.toColumnName(),
			jwtIDPKeysEndpointCol.toColumnName(),
			jwtIDPHeaderNameCol.toColumnName(),
			jwtIDPEndpointCol.toColumnName(),
		).From(idpTable).
			LeftJoin(idpOIDCConfigTable + " ON " + idpTable + "." + IDPIDCol.toColumnName() + " = " + idpOIDCConfigTable + "." + oidcIDPIDPIDCol.toColumnName()).
			LeftJoin(idpOIDCConfigTable + " ON " + idpTable + "." + IDPIDCol.toColumnName() + " = " + idpJWTConfigTable + "." + jwtIDPIDPIDCol.toColumnName()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IDP, error) {
			idp := new(IDP)
			oidc := new(OIDCIDP)
			jwt := new(JWTIDP)
			err := row.Scan(
				&idp.ID,
				&idp.State,
				&idp.Name,
				&idp.StylingType,
				&idp.OwnerType,
				&idp.AutoRegister,
				&oidc.IDPID,
				&oidc.ClientID,
				&oidc.ClientSecret,
				&oidc.Issuer,
				&oidc.Scopes,
				&oidc.DisplayNameMapping,
				&oidc.UsernameMapping,
				&oidc.AuthorizationEndpoint,
				&oidc.TokenEndpoint,
				&jwt.IDPID,
				&jwt.Issuer,
				&jwt.KeysEndpoint,
				&jwt.HeaderName,
				&jwt.Endpoint,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-rhR2o", "errors.idp.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-zE3Ro", "errors.internal")
			}

			if oidc.IDPID != "" {
				idp.OIDCIDP = oidc
			} else if jwt.IDPID != "" {
				idp.JWTIDP = jwt
			}

			return idp, nil
		}
}

func prepareIDPsQuery() (sq.SelectBuilder, func(*sql.Rows) (*IDPs, error)) {
	return sq.Select(
			IDPIDCol.toColumnName(),
			IDPStateCol.toColumnName(),
			IDPNameCol.toColumnName(),
			IDPStylingTypeCol.toColumnName(),
			IDPOwnerCol.toColumnName(),
			IDPAutoRegisterCol.toColumnName(),
			oidcIDPIDPIDCol.toColumnName(),
			oidcIDPClientIDCol.toColumnName(),
			oidcIDPClientSecretCol.toColumnName(),
			oidcIDPIssuerCol.toColumnName(),
			oidcIDPScopesCol.toColumnName(),
			oidcIDPDisplayNameMappingCol.toColumnName(),
			oidcIDPUsernameMappingCol.toColumnName(),
			oidcIDPAuthorizationEndpointCol.toColumnName(),
			oidcIDPTokenEndpointCol.toColumnName(),
			jwtIDPIDPIDCol.toColumnName(),
			jwtIDPIssuerCol.toColumnName(),
			jwtIDPKeysEndpointCol.toColumnName(),
			jwtIDPHeaderNameCol.toColumnName(),
			jwtIDPEndpointCol.toColumnName(),
			"COUNT(name) OVER ()").
			From(idpTable).
			LeftJoin(idpOIDCConfigTable + " ON " + idpTable + "." + IDPIDCol.toColumnName() + " = " + idpOIDCConfigTable + "." + oidcIDPIDPIDCol.toColumnName()).
			LeftJoin(idpOIDCConfigTable + " ON " + idpTable + "." + IDPIDCol.toColumnName() + " = " + idpJWTConfigTable + "." + jwtIDPIDPIDCol.toColumnName()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPs, error) {
			idps := make([]*IDP, 0)
			var count uint64
			for rows.Next() {
				idp := new(IDP)
				oidc := new(OIDCIDP)
				jwt := new(JWTIDP)
				err := rows.Scan(
					&idp.ID,
					&idp.State,
					&idp.Name,
					&idp.StylingType,
					&idp.OwnerType,
					&idp.AutoRegister,
					&oidc.IDPID,
					&oidc.ClientID,
					&oidc.ClientSecret,
					&oidc.Issuer,
					&oidc.Scopes,
					&oidc.DisplayNameMapping,
					&oidc.UsernameMapping,
					&oidc.AuthorizationEndpoint,
					&oidc.TokenEndpoint,
					&jwt.IDPID,
					&jwt.Issuer,
					&jwt.KeysEndpoint,
					&jwt.HeaderName,
					&jwt.Endpoint,
					&count,
				)
				if err != nil {
					return nil, err
				}

				if oidc.IDPID != "" {
					idp.OIDCIDP = oidc
				} else if jwt.IDPID != "" {
					idp.JWTIDP = jwt
				}

				idps = append(idps, idp)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-iiBgK", "unable to close rows")
			}

			return &IDPs{
				IDPs: idps,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

type idpColumn int32

const (
	IDPIDCol idpColumn = iota + 1
	IDPStateCol
	IDPNameCol
	IDPStylingTypeCol
	IDPOwnerCol
	IDPAutoRegisterCol
)

func (c idpColumn) toColumnName() string {
	switch c {
	case IDPIDCol:
		return "id"
	case IDPStateCol:
		return "state"
	case IDPNameCol:
		return "name"
	case IDPStylingTypeCol:
		return "styling_type"
	case IDPOwnerCol:
		return "owner"
	case IDPAutoRegisterCol:
		return "auto_register"
	default:
		return ""
	}
}

type oidcIDPColumn int32

const (
	oidcIDPIDPIDCol oidcIDPColumn = iota + 1
	oidcIDPClientIDCol
	oidcIDPClientSecretCol
	oidcIDPIssuerCol
	oidcIDPScopesCol
	oidcIDPDisplayNameMappingCol
	oidcIDPUsernameMappingCol
	oidcIDPAuthorizationEndpointCol
	oidcIDPTokenEndpointCol
)

func (c oidcIDPColumn) toColumnName() string {
	switch c {
	case oidcIDPIDPIDCol:
		return "idp_id"
	case oidcIDPClientIDCol:
		return "client_id"
	case oidcIDPClientSecretCol:
		return "client_secret"
	case oidcIDPIssuerCol:
		return "issuer"
	case oidcIDPScopesCol:
		return "scopes"
	case oidcIDPDisplayNameMappingCol:
		return "display_name_mapping"
	case oidcIDPUsernameMappingCol:
		return "username_mapping"
	case oidcIDPAuthorizationEndpointCol:
		return "authorization_endpoint"
	case oidcIDPTokenEndpointCol:
		return "token_endpoint"
	default:
		return ""
	}
}

type jwtIDPColumn int32

const (
	jwtIDPIDPIDCol jwtIDPColumn = iota + 1
	jwtIDPIssuerCol
	jwtIDPKeysEndpointCol
	jwtIDPHeaderNameCol
	jwtIDPEndpointCol
)

func (c jwtIDPColumn) toColumnName() string {
	switch c {
	case jwtIDPIDPIDCol:
		return "idp_id"
	case jwtIDPIssuerCol:
		return "issuer"
	case jwtIDPKeysEndpointCol:
		return "keys_endpoint"
	case jwtIDPHeaderNameCol:
		return "header_name"
	case jwtIDPEndpointCol:
		return "endpoint"
	default:
		return ""
	}
}
