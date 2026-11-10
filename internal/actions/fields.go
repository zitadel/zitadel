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
			// is the lowest field and can be set without further checks

			if !ok {
				// {
				//	"value": "some value"
				// }
				p.set(name, value)
				continue
			}

			var field fields
			if f, ok := p.fields[name]; ok {
				// check if the found field is an object
				if field, ok = f.(fields); !ok {
					// panic because overwriting fields is not allowed
					logging.WithFields("sub", name).Warn("sub is not an object")
					panic("unable to prepare parameter")
				}
			} else {
				// field does not exist so far.
				// sub object for field can be created
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
	switch v := value.(type) {
	case func(config *FieldConfig) interface{}:
		f.fields[name] = v(f)
	case func(config *FieldConfig) func(call goja.FunctionCall) goja.Value:
		f.fields[name] = v(f)
	default:
		f.fields[name] = value
	}
}
