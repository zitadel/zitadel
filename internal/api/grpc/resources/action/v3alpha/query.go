package action

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/durationpb"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
)

const (
	conditionIDAllSegmentCount                    = 0
	conditionIDRequestResponseServiceSegmentCount = 1
	conditionIDRequestResponseMethodSegmentCount  = 2
	conditionIDEventGroupSegmentCount             = 1
)

func (s *Server) GetTarget(ctx context.Context, req *action.GetTargetRequest) (*action.GetTargetResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}

	resp, err := s.query.GetTargetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &action.GetTargetResponse{
		Target: targetToPb(resp),
	}, nil
}

type InstanceContext interface {
	GetInstanceId() string
	GetInstanceDomain() string
}

type Context interface {
	GetOwner() InstanceContext
}

func (s *Server) SearchTargets(ctx context.Context, req *action.SearchTargetsRequest) (*action.SearchTargetsResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	queries, err := s.searchTargetsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchTargets(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.SearchTargetsResponse{
		Result:  targetsToPb(resp.Targets),
		Details: resource_object.ToSearchDetailsPb(queries.SearchRequest, resp.SearchResponse),
	}, nil
}

func (s *Server) SearchExecutions(ctx context.Context, req *action.SearchExecutionsRequest) (*action.SearchExecutionsResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	queries, err := s.searchExecutionsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchExecutions(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.SearchExecutionsResponse{
		Result:  executionsToPb(resp.Executions),
		Details: resource_object.ToSearchDetailsPb(queries.SearchRequest, resp.SearchResponse),
	}, nil
}

func targetsToPb(targets []*query.Target) []*action.GetTarget {
	t := make([]*action.GetTarget, len(targets))
	for i, target := range targets {
		t[i] = targetToPb(target)
	}
	return t
}

func targetToPb(t *query.Target) *action.GetTarget {
	target := &action.GetTarget{
		Details: resource_object.DomainToDetailsPb(&t.ObjectDetails, object.OwnerType_OWNER_TYPE_INSTANCE, t.ResourceOwner),
		Config: &action.Target{
			Name:     t.Name,
			Timeout:  durationpb.New(t.Timeout),
			Endpoint: t.Endpoint,
		},
		SigningKey: t.SigningKey,
	}
	switch t.TargetType {
	case domain.TargetTypeWebhook:
		target.Config.TargetType = &action.Target_RestWebhook{RestWebhook: &action.SetRESTWebhook{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeCall:
		target.Config.TargetType = &action.Target_RestCall{RestCall: &action.SetRESTCall{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeAsync:
		target.Config.TargetType = &action.Target_RestAsync{RestAsync: &action.SetRESTAsync{}}
	default:
		target.Config.TargetType = nil
	}
	return target
}

func (s *Server) searchTargetsRequestToModel(req *action.SearchTargetsRequest) (*query.TargetSearchQueries, error) {
	offset, limit, asc, err := resource_object.SearchQueryPbToQuery(s.systemDefaults, req.Query)
	if err != nil {
		return nil, err
	}
	queries, err := targetQueriesToQuery(req.Filters)
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

func targetQueriesToQuery(queries []*action.TargetSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, qry := range queries {
		q[i], err = targetQueryToQuery(qry)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func targetQueryToQuery(filter *action.TargetSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *action.TargetSearchFilter_TargetNameFilter:
		return targetNameQueryToQuery(q.TargetNameFilter)
	case *action.TargetSearchFilter_InTargetIdsFilter:
		return targetInTargetIdsQueryToQuery(q.InTargetIdsFilter)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func targetNameQueryToQuery(q *action.TargetNameFilter) (query.SearchQuery, error) {
	return query.NewTargetNameSearchQuery(resource_object.TextMethodPbToQuery(q.Method), q.GetTargetName())
}

func targetInTargetIdsQueryToQuery(q *action.InTargetIDsFilter) (query.SearchQuery, error) {
	return query.NewTargetInIDsSearchQuery(q.GetTargetIds())
}

// targetFieldNameToSortingColumn defaults to the creation date because this ensures deterministic pagination
func targetFieldNameToSortingColumn(field *action.TargetFieldName) query.Column {
	if field == nil {
		return query.TargetColumnCreationDate
	}
	switch *field {
	case action.TargetFieldName_TARGET_FIELD_NAME_UNSPECIFIED:
		return query.TargetColumnID
	case action.TargetFieldName_TARGET_FIELD_NAME_ID:
		return query.TargetColumnID
	case action.TargetFieldName_TARGET_FIELD_NAME_CREATED_DATE:
		return query.TargetColumnCreationDate
	case action.TargetFieldName_TARGET_FIELD_NAME_CHANGED_DATE:
		return query.TargetColumnChangeDate
	case action.TargetFieldName_TARGET_FIELD_NAME_NAME:
		return query.TargetColumnName
	case action.TargetFieldName_TARGET_FIELD_NAME_TARGET_TYPE:
		return query.TargetColumnTargetType
	case action.TargetFieldName_TARGET_FIELD_NAME_URL:
		return query.TargetColumnURL
	case action.TargetFieldName_TARGET_FIELD_NAME_TIMEOUT:
		return query.TargetColumnTimeout
	case action.TargetFieldName_TARGET_FIELD_NAME_INTERRUPT_ON_ERROR:
		return query.TargetColumnInterruptOnError
	default:
		return query.TargetColumnCreationDate
	}
}

// executionFieldNameToSortingColumn defaults to the creation date because this ensures deterministic pagination
func executionFieldNameToSortingColumn(field *action.ExecutionFieldName) query.Column {
	if field == nil {
		return query.ExecutionColumnCreationDate
	}
	switch *field {
	case action.ExecutionFieldName_EXECUTION_FIELD_NAME_UNSPECIFIED:
		return query.ExecutionColumnID
	case action.ExecutionFieldName_EXECUTION_FIELD_NAME_ID:
		return query.ExecutionColumnID
	case action.ExecutionFieldName_EXECUTION_FIELD_NAME_CREATED_DATE:
		return query.ExecutionColumnCreationDate
	case action.ExecutionFieldName_EXECUTION_FIELD_NAME_CHANGED_DATE:
		return query.ExecutionColumnChangeDate
	default:
		return query.ExecutionColumnCreationDate
	}
}

func (s *Server) searchExecutionsRequestToModel(req *action.SearchExecutionsRequest) (*query.ExecutionSearchQueries, error) {
	offset, limit, asc, err := resource_object.SearchQueryPbToQuery(s.systemDefaults, req.Query)
	if err != nil {
		return nil, err
	}
	queries, err := executionQueriesToQuery(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.ExecutionSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: executionFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func executionQueriesToQuery(queries []*action.ExecutionSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = executionQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func executionQueryToQuery(searchQuery *action.ExecutionSearchFilter) (query.SearchQuery, error) {
	switch q := searchQuery.Filter.(type) {
	case *action.ExecutionSearchFilter_InConditionsFilter:
		return inConditionsQueryToQuery(q.InConditionsFilter)
	case *action.ExecutionSearchFilter_ExecutionTypeFilter:
		return executionTypeToQuery(q.ExecutionTypeFilter)
	case *action.ExecutionSearchFilter_IncludeFilter:
		include, err := conditionToInclude(q.IncludeFilter.GetInclude())
		if err != nil {
			return nil, err
		}
		return query.NewIncludeSearchQuery(include)
	case *action.ExecutionSearchFilter_TargetFilter:
		return query.NewTargetSearchQuery(q.TargetFilter.GetTargetId())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func executionTypeToQuery(q *action.ExecutionTypeFilter) (query.SearchQuery, error) {
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

func inConditionsQueryToQuery(q *action.InConditionsFilter) (query.SearchQuery, error) {
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

func executionsToPb(executions []*query.Execution) []*action.GetExecution {
	e := make([]*action.GetExecution, len(executions))
	for i, execution := range executions {
		e[i] = executionToPb(execution)
	}
	return e
}

func executionToPb(e *query.Execution) *action.GetExecution {
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

	return &action.GetExecution{
		Details: resource_object.DomainToDetailsPb(&e.ObjectDetails, object.OwnerType_OWNER_TYPE_INSTANCE, e.ResourceOwner),
		Execution: &action.Execution{
			Targets: targets,
		},
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
	case conditionIDRequestResponseMethodSegmentCount:
		return &action.Condition{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Method{Method: id}}}}
	case conditionIDRequestResponseServiceSegmentCount:
		return &action.Condition{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_Service{Service: strings.TrimPrefix(id, "/")}}}}
	case conditionIDAllSegmentCount:
		return &action.Condition{ConditionType: &action.Condition_Request{Request: &action.RequestExecution{Condition: &action.RequestExecution_All{All: true}}}}
	default:
		return nil
	}
}
func includeResponseToCondition(id string) *action.Condition {
	switch strings.Count(id, "/") {
	case conditionIDRequestResponseMethodSegmentCount:
		return &action.Condition{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Method{Method: id}}}}
	case conditionIDRequestResponseServiceSegmentCount:
		return &action.Condition{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_Service{Service: strings.TrimPrefix(id, "/")}}}}
	case conditionIDAllSegmentCount:
		return &action.Condition{ConditionType: &action.Condition_Response{Response: &action.ResponseExecution{Condition: &action.ResponseExecution_All{All: true}}}}
	default:
		return nil
	}
}

func includeEventToCondition(id string) *action.Condition {
	switch strings.Count(id, "/") {
	case conditionIDEventGroupSegmentCount:
		if strings.HasSuffix(id, command.EventGroupSuffix) {
			return &action.Condition{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Group{Group: strings.TrimSuffix(strings.TrimPrefix(id, "/"), command.EventGroupSuffix)}}}}
		} else {
			return &action.Condition{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Event{Event: strings.TrimPrefix(id, "/")}}}}
		}
	case conditionIDAllSegmentCount:
		return &action.Condition{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_All{All: true}}}}
	default:
		return nil
	}
}

func includeFunctionToCondition(id string) *action.Condition {
	return &action.Condition{ConditionType: &action.Condition_Function{Function: &action.FunctionExecution{Name: strings.TrimPrefix(id, "/")}}}
}
