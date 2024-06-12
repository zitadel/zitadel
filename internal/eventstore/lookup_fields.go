package eventstore

import "golang.org/x/exp/constraints"

type LookupOperation struct {
	Add    *LookupField
	Remove map[LookupFieldType]*LookupValue
}

type LookupField struct {
	Aggregate *Aggregate
	IsUpsert  bool
	FieldName string
	Value     LookupValue
}

type LookupValue struct {
	valueType LookupValueType
	Value     any
}

func (v *LookupValue) IsNumeric() bool {
	return v.valueType == LookupValueTypeNumeric
}

func (v *LookupValue) IsText() bool {
	return v.valueType == LookupValueTypeString
}

type LookupValueType int8

const (
	LookupValueTypeString LookupValueType = iota
	LookupValueTypeNumeric
)

func UpsertLookupTextField[V ~string](aggregate *Aggregate, fieldName string, value V) *LookupOperation {
	return &LookupOperation{
		Add: &LookupField{
			Aggregate: aggregate,
			IsUpsert:  true,
			FieldName: fieldName,
			Value: LookupValue{
				valueType: LookupValueTypeString,
				Value:     value,
			},
		},
	}
}

func InsertLookupTextField[V ~string](aggregate *Aggregate, fieldName string, value V) *LookupOperation {
	return &LookupOperation{
		Add: &LookupField{
			Aggregate: aggregate,
			FieldName: fieldName,
			Value: LookupValue{
				valueType: LookupValueTypeString,
				Value:     value,
			},
		},
	}
}

func UpsertLookupNumericField[V constraints.Integer | constraints.Float](aggregate *Aggregate, fieldName string, value V) *LookupOperation {
	return &LookupOperation{
		Add: &LookupField{
			Aggregate: aggregate,
			IsUpsert:  true,
			FieldName: fieldName,
			Value: LookupValue{
				valueType: LookupValueTypeNumeric,
				Value:     value,
			},
		},
	}
}

func InsertLookupNumericField[V constraints.Integer | constraints.Float](aggregate *Aggregate, fieldName string, value V) *LookupOperation {
	return &LookupOperation{
		Add: &LookupField{
			Aggregate: aggregate,
			FieldName: fieldName,
			Value: LookupValue{
				valueType: LookupValueTypeNumeric,
				Value:     value,
			},
		},
	}
}

type LookupFieldType int8

const (
	LookupFieldTypeAggregateType LookupFieldType = iota
	LookupFieldTypeAggregateID
	LookupFieldTypeInstanceID
	LookupFieldTypeResourceOwner
	LookupFieldTypeFieldName
	LookupFieldTypeValue
)

func WithLookupFieldAggregateType(aggregateType string) *LookupValue {
	return &LookupValue{
		valueType: LookupValueTypeString,
		Value:     aggregateType,
	}
}

func WithLookupFieldAggregateID(aggregateID string) *LookupValue {
	return &LookupValue{
		valueType: LookupValueTypeString,
		Value:     aggregateID,
	}
}

func WithLookupFieldInstanceID(instanceID string) *LookupValue {
	return &LookupValue{
		valueType: LookupValueTypeString,
		Value:     instanceID,
	}
}

func WithLookupFieldResourceOwner(resourceOwner string) *LookupValue {
	return &LookupValue{
		valueType: LookupValueTypeString,
		Value:     resourceOwner,
	}
}

func WithLookupFieldFieldName(fieldName string) *LookupValue {
	return &LookupValue{
		valueType: LookupValueTypeString,
		Value:     fieldName,
	}
}

func WithLookupFieldFieldNumericValue[V constraints.Integer | constraints.Float](value V) *LookupValue {
	return &LookupValue{
		valueType: LookupValueTypeNumeric,
		Value:     value,
	}
}

func WithLookupFieldFieldTextValue[V ~string](value V) *LookupValue {
	return &LookupValue{
		valueType: LookupValueTypeString,
		Value:     value,
	}
}
