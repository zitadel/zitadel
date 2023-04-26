// Package integration provides helpers for integration testing.
package integration

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zitadel/zitadel/cmd"
	"github.com/zitadel/zitadel/cmd/start"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

var (
	//go:embed config/zitadel.yaml
	zitadelYAML []byte
	//go:embed config/cockroach.yaml
	cockroachYAML []byte
	//go:embed config/postgres.yaml
	postgresYAML []byte
)

type Tester struct {
	*start.Server
	GRPCClientConn *grpc.ClientConn
	wg             sync.WaitGroup // used for shutdown
}

const commandLine = `start --masterkey MasterkeyNeedsToHave32Characters`

func (s *Tester) createClientConn(ctx context.Context) {
	target := fmt.Sprintf("localhost:%d", s.Config.Port)
	cc, err := grpc.DialContext(ctx, target,
		grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.Shutdown <- os.Interrupt
		s.wg.Wait()
	}
	logging.OnError(err).Fatal("integration tester client dial")
	logging.New().WithField("target", target).Info("finished dialing grpc client conn")

	s.GRPCClientConn = cc
	err = s.pollHealth(ctx)
	logging.OnError(err).Fatal("integration tester health")
}

// pollHealth waits until a healthy status is reported.
// TODO: remove when we make the setup blocking on all
// projections completed.
func (s *Tester) pollHealth(ctx context.Context) (err error) {
	client := admin.NewAdminServiceClient(s.GRPCClientConn)

	for {
		err = func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			_, err := client.Healthz(ctx, &admin.HealthzRequest{})
			return err
		}(ctx)
		if err == nil {
			return nil
		}
		logging.WithError(err).Info("poll healthz")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			continue
		}
	}
}

func (s *Tester) Done() {
	err := s.GRPCClientConn.Close()
	logging.OnError(err).Error("integration tester client close")

	s.Shutdown <- os.Interrupt
	s.wg.Wait()
}

func NewTester(ctx context.Context) *Tester {
	args := strings.Split(commandLine, " ")

	sc := make(chan *start.Server)
	//nolint:contextcheck
	cmd := cmd.New(os.Stdout, os.Stdin, args, sc)
	cmd.SetArgs(args)
	err := viper.MergeConfig(bytes.NewBuffer(zitadelYAML))
	logging.OnError(err).Fatal()

	flavor := os.Getenv("INTEGRATION_DB_FLAVOR")
	switch flavor {
	case "cockroach", "":
		err = viper.MergeConfig(bytes.NewBuffer(cockroachYAML))
	case "postgres":
		err = viper.MergeConfig(bytes.NewBuffer(postgresYAML))
	default:
		logging.New().WithField("flavor", flavor).Fatal("unknown db flavor set in INTEGRATION_DB_FLAVOR")
	}
	logging.OnError(err).Fatal()

	tester := new(Tester)
	tester.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		logging.OnError(cmd.Execute()).Fatal()
		wg.Done()
	}(&tester.wg)

	select {
	case tester.Server = <-sc:
	case <-ctx.Done():
		logging.OnError(ctx.Err()).Fatal("waiting for integration tester server")
	}
	tester.createClientConn(ctx)

	return tester
}
