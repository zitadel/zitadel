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
		request, err := executeTargetsForGRPCFullMethod(ctx, queries, info.FullMethod, req, nil, domain.ExecutionTypeRequest)
		if err != nil {
			return nil, err
		}

		resp, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}

		return executeTargetsForGRPCFullMethod(ctx, queries, info.FullMethod, req, resp, domain.ExecutionTypeResponse)
	}
}

type ExecutionQueries interface {
	ExecutionTargetsRequestResponse(ctx context.Context, fullMethod, service, all string) (execution *query.ExecutionTargets, err error)
	SearchTargets(ctx context.Context, queries *query.TargetSearchQueries) (targets *query.Targets, err error)
}

func executeTargetsForGRPCFullMethod(
	ctx context.Context,
	queries ExecutionQueries,
	fullMethod string,
	req interface{},
	resp interface{},
	executionType domain.ExecutionType,
) (interface{}, error) {
	request := req

	exectargets, err := queries.ExecutionTargetsRequestResponse(ctx, exec_rp.ID(executionType, fullMethod), exec_rp.ID(executionType, serviceFromFullMethod(fullMethod)), exec_rp.IDAll(executionType))
	if err != nil {
		return nil, err
	}

	targetIDsQuery, err := query.NewTargetInIDsSearchQuery(exectargets.Targets)
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
		Response:   resp,
	}

	return execution.CallTargets(ctx, targets.Targets, info)
}

func serviceFromFullMethod(s string) string {
	parts := strings.Split(s, "/")
	return parts[1]
}
