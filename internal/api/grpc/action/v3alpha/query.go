package action

import (
	"context"
	"strings"

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
		Endpoint: t.Endpoint,
	}

	switch t.TargetType {
	case domain.TargetTypeWebhook:
		target.TargetType = &action.Target_RestWebhook{RestWebhook: &action.SetRESTWebhook{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeCall:
		target.TargetType = &action.Target_RestCall{RestCall: &action.SetRESTCall{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeAsync:
		target.TargetType = &action.Target_RestAsync{RestAsync: &action.SetRESTAsync{}}
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
	case *action.SearchQuery_IncludeQuery:
		include, err := conditionToInclude(q.IncludeQuery.GetInclude())
		if err != nil {
			return nil, err
		}
		return query.NewIncludeSearchQuery(include)
	case *action.SearchQuery_TargetQuery:
		return query.NewTargetSearchQuery(q.TargetQuery.GetTargetId())
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
		return command.ExecutionFunctionCondition(t.Function.GetName()).ID(), nil
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
	targets := make([]*action.ExecutionTargetType, len(e.Targets))
	for i := range e.Targets {
		switch e.Targets[i].Type {
		case domain.ExecutionTargetTypeInclude:
			targets[i] = &action.ExecutionTargetType{Type: &action.ExecutionTargetType_Include{Include: executionIDToCondition(e.Targets[i].Target)}}
		case domain.ExecutionTargetTypeTarget:
			targets[i] = &action.ExecutionTargetType{Type: &action.ExecutionTargetType_Target{Target: e.Targets[i].Target}}
		case domain.ExecutionTargetTypeUnspecified:
			continue
		default:
			continue
		}
	}

	return &action.Execution{
		Details:   object.DomainToDetailsPb(&e.ObjectDetails),
		Condition: executionIDToCondition(e.ID),
		Targets:   targets,
	}
}

func executionIDToCondition(include string) *action.Condition {
	if strings.HasPrefix(include, domain.ExecutionTypeRequest.String()) {
		return includeRequestToCondition(strings.TrimPrefix(include, domain.ExecutionTypeRequest.String()))
	}
	if strings.HasPrefix(include, domain.ExecutionTypeResponse.String()) {
		return includeResponseToCondition(strings.TrimPrefix(include, domain.ExecutionTypeResponse.String()))
	}
	if strings.HasPrefix(include, domain.ExecutionTypeEvent.String()) {
		return includeEventToCondition(strings.TrimPrefix(include, domain.ExecutionTypeEvent.String()))
	}
	if strings.HasPrefix(include, domain.ExecutionTypeFunction.String()) {
		return includeFunctionToCondition(strings.TrimPrefix(include, domain.ExecutionTypeFunction.String()))
	}
	return nil
}

func includeRequestToCondition(id string) *action.Condition {
	switch strings.Count(id, "/") {
	case 2:
		return &action.Condition{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Method{Method: id}}}}
	case 1:
		return &action.Condition{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Service{Service: strings.TrimPrefix(id, "/")}}}}
	case 0:
		return &action.Condition{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_All{All: true}}}}
	default:
		return nil
	}
}
func includeResponseToCondition(id string) *action.Condition {
	switch strings.Count(id, "/") {
	case 2:
		return &action.Condition{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Method{Method: id}}}}
	case 1:
		return &action.Condition{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Service{Service: strings.TrimPrefix(id, "/")}}}}
	case 0:
		return &action.Condition{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_All{All: true}}}}
	default:
		return nil
	}
}

func includeEventToCondition(id string) *action.Condition {
	switch strings.Count(id, "/") {
	case 1:
		if strings.HasSuffix(id, command.EventGroupSuffix) {
			return &action.Condition{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Group{Group: strings.TrimSuffix(strings.TrimPrefix(id, "/"), command.EventGroupSuffix)}}}}
		} else {
			return &action.Condition{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Event{Event: strings.TrimPrefix(id, "/")}}}}
		}
	case 0:
		return &action.Condition{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_All{All: true}}}}
	default:
		return nil
	}
}

func includeFunctionToCondition(id string) *action.Condition {
	return &action.Condition{ConditionType: &action.Condition_Function{Function: &action.FunctionExecution{Name: strings.TrimPrefix(id, "/")}}}
}
