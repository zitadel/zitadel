package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
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

func (l *IDPUserLinks) RemoveNoPermission(ctx context.Context, permissionCheck domain.PermissionCheck) {
	removableIndexes := make([]int, 0)
	for i := range l.Links {
		ctxData := authz.GetCtxData(ctx)
		if ctxData.UserID != l.Links[i].UserID {
			if err := permissionCheck(ctx, domain.PermissionUserRead, l.Links[i].ResourceOwner, l.Links[i].UserID); err != nil {
				removableIndexes = append(removableIndexes, i)
			}
		}
	}
	removed := 0
	for _, removeIndex := range removableIndexes {
		l.Links = removeIDPLink(l.Links, removeIndex-removed)
		removed++
	}
	// reset count as some users could be removed
	l.SearchResponse.Count = uint64(len(l.Links))
}

func removeIDPLink(slice []*IDPUserLink, s int) []*IDPUserLink {
	return append(slice[:s], slice[s+1:]...)
}

func (q *Queries) IDPUserLinks(ctx context.Context, queries *IDPUserLinksSearchQuery, withOwnerRemoved bool) (idps *IDPUserLinks, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareIDPUserLinksQuery(ctx, q.client)
	eq := sq.Eq{IDPUserLinkInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[IDPUserLinkOwnerRemovedCol.identifier()] = false
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
	return NewTextQuery(IDPUserLinkExternalUserIDCol, value, TextEquals)
}

func prepareIDPUserLinksQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*IDPUserLinks, error)) {
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
			LeftJoin(join(IDPTemplateIDCol, IDPUserLinkIDPIDCol) + db.Timetravel(call.Took(ctx))).
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
