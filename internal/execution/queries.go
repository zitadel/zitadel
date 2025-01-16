package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/query"
)

type Queries interface {
	SearchExecutions(ctx context.Context, queries *query.ExecutionSearchQueries) (executions *query.Executions, err error)
}
