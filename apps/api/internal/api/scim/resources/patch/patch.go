package patch

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type OperationRequest struct {
	Schemas    []schemas.ScimSchemaType `json:"Schemas"`
	Operations []*Operation             `json:"Operations"`
}

type Operation struct {
	Operation    OperationType   `json:"op"`
	Path         *filter.Path    `json:"path"`
	Value        json.RawMessage `json:"value"`
	valueIsArray bool
}

type OperationCollection []*Operation

type OperationType string

const (
	OperationTypeAdd     OperationType = "add"
	OperationTypeRemove  OperationType = "remove"
	OperationTypeReplace OperationType = "replace"

	fieldNamePrimary = "Primary"
	fieldNameValue   = "Value"
)

type ResourcePatcher interface {
	FilterEvaluator() *filter.Evaluator
	Added(attributePath []string) error
	Replaced(attributePath []string) error
	Removed(attributePath []string) error
}

func (req *OperationRequest) GetSchemas() []schemas.ScimSchemaType {
	return req.Schemas
}

func (req *OperationRequest) Validate() error {
	for _, op := range req.Operations {
		if err := op.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (op *Operation) validate() error {
	// ignore the casing, as some scim clients send these capitalized
	op.Operation = OperationType(strings.ToLower(string(op.Operation)))

	if !op.Operation.isValid() {
		return serrors.ThrowInvalidValue(zerrors.ThrowInvalidArgumentf(nil, "SCIM-opty1", "Patch op %s not supported", op.Operation))
	}

	// json deserialization initializes this field if an empty string is provided
	// to not special case this in the further code,
	// set it to nil here.
	if op.Path.IsZero() {
		op.Path = nil
	}

	op.valueIsArray = strings.HasPrefix(strings.TrimPrefix(string(op.Value), " "), "[")
	return nil
}

func (ops OperationCollection) Apply(patcher ResourcePatcher, value interface{}) error {
	for _, op := range ops {
		if err := op.validate(); err != nil {
			return err
		}

		if err := op.apply(patcher, value); err != nil {
			return err
		}
	}

	return nil
}

func (op *Operation) apply(patcher ResourcePatcher, value interface{}) error {
	switch op.Operation {
	case OperationTypeRemove:
		return applyRemovePatch(patcher, op, value)
	case OperationTypeReplace:
		return applyReplacePatch(patcher, op, value)
	case OperationTypeAdd:
		return applyAddPatch(patcher, op, value)
	}

	return zerrors.ThrowInvalidArgumentf(nil, "SCIM-opty3", "SCIM patch: Invalid operation %v", op.Operation)
}

func (o OperationType) isValid() bool {
	switch o {
	case OperationTypeAdd, OperationTypeRemove, OperationTypeReplace:
		return true
	default:
		return false
	}
}

func flattenAndApplyPatchOperations(patcher ResourcePatcher, op *Operation, value interface{}) error {
	ops, err := flattenPatchOperations(op)
	if err != nil {
		return err
	}

	for _, flattenedOperation := range ops {
		if err = flattenedOperation.apply(patcher, value); err != nil {
			return err
		}
	}

	return nil
}

// flattenPatchOperations flattens patch operations without a path
// it converts an op { "operation": "add", "value": { "path1": "value1", "path2": "value2" } }
// into [ { "operation": "add", "path": "path1", "value": "value1" }, { "operation": "add", "path": "path2", "value": "value2" } ]
func flattenPatchOperations(op *Operation) ([]*Operation, error) {
	if op.Path != nil {
		panic("Only operations without a path can be flattened")
	}

	patches := map[string]json.RawMessage{}
	if err := json.Unmarshal(op.Value, &patches); err != nil {
		logging.WithError(err).Error("SCIM: Invalid patch value while flattening")
		return nil, zerrors.ThrowInvalidArgument(err, "SCIM-ioyl1", "Invalid patch value")
	}

	result := make([]*Operation, 0, len(patches))
	for path, value := range patches {
		result = append(result, &Operation{
			Operation: op.Operation,
			Path: &filter.Path{
				AttrPath: &filter.AttrPath{
					AttrName: path,
				},
			},
			Value:        value,
			valueIsArray: strings.HasPrefix(string(value), "["),
		})
	}

	return result, nil
}

// unmarshalPatchValuesSlice unmarshal the raw json value (a scalar value, object or array) into a new slice
func unmarshalPatchValuesSlice(elementTypePtr reflect.Type, value json.RawMessage, valueIsArray bool) (reflect.Value, error) {
	if elementTypePtr.Kind() != reflect.Ptr {
		logging.Panicf("elementType must be a pointer to a struct, but is %s", elementTypePtr.Name())
		return reflect.Value{}, nil
	}

	if !valueIsArray {
		newElement := reflect.New(elementTypePtr.Elem())
		if err := unmarshalPatchValue(value, newElement); err != nil {
			return reflect.Value{}, err
		}

		newSlice := reflect.MakeSlice(reflect.SliceOf(elementTypePtr), 1, 1)
		newSlice.Index(0).Set(newElement)
		return newSlice, nil
	}

	newSlicePtr := reflect.New(reflect.SliceOf(elementTypePtr))
	newSlice := newSlicePtr.Elem()
	if err := json.Unmarshal(value, newSlicePtr.Interface()); err != nil {
		logging.WithError(err).Error("SCIM: Invalid patch values")
		return reflect.Value{}, zerrors.ThrowInvalidArgument(err, "SCIM-opxx8", "Invalid patch values")
	}
	return newSlice, nil
}

func unmarshalPatchValue(newValue json.RawMessage, targetElement reflect.Value) error {
	if targetElement.Kind() != reflect.Ptr {
		targetElement = targetElement.Addr()
	}

	if targetElement.IsNil() {
		targetElement.Set(reflect.New(targetElement.Type().Elem()))
	}

	if err := json.Unmarshal(newValue, targetElement.Interface()); err != nil {
		logging.WithError(err).Error("SCIM: Invalid patch value")
		return zerrors.ThrowInvalidArgument(err, "SCIM-opty9", "Invalid patch value")
	}

	return nil
}

// ensureSinglePrimary ensures the modification on a slice results in max one primary element.
// modifiedSlice contains the patched slice.
// modifiedElementsSlice contains only the modified elements.
// if a new element has Primary set to true and an existing is also Primary, the existing Primary flag is set to false.
// returns an error if multiple modifiedElements have a primary value of true.
func ensureSinglePrimary(modifiedSlice reflect.Value, modifiedElementsSlice []reflect.Value, modifiedElementIndexes map[int]bool) error {
	if len(modifiedElementsSlice) == 0 {
		return nil
	}

	hasPrimary, err := isAnyPrimary(modifiedElementsSlice)
	if err != nil || !hasPrimary {
		return err
	}

	for i := 0; i < modifiedSlice.Len(); i++ {
		if mod, ok := modifiedElementIndexes[i]; ok && mod {
			continue
		}

		sliceElement := modifiedSlice.Index(i)
		if sliceElement.Kind() == reflect.Ptr {
			sliceElement = sliceElement.Elem()
		}

		sliceElementPrimaryField := sliceElement.FieldByName(fieldNamePrimary)
		if !sliceElementPrimaryField.IsValid() || !sliceElementPrimaryField.Bool() {
			continue
		}

		sliceElementPrimaryField.SetBool(false)

		// we can stop at the first primary,
		// since there can only be one primary in a slice.
		return nil
	}

	return nil
}

func isAnyPrimary(elements []reflect.Value) (bool, error) {
	foundPrimary := false
	for _, element := range elements {
		if !isPrimary(element) {
			continue
		}

		if foundPrimary {
			return true, zerrors.ThrowInvalidArgument(nil, "SCIM-1d23", "Cannot add multiple primary values in one patch operation")
		}

		foundPrimary = true
	}

	return foundPrimary, nil
}

func isPrimary(element reflect.Value) bool {
	if element.Kind() == reflect.Ptr {
		element = element.Elem()
	}

	if element.Kind() != reflect.Struct {
		return false
	}

	primaryField := element.FieldByName(fieldNamePrimary)
	return primaryField.IsValid() && primaryField.Bool()
}
