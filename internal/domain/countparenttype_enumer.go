// Code generated by "enumer -type CountParentType -transform lower -trimprefix CountParentType -sql"; DO NOT EDIT.

package domain

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const _CountParentTypeName = "instanceorganization"

var _CountParentTypeIndex = [...]uint8{0, 8, 20}

const _CountParentTypeLowerName = "instanceorganization"

func (i CountParentType) String() string {
	if i < 0 || i >= CountParentType(len(_CountParentTypeIndex)-1) {
		return fmt.Sprintf("CountParentType(%d)", i)
	}
	return _CountParentTypeName[_CountParentTypeIndex[i]:_CountParentTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _CountParentTypeNoOp() {
	var x [1]struct{}
	_ = x[CountParentTypeInstance-(0)]
	_ = x[CountParentTypeOrganization-(1)]
}

var _CountParentTypeValues = []CountParentType{CountParentTypeInstance, CountParentTypeOrganization}

var _CountParentTypeNameToValueMap = map[string]CountParentType{
	_CountParentTypeName[0:8]:       CountParentTypeInstance,
	_CountParentTypeLowerName[0:8]:  CountParentTypeInstance,
	_CountParentTypeName[8:20]:      CountParentTypeOrganization,
	_CountParentTypeLowerName[8:20]: CountParentTypeOrganization,
}

var _CountParentTypeNames = []string{
	_CountParentTypeName[0:8],
	_CountParentTypeName[8:20],
}

// CountParentTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func CountParentTypeString(s string) (CountParentType, error) {
	if val, ok := _CountParentTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _CountParentTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to CountParentType values", s)
}

// CountParentTypeValues returns all values of the enum
func CountParentTypeValues() []CountParentType {
	return _CountParentTypeValues
}

// CountParentTypeStrings returns a slice of all String values of the enum
func CountParentTypeStrings() []string {
	strs := make([]string, len(_CountParentTypeNames))
	copy(strs, _CountParentTypeNames)
	return strs
}

// IsACountParentType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i CountParentType) IsACountParentType() bool {
	for _, v := range _CountParentTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

func (i CountParentType) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *CountParentType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of CountParentType: %[1]T(%[1]v)", value)
	}

	val, err := CountParentTypeString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
