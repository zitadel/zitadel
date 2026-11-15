package query

import (
	"context"
	"crypto/sha256"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var groupExternalIDsTable = table{name: projection.GroupExternalIDsProjectionTable, instanceIDCol: projection.GroupExternalIDsColumnInstanceID}

func groupExternalIDColumn(name string) Column {
	return Column{name: name, table: groupExternalIDsTable}
}

// GroupExternalIDs returns only bindings belonging to the caller's provisioning source and organization.
func (q *Queries) GroupExternalIDs(ctx context.Context, orgID, scope string, groupIDs []string) (map[string]string, error) {
	result := make(map[string]string, len(groupIDs))
	if len(groupIDs) == 0 {
		return result, nil
	}
	stmt, args, err := sq.Select(projection.GroupExternalIDsColumnGroupID, projection.GroupExternalIDsColumnValue).
		From(projection.GroupExternalIDsProjectionTable).
		Where(sq.Eq{projection.GroupExternalIDsColumnInstanceID: authz.GetInstance(ctx).InstanceID(), projection.GroupExternalIDsColumnResourceOwner: orgID, projection.GroupExternalIDsColumnScope: scope, projection.GroupExternalIDsColumnGroupID: groupIDs}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-GEXT1", "Errors.Query.SQLStatement")
	}
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var id, value string
			if err := rows.Scan(&id, &value); err != nil {
				return err
			}
			result[id] = value
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return rows.Close()
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-GEXT2", "Errors.Internal")
	}
	return result, nil
}

func NewGroupExternalIDSearchQuery(scope, value string) (SearchQuery, error) {
	instanceQuery, err := NewColumnComparisonQuery(groupExternalIDColumn(projection.GroupExternalIDsColumnInstanceID), GroupColumnInstanceID, ColumnEquals)
	if err != nil {
		return nil, err
	}
	orgQuery, err := NewColumnComparisonQuery(groupExternalIDColumn(projection.GroupExternalIDsColumnResourceOwner), GroupColumnResourceOwner, ColumnEquals)
	if err != nil {
		return nil, err
	}
	scopeQuery, err := NewTextQuery(groupExternalIDColumn(projection.GroupExternalIDsColumnScope), scope, TextEquals)
	if err != nil {
		return nil, err
	}
	valueQuery, err := NewTextQuery(groupExternalIDColumn(projection.GroupExternalIDsColumnValue), value, TextEquals)
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256([]byte(value))
	// Squirrel treats []byte as scalar bytea equality rather than expanding an IN list.
	digestQuery, err := NewListQuery(groupExternalIDColumn(projection.GroupExternalIDsColumnDigest), digest[:], ListIn)
	if err != nil {
		return nil, err
	}
	sub, err := NewSubSelect(groupExternalIDColumn(projection.GroupExternalIDsColumnGroupID), []SearchQuery{instanceQuery, orgQuery, scopeQuery, digestQuery, valueQuery})
	if err != nil {
		return nil, err
	}
	return NewListQuery(GroupColumnID, sub, ListIn)
}
