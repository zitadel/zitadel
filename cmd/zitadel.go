package cmd

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/admin"
	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/cmd/start"
)

var (
	configFiles []string

	//go:embed defaults.yaml
	defaultConfig []byte
)

func New(out io.Writer, in io.Reader, args []string, server chan<- *start.Server) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zitadel",
		Short: "The ZITADEL CLI lets you interact with ZITADEL",
		Long:  `The ZITADEL CLI lets you interact with ZITADEL`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
		Version: build.Version(),
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ZITADEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	logging.OnError(err).Fatal("unable to read default config")

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")

	cmd.AddCommand(
		admin.New(), //is now deprecated, remove later on
		initialise.New(),
		setup.New(),
		start.New(server),
		start.NewStartFromInit(server),
		start.NewStartFromSetup(server),
		key.New(),
	)

	cmd.InitDefaultVersionFlag()

	return cmd
}

func initConfig() {
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		err := viper.MergeInConfig()
		logging.WithFields("file", file).OnError(err).Warn("unable to read config file")
	}
}

type TestServer struct {
	*start.Server
	wg sync.WaitGroup
}

func (s *TestServer) Done() {
	s.Shutdown <- os.Interrupt
	s.wg.Wait()
}

func NewTestServer(args []string) *TestServer {
	testServer := new(TestServer)
	server := make(chan *start.Server, 1)

	testServer.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		cmd := New(os.Stdout, os.Stdin, args, server)
		cmd.SetArgs(args)
		logging.OnError(cmd.Execute()).Fatal()
	}(&testServer.wg)

	testServer.Server = <-server
	return testServer
}
