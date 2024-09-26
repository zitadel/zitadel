package mirror

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/zitadel/zitadel/internal/unixsocket"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/zitadel/cmd/key"
)

var (
	instanceIDs   []string
	isSystem      bool
	shouldReplace bool
)

func New(configFiles *[]string) *cobra.Command {
	closeSocket := func() error { return nil }
	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "mirrors all data of ZITADEL from one database to another",
		Long: `mirrors all data of ZITADEL from one database to another
ZITADEL needs to be initialized and set up with --for-mirror

The command does mirror all data needed and recomputes the projections.
For more details call the help functions of the sub commands.

Order of execution:
1. mirror system tables
2. mirror auth tables
3. mirror event store tables
4. recompute projections
5. verify`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			closeSocket, err = unixsocket.ListenAndIgnore()
			if err != nil {
				return fmt.Errorf("unable to listen on socket: %w", err)
			}
			if err = viper.MergeConfig(bytes.NewBuffer(defaultConfig)); err != nil {
				return errors.New("unable to read default config")
			}

			for _, file := range *configFiles {
				viper.SetConfigFile(file)
				if err = viper.MergeInConfig(); err != nil {
					return fmt.Errorf("unable to read config file: %w", err)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			defer closeSocket()
			config := mustNewMigrationConfig(viper.GetViper())
			projectionConfig := mustNewProjectionsConfig(viper.GetViper())

			masterKey, err := key.MasterKey(cmd)
			if err != nil {
				return fmt.Errorf("unable to read master key: %w", err)
			}
			if err = copySystem(cmd.Context(), config); err != nil {
				return fmt.Errorf("unable to copy system tables: %w", err)
			}
			if err = copyAuth(cmd.Context(), config); err != nil {
				return fmt.Errorf("unable to copy auth tables: %w", err)
			}
			if err = copyEventstore(cmd.Context(), config); err != nil {
				return fmt.Errorf("unable to copy eventstore tables: %w", err)
			}

			if err = projections(cmd.Context(), projectionConfig, masterKey); err != nil {
				return fmt.Errorf("unable to recompute projections: %w", err)
			}
			if err = verifyMigration(cmd.Context(), config); err != nil {
				return fmt.Errorf("unable to verify migration: %w", err)
			}
			return nil
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
		eventstoreCmd(closeSocket),
		systemCmd(closeSocket),
		projectionsCmd(closeSocket),
		authCmd(closeSocket),
		verifyCmd(closeSocket),
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
		return "WHERE instance_id <> ''"
	}
	for i := range instanceIDs {
		instanceIDs[i] = "'" + instanceIDs[i] + "'"
	}

	// COPY does not allow parameters so we need to set them directly
	return "WHERE instance_id IN (" + strings.Join(instanceIDs, ", ") + ")"
}
