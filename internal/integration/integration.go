// Package integration provides helpers for integration testing.
package integration

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zitadel/zitadel/cmd"
	"github.com/zitadel/zitadel/cmd/start"
)

type Tester struct {
	*start.Server
	ClientConn *grpc.ClientConn

	wg sync.WaitGroup // used for shutdown
}

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

	s.ClientConn = cc
}

func (s *Tester) Done() {
	err := s.ClientConn.Close()
	logging.OnError(err).Error("integration tester client close")

	s.Shutdown <- os.Interrupt
	s.wg.Wait()
}

func NewTester(ctx context.Context, args []string) *Tester {
	tester := new(Tester)
	sc := make(chan *start.Server)
	tester.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		cmd := cmd.New(os.Stdout, os.Stdin, args, sc)
		cmd.SetArgs(args)
		logging.OnError(cmd.Execute()).Fatal()
	}(&tester.wg)

	select {
	case tester.Server = <-sc:
	case <-ctx.Done():
		logging.OnError(ctx.Err()).Fatal("waiting for integration tester server")
	}
	tester.createClientConn(ctx)

	return tester
}
