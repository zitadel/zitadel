package bla2

import (
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
		// promptui.Select
		// err := handleStruct(v, reflect.ValueOf(config), configuration{})
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// err = v.WriteConfig()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
	}
}

const (
	tagName    = "configure"
	defaultKey = "default"
	oneOfKey   = "oneof"
)

type Field struct {
	Name         string
	DefaultValue any
	Prompt       prompt
}

type prompt interface {
	Run()
}

var (
	_ prompt = (*promptui.Select)(nil)
	_ prompt = (*promptui.SelectWithAdd)(nil)
	_ prompt = (*promptui.Prompt)(nil)
)
