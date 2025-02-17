package configure

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/zitadel/backend/cmd/config"
)

var (
	// ConfigureCmd represents the config command
	ConfigureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Guides you through configuring Zitadel for the specified command",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("config called")
		// 	// fmt.Println(viper.AllSettings())
		// 	// fmt.Println(viper.Sub("database").AllSettings())
		// 	// pool, err := config.Database.Connect(cmd.Context())
		// 	// _, _ = pool, err
		// },
		PersistentPreRun: configurePreRun,
		Run: func(cmd *cobra.Command, args []string) {
			t := new(test)
			// Update2(*t)
			Update("test", "test", t.Fields())(cmd, args)
		},
	}

	configuration Config
)

func configurePreRun(cmd *cobra.Command, args []string) {
	// cmd.InheritedFlags().Lookup("config").Hidden = true
	ReadConfigPreRun(viper.GetViper(), &configuration)(cmd, args)
}

func init() {
	// Here you will define your flags and configuration settings.
	ConfigureCmd.PersistentFlags().BoolVarP(&configuration.upgrade, "upgrade", "u", false, "Only changed configuration values since the previously used version will be asked for")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

type Config struct {
	upgrade bool `mapstructure:"-"`
}

func (c *Config) Hooks() []viper.DecoderConfigOption {
	return nil
}

type sub struct {
	F1 string
	F2 int
	F3 *bool
}

func (s sub) Fields() []Updater {
	return []Updater{
		Field[string]{
			FieldName:   "f1",
			Value:       &s.F1,
			Default:     "",
			Description: "field 1",
			Version:     config.V3,
		},
		Field[int]{
			FieldName:   "f2",
			Value:       &s.F2,
			Default:     0,
			Description: "field 2",
			Version:     config.V3,
		},
		Field[*bool]{
			FieldName:   "f3",
			Value:       &s.F3,
			Default:     nil,
			Description: "field 3",
			Version:     config.V3,
		},
	}
}

type test struct {
	F1  string
	Sub sub
}

func (t test) Fields() []Updater {
	return []Updater{
		Field[string]{
			FieldName:   "f1",
			Value:       &t.F1,
			Default:     "",
			Description: "field 1",
			Version:     config.V3,
		},
		Struct{
			FieldName:   "sub",
			Description: "sub field",
			SubFields:   t.Sub.Fields(),
		},
	}
}
