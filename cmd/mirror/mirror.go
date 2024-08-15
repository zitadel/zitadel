package mirror

import (
	"bytes"
	_ "embed"
	"os"
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
	isSrcFile     bool
	isDestFile    bool
	filePath      string
)

func New(configFiles *[]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "mirrors all data of ZITADEL between databases, or between a database and files",
		Long: `mirrors all data of ZITADEL between databases, or between a database and files
ZITADEL needs to be initialized and set up with --for-mirror

The command does mirror all data needed and recomputes the projections.
For more details call the help functions of the sub commands.

Order of execution:
1. mirror system tables
2. mirror auth tables
3. mirror event store tables
4. recompute projections
5. verify`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := viper.MergeConfig(bytes.NewBuffer(defaultConfig))
			logging.OnError(err).Fatal("unable to read default config")

			for _, file := range *configFiles {
				viper.SetConfigFile(file)
				err := viper.MergeInConfig()
				logging.WithFields("file", file).OnError(err).Warn("unable to read config file")
			}

			if isSrcFile = viper.IsSet("Source.file.path"); isSrcFile {
				filePath = viper.GetString("Source.file.path") + "/"
			}
			if isDestFile = viper.IsSet("Destination.file.path"); isDestFile {
				filePath = viper.GetString("Destination.file.path") + "/"
			}

			if isSrcFile || isDestFile {
				if isSrcFile && isDestFile {
					logging.Fatal("both source and destination cannot be files")
				}

				if !(shouldIgnorePrevious && shouldReplace) {
					logging.Fatal("both --ignore-previous and --replace flags must be set for mirroring files")
				}

				if stat, err := os.Stat(filePath); err != nil || !stat.IsDir() {
					if os.IsNotExist(err) {
						logging.Fatal("file path does not exist")
					}
					logging.Fatal("file path leads to a file not a directory")
				}
			}
		},
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

	mirrorFlags(cmd)
	cmd.Flags().BoolVar(&shouldIgnorePrevious, "ignore-previous", false, "ignores previous migrations of the events table")
	cmd.Flags().BoolVar(&shouldReplace, "replace", false, `replaces all data of the following tables for the provided instances or all if the "--system"-flag is set:
* system.assets
* auth.auth_requests
* eventstore.unique_constraints
The flag should be provided if you want to execute the mirror command multiple times so that the static data are also mirrored to prevent inconsistent states.`)
	migrateProjectionsFlags(cmd)

	cmd.AddCommand(
		eventstoreCmd(),
		systemCmd(),
		projectionsCmd(),
		authCmd(),
		verifyCmd(),
	)

	return cmd
}

func mirrorFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVar(&instanceIDs, "instance", nil, "id or comma separated ids of the instance(s) to migrate. Either this or the `--system`-flag must be set. Make sure to always use the same flag if you execute the command multiple times.")
	cmd.PersistentFlags().BoolVar(&isSystem, "system", false, "migrates the whole system. Either this or the `--instance`-flag must be set. Make sure to always use the same flag if you execute the command multiple times.")
	cmd.MarkFlagsOneRequired("system", "instance")
	cmd.MarkFlagsMutuallyExclusive("system", "instance")
}

func instanceClause() string {
	if isSystem {
		if isSrcFile || isDestFile {
			return ""
		}
		return "WHERE instance_id <> ''"
	}
	for i := range instanceIDs {
		instanceIDs[i] = "'" + instanceIDs[i] + "'"
	}

	// COPY does not allow parameters so we need to set them directly
	return "WHERE instance_id IN (" + strings.Join(instanceIDs, ", ") + ")"
}
