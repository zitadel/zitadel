package cmds

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/caos/orbos/pkg/kubernetes/cli"

	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/secrets"
	"github.com/spf13/cobra"
)

func WriteSecretCommand(getRv GetRootValues) *cobra.Command {

	var (
		value string
		file  string
		stdin bool
		cmd   = &cobra.Command{
			Use:   "writesecret [path]",
			Short: "Encrypt a secret and push it to the repository",
			Long:  "Encrypt a secret and push it to the repository.\nIf no path is provided, a secret can interactively be chosen from a list of all possible secrets",
			Args:  cobra.MaximumNArgs(1),
			Example: `zitadelctl writesecret database.bucket.serviceaccountjson.encrypted --file ~/googlecloudstoragesa.json
zitadelctl writesecret database.bucket.serviceaccountjson.encrypted --value "$(cat ~/googlecloudstoragesa.json)"
cat ~/googlecloudstoragesa.json | zitadelctl writesecret database.bucket.serviceaccountjson.encrypted --stdin`,
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&value, "value", "", "Secret value to encrypt")
	flags.StringVarP(&file, "file", "s", "", "File containing the value to encrypt")
	flags.BoolVar(&stdin, "stdin", false, "Value to encrypt is read from standard input")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
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

		s, err := key(value, file, stdin)
		if err != nil {
			monitor.Error(err)
			return nil
		}

		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		k8sClient, err := cli.Client(monitor, orbConfig, gitClient, rv.Kubeconfig, rv.Gitops, true)
		if err != nil && !rv.Gitops {
			return err
		}
		err = nil

		if err := secret.Write(
			monitor,
			k8sClient,
			path,
			s,
			"zitadelctl",
			fmt.Sprintf(rv.Version),
			secrets.GetAllSecretsFunc(monitor, path != "", rv.Gitops, gitClient, k8sClient, orbConfig),
			secrets.PushFunc(monitor, rv.Gitops, gitClient, k8sClient),
		); err != nil {
			monitor.Error(err)
		}
		return nil
	}
	return cmd
}

func key(value string, file string, stdin bool) (string, error) {

	channels := 0
	if value != "" {
		channels++
	}
	if file != "" {
		channels++
	}
	if stdin {
		channels++
	}

	if channels != 1 {
		return "", errors.New("Key must be provided eighter by value or by file path or by standard input")
	}

	if value != "" {
		return value, nil
	}

	readFunc := func() ([]byte, error) {
		return ioutil.ReadFile(file)
	}
	if stdin {
		readFunc = func() ([]byte, error) {
			return ioutil.ReadAll(os.Stdin)
		}
	}

	key, err := readFunc()
	if err != nil {
		panic(err)
	}
	return string(key), err
}
