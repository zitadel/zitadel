package mirror

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/cmd/key"
)

var (
	instanceIDs   []string
	isSystem      bool
	shouldReplace bool
)

func New(configFiles *[]string) *cobra.Command {
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
			defer func() {
				logging.OnError(cmd.Context(), err).ErrorContext(cmd.Context(), "zitadel mirror (sub)command failed")
			}()

			err = viper.MergeConfig(bytes.NewBuffer(defaultConfig))
			if err != nil {
				return fmt.Errorf("unable to read default config: %w", err)
			}
			for _, file := range *configFiles {
				viper.SetConfigFile(file)
				err := viper.MergeInConfig()
				logging.OnError(cmd.Context(), err).ErrorContext(cmd.Context(), "unable to read config file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).ErrorContext(cmd.Context(), "zitadel mirror command failed")
			}()

			config, shutdown, err := newMigrationConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return fmt.Errorf("unable to create migration config: %w", err)
			}
			// Set logger again to include changes from config
			cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()

			projectionConfig, _, err := newProjectionsConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return fmt.Errorf("unable to create projections config: %w", err)
			}

			masterKey, err := key.MasterKey(cmd)
			if err != nil {
				return fmt.Errorf("unable to read master key: %w", err)
			}

			copySystem(cmd.Context(), config)
			copyAuth(cmd.Context(), config)
			copyEventstore(cmd.Context(), config)

			defer func() {
				if recErr, ok := recover().(error); ok {
					err = recErr
				}
			}()
			projections(cmd.Context(), projectionConfig, masterKey)
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
		return "WHERE instance_id <> ''"
	}
	for i := range instanceIDs {
		instanceIDs[i] = "'" + instanceIDs[i] + "'"
	}

	// COPY does not allow parameters so we need to set them directly
	return "WHERE instance_id IN (" + strings.Join(instanceIDs, ", ") + ")"
}

func panicOnError(ctx context.Context, err error, logMsg string) {
	logging.OnError(ctx, err).ErrorContext(ctx, logMsg)
	if err != nil {
		panic(err)
	}
}
