package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
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
}

type OIDCIDP struct {
	IDP
}

type JWTIDP struct {
	IDP
}

func (q *Queries) IDPByID(ctx context.Context, id string) (*IDP, error) {
	return nil, nil
}

func prepareIDPByIDQuery() (sq.SelectBuilder, func(*sql.Row) (*IDP, error)) {
	return sq.Select(
		"hodor",
	).From(idpTable).
		LeftJoin(idpOIDCConfigTable + " ON " + idpTable + "." + idpIDCol.toColumnName() + " = " + idpOIDCConfigTable + "." + idpOIDCIDPIDCol.toColumnName()).
		LeftJoin(idpOIDCConfigTable + " ON " + idpTable + "." + idpIDCol.toColumnName() + " = " + idpJWTConfigTable + "." + idpJWTIDPIDCol.toColumnName()).
		PlaceholderFormat(sq.Dollar), nil
}

type idpColumn int32

const (
	idpIDCol idpColumn = iota + 1
	idpJWTIDPIDCol
	idpOIDCIDPIDCol
)

func (c idpColumn) toColumnName() string {
	return ""
}
