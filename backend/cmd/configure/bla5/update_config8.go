package bla5

import (
	"log/slog"
	"os"
	"reflect"

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

func objectToFlow(l *slog.Logger, value reflect.Value) {
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			value = reflect.New(value.Type().Elem())
		}
		value = value.Elem()
	}

	typeOfValue := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := typeOfValue.Field(i)

		l.Info("Processing field", slog.String("field", fieldType.Name))
	}
}

type field struct {
	set func(value reflect.Value) error
}

type structField struct {
	name string

	fields []field
}

type primitiveField struct {
}
