package bla3

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

var customType = reflect.TypeFor[Custom]()

type Custom interface {
	Configure() (value any, err error)
}

type Object struct {
	value reflect.Value

	tag fieldTag
}

func NewObject(value any) Object {
	return Object{
		value: reflect.Indirect(reflect.ValueOf(value)),
	}
}

func (o *Object) Configure(v *viper.Viper) error {
	if o.tag.label() != "" {
		fmt.Println("\n", o.tag.label())
	}
	for i := range o.value.NumField() {
		if !o.value.Type().Field(i).IsExported() {
			continue
		}
		structField := o.value.Field(i)

		f := Field{
			value:       structField,
			structField: o.value.Type().Field(i),
		}
		err := f.Configure(v)
		if err != nil {
			return err
		}
	}

	return nil
}

type Field struct {
	value reflect.Value

	structField reflect.StructField

	tag fieldTag
}

func NewField(value any) Field {
	return Field{
		value: reflect.ValueOf(value),
	}
}

func (f *Field) validate(input string) (err error) {
	configuredValue, err := mapValue(f.value.Type(), input)
	if err != nil {
		return err
	}
	f.value.Set(configuredValue)
	return nil
}

func (f *Field) callCustom(v *viper.Viper) (ok bool, err error) {
	if !f.value.Type().Implements(customType) {
		return false, nil
	}

	custom := f.value.Interface().(Custom)
	value, err := custom.Configure()
	if err != nil {
		return true, err
	}
	switch val := value.(type) {
	case Field:
		val.tag, err = newFieldTag(f.structField, val.value.Interface())
		val.structField = f.structField
		if err != nil {
			return false, err
		}
		return true, val.Configure(v)
	default:
		v.Set(f.tag.fieldName, val)
		return true, nil
	}
}

func (f *Field) Configure(v *viper.Viper) error {
	if f.value.IsNil() {
		f.value.Set(reflect.New(f.value.Type().Elem()))
	}

	tag, err := newFieldTag(f.structField, f.value.Interface())
	if err != nil {
		return err
	}
	if tag.skip {
		return nil
	}

	if ok, err := f.callCustom(v); ok || err != nil {
		return err
	}
	f.value = reflect.Indirect(f.value)
	if ok, err := f.callCustom(v); ok || err != nil {
		return err
	}

	kind := f.value.Kind()
	switch kind {
	case reflect.Bool:
		prompt := promptui.Prompt{
			Label:     f.tag.label(),
			IsConfirm: true,
		}
		_, err := prompt.Run()
		selected := true
		if err != nil {
			if err.Error() != "" {
				return err
			}
			selected = false
		}
		f.value.SetBool(selected)
		v.Set(f.tag.fieldName, f.value.Interface())
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:

		prompt := promptui.Prompt{
			Label:    f.tag.label(),
			Validate: f.validate,
			Default:  f.tag.defaultValue(),
		}
		_, err := prompt.Run()
		if err != nil {
			return err
		}
		v.Set(f.tag.fieldName, f.value.Interface())
		return nil
	case reflect.Struct:
		o := Object{
			value: f.value,
			tag:   f.tag,
		}
		if !v.IsSet(f.tag.fieldName) {
			v.Set(f.tag.fieldName, map[string]any{})
		}
		sub := v.Sub(f.tag.fieldName)
		err := o.Configure(sub)
		if err != nil {
			return err
		}
		v.Set(f.tag.fieldName, sub.AllSettings())
		return nil
	case reflect.Pointer:
		if f.value.IsNil() {
			f.value = reflect.New(f.value.Type().Elem())
		}
		f.value = f.value.Elem()
		return f.Configure(v)
	case reflect.Array, reflect.Slice, reflect.Map:
		slog.Warn("skipping because kind is unimplemented", slog.String("field", f.tag.fieldName), slog.String("kind", kind.String()))
		return nil
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer, reflect.Invalid:
		slog.Error("skipping because kind is unsupported", slog.String("field", f.tag.fieldName), slog.String("kind", kind.String()))
		return nil
	}
	return nil
}
