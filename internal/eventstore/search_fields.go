package eventstore

import (
	"encoding/json"

	"golang.org/x/exp/constraints"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type SearchOperation struct {
	Set    *SearchField
	Remove map[SearchFieldType]any
}

type SearchResult struct {
	Aggregate Aggregate
	Object    SearchObject
	FieldName string
	ValueType SearchValueType
	Value     []byte
}

func NumericResultValue[T constraints.Integer | constraints.Float](res *SearchResult) (T, error) {
	if res.ValueType != SearchValueTypeNumeric {
		return 0, zerrors.ThrowInvalidArgument(nil, "EVENT-JBhtu", "value is not numeric")
	}
	var value T
	if err := json.Unmarshal(res.Value, &value); err != nil {
		return 0, zerrors.ThrowInternal(err, "EVENT-2M9fs", "unable to unmarshal numeric value")
	}
	return value, nil
}

func TextResultValue[T ~string](res *SearchResult) (T, error) {
	if res.ValueType != SearchValueTypeString {
		return "", zerrors.ThrowInvalidArgument(nil, "EVENT-ywqg5", "value is not text")
	}
	return T(string(res.Value)), nil
}

type SearchField struct {
	Aggregate            *Aggregate
	Object               SearchObject
	UpsertConflictFields []SearchFieldType
	FieldName            string
	Value                SearchValue
}

type SearchObject struct {
	Type     string
	ID       string
	Revision uint8
}

type SearchValue struct {
	ValueType SearchValueType
	Value     any
}

func (v *SearchValue) IsNumeric() bool {
	return v.ValueType == SearchValueTypeNumeric
}

func (v *SearchValue) IsText() bool {
	return v.ValueType == SearchValueTypeString
}

type SearchValueType int8

const (
	SearchValueTypeString SearchValueType = iota
	SearchValueTypeNumeric
)

func SetSearchTextField[V ~string](aggregate *Aggregate, object SearchObject, fieldName string, value V, conflictFields ...SearchFieldType) *SearchOperation {
	return &SearchOperation{
		Set: &SearchField{
			Aggregate:            aggregate,
			Object:               object,
			UpsertConflictFields: conflictFields,
			FieldName:            fieldName,
			Value: SearchValue{
				ValueType: SearchValueTypeString,
				Value:     value,
			},
		},
	}
}

func SetSearchNumericField[V constraints.Integer | constraints.Float](aggregate *Aggregate, object SearchObject, fieldName string, value V, conflictFields ...SearchFieldType) *SearchOperation {
	return &SearchOperation{
		Set: &SearchField{
			Aggregate:            aggregate,
			Object:               object,
			UpsertConflictFields: conflictFields,
			FieldName:            fieldName,
			Value: SearchValue{
				ValueType: SearchValueTypeNumeric,
				Value:     value,
			},
		},
	}
}

func RemoveSearchFields(clause map[SearchFieldType]any) *SearchOperation {
	return &SearchOperation{
		Remove: clause,
	}
}

func RemoveSearchFieldsByAggregate(aggregate *Aggregate) *SearchOperation {
	return &SearchOperation{
		Remove: map[SearchFieldType]any{
			SearchFieldTypeInstanceID:    aggregate.InstanceID,
			SearchFieldTypeResourceOwner: aggregate.ResourceOwner,
			SearchFieldTypeAggregateType: aggregate.Type,
			SearchFieldTypeAggregateID:   aggregate.ID,
		},
	}
}

func RemoveSearchFieldsByAggregateAndObject(aggregate *Aggregate, object SearchObject) *SearchOperation {
	return &SearchOperation{
		Remove: map[SearchFieldType]any{
			SearchFieldTypeInstanceID:     aggregate.InstanceID,
			SearchFieldTypeResourceOwner:  aggregate.ResourceOwner,
			SearchFieldTypeAggregateType:  aggregate.Type,
			SearchFieldTypeAggregateID:    aggregate.ID,
			SearchFieldTypeObjectType:     object.Type,
			SearchFieldTypeObjectID:       object.ID,
			SearchFieldTypeObjectRevision: object.Revision,
		},
	}
}

func RemoveSearchFieldsByAggregateAndObjectAndField(aggregate *Aggregate, object SearchObject, field string) *SearchOperation {
	return &SearchOperation{
		Remove: map[SearchFieldType]any{
			SearchFieldTypeInstanceID:     aggregate.InstanceID,
			SearchFieldTypeResourceOwner:  aggregate.ResourceOwner,
			SearchFieldTypeAggregateType:  aggregate.Type,
			SearchFieldTypeAggregateID:    aggregate.ID,
			SearchFieldTypeObjectType:     object.Type,
			SearchFieldTypeObjectID:       object.ID,
			SearchFieldTypeObjectRevision: object.Revision,
			SearchFieldTypeFieldName:      field,
		},
	}
}

type SearchFieldType int8

const (
	SearchFieldTypeAggregateType SearchFieldType = iota
	SearchFieldTypeAggregateID
	SearchFieldTypeInstanceID
	SearchFieldTypeResourceOwner
	SearchFieldTypeObjectType
	SearchFieldTypeObjectID
	SearchFieldTypeObjectRevision
	SearchFieldTypeFieldName
	SearchFieldTypeTextValue
	SearchFieldTypeNumericValue
)

func WithSearchFieldAggregateType(aggregateType string) *SearchValue {
	return &SearchValue{
		ValueType: SearchValueTypeString,
		Value:     aggregateType,
	}
}

func WithSearchFieldAggregateID(aggregateID string) *SearchValue {
	return &SearchValue{
		ValueType: SearchValueTypeString,
		Value:     aggregateID,
	}
}

func WithSearchFieldInstanceID(instanceID string) *SearchValue {
	return &SearchValue{
		ValueType: SearchValueTypeString,
		Value:     instanceID,
	}
}

func WithSearchFieldResourceOwner(resourceOwner string) *SearchValue {
	return &SearchValue{
		ValueType: SearchValueTypeString,
		Value:     resourceOwner,
	}
}

func WithSearchFieldFieldName(fieldName string) *SearchValue {
	return &SearchValue{
		ValueType: SearchValueTypeString,
		Value:     fieldName,
	}
}

func WithSearchFieldFieldNumericValue[V constraints.Integer | constraints.Float](value V) *SearchValue {
	return &SearchValue{
		ValueType: SearchValueTypeNumeric,
		Value:     value,
	}
}

func WithSearchFieldFieldTextValue[V ~string](value V) *SearchValue {
	return &SearchValue{
		ValueType: SearchValueTypeString,
		Value:     value,
	}
}
