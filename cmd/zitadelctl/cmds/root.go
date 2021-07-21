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

type GetRootValues func(command string, tags map[string]interface{}, component string, moreComponents ...string) *RootValues

type errFunc func(err error) error

func RootCommand(version string, monitor mntr.Monitor) (*cobra.Command, GetRootValues) {

	var (
		ctx = context.Background()
		rv  = &RootValues{
			Ctx:     ctx,
			Version: version,
			ErrFunc: func(err error) error {
				if err != nil {
					monitor.Error(err)
				}
				return nil
			},
		}
		orbConfigPath    string
		verbose          bool
		disableIngestion bool
	)
	cmd := &cobra.Command{
		Use:   "zitadelctl [flags]",
		Short: "Interact with your IAM orbs",
		Long: `zitadelctl launches zitadel and simplifies common tasks such as deploying operators or reading and writing secrets.
Participate in our community on https://github.com/caos/orbos
and visit our website at https://caos.ch`,
		Example: `$ # For being able to use the --gitops flag, you need to create an orbconfig and add an SSH deploy key to your github project 
$ # Create an ssh key pair
$ ssh-keygen -b 2048 -t rsa -f ~/.ssh/myorbrepo -q -N ""
$ # Create the orbconfig
$ mkdir -p ~/.orb
$ cat > ~/.orb/myorb << EOF
> # this is the ssh URL to your git repository
> url: git@github.com:me/my-orb.git
> masterkey: "$(openssl rand -base64 21)" # used for encrypting and decrypting secrets
> # the repokey is used to connect to your git repository
> repokey: |
> $(cat ~/.ssh/myorbrepo | sed s/^/\ \ /g)
> EOF
$ zitadelctl --gitops -f ~/.orb/myorb [command]
`,
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&rv.Gitops, "gitops", false, "Run zitadelctl in gitops mode")
	flags.StringVarP(&orbConfigPath, "orbconfig", "f", "~/.orb/config", "Path to the file containing the orbs git repo URL, deploy key and the master key for encrypting and decrypting secrets")
	flags.StringVarP(&rv.Kubeconfig, "kubeconfig", "k", "~/.kube/config", "Path to the kubeconfig file to the cluster zitadelctl should target")
	flags.BoolVar(&verbose, "verbose", false, "Print debug levelled logs")
	flags.BoolVar(&disableIngestion, "disable-ingestion", false, "Don't help CAOS AG to improve ZITADEL by sending them errors and usage data")

	return cmd, func(command string, tags map[string]interface{}, component string, moreComponents ...string) *RootValues {

		if verbose {
			monitor = monitor.Verbose()
		}

		rv.Monitor = monitor
		rv.Kubeconfig = helpers.PruneHome(rv.Kubeconfig)
		rv.GitClient = git.New(ctx, monitor, "zitadel", "orbos@caos.ch")

		if rv.Gitops {
			prunedPath := helpers.PruneHome(orbConfigPath)
			rv.OrbConfig, _ = orb.ParseOrbConfig(prunedPath)
			if rv.OrbConfig == nil {
				rv.OrbConfig = &orb.Orb{Path: prunedPath}
			}
		}

		env := "unknown"
		if orbID, err := rv.OrbConfig.ID(); err == nil {
			env = orbID
		}

		if component == "" {
			component = "zitadelctl"
		}

		if !disableIngestion {
			if err := mntr.Ingest(rv.Monitor, "zitadel", version, env, component, moreComponents...); err != nil {
				panic(err)
			}
		}

		rv.Monitor.WithFields(map[string]interface{}{"command": command, "gitops": rv.Gitops}).WithFields(tags).CaptureMessage("zitadelctl invoked")

		return rv
	}
}
