package bla4

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
}))

type Configure func() (value any, err error)

type Configurer interface {
	// Configure is called to configure the value.
	// It must return the same type as itself. Otherwise [Update] will panic because it is not able to set the value.
	Configure() (value any, err error)
}

func Update(v *viper.Viper, config any) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		value := reflect.ValueOf(config)
		structConfigures := structToConfigureMap(Logger, v, value)
		for key, configure := range structConfigures {
			result, err := configure()
			if err != nil {
				Logger.Error("error configuring field", slog.String("field", key), slog.Any("cause", err))
				return
			}
			value.Elem().FieldByName(key).Set(reflect.ValueOf(result))
		}

		err := v.WriteConfig()
		if err != nil {
			Logger.Error("error writing config", slog.Any("cause", err))
			os.Exit(1)
		}
	}
}

func structToConfigureMap(l *slog.Logger, v *viper.Viper, object reflect.Value) map[string]Configure {
	if object.Kind() == reflect.Pointer {
		if object.IsNil() {
			l.Debug("initialize object")
			object = reflect.New(object.Type().Elem())
		}
		return structToConfigureMap(l, v, object.Elem())
	}
	if object.Kind() != reflect.Struct {
		panic("config must be a struct")
	}

	fields := make(map[string]Configure, object.NumField())

	for i := range object.NumField() {
		if !object.Type().Field(i).IsExported() {
			continue
		}

		tag, err := newTag(object.Type().Field(i), object.Field(i))
		if err != nil {
			l.Error("failed to parse field tag", slog.Any("cause", err))
			continue
		}
		if tag.skip {
			l.Debug("skipping field", slog.String("field", object.Type().Field(i).Name))
			continue
		}

		fields[object.Type().Field(i).Name] = fieldFunction(
			l.With(slog.String("field", object.Type().Field(i).Name)),
			&field{
				info:  object.Type().Field(i),
				tag:   tag,
				viper: v,
			},
		)
	}
	return fields
}

func fieldFunction(l *slog.Logger, f *field) (configure Configure) {
	for _, mapper := range []func(*slog.Logger, *field) Configure{
		fieldFunctionByConfigurer,
		fieldFunctionByIterator,
		fieldFunctionByReflection,
	} {
		if configure = mapper(l, f); configure != nil {
			return configure
		}
	}
	panic(fmt.Sprintf("unsupported field type: %s", f.info.Type.String()))
}

func fieldFunctionByConfigurer(l *slog.Logger, f *field) Configure {
	if f.value.Type().Implements(reflect.TypeFor[Configurer]()) {
		return func() (value any, err error) {
			f.printStructInfo()
			res, err := f.value.Interface().(Configurer).Configure()
			if err != nil {
				return nil, err
			}
			f.viper.Set(f.tag.fieldName, res)
			return res, nil
		}
	}
	return nil
}

func fieldFunctionByReflection(l *slog.Logger, f *field) Configure {
	kind := f.value.Kind()

	//nolint:exhaustive // only types that require special treatment are covered
	switch kind {
	case reflect.Pointer:
		if f.value.IsNil() {
			f.value.Set(reflect.New(f.value.Type().Elem()))
		}
		sub := f.sub()
		m := structToConfigureMap(l, sub, f.value)
		return func() (value any, err error) {
			f.printStructInfo()
			for key, configure := range m {
				value, err = configure()
				if err != nil {
					return nil, err
				}
				f.value.Elem().FieldByName(key).Set(reflect.ValueOf(value))
			}
			f.viper.Set(f.tag.fieldName, sub.AllSettings())
			return f.value.Interface(), nil
		}
	case reflect.Struct:
		sub := f.sub()
		m := structToConfigureMap(l, sub, f.value)
		return func() (value any, err error) {
			f.printStructInfo()
			for key, configure := range m {
				value, err = configure()
				if err != nil {
					return nil, err
				}
				f.value.FieldByName(key).Set(reflect.ValueOf(value))
			}
			f.viper.Set(f.tag.fieldName, sub.AllSettings())
			return f.value.Interface(), nil
		}

	case reflect.Array, reflect.Slice, reflect.Map:
		l.Warn("skipping because kind is unimplemented", slog.String("kind", kind.String()))
		return nil
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer, reflect.Invalid:
		slog.Error("skipping because kind is unsupported", slog.String("kind", kind.String()))
		return nil
	}

	mapper := typeMapper(f.info.Type)
	if mapper == nil {
		mapper = kindMapper(kind)
	}
	if mapper == nil {
		l.Error("unsupported kind", slog.String("kind", kind.String()))
		panic(fmt.Sprintf("unsupported kind: %s", kind.String()))
	}

	return func() (value any, err error) {
		prompt := promptui.Prompt{
			Label: f.label(),
			Validate: func(input string) error {
				value, err = mapper(input)
				return err
			},
			Default: f.defaultValue(),
		}
		_, err = prompt.Run()
		if err != nil {
			return nil, err
		}
		f.viper.Set(f.tag.fieldName, value)
		return value, nil
	}
}
