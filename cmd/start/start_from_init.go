package start

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/cmd/tls"
)

func NewStartFromInit(server chan<- *Server) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-from-init",
		Short: "cold starts zitadel",
		Long: `cold starts ZITADEL.
First the minimum requirements to start ZITADEL are set up.
Second the initial events are created.
Last ZITADEL starts.

Requirements:
- postgreSQL`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel start-from-init command failed")
			}()

			err = tls.ModeFromFlag(cmd)
			if err != nil {
				return fmt.Errorf("invalid tlsMode: %w", err)
			}

			masterKey, err := key.MasterKey(cmd)
			if err != nil {
				return fmt.Errorf("no master key provided: %w", err)
			}

			initCtx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			initConfig, shutdown, err := initialise.NewConfig(initCtx, viper.GetViper())
			if err != nil {
				return err
			}
			// Set logger again to include changes from config
			cmd.SetContext(logging.NewCtx(cmd.Context(), logging.StreamRuntime))
			defer func() {
				err = errors.Join(err, shutdown(initCtx))
			}()

			err = initialise.InitAll(initCtx, initConfig)
			if err != nil {
				return err
			}

			err = setup.BindInitProjections(cmd)
			if err != nil {
				return fmt.Errorf("unable to bind \"init-projections\" flag: %w", err)
			}

			setupConfig, _, err := setup.NewConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return err
			}

			setupSteps, err := setup.NewSteps(cmd.Context(), viper.New())
			if err != nil {
				return err
			}

			setupCtx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			err = setup.Setup(setupCtx, setupConfig, setupSteps, masterKey)
			if err != nil {
				return err
			}

			startConfig, _, err := NewConfig(cmd.Context(), viper.GetViper())
			if err != nil {
				return err
			}

			return startZitadel(cmd.Context(), startConfig, masterKey, server)
		},
	}

	startFlags(cmd)
	setup.Flags(cmd)

	return cmd
}
