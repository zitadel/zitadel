package repository

import (
	"database/sql"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// SearchQuery defines the which and how data are queried
type SearchQuery struct {
	Columns               eventstore.Columns
	Limit                 uint64
	Desc                  bool
	Filters               [][]*Filter
	Tx                    *sql.Tx
	AllowTimeTravel       bool
	AwaitOpenTransactions bool
}

// Filter represents all fields needed to compare a field of an event with a value
type Filter struct {
	Field     Field
	Value     interface{}
	Operation Operation
}

// Operation defines how fields are compared
type Operation int32

const (
	// OperationEquals compares two values for equality
	OperationEquals Operation = iota + 1
	// OperationGreater compares if the given values is greater than the stored one
	OperationGreater
	// OperationLess compares if the given values is less than the stored one
	OperationLess
	//OperationIn checks if a stored value matches one of the passed value list
	OperationIn
	//OperationJSONContains checks if a stored value matches the given json
	OperationJSONContains
	//OperationNotIn checks if a stored value does not match one of the passed value list
	OperationNotIn

	operationCount
)

// Field is the representation of a field from the event
type Field int32

const (
	//FieldAggregateType represents the aggregate type field
	FieldAggregateType Field = iota + 1
	//FieldAggregateID represents the aggregate id field
	FieldAggregateID
	//FieldSequence represents the sequence field
	FieldSequence
	//FieldResourceOwner represents the resource owner field
	FieldResourceOwner
	//FieldInstanceID represents the instance id field
	FieldInstanceID
	//FieldEditorService represents the editor service field
	FieldEditorService
	//FieldEditorUser represents the editor user field
	FieldEditorUser
	//FieldEventType represents the event type field
	FieldEventType
	//FieldEventData represents the event data field
	FieldEventData
	//FieldCreationDate represents the creation date field
	FieldCreationDate
	// FieldPosition represents the field of the global sequence
	FieldPosition

	fieldCount
)

// NewFilter is used in tests. Use searchQuery.*Filter() instead
func NewFilter(field Field, value interface{}, operation Operation) *Filter {
	return &Filter{
		Field:     field,
		Value:     value,
		Operation: operation,
	}
}

// Validate checks if the fields of the filter have valid values
func (f *Filter) Validate() error {
	if f == nil {
		return errors.ThrowPreconditionFailed(nil, "REPO-z6KcG", "filter is nil")
	}
	if f.Field <= 0 || f.Field >= fieldCount {
		return errors.ThrowPreconditionFailed(nil, "REPO-zw62U", "field not definded")
	}
	if f.Value == nil {
		return errors.ThrowPreconditionFailed(nil, "REPO-GJ9ct", "no value definded")
	}
	if f.Operation <= 0 || f.Operation >= operationCount {
		return errors.ThrowPreconditionFailed(nil, "REPO-RrQTy", "operation not definded")
	}
	return nil
}

func QueryFromBuilder(builder *eventstore.SearchQueryBuilder) (*SearchQuery, error) {
	if builder == nil ||
		builder.GetColumns().Validate() != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-4m9gs", "builder invalid")
	}

	filters := make([][]*Filter, len(builder.GetQueries()))

	var builderFilters []*Filter
	for _, f := range []func(builder *eventstore.SearchQueryBuilder) *Filter{
		instanceIDFilterFromBuilder,
		editorUserFilter,
		resourceOwnerFilter,
		positionAfterFilter,
	} {
		filter := f(builder)
		if filter == nil {
			continue
		}
		if err := filter.Validate(); err != nil {
			return nil, err
		}
		builderFilters = append(builderFilters, filter)
	}

	for i, query := range builder.GetQueries() {
		for _, f := range []func(query *eventstore.SearchQuery) *Filter{
			instanceIDFilterFromQuery,
			aggregateTypeFilter,
			aggregateIDFilter,
			eventTypeFilter,
			eventDataFilter,
			eventSequenceGreaterFilter,
			eventSequenceLessFilter,
			excludedInstanceIDFilter,
			creationDateAfterFilter,
		} {
			filter := f(query)
			if filter == nil {
				continue
			}
			if err := filter.Validate(); err != nil {
				return nil, err
			}
			filters[i] = append(filters[i], filter)
		}
		filters[i] = append(builderFilters, filters[i]...)
	}

	if len(filters) == 0 {
		filters = append(filters, builderFilters)
	}

	return &SearchQuery{
		Columns:               builder.GetColumns(),
		Limit:                 builder.GetLimit(),
		Desc:                  builder.GetDesc(),
		Filters:               filters,
		Tx:                    builder.GetTx(),
		AllowTimeTravel:       builder.GetAllowTimeTravel(),
		AwaitOpenTransactions: builder.GetAwaitOpenTransactions(),
	}, nil
}

func aggregateIDFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetAggregateIDs()) < 1 {
		return nil
	}
	if len(query.GetAggregateIDs()) == 1 {
		return NewFilter(FieldAggregateID, query.GetAggregateIDs()[0], OperationEquals)
	}
	return NewFilter(FieldAggregateID, database.TextArray[string](query.GetAggregateIDs()), OperationIn)
}

func eventTypeFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetEventTypes()) < 1 {
		return nil
	}
	if len(query.GetEventTypes()) == 1 {
		return NewFilter(FieldEventType, query.GetEventTypes()[0], OperationEquals)
	}
	eventTypes := make(database.TextArray[eventstore.EventType], len(query.GetEventTypes()))
	for i, eventType := range query.GetEventTypes() {
		eventTypes[i] = eventType
	}
	return NewFilter(FieldEventType, eventTypes, OperationIn)
}

func aggregateTypeFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetAggregateTypes()) < 1 {
		return nil
	}
	if len(query.GetAggregateTypes()) == 1 {
		return NewFilter(FieldAggregateType, query.GetAggregateTypes()[0], OperationEquals)
	}
	aggregateTypes := make(database.TextArray[eventstore.AggregateType], len(query.GetAggregateTypes()))
	for i, aggregateType := range query.GetAggregateTypes() {
		aggregateTypes[i] = aggregateType
	}
	return NewFilter(FieldAggregateType, aggregateTypes, OperationIn)
}

func eventSequenceGreaterFilter(query *eventstore.SearchQuery) *Filter {
	if query.GetEventSequenceGreater() == 0 {
		return nil
	}
	sortOrder := OperationGreater
	if query.Builder().GetDesc() {
		sortOrder = OperationLess
	}
	return NewFilter(FieldSequence, query.GetEventSequenceGreater(), sortOrder)
}

func eventSequenceLessFilter(query *eventstore.SearchQuery) *Filter {
	if query.GetEventSequenceLess() == 0 {
		return nil
	}
	sortOrder := OperationLess
	if query.Builder().GetDesc() {
		sortOrder = OperationGreater
	}
	return NewFilter(FieldSequence, query.GetEventSequenceLess(), sortOrder)
}

func instanceIDFilterFromQuery(query *eventstore.SearchQuery) *Filter {
	if query.GetInstanceID() != "" {
		return NewFilter(FieldInstanceID, query.GetInstanceID(), OperationEquals)
	}
	if query.Builder().GetInstanceID() != "" {
		NewFilter(FieldInstanceID, query.Builder().GetInstanceID(), OperationEquals)
	}
	return nil
}

func excludedInstanceIDFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetExcludedInstanceIDs()) == 0 {
		return nil
	}
	return NewFilter(FieldInstanceID, database.TextArray[string](query.GetExcludedInstanceIDs()), OperationNotIn)
}

func resourceOwnerFilter(builder *eventstore.SearchQueryBuilder) *Filter {
	if builder.GetResourceOwner() == "" {
		return nil
	}
	return NewFilter(FieldResourceOwner, builder.GetResourceOwner(), OperationEquals)
}

func editorUserFilter(builder *eventstore.SearchQueryBuilder) *Filter {
	if builder.GetEditorUser() == "" {
		return nil
	}
	return NewFilter(FieldEditorUser, builder.GetEditorUser(), OperationEquals)
}

func instanceIDFilterFromBuilder(builder *eventstore.SearchQueryBuilder) *Filter {
	if builder.GetInstanceID() == "" {
		return nil
	}
	return NewFilter(FieldInstanceID, builder.GetInstanceID(), OperationEquals)
}

func positionAfterFilter(query *eventstore.SearchQueryBuilder) *Filter {
	if query.GetPositionAfter() == 0 {
		return nil
	}
	return NewFilter(FieldPosition, query.GetPositionAfter(), OperationGreater)
}

func creationDateAfterFilter(query *eventstore.SearchQuery) *Filter {
	if query.GetCreationDateAfter().IsZero() {
		return nil
	}
	return NewFilter(FieldCreationDate, query.GetCreationDateAfter(), OperationGreater)
}

func eventDataFilter(query *eventstore.SearchQuery) *Filter {
	if len(query.GetEventData()) == 0 {
		return nil
	}
	return NewFilter(FieldEventData, query.GetEventData(), OperationJSONContains)
}
