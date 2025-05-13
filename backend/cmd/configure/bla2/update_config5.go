package bla2

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

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
		s := object{
			value: reflect.Indirect(reflect.ValueOf(config)),
			viper: v,
		}
		fields, err := s.fields()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, field := range fields {
			err = field.execute()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		err = v.WriteConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

type Field struct {
	// Either Prompt or Sub must be set
	Prompt prompt
	Sub    []*Field
	Set    func()
}

func (f *Field) execute() error {
	if f.Prompt != nil {
		if err := f.Prompt.Run(); err != nil {
			return err
		}
		f.Set()
		return nil
	}
	for _, sub := range f.Sub {
		err := sub.execute()
		if err != nil {
			return err
		}
		sub.Set()
	}
	f.Set()
	return nil
}

func mapValue(typ reflect.Type, value string) (v reflect.Value, err error) {
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
	case reflect.Complex64:
		val, err = strconv.ParseComplex(value, 64)
		val = complex64(val.(complex128))
	case reflect.Complex128:
		val, err = strconv.ParseComplex(value, 128)
	default:
		return v, fmt.Errorf("unsupported type: %s", typ.Kind().String())
	}
	if err != nil {
		return v, err
	}

	res := reflect.ValueOf(val)
	if !res.CanConvert(typ) {
		return v, fmt.Errorf("cannot convert %T to %s", val, typ.Kind().String())
	}
	return res.Convert(typ), nil
}
