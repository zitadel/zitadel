package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type IDPLoginPolicyLink struct {
	IDPID             string
	IDPName           string
	IDPType           domain.IDPType
	OwnerType         domain.IdentityProviderType
	IsCreationAllowed bool
	IsLinkingAllowed  bool
	IsAutoCreation    bool
	IsAutoUpdate      bool
	AutoLinking       domain.AutoLinkingOption
}

type IDPLoginPolicyLinks struct {
	SearchResponse
	Links []*IDPLoginPolicyLink
}

type IDPLoginPolicyLinksSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *IDPLoginPolicyLinksSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

var (
	idpLoginPolicyLinkTable = table{
		name:          projection.IDPLoginPolicyLinkTable,
		instanceIDCol: projection.IDPLoginPolicyLinkInstanceIDCol,
	}
	IDPLoginPolicyLinkIDPIDCol = Column{
		name:  projection.IDPLoginPolicyLinkIDPIDCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkCreationDateCol = Column{
		name:  projection.IDPLoginPolicyLinkCreationDateCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkChangeDateCol = Column{
		name:  projection.IDPLoginPolicyLinkChangeDateCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkSequenceCol = Column{
		name:  projection.IDPLoginPolicyLinkSequenceCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkAggregateIDCol = Column{
		name:  projection.IDPLoginPolicyLinkAggregateIDCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkResourceOwnerCol = Column{
		name:  projection.IDPLoginPolicyLinkResourceOwnerCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkInstanceIDCol = Column{
		name:  projection.IDPLoginPolicyLinkInstanceIDCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkProviderTypeCol = Column{
		name:  projection.IDPLoginPolicyLinkProviderTypeCol,
		table: idpLoginPolicyLinkTable,
	}
	IDPLoginPolicyLinkOwnerRemovedCol = Column{
		name:  projection.IDPLoginPolicyLinkOwnerRemovedCol,
		table: idpLoginPolicyLinkTable,
	}

	idpLoginPolicyOwnerTable           = loginPolicyTable.setAlias("login_policy_owner")
	idpLoginPolicyOwnerIDCol           = LoginPolicyColumnOrgID.setTable(idpLoginPolicyOwnerTable)
	idpLoginPolicyOwnerInstanceIDCol   = LoginPolicyColumnInstanceID.setTable(idpLoginPolicyOwnerTable)
	idpLoginPolicyOwnerIsDefaultCol    = LoginPolicyColumnIsDefault.setTable(idpLoginPolicyOwnerTable)
	idpLoginPolicyOwnerOwnerRemovedCol = LoginPolicyColumnOwnerRemoved.setTable(idpLoginPolicyOwnerTable)
)

func (q *Queries) IDPLoginPolicyLinks(ctx context.Context, resourceOwner string, queries *IDPLoginPolicyLinksSearchQuery, withOwnerRemoved bool) (idps *IDPLoginPolicyLinks, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareIDPLoginPolicyLinksQuery(ctx, resourceOwner)
	eq := sq.Eq{
		IDPLoginPolicyLinkInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[IDPLoginPolicyLinkOwnerRemovedCol.identifier()] = false
		eq[idpLoginPolicyOwnerOwnerRemovedCol.identifier()] = false
	}

	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-FDbKW", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		idps, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-ZkKUc", "Errors.Internal")
	}
	idps.State, err = q.latestState(ctx, idpLoginPolicyLinkTable)
	return idps, err
}

//nolint:gocognit
func prepareIDPLoginPolicyLinksQuery(ctx context.Context, resourceOwner string) (sq.SelectBuilder, func(*sql.Rows) (*IDPLoginPolicyLinks, error)) {
	resourceOwnerQuery, resourceOwnerArgs, err := prepareIDPLoginPolicyLinksResourceOwnerQuery(ctx, resourceOwner)
	if err != nil {
		return sq.SelectBuilder{}, nil
	}
	return sq.Select(
			IDPLoginPolicyLinkIDPIDCol.identifier(),
			IDPTemplateNameCol.identifier(),
			IDPTemplateTypeCol.identifier(),
			IDPTemplateOwnerTypeCol.identifier(),
			IDPTemplateIsCreationAllowedCol.identifier(),
			IDPTemplateIsLinkingAllowedCol.identifier(),
			IDPTemplateIsAutoCreationCol.identifier(),
			IDPTemplateIsAutoUpdateCol.identifier(),
			IDPTemplateAutoLinkingCol.identifier(),
			countColumn.identifier()).
			From(idpLoginPolicyLinkTable.identifier()).
			LeftJoin(join(IDPTemplateIDCol, IDPLoginPolicyLinkIDPIDCol)).
			RightJoin("("+resourceOwnerQuery+") AS "+idpLoginPolicyOwnerTable.alias+" ON "+
				idpLoginPolicyOwnerIDCol.identifier()+" = "+IDPLoginPolicyLinkResourceOwnerCol.identifier()+" AND "+
				idpLoginPolicyOwnerInstanceIDCol.identifier()+" = "+IDPLoginPolicyLinkInstanceIDCol.identifier(),
				resourceOwnerArgs...).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPLoginPolicyLinks, error) {
			links := make([]*IDPLoginPolicyLink, 0)
			var count uint64
			for rows.Next() {
				var (
					idpName           = sql.NullString{}
					idpType           = sql.NullInt16{}
					idpOwnerType      = sql.NullInt16{}
					link              = new(IDPLoginPolicyLink)
					isCreationAllowed = sql.NullBool{}
					isLinkingAllowed  = sql.NullBool{}
					isAutoCreation    = sql.NullBool{}
					isAutoUpdate      = sql.NullBool{}
					autoLinking       = sql.NullInt16{}
				)
				err := rows.Scan(
					&link.IDPID,
					&idpName,
					&idpType,
					&idpOwnerType,
					&isCreationAllowed,
					&isLinkingAllowed,
					&isAutoCreation,
					&isAutoUpdate,
					&autoLinking,
					&count,
				)
				if err != nil {
					return nil, err
				}
				if idpName.Valid {
					link.IDPName = idpName.String
				}
				//IDPType 0 is oidc so we have to set unspecified manually
				if idpType.Valid {
					link.IDPType = domain.IDPType(idpType.Int16)
				} else {
					link.IDPType = domain.IDPTypeUnspecified
				}
				if idpOwnerType.Valid {
					link.OwnerType = domain.IdentityProviderType(idpOwnerType.Int16)
				}
				if isCreationAllowed.Valid {
					link.IsCreationAllowed = isCreationAllowed.Bool
				}
				if isLinkingAllowed.Valid {
					link.IsLinkingAllowed = isLinkingAllowed.Bool
				}
				if isAutoCreation.Valid {
					link.IsAutoCreation = isAutoCreation.Bool
				}
				if isAutoUpdate.Valid {
					link.IsAutoUpdate = isAutoUpdate.Bool
				}
				if autoLinking.Valid {
					link.AutoLinking = domain.AutoLinkingOption(autoLinking.Int16)
				} else {
					link.AutoLinking = domain.AutoLinkingOptionUnspecified
				}
				links = append(links, link)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-vOLFG", "Errors.Query.CloseRows")
			}

			return &IDPLoginPolicyLinks{
				Links: links,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareIDPLoginPolicyLinksResourceOwnerQuery(ctx context.Context, resourceOwner string) (string, []interface{}, error) {
	eqPolicy := sq.Eq{idpLoginPolicyOwnerInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID()}
	return sq.Select(
		idpLoginPolicyOwnerIDCol.identifier(),
		idpLoginPolicyOwnerInstanceIDCol.identifier(),
		idpLoginPolicyOwnerOwnerRemovedCol.identifier(),
	).
		From(idpLoginPolicyOwnerTable.identifier()).
		Where(
			sq.And{
				eqPolicy,
				sq.Or{
					sq.Eq{idpLoginPolicyOwnerIDCol.identifier(): resourceOwner},
					sq.Eq{idpLoginPolicyOwnerIDCol.identifier(): authz.GetInstance(ctx).InstanceID()},
				},
			}).
		Limit(1).OrderBy(idpLoginPolicyOwnerIsDefaultCol.identifier()).ToSql()
}
