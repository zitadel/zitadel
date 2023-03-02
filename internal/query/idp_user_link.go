package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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
		return nil, errors.ThrowInvalidArgument(err, "QUERY-4zzFK", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-C1E4D", "Errors.Internal")
	}
	idps, err = scan(rows)
	if err != nil {
		return nil, err
	}
	idps.LatestSequence, err = q.latestSequence(ctx, idpUserLinkTable)
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
				return nil, errors.ThrowInternal(err, "QUERY-nwx6U", "Errors.Query.CloseRows")
			}

			return &IDPUserLinks{
				Links: idps,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
