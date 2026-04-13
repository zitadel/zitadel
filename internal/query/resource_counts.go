package query

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed resource_counts_list.sql
	resourceCountsListQuery string
)

type ResourceCount struct {
	ID         int // Primary key, used for pagination
	InstanceID string
	TableName  string
	ParentType domain.CountParentType
	ParentID   string
	Resource   string
	UpdatedAt  time.Time
	Amount     int
}

// ListResourceCounts retrieves all resource counts.
// It supports pagination using lastID and limit parameters.
//
// TODO: Currently only a proof of concept, filters may be implemented later if required.
func (q *Queries) ListResourceCounts(ctx context.Context, lastID, limit int) (result []ResourceCount, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var count ResourceCount
			err := rows.Scan(
				&count.ID,
				&count.InstanceID,
				&count.TableName,
				&count.ParentType,
				&count.ParentID,
				&count.Resource,
				&count.UpdatedAt,
				&count.Amount)
			if err != nil {
				return zerrors.ThrowInternal(err, "QUERY-2f4g5", "Errors.Internal")
			}
			result = append(result, count)
		}
		return nil
	}, resourceCountsListQuery, lastID, limit)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-3f4g5", "Errors.Internal")
	}
	return result, nil
}
