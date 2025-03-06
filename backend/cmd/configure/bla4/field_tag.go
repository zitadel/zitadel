package bla4

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type field struct {
	info  reflect.StructField
	viper *viper.Viper

	tag
}

func (f *field) label() string {
	var builder strings.Builder
	builder.WriteString(f.info.Name)
	builder.WriteString(" (")
	builder.WriteString(f.info.Type.Kind().String())
	builder.WriteString(")")
	if f.description != "" {
		builder.WriteString(": ")
		builder.WriteString(f.description)
	}
	return builder.String()
}

func (f *field) sub() *viper.Viper {
	if !f.viper.IsSet(f.fieldName) {
		f.viper.Set(f.fieldName, map[string]any{})
	}
	return f.viper.Sub(f.fieldName)
}

func (f *field) printStructInfo() {
	var builder strings.Builder
	builder.WriteString("------- ")
	builder.WriteString(f.info.Name)
	builder.WriteString(" -------")
	if f.description != "" {
		builder.WriteString(": ")
		builder.WriteString(f.description)
	}
	fmt.Println(builder.String())
}

type tag struct {
	skip bool

	fieldName   string
	description string

	value reflect.Value
}

const (
	tagName        = "configure"
	defaultKey     = "default"
	descriptionKey = "description"
)

func newTag(field reflect.StructField, current reflect.Value) (config tag, err error) {
	config.fieldName = field.Name
	defer func() {
		if !config.value.IsValid() {
			if field.Type.Kind() == reflect.Pointer {
				config.value = reflect.New(field.Type.Elem())
				return
			}
			config.value = reflect.New(field.Type).Elem()
		}
	}()
	if !current.IsZero() {
		config.value = current
	}

	value, ok := field.Tag.Lookup(tagName)
	if !ok {
		return config, nil
	}
	if value == "-" {
		config.skip = true
		return config, nil
	}

	fields := strings.Split(value, ",")
	for _, f := range fields {
		configSplit := strings.Split(f, "=")
		switch strings.ToLower(configSplit[0]) {
		case defaultKey:
			if !config.value.IsZero() {
				continue
			}
			value, err := kindMapper(field.Type.Kind())(configSplit[1])
			if err != nil {
				return config, err
			}
			config.value.Set(reflect.ValueOf(value))
		case descriptionKey:
			config.description = configSplit[1]
		}
	}

	return config, nil
}

func (tag tag) defaultValue() string {
	if tag.value.IsZero() {
		return ""
	}
	return fmt.Sprintf("%v", tag.value.Interface())
}
