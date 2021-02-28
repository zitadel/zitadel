package test

import (
	"reflect"
	"sort"
	"strings"
)

//testingT is a wrapper for testing.T
//
// this wrapper is needed for internal testing
type testingT interface {
	Errorf(format string, args ...interface{})
	Helper()
}

func AssertFieldsMapped(t testingT, object interface{}, ignoreFields ...string) (failed bool) {
	t.Helper()
	val := reflect.ValueOf(object)

	fields := BuildList(val)

	notEmptyFields := validateEmptyFields(fields, ignoreFields)
	if len(notEmptyFields) > 0 {
		sort.Strings(notEmptyFields)
		t.Errorf("expected fields are not empty:\n\t%s\n", strings.Join(notEmptyFields, ",\n\t"))
	}

	notFilledFields := validateFilledFields(fields)
	if len(notFilledFields) > 0 {
		sort.Strings(notFilledFields)
		t.Errorf("unexpected empty fields:\n\t%s\n", strings.Join(notFilledFields, ",\n\t"))
	}

	return len(notEmptyFields) > 0 || len(notFilledFields) > 0
}

func BuildList(val reflect.Value) map[string]bool {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = reflect.Indirect(val)
	}

	fields := map[string]bool{}

	for i := 0; i < val.NumField(); i++ {
		if !isPublicField(val.Type().Field(i).Name) {
			continue
		}

		if val.Field(i).Kind() == reflect.Struct || val.Field(i).Kind() == reflect.Ptr {
			fieldName := val.Type().Field(i).Name
			fields[fieldName] = false
			subFields := BuildList(val.Field(i))
			for k, v := range subFields {
				fields[fieldName+"."+k] = v
				fields[fieldName] = fields[val.Type().Field(i).Name] || v
			}
			if len(subFields) == 0 &&
				((val.Field(i).Kind() == reflect.Ptr && !val.Field(i).IsNil()) ||
					val.Field(i).Kind() == reflect.Struct && !val.Field(i).IsZero()) {

				fields[fieldName] = true
			}
		} else {
			fields[val.Type().Field(i).Name] = isFieldFilled(val.Field(i))
		}
	}

	return fields
}

func isPublicField(fieldName string) bool {
	return fieldName[0] >= 'A' && fieldName[0] <= 'Z'
}

func isFieldFilled(val reflect.Value) bool {
	if isLengthMeasurable(val) {
		switch val.Kind() {
		case reflect.Slice:
			if val.IsNil() {
				return false
			}
			fallthrough
		case reflect.Array:
			if val.Len() == 0 {
				return false
			}
			for i := 0; i < val.Len(); i++ {
				if val.Index(i).IsZero() {
					return false
				}
			}
			return true
		case reflect.Map:
			if val.Len() == 0 {
				return false
			}
			for _, key := range val.MapKeys() {
				if val.MapIndex(key).IsZero() {
					return false
				}
			}
			return true
		}
	}

	return !val.IsZero()
}

func isLengthMeasurable(val reflect.Value) bool {
	return val.Kind() == reflect.Slice ||
		val.Kind() == reflect.Array ||
		val.Kind() == reflect.Map
}

func validateEmptyFields(fields map[string]bool, emptyFields []string) (notEmptyFields []string) {
	for _, emptyField := range emptyFields {
		isFilled := fields[emptyField]
		subs := subFields(fields, emptyField)
		if isFilled {
			filledFields := filledSubFields(fields, subs)
			if len(filledFields) == 0 {
				filledFields = append(filledFields, emptyField)
			}
			notEmptyFields = append(notEmptyFields, filledFields...)
		}
		fields = removeFields(fields, append([]string{emptyField}, subs...))
	}
	return notEmptyFields
}

func filledSubFields(fields map[string]bool, subFields []string) (filledSubs []string) {
	for _, subField := range subFields {
		if fields[subField] {
			filledSubs = append(filledSubs, subField)
		}
	}
	return filledSubs
}

func subFields(fields map[string]bool, parentName string) (subs []string) {
	for fieldName := range fields {
		if strings.HasPrefix(fieldName, parentName+".") {
			subs = append(subs, fieldName)
		}
	}
	return subs
}

func removeFields(fields map[string]bool, fieldNames []string) map[string]bool {
	for _, fieldName := range fieldNames {
		delete(fields, fieldName)
	}
	return fields
}

func validateFilledFields(fields map[string]bool) (emptyFields []string) {
	for fieldName, isFilled := range fields {
		if !isFilled {
			emptyFields = append(emptyFields, fieldName)
		}
	}

	return emptyFields
}
