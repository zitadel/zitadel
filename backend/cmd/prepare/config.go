package prepare

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/cmd/config"
	"github.com/zitadel/zitadel/backend/cmd/configure"
	"github.com/zitadel/zitadel/backend/cmd/configure/bla2"
	step001 "github.com/zitadel/zitadel/backend/cmd/prepare/001"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/database/dialect"
)

var (
	configuration Config

	// configurePrepare represents the prepare command
	configurePrepare = &cobra.Command{
		Use:   "prepare",
		Short: "Writes the configuration for the prepare command",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		// Run: func(cmd *cobra.Command, args []string) {
		// 	var err error
		// 	config.Client, err = config.Database.Connect(cmd.Context())
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	defer config.Client.Close(cmd.Context())
		// 	if err := (&step001.Step001{Database: config.Client}).Migrate(cmd.Context()); err != nil {
		// 		panic(err)
		// 	}
		// },
		// Run: configure.Update(
		// 	"prepare",
		// 	"Writes the configuration for the prepare command",
		// 	configuration.Fields(),
		// ),
		Run: func(cmd *cobra.Command, args []string) {
			bla2.Update(viper.GetViper(), &configuration)(cmd, args)
		},
		PreRun: configure.ReadConfigPreRun(viper.GetViper(), &configuration),
	}
)

type Config struct {
	config.Config `mapstructure:",squash" configure:"-"`

	Database dialect.Config `configure:"-"`
	Step001  step001.Step001

	// runtime config
	Client database.Pool `mapstructure:"-" configure:"-"`
}

func (c *Config) Hooks() (decoders []viper.DecoderConfigOption) {
	for _, hooks := range []configure.Unmarshaller{
		c.Config,
		c.Database,
	} {
		decoders = append(decoders, hooks.Hooks()...)
	}
	return decoders
}

func init() {
	configure.ConfigureCmd.AddCommand(configurePrepare)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prepareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prepareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
