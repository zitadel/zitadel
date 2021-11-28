package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	idpUserLinkTable = table{
		name: projection.IDPUserLinkTable,
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
	IDPUserLinkDisplayNameCol = Column{
		name:  projection.IDPUserLinkDisplayNameCol,
		table: idpUserLinkTable,
	}
)

type LinkedIDP struct {
	IDPID            string
	UserID           string
	IDPName          string
	ProvidedUserID   string
	ProvidedUsername string
	IDPType          domain.IDPConfigType
}

type LinkedIDPs struct {
	SearchResponse
	IDPs []*LinkedIDP
}

type LinkedIDPsSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *LinkedIDPsSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) LinkedIDPsByUser(ctx context.Context, queries *LinkedIDPsSearchQuery) (idps *LinkedIDPs, err error) {
	query, scan := prepareLinkedIDPsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
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

func NewLinkedIDPsUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(OrgColumnDomain, value, TextEquals)
}

func prepareLinkedIDPsQuery() (sq.SelectBuilder, func(*sql.Rows) (*LinkedIDPs, error)) {
	return sq.Select(
			IDPUserLinkIDPIDCol.identifier(),
			IDPUserLinkUserIDCol.identifier(),
			IDPNameCol.identifier(),
			IDPUserLinkExternalUserIDCol.identifier(),
			IDPUserLinkDisplayNameCol.identifier(),
			IDPTypeCol.identifier(),
			countColumn.identifier()).
			From(idpUserLinkTable.identifier()).
			LeftJoin(join(IDPUserLinkIDPIDCol, IDPIDCol)).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*LinkedIDPs, error) {
			idps := make([]*LinkedIDP, 0)
			var count uint64
			for rows.Next() {
				idp := new(LinkedIDP)
				err := rows.Scan(
					&idp.IDPID,
					&idp.UserID,
					&idp.IDPName,
					&idp.ProvidedUserID,
					&idp.ProvidedUsername,
					&idp.IDPType,
					&count,
				)
				if err != nil {
					return nil, err
				}
				idps = append(idps, idp)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-nwx6U", "Errors.Query.CloseRows")
			}

			return &LinkedIDPs{
				IDPs: idps,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
