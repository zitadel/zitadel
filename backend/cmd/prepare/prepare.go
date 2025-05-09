package prepare

import (
	"github.com/spf13/cobra"

	step001 "github.com/zitadel/zitadel/backend/cmd/prepare/001"
)

var (
	// PrepareCmd represents the prepare command
	PrepareCmd = &cobra.Command{
		Use:   "prepare",
		Short: "Prepares external services before starting Zitadel",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			// var err error
			// configuration.Client, err = configuration.Database.Connect(cmd.Context())
			// if err != nil {
			// 	panic(err)
			// }
			defer configuration.Client.Close(cmd.Context())
			if err := (&step001.Step001{Database: configuration.Client}).Migrate(cmd.Context()); err != nil {
				panic(err)
			}
		},
	}
)

type Migration interface {
}
