package bla2

import (
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
			config.currentValue, err = mapValue(field.Type, configSplit[1])
			if err != nil {
				return config, err
			}
		case descriptionKey:
			config.description = configSplit[1]
		}
	}
	if current != nil {
		config.currentValue = current
	}

	return config, nil
}
