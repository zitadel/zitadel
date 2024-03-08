package execution

import (
	"context"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func (s *Server) ListTargets(ctx context.Context, req *execution.ListTargetsRequest) (*execution.ListTargetsResponse, error) {
	queries, err := listTargetsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchTargets(ctx, queries, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &execution.ListTargetsResponse{
		Result:  targetsToPb(resp.Targets),
		Details: object.ToListDetails(resp.SearchResponse),
	}, nil
}

func listTargetsRequestToModel(req *execution.ListTargetsRequest) (*query.TargetSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := targetQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.TargetSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: targetFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func targetFieldNameToSortingColumn(field execution.TargetFieldName) query.Column {
	switch field {
	case execution.TargetFieldName_FIELD_NAME_UNSPECIFIED:
		return query.TargetColumnID
	case execution.TargetFieldName_FIELD_NAME_ID:
		return query.TargetColumnID
	case execution.TargetFieldName_FIELD_NAME_CREATION_DATE:
		return query.TargetColumnCreationDate
	case execution.TargetFieldName_FIELD_NAME_CHANGE_DATE:
		return query.TargetColumnChangeDate
	case execution.TargetFieldName_FIELD_NAME_NAME:
		return query.TargetColumnName
	case execution.TargetFieldName_FIELD_NAME_TARGET_TYPE:
		return query.TargetColumnTargetType
	case execution.TargetFieldName_FIELD_NAME_URL:
		return query.TargetColumnURL
	case execution.TargetFieldName_FIELD_NAME_TIMEOUT:
		return query.TargetColumnTimeout
	case execution.TargetFieldName_FIELD_NAME_ASYNC:
		return query.TargetColumnAsync
	case execution.TargetFieldName_FIELD_NAME_INTERRUPT_ON_ERROR:
		return query.TargetColumnInterruptOnError
	default:
		return query.TargetColumnID
	}
}

func targetQueriesToQuery(queries []*execution.TargetSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = targetQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func targetQueryToQuery(query *execution.TargetSearchQuery) (query.SearchQuery, error) {
	switch q := query.Query.(type) {
	case *execution.TargetSearchQuery_TargetNameQuery:
		return targetNameQueryToQuery(q.TargetNameQuery)
	case *execution.TargetSearchQuery_InTargetIdsQuery:
		return targetInTargetIdsQueryToQuery(q.InTargetIdsQuery)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func targetNameQueryToQuery(q *execution.TargetNameQuery) (query.SearchQuery, error) {
	return query.NewTargetNameSearchQuery(object.TextMethodToQuery(q.Method), q.GetTargetName())
}

func targetInTargetIdsQueryToQuery(q *execution.InTargetIDsQuery) (query.SearchQuery, error) {
	return query.NewTargetInIDsSearchQuery(q.GetTargetIds())
}

func (s *Server) GetTargetByID(ctx context.Context, req *execution.GetTargetByIDRequest) (_ *execution.GetTargetByIDResponse, err error) {
	resp, err := s.query.GetTargetByID(ctx, req.GetTargetId(), authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &execution.GetTargetByIDResponse{
		Target: targetToPb(resp),
	}, nil
}

func targetsToPb(targets []*query.Target) []*execution.Target {
	t := make([]*execution.Target, len(targets))
	for i, target := range targets {
		t[i] = targetToPb(target)
	}
	return t
}

func targetToPb(t *query.Target) *execution.Target {
	target := &execution.Target{
		Details:  object.DomainToDetailsPb(t.ObjectDetails),
		TargetId: t.ID,
		Name:     t.Name,
		Timeout:  durationpb.New(t.Timeout()),
	}
	if t.Async {
		target.ExecutionType = &execution.Target_IsAsync{IsAsync: t.Async}
	}
	if t.InterruptOnError {
		target.ExecutionType = &execution.Target_InterruptOnError{InterruptOnError: t.InterruptOnError}
	}

	switch t.TargetType {
	case domain.TargetTypeWebhook:
		target.TargetType = &execution.Target_RestWebhook{RestWebhook: &execution.SetRESTWebhook{Url: t.URL}}
	case domain.TargetTypeRequestResponse:
		target.TargetType = &execution.Target_RestRequestResponse{RestRequestResponse: &execution.SetRESTRequestResponse{Url: t.URL}}
	case domain.TargetTypeUnspecified:
		target.TargetType = nil
	default:
		target.TargetType = nil
	}
	return target
}

func (s *Server) ListExecutions(ctx context.Context, req *execution.ListExecutionsRequest) (*execution.ListExecutionsResponse, error) {
	queries, err := listExecutionsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchExecutions(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &execution.ListExecutionsResponse{
		Result:  executionsToPb(resp.Executions),
		Details: object.ToListDetails(resp.SearchResponse),
	}, nil
}

func listExecutionsRequestToModel(req *execution.ListExecutionsRequest) (*query.ExecutionSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := executionQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.ExecutionSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func executionQueriesToQuery(queries []*execution.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = executionQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func executionQueryToQuery(searchQuery *execution.SearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *execution.SearchQuery_ConditionQuery:
		return conditionQueryToQuery(q.ConditionQuery)
	case *execution.SearchQuery_InConditionsQuery:
		return inConditionsQueryToQuery(q.InConditionsQuery)
	case *execution.SearchQuery_ExecutionTypeQuery:
		return executionTypeToQuery(q.ExecutionTypeQuery)
	case *execution.SearchQuery_TargetQuery:
		return query.NewExecutionTargetSearchQuery(q.TargetQuery.GetTargetId())
	case *execution.SearchQuery_IncludeQuery:
		return query.NewExecutionIncludeSearchQuery(q.IncludeQuery.GetInclude())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func executionTypeToQuery(q *execution.ExecutionTypeQuery) (query.SearchQuery, error) {
	switch q.ExecutionType {
	case execution.ExecutionType_EXECUTION_TYPE_UNSPECIFIED:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeUnspecified)
	case execution.ExecutionType_EXECUTION_TYPE_REQUEST:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeRequest)
	case execution.ExecutionType_EXECUTION_TYPE_RESPONSE:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeResponse)
	case execution.ExecutionType_EXECUTION_TYPE_EVENT:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeEvent)
	case execution.ExecutionType_EXECUTION_TYPE_FUNCTION:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeFunction)
	default:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeUnspecified)
	}
}

func inConditionsQueryToQuery(q *execution.InConditionsQuery) (query.SearchQuery, error) {
	values := make([]string, len(q.GetConditions()))
	for i, condition := range q.GetConditions() {
		id, err := conditionToID(condition)
		if err != nil {
			return nil, err
		}
		values[i] = id
	}
	return query.NewExecutionInIDsSearchQuery(values)
}

func conditionQueryToQuery(q *execution.ConditionQuery) (query.SearchQuery, error) {
	id, err := conditionToID(q.GetCondition())
	if err != nil {
		return nil, err
	}
	return query.NewExecutionIDSearchQuery(id)
}

func conditionToID(q *execution.SetConditions) (string, error) {
	switch t := q.GetConditionType().(type) {
	case *execution.SetConditions_Request:
		cond := &command.ExecutionAPICondition{
			Method:  t.Request.GetMethod(),
			Service: t.Request.GetService(),
			All:     t.Request.GetAll(),
		}
		return cond.ID(domain.ExecutionTypeRequest), nil
	case *execution.SetConditions_Response:
		cond := &command.ExecutionAPICondition{
			Method:  t.Response.GetMethod(),
			Service: t.Response.GetService(),
			All:     t.Response.GetAll(),
		}
		return cond.ID(domain.ExecutionTypeResponse), nil
	case *execution.SetConditions_Event:
		cond := &command.ExecutionEventCondition{
			Event: t.Event.GetEvent(),
			Group: t.Event.GetGroup(),
			All:   t.Event.GetAll(),
		}
		return cond.ID(), nil
	case *execution.SetConditions_Function:
		return t.Function, nil
	default:
		return "", zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func executionsToPb(executions []*query.Execution) []*execution.Execution {
	e := make([]*execution.Execution, len(executions))
	for i, execution := range executions {
		e[i] = executionToPb(execution)
	}
	return e
}

func executionToPb(e *query.Execution) *execution.Execution {
	var targets, includes []string
	if len(e.Targets) > 0 {
		targets = e.Targets
	}
	if len(e.Includes) > 0 {
		includes = e.Includes
	}
	return &execution.Execution{
		Details:     object.DomainToDetailsPb(e.ObjectDetails),
		ExecutionId: e.ID,
		Targets:     targets,
		Includes:    includes,
	}
}
