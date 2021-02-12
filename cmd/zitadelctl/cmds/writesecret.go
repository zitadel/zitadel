package cmds

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/secrets"
	"github.com/spf13/cobra"
)

func WriteSecretCommand(rv RootValues) *cobra.Command {

	var (
		value string
		file  string
		stdin bool
		cmd   = &cobra.Command{
			Use:   "writesecret [path]",
			Short: "Encrypt a secret and push it to the repository",
			Long:  "Encrypt a secret and push it to the repository.\nIf no path is provided, a secret can interactively be chosen from a list of all possible secrets",
			Args:  cobra.MaximumNArgs(1),
			Example: `orbctl writesecret mystaticprovider.bootstrapkey --file ~/.ssh/my-orb-bootstrap
orbctl writesecret mygceprovider.google_application_credentials_value --value "$(cat $GOOGLE_APPLICATION_CREDENTIALS)" `,
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&value, "value", "", "Secret value to encrypt")
	flags.StringVarP(&file, "file", "s", "", "File containing the value to encrypt")
	flags.BoolVar(&stdin, "stdin", false, "Value to encrypt is read from standard input")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		_, monitor, orbConfig, gitClient, _, errFunc, err := rv()
		if err != nil {
			return err
		}
		defer func() {
			err = errFunc(err)
		}()

		s, err := key(value, file, stdin)
		if err != nil {
			monitor.Error(err)
			return nil
		}

		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			monitor.Error(err)
			return nil
		}

		if err := gitClient.Clone(); err != nil {
			monitor.Error(err)
			return nil
		}

		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		if err := secret.Write(
			monitor,
			gitClient,
			path,
			s,
			secrets.GetAllSecretsFunc(orbConfig),
			secrets.PushFunc(),
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
