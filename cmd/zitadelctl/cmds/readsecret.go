package cmds

import (
	"os"

	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/secrets"
	"github.com/spf13/cobra"
)

func ReadSecretCommand(getRv GetRootValues) *cobra.Command {
	return &cobra.Command{
		Use:     "readsecret [path]",
		Short:   "Print a secrets decrypted value to stdout",
		Long:    "Print a secrets decrypted value to stdout.\nIf no path is provided, a secret can interactively be chosen from a list of all possible secrets",
		Args:    cobra.MaximumNArgs(1),
		Example: `zitadelctl readsecret zitadel.emailappkey > ~/emailappkey`,
		RunE: func(cmd *cobra.Command, args []string) error {
			rv, err := getRv()
			if err != nil {
				return err
			}
			defer func() {
				err = rv.ErrFunc(err)
			}()

			monitor := rv.Monitor
			orbConfig := rv.OrbConfig
			gitClient := rv.GitClient

			if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
				return err
			}

			if err := gitClient.Clone(); err != nil {
				return err
			}

			path := ""
			if len(args) > 0 {
				path = args[0]
			}

			value, err := secret.Read(
				monitor,
				gitClient,
				path,
				secrets.GetAllSecretsFunc(orbConfig))
			if err != nil {
				monitor.Error(err)
				return nil
			}

			if _, err := os.Stdout.Write([]byte(value)); err != nil {
				monitor.Error(err)
				return nil
			}
			return nil
		},
	}
}
