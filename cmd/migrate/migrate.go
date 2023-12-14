package migrate

import (
	_ "embed"

	"github.com/spf13/cobra"
)

var (
	instanceID string
	system     bool
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrates the events of an instance from one database to another",
	}

	migrateFlags(cmd)
	cmd.AddCommand(
		eventsCmd(),
		projectionsCmd(),
	)

	return cmd
}

func migrateFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&instanceID, "instance", "", "id of the instance to migrate")
	cmd.PersistentFlags().BoolVar(&system, "system", false, "migrates the whole system")
}
