package bla

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TestConfig struct {
	API      APIConfig     `configure:""`
	Database DatabaseOneOf `configure:"type=oneof"`
}

type APIConfig struct {
	Host string `configure:""`
	Port uint16 `configure:""`
}

type DatabaseOneOf struct {
	ConnectionString *string         `configure:""`
	Config           *DatabaseConfig `configure:""`
}

type DatabaseConfig struct {
	Host    string `configure:""`
	Port    uint16 `configure:""`
	SSLMode string `configure:""`
}

type Config interface {
	Hooks() []viper.DecoderConfigOption
}

func Update(v *viper.Viper, config any) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := handleStruct(v, reflect.ValueOf(config), configuration{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = v.WriteConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

const (
	tagName    = "configure"
	defaultKey = "default"
	oneOfKey   = "oneof"
)

func handleStruct(v *viper.Viper, val reflect.Value, config configuration) error {
	val = reflect.Indirect(val)
	if val.Type().Kind() != reflect.Struct {
		return errors.New("object is not a struct")
	}

	for i := range val.NumField() {
		if err := handleField(v, val.Type().Field(i), val.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

func handlePrimitive(field reflect.StructField, primitive reflect.Value, config configuration) error {
	if config.customPrimitiveHandler != nil {
		return config.customPrimitiveHandler(config.defaultValue)
	}
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s: %s", field.Name, config.description),
		Validate: func(s string) error {
			_, err := mapValue(primitive.Type(), s)
			return err
		},
	}
	if config.defaultValue != nil {
		prompt.Default = fmt.Sprintf("%v", config.defaultValue)
	}
	res, err := prompt.Run()
	if err != nil {
		return err
	}
	value, _ := mapValue(primitive.Type(), res)
	primitive.Set(reflect.ValueOf(value))

	return nil
}

func handleField(v *viper.Viper, field reflect.StructField, value reflect.Value) error {
	fmt.Println(field.Name)
	fieldConfig, err := fieldConfiguration(field, v.Get(field.Name))
	if err != nil {
		return err
	}
	if fieldConfig.skip {
		return nil
	}

	// check if potential pointer implements [OneOfField]
	if value.Type().Implements(reflect.TypeFor[OneOfField]()) {
		return handleOneOf(v, field, value, fieldConfig)
	}

	value = reflect.Indirect(value)

	if value.IsZero() {
		value = reflect.New(value.Type()).Elem()
	}

	// Check if potential non pointer value implements [OneOfField]
	if value.Type().Implements(reflect.TypeFor[OneOfField]()) {
		return handleOneOf(v, field, value, fieldConfig)
	}
	//nolint: exhaustive
	// Default are all primitive types
	switch value.Kind() {
	case reflect.Struct:
		if v.Get(field.Name) == nil {
			defaultValue := fieldConfig.defaultValue
			if defaultValue == nil {
				defaultValue = map[string]any{}
			}
			v.Set(field.Name, defaultValue)
		}
		sub := v.Sub(field.Name)
		if err := handleStruct(sub, value, fieldConfig); err != nil {
			return err
		}
		v.Set(field.Name, sub.AllSettings())
		return nil
	case reflect.Uintptr, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice, reflect.UnsafePointer:
		return fmt.Errorf("unsupported kind: %s", value.Type().Kind().String())
	default:
		if err := handlePrimitive(field, value, fieldConfig); err != nil {
			return err
		}
		v.Set(field.Name, value.Interface())
		return nil
	}
}

func handleOneOf(v *viper.Viper, field reflect.StructField, value reflect.Value, config configuration) error {
	possibilities := value.MethodByName("Possibilities").Call(nil)[0].Interface().([]string)
	prompt := promptui.Select{
		Label: config.description,
		Items: possibilities,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return err
	}

	selectedValue := value.MethodByName("ConfigForIndex").Call([]reflect.Value{reflect.ValueOf(idx)})[0].Interface()
	return handleField(v, field, reflect.ValueOf(selectedValue))
}

type OneOfField interface {
	Possibilities() []string
	// ConfigForIndex returns a pointer to the field that was selected
	ConfigForIndex(int) any
}

type configuration struct {
	skip         bool
	defaultValue any
	description  string

	isOneOf     bool
	oneOfValues []any

	customPrimitiveHandler func(current any) error
	customStructHandler    func(*viper.Viper) error
	customRepeatedHandler  func(*viper.Viper) error
	customOneOfHandler     func(*viper.Viper) error
}

func fieldConfiguration(field reflect.StructField, current any) (config configuration, err error) {
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
			config.defaultValue, err = mapValue(field.Type, configSplit[1])
			if err != nil {
				return config, err
			}
		case oneOfKey:
			config.isOneOf = true
			oneOfValues := strings.Split(configSplit[1], "|")
			for _, oneOfValue := range oneOfValues {
				value, err := mapValue(field.Type, oneOfValue)
				if err != nil {
					return config, err
				}
				config.oneOfValues = append(config.oneOfValues, value)
			}
		}
	}
	if current != nil {
		config.defaultValue = current
	}

	return config, nil
}

func mapValue(typ reflect.Type, value string) (_ any, err error) {
	var val any
	switch typ.Kind() {
	case reflect.String:
		val = value
	case reflect.Bool:
		val, err = strconv.ParseBool(value)
	case reflect.Int:
		val, err = strconv.Atoi(value)
	case reflect.Int8:
		val, err = strconv.ParseInt(value, 10, 8)
		val = int8(val.(int64))
	case reflect.Int16:
		val, err = strconv.ParseInt(value, 10, 16)
		val = int16(val.(int64))
	case reflect.Int32:
		val, err = strconv.ParseInt(value, 10, 32)
		val = int32(val.(int64))
	case reflect.Int64:
		val, err = strconv.ParseInt(value, 10, 64)
	case reflect.Uint:
		val, err = strconv.ParseUint(value, 10, 0)
		val = uint(val.(uint64))
	case reflect.Uint8:
		val, err = strconv.ParseUint(value, 10, 8)
		val = uint8(val.(uint64))
	case reflect.Uint16:
		val, err = strconv.ParseUint(value, 10, 16)
		val = uint16(val.(uint64))
	case reflect.Uint32:
		val, err = strconv.ParseUint(value, 10, 32)
		val = uint32(val.(uint64))
	case reflect.Uint64:
		val, err = strconv.ParseUint(value, 10, 64)
	case reflect.Float32:
		val, err = strconv.ParseFloat(value, 32)
		val = float32(val.(float64))
	case reflect.Float64:
		val, err = strconv.ParseFloat(value, 64)
	default:
		return nil, fmt.Errorf("unsupported type: %s", typ.Kind().String())
	}
	if err != nil {
		return nil, err
	}

	res := reflect.ValueOf(val)
	if !res.CanConvert(typ) {
		return nil, fmt.Errorf("cannot convert %T to %s", val, typ.Kind().String())
	}
	return res.Convert(typ).Interface(), nil
}
