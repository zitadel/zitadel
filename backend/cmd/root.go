package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zitadel/zitadel/backend/cmd/config"
	"github.com/zitadel/zitadel/backend/cmd/configure"
	"github.com/zitadel/zitadel/backend/cmd/prepare"
	"github.com/zitadel/zitadel/backend/cmd/start"
	"github.com/zitadel/zitadel/backend/cmd/upgrade"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "zitadel [subcommand]",
	Short: "A brief description of your application",
	Long:  `zitadel`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(
		configure.ConfigureCmd,
		prepare.PrepareCmd,
		start.StartCmd,
		upgrade.UpgradeCmd,
	)

	cobra.OnInitialize(config.InitConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&config.Path, "config", "", "config file (default is $HOME/.zitadel.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
