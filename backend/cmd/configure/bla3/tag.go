package bla3

import (
	"fmt"
	"reflect"
	"strings"
)

type fieldTag struct {
	skip bool

	fieldName   string
	description string

	currentValue any
}

const (
	tagName        = "configure"
	defaultKey     = "default"
	descriptionKey = "description"
)

func newFieldTag(field reflect.StructField, current any) (config fieldTag, err error) {
	config.fieldName = field.Name
	if current != nil {
		config.currentValue = current
	}

	value, ok := field.Tag.Lookup(tagName)
	if !ok {
		if config.currentValue == nil {
			config.currentValue = reflect.New(field.Type).Elem().Interface()
		}
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
			config.currentValue, err = mapValue(field.Type, configSplit[1])
			if err != nil {
				return config, err
			}
		case descriptionKey:
			config.description = configSplit[1]
		}
	}
	if config.currentValue == nil {
		config.currentValue = reflect.New(field.Type).Elem().Interface()
	}

	return config, nil
}

func (tag fieldTag) label() string {
	if tag.fieldName == "" {
		return ""
	}
	valueType := reflect.TypeOf(tag.currentValue)
	typ := valueType.Kind().String()
	if typeParsers[valueType] != nil {
		typ = valueType.Name()
	}

	label := fmt.Sprintf("%s (%s)", tag.fieldName, typ)
	if tag.description == "" {
		return label
	}
	return fmt.Sprintf("%s (%s)", label, tag.description)
}

func (tag fieldTag) defaultValue() string {
	return fmt.Sprintf("%v", tag.currentValue)
}
