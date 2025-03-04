package bla2

import (
	"log/slog"
	"reflect"

	"github.com/spf13/viper"
)

type object struct {
	value reflect.Value
	viper *viper.Viper
}

func (o *object) fields() ([]*Field, error) {
	fields := make([]*Field, 0, o.value.NumField())
	for i := range o.value.NumField() {
		if !o.value.CanSet() {
			slog.Info("skipping field because it is not settable", slog.String("field", o.value.Type().Field(i).Name))
			continue
		}
		var currentValue any
		if !o.value.Field(i).IsZero() {
			currentValue = o.value.Field(i).Interface()
		}
		tag, err := newFieldTag(o.value.Type().Field(i), currentValue)
		if err != nil {
			return nil, err
		}

		if tag.skip {
			continue
		}

		field, err := o.field(o.value.Field(i), tag)
		if err != nil {
			return nil, err
		}
		// field can be nil if the fields kind is not supported
		if field != nil {
			fields = append(fields, field)
		}
	}

	return fields, nil
}

func (o *object) field(field reflect.Value, tag fieldTag) (*Field, error) {
	switch field.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String:

		prompt := promptFromPrimitive(&primitive{
			typ: field.Type(),
			tag: tag,
		})
		return &Field{
			Prompt: prompt,
			Set: func() {
				o.viper.Set(tag.fieldName, prompt.Result())
			},
		}, nil
	case reflect.Struct:
		if o.viper.Get(tag.fieldName) == nil {
			o.viper.Set(tag.fieldName, make(map[string]any))
		}
		sub := &object{
			value: reflect.Indirect(field),
			viper: o.viper.Sub(tag.fieldName),
		}
		subFields, err := sub.fields()
		if err != nil {
			return nil, err
		}
		return &Field{
			Sub: subFields,
			Set: func() {
				o.viper.Set(tag.fieldName, sub.viper.AllSettings())
			},
		}, nil
	case reflect.Pointer:
		return o.field(field.Elem(), tag)
	case reflect.Array, reflect.Interface, reflect.Map, reflect.Slice:
		slog.Info(
			"skipping field because type is not implemented",
			slog.String("field", tag.fieldName),
			slog.String("kind", field.Kind().String()),
		)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Invalid:
		slog.Info(
			"skipping field because field type is invalid, add `configure=\"-\"` to skip",
			slog.String("field", tag.fieldName),
			slog.String("kind", field.Kind().String()),
		)
	}
	return nil, nil
}
