package execution

import (
	"context"
	_ "embed"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
)

type Queries interface {
	ActiveInstances() []string
	GetNotifyUserByID(ctx context.Context, shouldTriggered bool, userID string) (*query.NotifyUser, error)
	TargetsByExecutionID(ctx context.Context, ids []string) (execution []*query.ExecutionTarget, err error)
	GetInstanceFeatures(ctx context.Context, cascade bool) (_ *query.InstanceFeatures, err error)

	InstanceByID(ctx context.Context, id string) (instance authz.Instance, err error)
}

type ExecutionsQueries struct {
	Queries
}

func NewExecutionsQueries(
	baseQueries Queries,
) *ExecutionsQueries {
	return &ExecutionsQueries{
		Queries: baseQueries,
	}
}

func (q *ExecutionsQueries) ActiveInstancesWithFeatureFlag(ctx context.Context) []string {
	return slices.DeleteFunc(q.Queries.ActiveInstances(), func(s string) bool {
		features, err := q.Queries.GetInstanceFeatures(ctx, true)
		if err != nil {
			return true
		}
		if features == nil || !features.Actions.Value {
			return true
		}
		return false
	})
}
