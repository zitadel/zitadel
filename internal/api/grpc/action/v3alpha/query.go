package action

import (
	"context"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
)

func (s *Server) ListTargets(ctx context.Context, req *action.ListTargetsRequest) (*action.ListTargetsResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	queries, err := listTargetsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchTargets(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.ListTargetsResponse{
		Result:  targetsToPb(resp.Targets),
		Details: object.ToListDetails(resp.SearchResponse),
	}, nil
}

func listTargetsRequestToModel(req *action.ListTargetsRequest) (*query.TargetSearchQueries, error) {
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

func targetFieldNameToSortingColumn(field action.TargetFieldName) query.Column {
	switch field {
	case action.TargetFieldName_FIELD_NAME_UNSPECIFIED:
		return query.TargetColumnID
	case action.TargetFieldName_FIELD_NAME_ID:
		return query.TargetColumnID
	case action.TargetFieldName_FIELD_NAME_CREATION_DATE:
		return query.TargetColumnCreationDate
	case action.TargetFieldName_FIELD_NAME_CHANGE_DATE:
		return query.TargetColumnChangeDate
	case action.TargetFieldName_FIELD_NAME_NAME:
		return query.TargetColumnName
	case action.TargetFieldName_FIELD_NAME_TARGET_TYPE:
		return query.TargetColumnTargetType
	case action.TargetFieldName_FIELD_NAME_URL:
		return query.TargetColumnURL
	case action.TargetFieldName_FIELD_NAME_TIMEOUT:
		return query.TargetColumnTimeout
	case action.TargetFieldName_FIELD_NAME_ASYNC:
		return query.TargetColumnAsync
	case action.TargetFieldName_FIELD_NAME_INTERRUPT_ON_ERROR:
		return query.TargetColumnInterruptOnError
	default:
		return query.TargetColumnID
	}
}

func targetQueriesToQuery(queries []*action.TargetSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = targetQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func targetQueryToQuery(query *action.TargetSearchQuery) (query.SearchQuery, error) {
	switch q := query.Query.(type) {
	case *action.TargetSearchQuery_TargetNameQuery:
		return targetNameQueryToQuery(q.TargetNameQuery)
	case *action.TargetSearchQuery_InTargetIdsQuery:
		return targetInTargetIdsQueryToQuery(q.InTargetIdsQuery)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func targetNameQueryToQuery(q *action.TargetNameQuery) (query.SearchQuery, error) {
	return query.NewTargetNameSearchQuery(object.TextMethodToQuery(q.Method), q.GetTargetName())
}

func targetInTargetIdsQueryToQuery(q *action.InTargetIDsQuery) (query.SearchQuery, error) {
	return query.NewTargetInIDsSearchQuery(q.GetTargetIds())
}

func (s *Server) GetTargetByID(ctx context.Context, req *action.GetTargetByIDRequest) (_ *action.GetTargetByIDResponse, err error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	resp, err := s.query.GetTargetByID(ctx, req.GetTargetId())
	if err != nil {
		return nil, err
	}
	return &action.GetTargetByIDResponse{
		Target: targetToPb(resp),
	}, nil
}

func targetsToPb(targets []*query.Target) []*action.Target {
	t := make([]*action.Target, len(targets))
	for i, target := range targets {
		t[i] = targetToPb(target)
	}
	return t
}

func targetToPb(t *query.Target) *action.Target {
	target := &action.Target{
		Details:  object.DomainToDetailsPb(&t.ObjectDetails),
		TargetId: t.ID,
		Name:     t.Name,
		Timeout:  durationpb.New(t.Timeout),
	}
	if t.Async {
		target.ExecutionType = &action.Target_IsAsync{IsAsync: t.Async}
	}
	if t.InterruptOnError {
		target.ExecutionType = &action.Target_InterruptOnError{InterruptOnError: t.InterruptOnError}
	}

	switch t.TargetType {
	case domain.TargetTypeWebhook:
		target.TargetType = &action.Target_RestWebhook{RestWebhook: &action.SetRESTWebhook{Url: t.URL}}
	case domain.TargetTypeRequestResponse:
		target.TargetType = &action.Target_RestRequestResponse{RestRequestResponse: &action.SetRESTRequestResponse{Url: t.URL}}
	default:
		target.TargetType = nil
	}
	return target
}

func (s *Server) ListExecutions(ctx context.Context, req *action.ListExecutionsRequest) (*action.ListExecutionsResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	queries, err := listExecutionsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchExecutions(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.ListExecutionsResponse{
		Result:  executionsToPb(resp.Executions),
		Details: object.ToListDetails(resp.SearchResponse),
	}, nil
}

func listExecutionsRequestToModel(req *action.ListExecutionsRequest) (*query.ExecutionSearchQueries, error) {
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

func executionQueriesToQuery(queries []*action.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = executionQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func executionQueryToQuery(searchQuery *action.SearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *action.SearchQuery_InConditionsQuery:
		return inConditionsQueryToQuery(q.InConditionsQuery)
	case *action.SearchQuery_ExecutionTypeQuery:
		return executionTypeToQuery(q.ExecutionTypeQuery)
	case *action.SearchQuery_TargetQuery:
		return query.NewExecutionTargetSearchQuery(q.TargetQuery.GetTargetId())
	case *action.SearchQuery_IncludeQuery:
		return query.NewExecutionIncludeSearchQuery(q.IncludeQuery.GetInclude())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func executionTypeToQuery(q *action.ExecutionTypeQuery) (query.SearchQuery, error) {
	switch q.ExecutionType {
	case action.ExecutionType_EXECUTION_TYPE_UNSPECIFIED:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeUnspecified)
	case action.ExecutionType_EXECUTION_TYPE_REQUEST:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeRequest)
	case action.ExecutionType_EXECUTION_TYPE_RESPONSE:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeResponse)
	case action.ExecutionType_EXECUTION_TYPE_EVENT:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeEvent)
	case action.ExecutionType_EXECUTION_TYPE_FUNCTION:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeFunction)
	default:
		return query.NewExecutionTypeSearchQuery(domain.ExecutionTypeUnspecified)
	}
}

func inConditionsQueryToQuery(q *action.InConditionsQuery) (query.SearchQuery, error) {
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

func conditionToID(q *action.Condition) (string, error) {
	switch t := q.GetConditionType().(type) {
	case *action.Condition_Request:
		cond := &command.ExecutionAPICondition{
			Method:  t.Request.GetMethod(),
			Service: t.Request.GetService(),
			All:     t.Request.GetAll(),
		}
		return cond.ID(domain.ExecutionTypeRequest), nil
	case *action.Condition_Response:
		cond := &command.ExecutionAPICondition{
			Method:  t.Response.GetMethod(),
			Service: t.Response.GetService(),
			All:     t.Response.GetAll(),
		}
		return cond.ID(domain.ExecutionTypeResponse), nil
	case *action.Condition_Event:
		cond := &command.ExecutionEventCondition{
			Event: t.Event.GetEvent(),
			Group: t.Event.GetGroup(),
			All:   t.Event.GetAll(),
		}
		return cond.ID(), nil
	case *action.Condition_Function:
		return t.Function, nil
	default:
		return "", zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func executionsToPb(executions []*query.Execution) []*action.Execution {
	e := make([]*action.Execution, len(executions))
	for i, execution := range executions {
		e[i] = executionToPb(execution)
	}
	return e
}

func executionToPb(e *query.Execution) *action.Execution {
	var targets, includes []string
	if len(e.Targets) > 0 {
		targets = e.Targets
	}
	if len(e.Includes) > 0 {
		includes = e.Includes
	}
	return &action.Execution{
		Details:     object.DomainToDetailsPb(&e.ObjectDetails),
		ExecutionId: e.ID,
		Targets:     targets,
		Includes:    includes,
	}
}
