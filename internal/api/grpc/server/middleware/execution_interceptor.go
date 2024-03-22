package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/query"
	exec_rp "github.com/zitadel/zitadel/internal/repository/execution"
)

func ExecutionHandler(queries *query.Queries) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		request, err := executeTargetsForGRPCFullMethod(ctx, queries, info.FullMethod, req, domain.ExecutionTypeRequest)
		if err != nil {
			return nil, err
		}

		resp, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}

		return executeTargetsForGRPCFullMethod(ctx, queries, info.FullMethod, resp, domain.ExecutionTypeResponse)
	}
}

func executeTargetsForGRPCFullMethod(ctx context.Context, queries *query.Queries, fullMethod string, req interface{}, executionType domain.ExecutionType) (interface{}, error) {
	request := req
	typeQuery, err := query.NewExecutionTypeSearchQuery(executionType)
	if err != nil {
		return nil, err
	}

	idQuery, err := inIDsQuery(executionType, fullMethod)
	if err != nil {
		return nil, err
	}

	execs, err := queries.SearchExecutions(ctx, &query.ExecutionSearchQueries{Queries: []query.SearchQuery{typeQuery, idQuery}})
	if err != nil {
		return nil, err
	}

	for _, exec := range execs.Executions {
		exectargets, err := queries.ExecutionTargets(ctx, exec.ID)
		if err != nil {
			return nil, err
		}

		targetIDsQuery, err := query.NewTargetInIDsSearchQuery(exectargets.Targets())
		if err != nil {
			return nil, err
		}

		targets, err := queries.SearchTargets(ctx, &query.TargetSearchQueries{Queries: []query.SearchQuery{targetIDsQuery}})
		if err != nil {
			return nil, err
		}

		ctxData := authz.GetCtxData(ctx)
		info := &execution.ContextInfo{
			FullMethod: fullMethod,
			InstanceID: authz.GetInstance(ctx).InstanceID(),
			ProjectID:  ctxData.ProjectID,
			OrgID:      ctxData.OrgID,
			UserID:     ctxData.UserID,
			Request:    request,
		}

		request, err = execution.CallTargets(ctx, targets.Targets, info)
		if err != nil {
			return nil, err
		}
	}
	return request, nil
}

func inIDsQuery(t domain.ExecutionType, fullMethod string) (query.SearchQuery, error) {
	return query.NewExecutionInIDsSearchQuery(
		[]string{
			exec_rp.IDAll(t),
			exec_rp.ID(t, serviceFromFullMethod(fullMethod)),
			exec_rp.ID(t, fullMethod),
		},
	)
}

func serviceFromFullMethod(s string) string {
	parts := strings.Split(s, "/")
	return parts[1]
}
