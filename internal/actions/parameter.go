package actions

import (
	"github.com/dop251/goja"
	"github.com/zitadel/logging"
)

type fields map[string]interface{}

type FieldOption func(*FieldConfig)

type FieldConfig struct {
	fields

	Runtime *goja.Runtime
}

func SetFields(name string, values ...interface{}) FieldOption {
	return func(p *FieldConfig) {
		if len(values) == 0 {
			return
		}

		for _, value := range values {
			val, ok := value.(FieldOption)
			// should be a primitive type or function
			if !ok {
				p.set(name, value)
				continue
			}

			// another SetFields call
			// type of current fields value must be parameter{}
			var field fields
			if f, ok := p.fields[name]; ok {
				if field, ok = f.(fields); !ok {
					logging.WithFields("sub", name).Warn("sub is not a parameter{}")
					panic("unable to prepare parameter")
				}
			} else {
				field = fields{}
				p.fields[name] = field
			}

			fieldParam := FieldConfig{
				Runtime: p.Runtime,
				fields:  field,
			}
			val(&fieldParam)
		}
	}
}

func (f *FieldConfig) set(name string, value interface{}) {
	if _, ok := f.fields[name]; ok {
		logging.WithFields("name", name).Error("tried to overwrite field")
		panic("tried to overwrite field")
	}
	v, ok := value.(func(*FieldConfig) interface{})
	if ok {
		f.fields[name] = v(f)
		return
	}
	f.fields[name] = value
}
