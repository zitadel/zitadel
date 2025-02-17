package configure

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/storage/database/dialect"
)

var (
	// ConfigureCmd represents the config command
	ConfigureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Guides you through configuring Zitadel",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("config called")
			fmt.Println(viper.AllSettings())
			fmt.Println(viper.Sub("database").AllSettings())
			pool, err := config.Database.Connect(cmd.Context())
			_, _ = pool, err
		},
		PreRun: ReadConfigPreRun[Config](viper.GetViper(), &config),
	}

	config Config
)

func init() {
	// Here you will define your flags and configuration settings.
	ConfigureCmd.Flags().BoolVarP(&config.upgrade, "upgrade", "u", false, "Only changed configuration values since the previously used version will be asked for")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

type Config struct {
	Database dialect.Config

	upgrade bool
}

func (c Config) Hooks() []viper.DecoderConfigOption {
	return c.Database.Hooks()
}
