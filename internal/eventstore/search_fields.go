package eventstore

import (
	"encoding/json"

	"golang.org/x/exp/constraints"

	"github.com/zitadel/zitadel/internal/zerrors"
)

// SearchOperation if the definition of the operation to be executed on the search index
type SearchOperation struct {
	// Set a field in the search index
	// if [SearchField.UpsertConflictFields] are set the field will be updated if the conflict fields match
	// if no [SearchField.UpsertConflictFields] are set the field will be inserted
	Set *SearchField
	// Remove fields from the search index using the map as `AND`ed conditions
	Remove map[SearchFieldType]any
}

type SearchResult struct {
	Aggregate Aggregate
	Object    SearchObject
	FieldName string
	ValueType SearchValueType
	// Value is either a string or a number
	// use the helper functions ([NumericResultValue] or [TextResultValue]) to convert the value to the correct type
	Value []byte
}

// NumericResultValue marshals the value to the given type
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

// TextResultValue marshals the value to the given type
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
	// Type of the object
	Type string
	// ID of the object
	ID string
	// Revision of the object, if an object evolves the revision should be increased
	// analog to current projection versioning
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

// SetSearchTextField sets a text field in the search index
// if conflictFields are set the field will be updated if the conflict fields match
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

// SetSearchNumericField sets a text field in the search index
// if conflictFields are set the field will be updated if the conflict fields match
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

// RemoveSearchFields removes fields from the search index using the map as `AND`ed conditions
func RemoveSearchFields(clause map[SearchFieldType]any) *SearchOperation {
	return &SearchOperation{
		Remove: clause,
	}
}

// RemoveSearchFieldsByAggregate removes fields from the search index using the aggregate as `AND`ed conditions
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

// RemoveSearchFieldsByAggregateAndObject removes fields from the search index using the aggregate and object as `AND`ed conditions
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

// RemoveSearchFieldsByAggregateAndObjectAndField removes fields from the search index using the aggregate, object and field as `AND`ed conditions
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
