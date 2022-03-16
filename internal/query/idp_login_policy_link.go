package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type IDPLoginPolicyLink struct {
	IDPID   string
	IDPName string
	IDPType domain.IDPConfigType
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
		name: projection.IDPLoginPolicyLinkTable,
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
)

func (q *Queries) IDPLoginPolicyLinks(ctx context.Context, resourceOwner string, queries *IDPLoginPolicyLinksSearchQuery) (idps *IDPLoginPolicyLinks, err error) {
	query, scan := prepareIDPLoginPolicyLinksQuery()
	stmt, args, err := queries.toQuery(query).Where(
		sq.Eq{
			IDPLoginPolicyLinkResourceOwnerCol.identifier(): resourceOwner,
			IDPLoginPolicyLinkInstanceIDCol.identifier():    authz.GetInstance(ctx).ID,
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-FDbKW", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-ZkKUc", "Errors.Internal")
	}
	idps, err = scan(rows)
	if err != nil {
		return nil, err
	}
	idps.LatestSequence, err = q.latestSequence(ctx, idpLoginPolicyLinkTable)
	return idps, err
}

func prepareIDPLoginPolicyLinksQuery() (sq.SelectBuilder, func(*sql.Rows) (*IDPLoginPolicyLinks, error)) {
	return sq.Select(
			IDPLoginPolicyLinkIDPIDCol.identifier(),
			IDPNameCol.identifier(),
			IDPTypeCol.identifier(),
			countColumn.identifier()).
			From(idpLoginPolicyLinkTable.identifier()).
			LeftJoin(join(IDPIDCol, IDPLoginPolicyLinkIDPIDCol)).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPLoginPolicyLinks, error) {
			links := make([]*IDPLoginPolicyLink, 0)
			var count uint64
			for rows.Next() {
				var (
					idpName = sql.NullString{}
					idpType = sql.NullInt16{}
					link    = new(IDPLoginPolicyLink)
				)
				err := rows.Scan(
					&link.IDPID,
					&idpName,
					&idpType,
					&count,
				)
				if err != nil {
					return nil, err
				}
				link.IDPName = idpName.String
				//IDPType 0 is oidc so we have to set unspecified manually
				if idpType.Valid {
					link.IDPType = domain.IDPConfigType(idpType.Int16)
				} else {
					link.IDPType = domain.IDPConfigTypeUnspecified
				}
				links = append(links, link)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-vOLFG", "Errors.Query.CloseRows")
			}

			return &IDPLoginPolicyLinks{
				Links: links,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
