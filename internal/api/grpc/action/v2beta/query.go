package action

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
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

func (s *Server) ListTargets(ctx context.Context, req *action.ListTargetsRequest) (*action.ListTargetsResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	queries, err := s.ListTargetsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchTargets(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.ListTargetsResponse{
		Result:     targetsToPb(resp.Targets),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}, nil
}

func (s *Server) ListExecutions(ctx context.Context, req *action.ListExecutionsRequest) (*action.ListExecutionsResponse, error) {
	if err := checkActionsEnabled(ctx); err != nil {
		return nil, err
	}
	queries, err := s.ListExecutionsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchExecutions(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.ListExecutionsResponse{
		Result:     executionsToPb(resp.Executions),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
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
		Id:         t.ObjectDetails.ID,
		Name:       t.Name,
		Timeout:    durationpb.New(t.Timeout),
		Endpoint:   t.Endpoint,
		SigningKey: t.SigningKey,
	}
	switch t.TargetType {
	case domain.TargetTypeWebhook:
		target.TargetType = &action.Target_RestWebhook{RestWebhook: &action.RESTWebhook{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeCall:
		target.TargetType = &action.Target_RestCall{RestCall: &action.RESTCall{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeAsync:
		target.TargetType = &action.Target_RestAsync{RestAsync: &action.RESTAsync{}}
	default:
		target.TargetType = nil
	}

	if !t.ObjectDetails.EventDate.IsZero() {
		target.ChangeDate = timestamppb.New(t.ObjectDetails.EventDate)
	}
	if !t.ObjectDetails.CreationDate.IsZero() {
		target.CreationDate = timestamppb.New(t.ObjectDetails.CreationDate)
	}
	return target
}

func (s *Server) ListTargetsRequestToModel(req *action.ListTargetsRequest) (*query.TargetSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
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
	return query.NewTargetNameSearchQuery(filter.TextMethodPbToQuery(q.Method), q.GetTargetName())
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

func (s *Server) ListExecutionsRequestToModel(req *action.ListExecutionsRequest) (*query.ExecutionSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
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

	exec := &action.Execution{
		Condition: executionIDToCondition(e.ID),
		Targets:   targets,
	}
	if !e.ObjectDetails.EventDate.IsZero() {
		exec.ChangeDate = timestamppb.New(e.ObjectDetails.EventDate)
	}
	if !e.ObjectDetails.CreationDate.IsZero() {
		exec.CreationDate = timestamppb.New(e.ObjectDetails.CreationDate)
	}
	return exec
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
