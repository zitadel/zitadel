package query

import (
	"context"
	"database/sql"
	"slices"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type IDPUserLink struct {
	IDPID            string
	UserID           string
	IDPName          string
	ProvidedUserID   string
	ProvidedUsername string
	ResourceOwner    string
	IDPType          domain.IDPType
}

type IDPUserLinks struct {
	SearchResponse
	Links []*IDPUserLink
}

type IDPUserLinksSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *IDPUserLinksSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *IDPUserLinksSearchQuery) hasUserID() bool {
	for _, query := range q.Queries {
		if query.Col() == IDPUserLinkUserIDCol {
			return true
		}
	}
	return false
}

var (
	idpUserLinkTable = table{
		name:          projection.IDPUserLinkTable,
		instanceIDCol: projection.IDPUserLinkInstanceIDCol,
	}
	IDPUserLinkIDPIDCol = Column{
		name:  projection.IDPUserLinkIDPIDCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkUserIDCol = Column{
		name:  projection.IDPUserLinkUserIDCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkExternalUserIDCol = Column{
		name:  projection.IDPUserLinkExternalUserIDCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkCreationDateCol = Column{
		name:  projection.IDPUserLinkCreationDateCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkChangeDateCol = Column{
		name:  projection.IDPUserLinkChangeDateCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkSequenceCol = Column{
		name:  projection.IDPUserLinkSequenceCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkResourceOwnerCol = Column{
		name:  projection.IDPUserLinkResourceOwnerCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkInstanceIDCol = Column{
		name:  projection.IDPUserLinkInstanceIDCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkDisplayNameCol = Column{
		name:  projection.IDPUserLinkDisplayNameCol,
		table: idpUserLinkTable,
	}
	IDPUserLinkOwnerRemovedCol = Column{
		name:  projection.IDPUserLinkOwnerRemovedCol,
		table: idpUserLinkTable,
	}
)

func idpLinksCheckPermission(ctx context.Context, links *IDPUserLinks, permissionCheck domain.PermissionCheck) {
	links.Links = slices.DeleteFunc(links.Links,
		func(link *IDPUserLink) bool {
			return userCheckPermission(ctx, link.ResourceOwner, link.UserID, permissionCheck) != nil
		},
	)
}

func idpLinksPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *IDPUserLinksSearchQuery) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		IDPUserLinkResourceOwnerCol,
		domain.PermissionUserRead,
		SingleOrgPermissionOption(queries.Queries),
		OwnedRowsPermissionOption(IDPUserLinkUserIDCol),
	)
	return query.JoinClause(join, args...)
}

func (q *Queries) IDPUserLinks(ctx context.Context, queries *IDPUserLinksSearchQuery, permissionCheck domain.PermissionCheck) (idps *IDPUserLinks, err error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	links, err := q.idpUserLinks(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && len(links.Links) > 0 && !permissionCheckV2 {
		// when userID for query is provided, only one check has to be done
		if queries.hasUserID() {
			if err := userCheckPermission(ctx, links.Links[0].ResourceOwner, links.Links[0].UserID, permissionCheck); err != nil {
				return nil, err
			}
		} else {
			idpLinksCheckPermission(ctx, links, permissionCheck)
		}
	}
	return links, nil
}

func (q *Queries) idpUserLinks(ctx context.Context, queries *IDPUserLinksSearchQuery, permissionCheckV2 bool) (idps *IDPUserLinks, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareIDPUserLinksQuery()
	query = idpLinksPermissionCheckV2(ctx, query, permissionCheckV2, queries)
	eq := sq.Eq{
		IDPUserLinkInstanceIDCol.identifier():   authz.GetInstance(ctx).InstanceID(),
		IDPUserLinkOwnerRemovedCol.identifier(): false,
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-4zzFK", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		idps, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-C1E4D", "Errors.Internal")
	}
	idps.State, err = q.latestState(ctx, idpUserLinkTable)
	return idps, err
}

func NewIDPUserLinkIDPIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(IDPUserLinkIDPIDCol, value, TextEquals)
}

func NewIDPUserLinksUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(IDPUserLinkUserIDCol, value, TextEquals)
}

func NewIDPUserLinksResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(IDPUserLinkResourceOwnerCol, value, TextEquals)
}

func NewIDPUserLinksExternalIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(IDPUserLinkExternalUserIDCol, value, TextEqualsIgnoreCase)
}

func prepareIDPUserLinksQuery() (sq.SelectBuilder, func(*sql.Rows) (*IDPUserLinks, error)) {
	return sq.Select(
			IDPUserLinkIDPIDCol.identifier(),
			IDPUserLinkUserIDCol.identifier(),
			IDPTemplateNameCol.identifier(),
			IDPUserLinkExternalUserIDCol.identifier(),
			IDPUserLinkDisplayNameCol.identifier(),
			IDPTemplateTypeCol.identifier(),
			IDPUserLinkResourceOwnerCol.identifier(),
			countColumn.identifier()).
			From(idpUserLinkTable.identifier()).
			LeftJoin(join(IDPTemplateIDCol, IDPUserLinkIDPIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPUserLinks, error) {
			idps := make([]*IDPUserLink, 0)
			var count uint64
			for rows.Next() {
				var (
					idpName = sql.NullString{}
					idpType = sql.NullInt16{}
					idp     = new(IDPUserLink)
				)
				err := rows.Scan(
					&idp.IDPID,
					&idp.UserID,
					&idpName,
					&idp.ProvidedUserID,
					&idp.ProvidedUsername,
					&idpType,
					&idp.ResourceOwner,
					&count,
				)
				if err != nil {
					return nil, err
				}
				idp.IDPName = idpName.String
				//IDPType 0 is oidc so we have to set unspecified manually
				if idpType.Valid {
					idp.IDPType = domain.IDPType(idpType.Int16)
				} else {
					idp.IDPType = domain.IDPTypeUnspecified
				}
				idps = append(idps, idp)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-nwx6U", "Errors.Query.CloseRows")
			}

			return &IDPUserLinks{
				Links: idps,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
