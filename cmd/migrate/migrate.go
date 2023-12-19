package migrate

import (
	_ "embed"
	"strings"

	"github.com/spf13/cobra"
)

var (
	instanceIDs []string
	system      bool
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrates the events of an instance from one database to another",
	}

	migrateFlags(cmd)
	cmd.AddCommand(
		eventstoreCmd(),
		systemCmd(),
		projectionsCmd(),
		authCmd(),
		verifyCmd(),
	)

	return cmd
}

func migrateFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVar(&instanceIDs, "instance", nil, "id of the instance to migrate")
	cmd.PersistentFlags().BoolVar(&system, "system", false, "migrates the whole system")
	cmd.MarkFlagsOneRequired("system", "instance")
}

func instanceClause() string {
	if system {
		return "WHERE instance_id <> ''"
	}
	for i := range instanceIDs {
		instanceIDs[i] = "'" + instanceIDs[i] + "'"
	}

	// COPY does not allow parameters so we need to set them directly
	return "WHERE instance_id IN (" + strings.Join(instanceIDs, ", ") + ")"
}
