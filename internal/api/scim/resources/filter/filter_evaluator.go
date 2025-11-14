package filter

import (
	"reflect"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Evaluator struct {
	schema            schemas.ScimSchemaType
	attributeResolver *AttributeResolver
}

type SimpleValueEvaluationResult struct {
	PathSegments []string
	Value        reflect.Value
}

type FilteredValuesEvaluationResult struct {
	PathSegments []string
	Source       reflect.Value
	Matches      []*FilteredValuesEvaluationResultMatch
}

type FilteredValuesEvaluationResultMatch struct {
	// Value the selected value of the Element. Can be the Element itself or an attribute of the Element.
	Value       reflect.Value
	Element     reflect.Value
	SourceIndex int
}

type EvaluationResult interface{}

func NewEvaluator(schema schemas.ScimSchemaType) *Evaluator {
	return &Evaluator{
		schema:            schema,
		attributeResolver: newAttributeResolver(schema),
	}
}

func (e *Evaluator) Evaluate(v reflect.Value, path *Path) (EvaluationResult, error) {
	switch {
	case path == nil:
		return &SimpleValueEvaluationResult{
			PathSegments: nil,
			Value:        v,
		}, nil
	case path.AttrPath != nil:
		segments, value, err := e.attributeResolver.resolveAttrPath(v, path.AttrPath)
		if err != nil {
			return nil, err
		}

		return &SimpleValueEvaluationResult{
			PathSegments: segments,
			Value:        value,
		}, nil
	}

	return e.filterByValuePath(v, path.ValuePath)
}

func (e *Evaluator) filterByValuePath(v reflect.Value, valuePath *ValuePathWithSubAttr) (*FilteredValuesEvaluationResult, error) {
	segments, value, err := e.attributeResolver.resolveAttrPath(v, &valuePath.ValuePath.AttrPath)
	if err != nil {
		return nil, err
	}

	// filtering does only work on slices
	if value.Kind() != reflect.Slice {
		return nil, zerrors.ThrowInvalidArgument(nil, "SCIM-21x2", "value path target must be a slice")
	}

	// apply the filter
	var matches []*FilteredValuesEvaluationResultMatch
	for i := 0; i < value.Len(); i++ {
		element := value.Index(i)
		if match, err := e.evaluateOr(element, &valuePath.ValuePath.ValFilter); err != nil {
			return nil, err
		} else if match && valuePath.SubAttr != nil {
			subElement, err := e.attributeResolver.resolveField(element, *valuePath.SubAttr)
			if err != nil {
				return nil, err
			}

			matches = append(matches, &FilteredValuesEvaluationResultMatch{
				Value:       subElement,
				Element:     element,
				SourceIndex: i,
			})
		} else if match {
			matches = append(matches, &FilteredValuesEvaluationResultMatch{
				Value:       element,
				Element:     element,
				SourceIndex: i,
			})
		}
	}

	if valuePath.SubAttr != nil {
		segments = append(segments, *valuePath.SubAttr)
	}

	return &FilteredValuesEvaluationResult{
		PathSegments: segments,
		Source:       value,
		Matches:      matches,
	}, nil
}

func (e *Evaluator) evaluateOr(item reflect.Value, or *OrLogExp) (bool, error) {
	if match, err := e.evaluateAnd(item, &or.Left); err != nil {
		return false, err
	} else if match {
		return true, nil
	}

	if or.Right != nil {
		return e.evaluateOr(item, or.Right)
	}

	return false, nil
}

func (e *Evaluator) evaluateAnd(item reflect.Value, and *AndLogExp) (bool, error) {
	if match, err := e.evaluateValueAtom(item, &and.Left); err != nil {
		return false, err
	} else if !match {
		return false, nil
	}

	if and.Right != nil {
		return e.evaluateAnd(item, and.Right)
	}

	return true, nil
}

func (e *Evaluator) evaluateValueAtom(item reflect.Value, atom *ValueAtom) (bool, error) {
	switch {
	case atom.SubFilter != nil:
		return e.evaluateOr(item, &atom.SubFilter.OrExp)
	case atom.Negation != nil:
		if match, err := e.evaluateOr(item, &atom.Negation.OrExp); err != nil {
			return false, err
		} else {
			return !match, nil
		}
	case atom.AttrExp != nil:
		return e.evaluateAttrExp(item, atom.AttrExp)
	}

	// atom.ValuePath is not evaluated since nested valuePaths are not supported
	logging.Warn("Encountered unsupported nested value path")
	return false, zerrors.ThrowInvalidArgument(nil, "SCIM-21x7", "nested value paths are not supported")
}

func (e *Evaluator) evaluateAttrExp(item reflect.Value, attrExp *AttrExp) (bool, error) {
	if attrExp.UnaryCondition != nil {
		return e.evaluateUnaryCondition(item, attrExp.UnaryCondition)
	}

	return e.evaluateBinaryCondition(item, attrExp.BinaryCondition)
}

func (e *Evaluator) evaluateUnaryCondition(item reflect.Value, condition *UnaryCondition) (bool, error) {
	_, field, err := e.attributeResolver.resolveAttrPath(item, &condition.Left)
	if err != nil {
		return false, err
	}

	if field.Kind() == reflect.Ptr {
		return !field.IsZero(), nil
	}

	return !field.IsZero(), nil
}

func (e *Evaluator) evaluateBinaryCondition(item reflect.Value, condition *BinaryCondition) (bool, error) {
	_, field, err := e.attributeResolver.resolveAttrPath(item, &condition.Left)
	if err != nil {
		return false, err
	}

	return e.compareValues(field, &condition.Operator, &condition.Right)
}

func (e *Evaluator) compareValues(value reflect.Value, op *CompareOp, compValue *CompValue) (bool, error) {
	if (op.Equal) && compValue.Null {
		return value.IsNil() || value.IsZero(), nil
	}

	if op.NotEqual && compValue.Null {
		return !value.IsNil() && !value.IsZero(), nil
	}

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.compareUIntValue(value, op, compValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.compareIntValue(value, op, compValue)
	case reflect.Float64, reflect.Float32:
		return e.compareFloatValue(value, op, compValue)
	case reflect.String:
		return e.compareString(value, op, compValue)
	case reflect.Bool:
		return e.compareBool(value, op, compValue)

	// for complex attributes, compare the "value" by default as required by the rfc
	case reflect.Struct:
		if defaultValue, err := e.attributeResolver.resolveField(value, "value"); err == nil {
			return e.compareValues(defaultValue, op, compValue)
		}
	}

	logging.WithFields("kind", value.Kind()).Warn("Encountered unsupported nested value path")
	return false, zerrors.ThrowInvalidArgument(nil, "SCIM-23x8", "unsupported filter path value type")
}

func (e *Evaluator) compareUIntValue(value reflect.Value, op *CompareOp, compValue *CompValue) (bool, error) {
	if compValue.Int == nil {
		return false, zerrors.ThrowInvalidArgument(nil, "SCIM-23x8", "Invalid comparison type, expected integer")
	}

	intCompValue := uint64(*compValue.Int)
	v := value.Uint()
	switch {
	case op.Equal:
		return v == intCompValue, nil
	case op.NotEqual:
		return v != intCompValue, nil
	case op.GreaterThan:
		return v > intCompValue, nil
	case op.GreaterThanOrEqual:
		return v >= intCompValue, nil
	case op.LessThan:
		return v < intCompValue, nil
	case op.LessThanOrEqual:
		return v <= intCompValue, nil
	}

	return false, zerrors.ThrowInvalidArgumentf(nil, "SCIM-24x8", "Invalid comparison operator %s, not supported for number comparisons", op.String())
}

func (e *Evaluator) compareIntValue(value reflect.Value, op *CompareOp, compValue *CompValue) (bool, error) {
	if compValue.Int == nil {
		return false, zerrors.ThrowInvalidArgument(nil, "SCIM-23x9", "Invalid comparison type, expected integer")
	}

	intCompValue := int64(*compValue.Int)
	v := value.Int()
	switch {
	case op.Equal:
		return v == intCompValue, nil
	case op.NotEqual:
		return v != intCompValue, nil
	case op.GreaterThan:
		return v > intCompValue, nil
	case op.GreaterThanOrEqual:
		return v >= intCompValue, nil
	case op.LessThan:
		return v < intCompValue, nil
	case op.LessThanOrEqual:
		return v <= intCompValue, nil
	}

	return false, zerrors.ThrowInvalidArgumentf(nil, "SCIM-24x7", "Invalid comparison operator %s, not supported for number comparisons", op.String())
}

func (e *Evaluator) compareFloatValue(value reflect.Value, op *CompareOp, compValue *CompValue) (bool, error) {
	var floatCompValue float64
	switch {
	case compValue.Float != nil:
		floatCompValue = *compValue.Float
	case compValue.Int != nil:
		floatCompValue = float64(*compValue.Int)
	default:
		return false, zerrors.ThrowInvalidArgument(nil, "SCIM-24x1", "Invalid comparison type, expected number")
	}

	v := value.Float()
	switch {
	case op.Equal:
		return v == floatCompValue, nil
	case op.NotEqual:
		return v != floatCompValue, nil
	case op.GreaterThan:
		return v > floatCompValue, nil
	case op.GreaterThanOrEqual:
		return v >= floatCompValue, nil
	case op.LessThan:
		return v < floatCompValue, nil
	case op.LessThanOrEqual:
		return v <= floatCompValue, nil
	}

	return false, zerrors.ThrowInvalidArgumentf(nil, "SCIM-24x6", "Invalid comparison operator %s, not supported for number comparisons", op.String())
}

func (e *Evaluator) compareString(value reflect.Value, op *CompareOp, compValue *CompValue) (bool, error) {
	if compValue.StringValue == nil {
		return false, zerrors.ThrowInvalidArgument(nil, "SCIM-24x2", "Invalid comparison type, expected string")
	}

	strCompValue := *compValue.StringValue
	v := value.String()
	switch {
	case op.Equal:
		return v == strCompValue, nil
	case op.NotEqual:
		return v != strCompValue, nil
	case op.Contains:
		return strings.Contains(v, strCompValue), nil
	case op.StartsWith:
		return strings.HasPrefix(v, strCompValue), nil
	case op.EndsWith:
		return strings.HasSuffix(v, strCompValue), nil
	}

	return false, zerrors.ThrowInvalidArgumentf(nil, "SCIM-24x5", "Invalid comparison operator %s, not supported for string comparisons", op.String())
}

func (e *Evaluator) compareBool(value reflect.Value, op *CompareOp, compValue *CompValue) (bool, error) {
	if !compValue.BooleanTrue && !compValue.BooleanFalse {
		return false, zerrors.ThrowInvalidArgument(nil, "SCIM-24x3", "Invalid comparison type, expected bool")
	}

	if !op.Equal && !op.NotEqual {
		return false, zerrors.ThrowInvalidArgumentf(nil, "SCIM-24x4", "Invalid comparison operator %s, expected eq or ne", op.String())
	}

	return value.Bool() == (op.Equal && compValue.BooleanTrue), nil
}
