package cmds

import (
	"context"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/spf13/cobra"
)

type RootValues struct {
	Ctx        context.Context
	Monitor    mntr.Monitor
	Version    string
	Gitops     bool
	OrbConfig  *orb.Orb
	GitClient  *git.Client
	Kubeconfig string
	ErrFunc    errFunc
}

type GetRootValues func() (*RootValues, error)

type errFunc func(err error) error

func RootCommand(version string) (*cobra.Command, GetRootValues) {

	var (
		ctx     = context.Background()
		monitor = mntr.Monitor{
			OnInfo:   mntr.LogMessage,
			OnChange: mntr.LogMessage,
			OnError:  mntr.LogError,
		}
		rv = &RootValues{
			Ctx:     ctx,
			Version: version,
			ErrFunc: func(err error) error {
				if err != nil {
					monitor.Error(err)
				}
				return nil
			},
		}
		orbConfigPath string
		verbose       bool
	)
	cmd := &cobra.Command{
		Use:   "zitadelctl [flags]",
		Short: "Interact with your IAM orbs",
		Long: `zitadelctl launches zitadel and simplifies common tasks such as updating your kubeconfig.
Participate in our community on https://github.com/caos/orbos
and visit our website at https://caos.ch`,
		Example: `$ mkdir -p ~/.orb
$ cat > ~/.orb/myorb << EOF
> url: git@github.com:me/my-orb.git
> masterkey: "$(gopass my-secrets/orbs/myorb/masterkey)"
> repokey: |
> $(cat ~/.ssh/myorbrepo | sed s/^/\ \ /g)
> EOF
$ orbctl -f ~/.orb/myorb [command]
`,
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&rv.Gitops, "gitops", false, "Run orbctl in gitops mode. Not specifying this flag is only supported for BOOM and Networking Operator")
	flags.StringVarP(&orbConfigPath, "orbconfig", "f", "~/.orb/config", "Path to the file containing the orbs git repo URL, deploy key and the master key for encrypting and decrypting secrets")
	flags.StringVarP(&rv.Kubeconfig, "kubeconfig", "k", "~/.kube/config", "Path to the kubeconfig file to the cluster orbctl should target")
	flags.BoolVar(&verbose, "verbose", false, "Print debug levelled logs")

	return cmd, func() (*RootValues, error) {

		if verbose {
			monitor = monitor.Verbose()
		}

		rv.Monitor = monitor
		rv.Kubeconfig = helpers.PruneHome(rv.Kubeconfig)
		rv.GitClient = git.New(ctx, monitor, "orbos", "orbos@caos.ch")

		if rv.Gitops {
			prunedPath := helpers.PruneHome(orbConfigPath)
			orbConfig, err := orb.ParseOrbConfig(prunedPath)
			if err != nil {
				orbConfig = &orb.Orb{Path: prunedPath}
				return nil, err
			}
			rv.OrbConfig = orbConfig
		}

		return rv, nil
	}
}
