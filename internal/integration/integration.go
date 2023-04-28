// Package integration provides helpers for integration testing.
package integration

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/cmd"
	"github.com/zitadel/zitadel/cmd/start"
	"github.com/zitadel/zitadel/internal/api/authz"
	z_oidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
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

type UserType int

//go:generate stringer -type=UserType
const (
	Unspecified UserType = iota
	OrgOwner
)

type User struct {
	*query.User
	Token string
}

type Tester struct {
	*start.Server

	Instance     authz.Instance
	Organisation *query.Org
	Users        map[UserType]User

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

const (
	SystemUser = "integration1"
)

func (s *Tester) createSystemUser(ctx context.Context) {
	var err error

	s.Instance, err = s.Queries.InstanceByHost(ctx, "localhost:8080")
	logging.OnError(err).Fatal("query instance")
	ctx = authz.WithInstance(ctx, s.Instance)

	s.Organisation, err = s.Queries.OrgByID(ctx, true, s.Instance.DefaultOrganisationID())
	logging.OnError(err).Fatal("query organisation")

	query, err := query.NewUserUsernameSearchQuery(SystemUser, query.TextEquals)
	logging.OnError(err).Fatal("user query")
	user, err := s.Queries.GetUser(ctx, true, true, query)

	if errors.Is(err, sql.ErrNoRows) {
		_, err = s.Commands.AddMachine(ctx, &command.Machine{
			ObjectRoot: models.ObjectRoot{
				ResourceOwner: s.Organisation.ID,
			},
			Username:        SystemUser,
			Name:            SystemUser,
			Description:     "who cares?",
			AccessTokenType: domain.OIDCTokenTypeJWT,
		})
		logging.OnError(err).Fatal("add machine user")
		user, err = s.Queries.GetUser(ctx, true, true, query)

	}
	logging.OnError(err).Fatal("get user")

	_, err = s.Commands.AddOrgMember(ctx, s.Organisation.ID, user.ID, "ORG_OWNER")
	target := new(caos_errs.AlreadyExistsError)
	if !errors.As(err, &target) {
		logging.OnError(err).Fatal("add org member")
	}

	scopes := []string{oidc.ScopeOpenID, z_oidc.ScopeUserMetaData, z_oidc.ScopeResourceOwner}
	pat := command.NewPersonalAccessToken(user.ResourceOwner, user.ID, time.Now().Add(time.Hour), scopes, domain.UserTypeMachine)
	_, err = s.Commands.AddPersonalAccessToken(ctx, pat)
	logging.OnError(err).Fatal("add pat")

	s.Users = map[UserType]User{
		OrgOwner: {
			User:  user,
			Token: pat.Token,
		},
	}
}

func (s *Tester) WithSystemAuthorization(ctx context.Context, u UserType) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", s.Users[u].Token))
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
	tester.createSystemUser(ctx)

	return tester
}

func Contexts(timeout time.Duration) (ctx, errCtx context.Context, cancel context.CancelFunc) {
	errCtx, cancel = context.WithCancel(context.Background())
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	return ctx, errCtx, cancel
}
