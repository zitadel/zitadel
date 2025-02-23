package configure

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/cmd/info"
)

var lastConfiguredVersion *semver.Version

func Update(name, description string, fields []Updater) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		configVersion := viper.GetString("configuredVersion")
		if configVersion != "" {
			lastConfiguredVersion, _ = semver.NewVersion(configVersion)
		}
		for _, field := range fields {
			if err := setField(field, "", 1); err != nil {
				panic(err)
			}
		}

		fmt.Println("Using config file:", viper.ConfigFileUsed(), "fields:", fields)
		fmt.Println(viper.AllSettings())
		viper.Set("configuredVersion", info.Version().String())
		// viper.MergeConfigMap(fields)
		// if err := viper.WriteConfig(); err != nil {
		// 	panic(err)
		// }
		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}
}

func setField(field Updater, parentPath string, depth int) error {
	if !field.ShouldUpdate(lastConfiguredVersion) {
		prompt := promptui.Prompt{
			Label:     field.Name() + " did not change since last configure, skip?",
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err == nil {
			fmt.Println("skip")
			return nil
		}
	}
	fieldPath := path(parentPath, field.Name())
	switch f := field.(type) {
	case StructUpdater:
		fmt.Printf("\n%.*s %s: %s\n\n", depth*2, "-", f.Name(), f.Describe())
		for _, subField := range f.Fields() {
			err := setField(subField, fieldPath, depth+1)
			if err != nil {
				return err
			}
		}
	case FieldUpdater:
		value := viper.Get(fieldPath)
		if value == nil {
			value = f.ValueOrDefault()
		}
		prompt := promptui.Prompt{
			Label:   fmt.Sprintf("%s (%s) (%T)", f.Name(), f.Describe(), f.ValueOrDefault()),
			Default: fmt.Sprintf("%v", value),
			Validate: func(s string) error {
				if isConfirm(reflect.TypeOf(value)) {
					return nil
				}
				return f.Set(s)
			},
			HideEntered: f.ShouldHide(),
			IsConfirm:   isConfirm(reflect.TypeOf(value)),
		}
		_, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		viper.Set(fieldPath, f.ValueOrDefault())
	case OneOfUpdater:
		var possibilities []string
		for _, subField := range f.Fields() {
			possibility := subField.Name()
			if possibility == "" {
				possibility = subField.Describe()
			}
			possibilities = append(possibilities, possibility)
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("Select one of %s: (%s)", f.Name(), f.Describe()),
			Items: possibilities,
		}
		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("panic", err)
			panic(err)
		}
		setField(f.Fields()[i], fieldPath, depth-1)
		// TODO: fv is result of decoder (postgres.config in this case)
	case ConstantUpdater:
		viper.Set(fieldPath, f.Value())
	}
	return nil
}

func path(parent, field string) string {
	if parent == "" {
		return field
	}
	if strings.HasSuffix(parent, field) {
		return parent
	}
	if field == "" {
		return parent
	}
	return parent + "." + field
}

func isConfirm(t reflect.Type) bool {
	if t == nil {
		return false
	}
	if t.Kind() == reflect.Ptr {
		return isConfirm(t.Elem())
	}
	return t.Kind() == reflect.Bool
}

type Updater interface {
	Name() string
	Describe() string
	ShouldUpdate(version *semver.Version) bool
}

type FieldUpdater interface {
	Updater
	// DefaultValue() any
	Set(value any) error
	ValueOrDefault() any
	_field()
	ShouldHide() bool
}

type StructUpdater interface {
	Updater
	Fields() []Updater
	_struct()
}

type OneOfUpdater interface {
	Updater
	Fields() []Updater
	_oneOf()
}

type ConstantUpdater interface {
	Updater
	Value() any
}

var _ FieldUpdater = (*Field[string])(nil)

type Field[T any] struct {
	FieldName   string
	Value       T
	HideInput   bool
	Description string
	Version     *semver.Version
	Validate    func(T) error
}

// Describe implements [FieldUpdater].
func (f Field[T]) Describe() string {
	return f.Description
}

// Set implements [FieldUpdater].
func (f *Field[T]) Set(value any) (err error) {
	switch v := value.(type) {
	case string:
		f.Value, err = mapString[T](v)
	case map[string]any:
		f.Value, err = mapMap[T](v)
	case T:
		f.Value = v
	default:
		err = fmt.Errorf("unsupported type %T", v)
	}

	if err != nil || f.Validate == nil {
		return err
	}

	return f.Validate(f.Value)
}

func (f Field[T]) ValueOrDefault() any {
	return f.Value
}

func mapMap[T any](value map[string]any) (t T, err error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return t, err
	}
	return t, json.Unmarshal(jsonValue, &t)
}

func mapString[T any](value string) (t T, err error) {
	var v any
	switch reflect.TypeFor[T]().Kind() {
	case reflect.Bool:
		v, err = strconv.ParseBool(value)
		if err != nil {
			return t, err
		}
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return t, err
		}
		v = int(i)
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return t, err
		}
		v = int8(i)
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return t, err
		}
		v = int16(i)
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return t, err
		}
		v = int32(i)
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return t, err
		}
		v = int64(i)
	case reflect.Uint:
		i, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return t, err
		}
		v = uint(i)
	case reflect.Uint8:
		i, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return t, err
		}
		v = uint8(i)
	case reflect.Uint16:
		i, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return t, err
		}
		v = uint16(i)
	case reflect.Uint32:
		i, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return t, err
		}
		v = uint32(i)
	case reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return t, err
		}
		v = uint64(i)
	case reflect.Float32:
		i, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return t, err
		}
		v = float32(i)
	case reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return t, err
		}
		v = float64(i)
	case reflect.Complex64:
		i, err := strconv.ParseComplex(value, 64)
		if err != nil {
			return t, err
		}
		v = complex64(i)
	case reflect.Complex128:
		v, err = strconv.ParseComplex(value, 128)
	case reflect.String:
		v = value
	default:
		k := reflect.TypeFor[T]().Kind()
		_ = k
		return t, errors.New("not implemented")
	}
	return v.(T), nil
}

// Field implements [FieldUpdater].
func (f Field[T]) Name() string {
	return f.FieldName
}

// ShouldUpdate implements [FieldUpdater].
func (f Field[T]) ShouldUpdate(version *semver.Version) bool {
	if version == nil {
		return true
	}
	return f.Version.GreaterThan(version)
}

// ShouldHide implements [FieldUpdater].
func (f Field[T]) ShouldHide() bool {
	return f.HideInput
}

func (f Field[T]) _field() {}

var _ StructUpdater = (*Struct)(nil)

type Struct struct {
	FieldName   string
	Description string
	SubFields   []Updater
}

// Describe implements [StructUpdater].
func (s Struct) Describe() string {
	return s.Description
}

func (s Struct) Name() string {
	return s.FieldName
}

func (s Struct) Fields() []Updater {
	return s.SubFields
}

func (s Struct) ShouldUpdate(version *semver.Version) bool {
	if version == nil {
		return true
	}
	for _, field := range s.SubFields {
		if !field.ShouldUpdate(version) {
			continue
		}
		return true
	}
	return false
}

func (Struct) _struct() {}

type OneOf struct {
	FieldName   string
	Description string
	SubFields   []Updater
}

// Describe implements [OneOfUpdater].
func (o OneOf) Describe() string {
	return o.Description
}

// Fields implements [OneOfUpdater].
func (o OneOf) Fields() []Updater {
	return o.SubFields
}

// Name implements [FieldUpdater].
func (o OneOf) Name() string {
	return o.FieldName
}

func (OneOf) _oneOf() {}

// ShouldUpdate implements [OneOfUpdater].
func (o OneOf) ShouldUpdate(version *semver.Version) bool {
	if version == nil {
		return true
	}

	for _, field := range o.SubFields {
		if !field.ShouldUpdate(version) {
			continue
		}
		return true
	}
	return false
}

var _ OneOfUpdater = (*OneOf)(nil)

func FieldName(parent, field string) string {
	if parent == "" {
		return field
	}
	return parent + "." + field
}

type Constant[T any] struct {
	Description string
	Constant    T
	Version     *semver.Version
}

var _ ConstantUpdater = (*Constant[any])(nil)

// Describe implements [ConstantUpdater].
func (c Constant[T]) Describe() string {
	return c.Description
}

// Name implements [ConstantUpdater].
func (c Constant[T]) Name() string {
	return ""
}

// ShouldUpdate implements [ConstantUpdater].
func (c Constant[T]) ShouldUpdate(version *semver.Version) bool {
	if version == nil {
		return true
	}
	return c.Version.GreaterThan(version)
}

// Value implements [ConstantUpdater].
func (c Constant[T]) Value() any {
	return c.Constant
}
