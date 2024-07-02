package eventstore

// FieldOperation if the definition of the operation to be executed on the field
type FieldOperation struct {
	// Set a field in the field table
	// if [SearchField.UpsertConflictFields] are set the field will be updated if the conflict fields match
	// if no [SearchField.UpsertConflictFields] are set the field will be inserted
	Set *Field
	// Remove fields using the map as `AND`ed conditions
	Remove map[FieldType]any
}

type SearchResult struct {
	Aggregate Aggregate
	Object    Object
	FieldName string
	// Value represents the stored value
	// use the Unmarshal method to parse the value to the desired type
	Value interface {
		// Unmarshal parses the value to ptr
		Unmarshal(ptr any) error
	}
}

// // NumericResultValue marshals the value to the given type

type Object struct {
	// Type of the object
	Type string
	// ID of the object
	ID string
	// Revision of the object, if an object evolves the revision should be increased
	// analog to current projection versioning
	Revision uint8
}

type Field struct {
	Aggregate            *Aggregate
	Object               Object
	UpsertConflictFields []FieldType
	FieldName            string
	Value                Value
}

type Value struct {
	Value any
	// MustBeUnique defines if the field must be unique
	// This field will replace unique constraints in the future
	// If MustBeUnique is true the value must be a primitive type
	MustBeUnique bool
	// ShouldIndex defines if the field should be indexed
	// If the field is not indexed it can not be used in search queries
	// If ShouldIndex is true the value must be a primitive type
	ShouldIndex bool
}

type SearchValueType int8

const (
	SearchValueTypeString SearchValueType = iota
	SearchValueTypeNumeric
)

// SetSearchField sets the field based on the defined parameters
// if conflictFields are set the field will be updated if the conflict fields match
func SetField(aggregate *Aggregate, object Object, fieldName string, value *Value, conflictFields ...FieldType) *FieldOperation {
	return &FieldOperation{
		Set: &Field{
			Aggregate:            aggregate,
			Object:               object,
			UpsertConflictFields: conflictFields,
			FieldName:            fieldName,
			Value:                *value,
		},
	}
}

// RemoveSearchFields removes fields using the map as `AND`ed conditions
func RemoveSearchFields(clause map[FieldType]any) *FieldOperation {
	return &FieldOperation{
		Remove: clause,
	}
}

// RemoveSearchFieldsByAggregate removes fields using the aggregate as `AND`ed conditions
func RemoveSearchFieldsByAggregate(aggregate *Aggregate) *FieldOperation {
	return &FieldOperation{
		Remove: map[FieldType]any{
			FieldTypeInstanceID:    aggregate.InstanceID,
			FieldTypeResourceOwner: aggregate.ResourceOwner,
			FieldTypeAggregateType: aggregate.Type,
			FieldTypeAggregateID:   aggregate.ID,
		},
	}
}

// RemoveSearchFieldsByAggregateAndObject removes fields using the aggregate and object as `AND`ed conditions
func RemoveSearchFieldsByAggregateAndObject(aggregate *Aggregate, object Object) *FieldOperation {
	return &FieldOperation{
		Remove: map[FieldType]any{
			FieldTypeInstanceID:     aggregate.InstanceID,
			FieldTypeResourceOwner:  aggregate.ResourceOwner,
			FieldTypeAggregateType:  aggregate.Type,
			FieldTypeAggregateID:    aggregate.ID,
			FieldTypeObjectType:     object.Type,
			FieldTypeObjectID:       object.ID,
			FieldTypeObjectRevision: object.Revision,
		},
	}
}

// RemoveSearchFieldsByAggregateAndObjectAndField removes fields using the aggregate, object and field as `AND`ed conditions
func RemoveSearchFieldsByAggregateAndObjectAndField(aggregate *Aggregate, object Object, field string) *FieldOperation {
	return &FieldOperation{
		Remove: map[FieldType]any{
			FieldTypeInstanceID:     aggregate.InstanceID,
			FieldTypeResourceOwner:  aggregate.ResourceOwner,
			FieldTypeAggregateType:  aggregate.Type,
			FieldTypeAggregateID:    aggregate.ID,
			FieldTypeObjectType:     object.Type,
			FieldTypeObjectID:       object.ID,
			FieldTypeObjectRevision: object.Revision,
			FieldTypeFieldName:      field,
		},
	}
}

type FieldType int8

const (
	FieldTypeAggregateType FieldType = iota
	FieldTypeAggregateID
	FieldTypeInstanceID
	FieldTypeResourceOwner
	FieldTypeObjectType
	FieldTypeObjectID
	FieldTypeObjectRevision
	FieldTypeFieldName
	FieldTypeValue
)
