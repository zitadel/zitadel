package migrate

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/key"
)

var (
	instanceIDs   []string
	isSystem      bool
	shouldReplace bool
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "mirrors all data of ZITADEL from one database to another",
		Long: `mirrors all data of ZITADEL from one database to another
ZITADEL needs to be initialized

The command does mirror all data needed and recomputes the projections.
For more details call the help functions of the sub commands.

Order of execution:
1. mirror system tables
2. mirror auth tables
3. mirror event store tables
4. recompute projections
5. verify`,
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewMigrationConfig(viper.GetViper())
			projectionConfig := mustNewProjectionsConfig(viper.GetViper())

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Fatal("unable to read master key")

			copySystem(cmd.Context(), config)
			copyAuth(cmd.Context(), config)
			copyEventstore(cmd.Context(), config)

			projections(cmd.Context(), projectionConfig, masterKey)
			verifyMigration(cmd.Context(), config)
		},
	}

	migrateFlags(cmd)
	cmd.Flags().BoolVar(&shouldIgnorePrevious, "ignore-previous", false, "ignores previous migrations of the events table")
	cmd.Flags().BoolVar(&shouldReplace, "replace", false, "replaces all data except events and projections")
	migrateProjectionsFlags(cmd)

	err := viper.MergeConfig(bytes.NewBuffer(defaultConfig))
	logging.OnError(err).Fatal("unable to read default config")

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
	cmd.PersistentFlags().BoolVar(&isSystem, "system", false, "migrates the whole system")
	cmd.MarkFlagsOneRequired("system", "instance")
}

func instanceClause() string {
	if isSystem {
		return "WHERE instance_id <> ''"
	}
	for i := range instanceIDs {
		instanceIDs[i] = "'" + instanceIDs[i] + "'"
	}

	// COPY does not allow parameters so we need to set them directly
	return "WHERE instance_id IN (" + strings.Join(instanceIDs, ", ") + ")"
}
