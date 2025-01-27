package patch

import (
	"reflect"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func applyRemovePatch(patcher ResourcePatcher, op *Operation, value interface{}) error {
	// the root cannot be removed
	if op.Path == nil {
		logging.Info("SCIM: remove patch without path")
		return serrors.ThrowNoTarget(zerrors.ThrowInvalidArgument(nil, "SCIM-ozzy54", "Remove patch without path is not supported"))
	}

	result, err := patcher.FilterEvaluator().Evaluate(reflect.ValueOf(value), op.Path)
	if err != nil {
		return serrors.ThrowInvalidPath(zerrors.ThrowInvalidArgument(err, "SCIM-sd41", "Failed to evaluate path"))
	}

	switch filterResult := result.(type) {
	case *filter.SimpleValueEvaluationResult:
		return applyRemovePatchSimple(patcher, filterResult)
	case *filter.FilteredValuesEvaluationResult:
		return applyRemovePatchFiltered(patcher, filterResult)
	}

	logging.Errorf("SCIM remove patch: unsupported filter type %T", result)
	return serrors.ThrowInvalidPath(zerrors.ThrowInvalidArgument(err, "SCIM-12syw", "Invalid patch path"))
}

func applyRemovePatchSimple(patcher ResourcePatcher, filterResult *filter.SimpleValueEvaluationResult) error {
	filterResult.Value.Set(reflect.Zero(filterResult.Value.Type()))
	return patcher.Removed(filterResult.PathSegments)
}

func applyRemovePatchFiltered(patcher ResourcePatcher, filterResult *filter.FilteredValuesEvaluationResult) error {
	if len(filterResult.Matches) == 0 {
		return nil
	}

	// if a subattribute is selected, set that one to nil instead of removing the elements from the slice
	if len(filterResult.PathSegments) > 1 {
		for _, match := range filterResult.Matches {
			match.Value.Set(reflect.Zero(match.Value.Type()))
		}

		return patcher.Removed(filterResult.PathSegments)
	}

	slice := filterResult.Source
	sliceLen := slice.Len()

	// if all elements are matched, set the field to nil
	if sliceLen == len(filterResult.Matches) {
		filterResult.Source.Set(reflect.Zero(slice.Type()))
		return patcher.Removed(filterResult.PathSegments)
	}

	// start at the very last matched value to keep correct indexing
	for i := len(filterResult.Matches) - 1; i >= 0; i-- {
		match := filterResult.Matches[i]
		slice = reflect.AppendSlice(slice.Slice(0, match.SourceIndex), slice.Slice(match.SourceIndex+1, sliceLen))
		sliceLen--
	}

	filterResult.Source.Set(slice)
	return patcher.Removed(filterResult.PathSegments)
}
