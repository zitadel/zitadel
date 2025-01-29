package patch

import (
	"encoding/json"
	"reflect"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func applyReplacePatch(patcher ResourcePatcher, op *Operation, value interface{}) error {
	if op.Path == nil {
		return flattenAndApplyPatchOperations(patcher, op, value)
	}

	result, err := patcher.FilterEvaluator().Evaluate(reflect.ValueOf(value), op.Path)
	if err != nil {
		return serrors.ThrowInvalidPath(zerrors.ThrowInvalidArgument(err, "SCIM-i2o3", "Failed to evaluate path"))
	}

	switch filterResult := result.(type) {
	case *filter.SimpleValueEvaluationResult:
		return applyReplacePatchSimple(patcher, filterResult, op.Value, op.valueIsArray)
	case *filter.FilteredValuesEvaluationResult:
		return applyReplacePatchFiltered(patcher, filterResult, op.Value)
	}

	logging.Errorf("SCIM replace patch: unsupported filter type %T", result)
	return serrors.ThrowInvalidPath(zerrors.ThrowInvalidArgument(err, "SCIM-optu9", "Invalid patch path"))
}

func applyReplacePatchSimple(patcher ResourcePatcher, evaluationResult *filter.SimpleValueEvaluationResult, newValueRaw json.RawMessage, valueIsArray bool) error {
	// patch value is an array
	// or it is a scalar or an object but the target is a slice
	// unmarshal it as a slice and set it on the target, clearing all existing entries
	if valueIsArray || evaluationResult.Value.Kind() == reflect.Slice {
		return applyReplacePatchSimpleSlice(patcher, evaluationResult, newValueRaw, valueIsArray)
	}

	// patch value and target is a scalar or object
	if err := unmarshalPatchValue(newValueRaw, evaluationResult.Value); err != nil {
		return err
	}

	return patcher.Replaced(evaluationResult.PathSegments)
}

func applyReplacePatchSimpleSlice(patcher ResourcePatcher, evaluationResult *filter.SimpleValueEvaluationResult, newValueRaw json.RawMessage, valueIsArray bool) error {
	if evaluationResult.Value.Kind() != reflect.Slice {
		return zerrors.ThrowInvalidArgument(nil, "SCIM-E345X", "Cannot apply array patch value to single value attribute")
	}

	values, err := unmarshalPatchValuesSlice(evaluationResult.Value.Type().Elem(), newValueRaw, valueIsArray)
	if err != nil {
		return err
	}

	evaluationResult.Value.Set(values)
	modifiedIndexes := make(map[int]bool, values.Len())
	for i := 0; i < values.Len(); i++ {
		modifiedIndexes[i] = true
	}

	if err = ensureSinglePrimaryAdded(values, values, modifiedIndexes); err != nil {
		return err
	}

	return patcher.Replaced(evaluationResult.PathSegments)
}

func applyReplacePatchFiltered(patcher ResourcePatcher, result *filter.FilteredValuesEvaluationResult, newValueRaw json.RawMessage) error {
	if len(result.Matches) == 0 {
		return serrors.ThrowNoTarget(zerrors.ThrowInvalidArgument(nil, "SCIM-4513", "Path evaluation resulted in no matches"))
	}

	for _, match := range result.Matches {
		if err := unmarshalPatchValue(newValueRaw, match.Value); err != nil {
			return err
		}
	}

	if err := ensureSinglePrimaryBasedOnMatches(result.Source, result.Matches); err != nil {
		return err
	}

	return patcher.Replaced(result.PathSegments)
}

func ensureSinglePrimaryBasedOnMatches(slice reflect.Value, matches []*filter.FilteredValuesEvaluationResultMatch) error {
	if len(matches) == 0 {
		return nil
	}

	modifiedElements := make([]reflect.Value, 0, len(matches))
	modifiedIndexes := make(map[int]bool, len(matches))
	for _, match := range matches {
		modifiedElements = append(modifiedElements, match.Element)
		modifiedIndexes[match.SourceIndex] = true
	}

	return ensureSinglePrimary(slice, modifiedElements, modifiedIndexes)
}
