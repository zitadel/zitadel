package patch

import (
	"reflect"

	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func applyAddPatch(patcher ResourcePatcher, op *Operation, value interface{}) error {
	if op.Path == nil {
		return flattenAndApplyPatchOperations(patcher, op, value)
	}

	result, err := patcher.FilterEvaluator().Evaluate(reflect.ValueOf(value), op.Path)
	if err != nil {
		return serrors.ThrowInvalidPath(zerrors.ThrowInvalidArgument(err, "SCIM-opzz8", "Failed to evaluate path"))
	}

	evaluationResult, ok := result.(*filter.SimpleValueEvaluationResult)
	if !ok {
		return serrors.ThrowInvalidPath(zerrors.ThrowInvalidArgument(nil, "SCIM-opty8", "Filtered paths are not allowed for add patch operations"))
	}

	if evaluationResult.Value.Kind() != reflect.Slice {
		return applyReplacePatchSimple(patcher, evaluationResult, op.Value, op.valueIsArray)
	}

	elementType := evaluationResult.Value.Type().Elem()
	newElementsSlice, err := unmarshalPatchValuesSlice(elementType, op.Value, op.valueIsArray)
	if err != nil {
		return err
	}

	oldLen := evaluationResult.Value.Len()
	newSlice := reflect.MakeSlice(reflect.SliceOf(elementType), oldLen, oldLen+newElementsSlice.Len())

	// copy over existing values
	reflect.Copy(newSlice, evaluationResult.Value)

	// according to the RFC only "new" values should be added
	// existing values should be replaced
	newSlice, modifiedIndexes := addOrReplaceByValue(newSlice, newElementsSlice)

	evaluationResult.Value.Set(newSlice)
	if err = ensureSinglePrimaryAdded(evaluationResult.Value, newElementsSlice, modifiedIndexes); err != nil {
		return err
	}

	return patcher.Added(evaluationResult.PathSegments)
}

func ensureSinglePrimaryAdded(resultSlice, newSlice reflect.Value, modifiedIndexes map[int]bool) error {
	modifiedValues := make([]reflect.Value, newSlice.Len())
	for i := 0; i < newSlice.Len(); i++ {
		modifiedValues[i] = newSlice.Index(i)
	}

	return ensureSinglePrimary(resultSlice, modifiedValues, modifiedIndexes)
}

func addOrReplaceByValue(entries, newEntries reflect.Value) (reflect.Value, map[int]bool) {
	modifiedIndexes := make(map[int]bool, newEntries.Len())
	if entries.Len() == 0 {
		for i := 0; i < newEntries.Len(); i++ {
			modifiedIndexes[i] = true
		}

		return newEntries, modifiedIndexes
	}

	valueField := entries.Index(0).Elem().FieldByName(fieldNameValue)
	if !valueField.IsValid() || valueField.Kind() != reflect.String {
		entriesLen := entries.Len()
		for i := 0; i < newEntries.Len(); i++ {
			modifiedIndexes[i+entriesLen] = true
		}

		return reflect.AppendSlice(entries, newEntries), modifiedIndexes
	}

	existingValueIndexes := make(map[string]int, entries.Len())
	for i := 0; i < entries.Len(); i++ {
		value := entries.Index(i).Elem().FieldByName(fieldNameValue).String()
		if _, ok := existingValueIndexes[value]; ok {
			continue
		}

		existingValueIndexes[value] = i
	}

	entriesLen := entries.Len()
	for i := 0; i < newEntries.Len(); i++ {
		newEntry := newEntries.Index(i)
		value := newEntry.Elem().FieldByName(fieldNameValue).String()
		index, valueExists := existingValueIndexes[value]

		// according to the rfc if the value already exists it should be replaced
		if valueExists {
			entries.Index(index).Set(newEntry)
			modifiedIndexes[index] = true
			continue
		}

		entries = reflect.Append(entries, newEntry)
		modifiedIndexes[entriesLen] = true
		entriesLen++
	}

	return entries, modifiedIndexes
}
