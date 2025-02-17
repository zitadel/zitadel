package configure

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/zitadel/backend/cmd/config"
)

func Update(name, description string, fields []Updater) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		u := Struct{
			FieldName:   name,
			Description: description,
			SubFields:   fields,
		}
		fields := viper.AllSettings()
		updateStruct(u, fields, 1)

		fmt.Println("Using config file:", viper.ConfigFileUsed(), "fields:", fields)
		// viper.MergeConfigMap(fields)
		// if err := viper.WriteConfig(); err != nil {
		// 	panic(err)
		// }
	}
}

func updateStruct(updater StructUpdater, fields map[string]any, depth int) {
	for _, field := range updater.Fields() {
		// fmt.Printf("field: %s.%s\n       %s\n", parent, field.Name(), field.Describe())
		setField(field, fields, depth)
	}
}

func setField(field Updater, fields map[string]any, depth int) {
	if !field.ShouldUpdate(config.V3) {
		prompt := promptui.Prompt{
			Label:     field.Name() + " did not change since last configure, skip?",
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err == nil {
			fmt.Println("skip")
			return
		}
		return
	}
	fmt.Println(field.Name(), field.Describe())
	switch f := field.(type) {
	case StructUpdater:
		if fields[f.Name()] == nil {
			fields[f.Name()] = map[string]any{}
		}
		fmt.Printf("%.*s %s: %s\n", depth*2, "-", f.Name(), f.Describe())
		updateStruct(f, fields[f.Name()].(map[string]any), depth+1)
	case FieldUpdater:
		prompt := promptui.Prompt{
			Label:   fmt.Sprintf("%s (%s) (%T)", f.Name(), f.Describe(), f.DefaultValue()),
			Default: fmt.Sprintf("%v", f.DefaultValue()),
			Validate: func(s string) error {
				if isConfirm(reflect.TypeOf(f.DefaultValue())) {
					return nil
				}
				return f.Set(s)
			},
			IsConfirm: isConfirm(reflect.TypeOf(f.DefaultValue())),
		}
		val, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		fields[f.Name()] = f.Set(val)
	case OneOfUpdater:
		var possibilities []string
		for _, subField := range f.Fields() {
			possibilities = append(possibilities, subField.Name())
			fields[subField.Name()] = subField
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("Select one of %s: (%s)", f.Name(), f.Describe()),
			Items: possibilities,
		}
		i, _, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		setField(f.Fields()[i], fields, depth+1)
	}
}

func isConfirm(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		return isConfirm(t.Elem())
	}
	return t.Kind() == reflect.Bool
}

type Updater interface {
	Name() string
	Describe() string
	ShouldUpdate(version config.Version) bool
}

type FieldUpdater interface {
	Updater
	DefaultValue() any
	Set(value string) error
	_field()
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

var _ FieldUpdater = (*Field[string])(nil)

type Field[T any] struct {
	FieldName   string
	Default     T
	Value       *T
	Description string
	Version     config.Version
}

// DefaultValue implements [FieldUpdater].
func (uf Field[T]) DefaultValue() any {
	return uf.Default
}

// Describe implements [FieldUpdater].
func (uf Field[T]) Describe() string {
	return uf.Description
}

// Set implements [FieldUpdater].
func (uf Field[T]) Set(value string) error {
	var v T
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return fmt.Errorf("failed to unmarshal value: %v", err)
	}
	*uf.Value = v
	return nil
}

// Field implements [FieldUpdater].
func (uf Field[T]) Name() string {
	return uf.FieldName
}

// ShouldUpdate implements [FieldUpdater].
func (uf Field[T]) ShouldUpdate(version config.Version) bool {
	return uf.Version <= version
}

func (f Field[T]) _field() {}

var _ StructUpdater = (*Struct)(nil)

type Struct struct {
	FieldName   string
	Description string
	SubFields   []Updater
}

// Describe implements [StructUpdater].
func (us Struct) Describe() string {
	return us.Description
}

func (us Struct) Name() string {
	return us.FieldName
}

func (us Struct) Fields() []Updater {
	return us.SubFields
}

func (f Struct) _struct() {}

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

func (f OneOf) _oneOf() {}

// ShouldUpdate implements [OneOfUpdater].
func (o OneOf) ShouldUpdate(version config.Version) bool {
	for _, field := range o.SubFields {
		if !field.ShouldUpdate(version) {
			continue
		}
		return true
	}
	return false
}

var _ OneOfUpdater = (*OneOf)(nil)

func (us Struct) ShouldUpdate(version config.Version) bool {
	for _, field := range us.SubFields {
		if !field.ShouldUpdate(version) {
			continue
		}
		return true
	}
	return false
}

func FieldName(parent, field string) string {
	if parent == "" {
		return field
	}
	return parent + "." + field
}
