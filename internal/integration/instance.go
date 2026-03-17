// Package integration provides helpers for integration testing.
package integration

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/zitadel/logging"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/webauthn"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/instance"
	internal_permission_v2 "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/org"
	"github.com/zitadel/zitadel/pkg/grpc/system"
	"github.com/zitadel/zitadel/pkg/grpc/user"
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

// NotEmpty can be used as placeholder, when the returned values is unknown.
// It can be used in tests to assert whether a value should be empty or not.
const NotEmpty = "not empty"

// newInstanceSem limits how many NewInstance calls run concurrently.
// Each CreateInstance call generates ~50 projection events; bursting all
// test packages simultaneously overwhelms the projection pipeline and causes
// "not found" errors. Limiting to a small number lets workers keep pace.
// Override with INTEGRATION_INSTANCE_CONCURRENCY env var.
var newInstanceSem = make(chan struct{}, newInstanceConcurrency())

func newInstanceConcurrency() int {
	if s := os.Getenv("INTEGRATION_INSTANCE_CONCURRENCY"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			return n
		}
	}
	return 4
}

const (
	adminPATFile = "admin-pat.txt"
)

// UserType provides constants that give
// a short explanation with the purpose
// a service account.
// This allows to pre-create users with
// different permissions and reuse them.
type UserType int

//go:generate enumer -type UserType -transform snake -trimprefix UserType
const (
	UserTypeUnspecified UserType = iota
	UserTypeIAMOwner
	UserTypeOrgOwner
	UserTypeLogin
	UserTypeNoPermission
)

const (
	UserPassword = "VeryS3cret!"
)

const (
	PortMilestoneServer = "8081"
	PortQuotaServer     = "8082"
)

// User information with a Personal Access Token.
type User struct {
	ID       string
	Username string
	Token    string
}

type UserMap map[UserType]*User

func (m UserMap) Set(typ UserType, user *User) {
	m[typ] = user
}

func (m UserMap) Get(typ UserType) *User {
	return m[typ]
}

// Host returns the primary host of zitadel, on which the first instance is served.
// http://localhost:8082 by default
func (c *Config) Host() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

// Instance is a Zitadel server and client with all resources available for testing.
type Instance struct {
	Config      Config
	Domain      string
	Instance    *instance.InstanceDetail
	DefaultOrg  *org.Org
	Users       UserMap
	AdminUserID string // First User (Human) for password login

	Client   *Client
	WebAuthN *webauthn.Client
	SAML     struct {
		PrivateKey  *rsa.PrivateKey
		Certificate *x509.Certificate
	}
}

// NewInstance returns a new instance that can be used for integration tests.
// The instance contains a gRPC client connected to the domain of this instance.
// The included users are the IAM_OWNER, ORG_OWNER of the default org and
// a Login client user.
//
// The instance is isolated and is safe for parallel testing.
func NewInstance(ctx context.Context) *Instance {
	inst, err := CreateInstance(ctx)
	if err != nil {
		panic(err)
	}
	return inst
}

// CreateInstance returns a new instance that can be used for integration tests.
// It returns an error instead of panicking.
func CreateInstance(ctx context.Context) (*Instance, error) {
	// Acquire semaphore to prevent burst-creating too many instances at once,
	// which overwhelms read-model projection workers.
	select {
	case newInstanceSem <- struct{}{}:
		defer func() { <-newInstanceSem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	primaryDomain := RandString(5) + ".integration.localhost"

	ctx = WithSystemAuthorization(ctx)
	resp, err := SystemClient().CreateInstance(ctx, &system.CreateInstanceRequest{
		InstanceName: "testinstance",
		CustomDomain: primaryDomain,
		Owner: &system.CreateInstanceRequest_Machine_{
			Machine: &system.CreateInstanceRequest_Machine{
				UserName:            "owner",
				Name:                "owner",
				PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	i := &Instance{
		Config: loadedConfig,
		Domain: primaryDomain,
	}
	if err := i.setClient(ctx); err != nil {
		return nil, err
	}
	if err := i.awaitFirstUser(WithAuthorizationToken(ctx, resp.GetPat())); err != nil {
		return nil, err
	}
	if err := i.setupInstance(ctx, resp.GetPat()); err != nil {
		return nil, err
	}
	if err := i.awaitReadProjections(ctx); err != nil {
		return nil, err
	}
	return i, nil
}

func (i *Instance) ID() string {
	return i.Instance.GetId()
}

func (i *Instance) awaitFirstUser(ctx context.Context) error {
	var allErrs []error
	for {
		resp, err := i.Client.UserV2.AddHumanUser(ctx, &user_v2.AddHumanUserRequest{
			Username: proto.String("zitadel-admin@zitadel.localhost"),
			Profile: &user_v2.SetHumanProfile{
				GivenName:  "hodor",
				FamilyName: "hodor",
				NickName:   proto.String("hodor"),
			},
			Email: &user_v2.SetHumanEmail{
				Email: "zitadel-admin@zitadel.localhost",
				Verification: &user_v2.SetHumanEmail_IsVerified{
					IsVerified: true,
				},
			},
			PasswordType: &user_v2.AddHumanUserRequest_Password{
				Password: &user_v2.Password{
					Password:       "Password1!",
					ChangeRequired: false,
				},
			},
		})
		if err == nil {
			i.AdminUserID = resp.GetUserId()
			return nil
		}
		logging.WithError(err).Debug("await first instance user")
		allErrs = append(allErrs, err)
		select {
		case <-ctx.Done():
			return errors.Join(append(allErrs, ctx.Err())...)
		case <-time.After(time.Second):
			continue
		}
	}
}

func (i *Instance) setupInstance(ctx context.Context, token string) error {
	i.Users = make(UserMap)
	ctx = WithAuthorizationToken(ctx, token)
	if err := i.setInstance(ctx); err != nil {
		return err
	}
	if err := i.setOrganization(ctx); err != nil {
		return err
	}
	if err := i.createMachineUserInstanceOwner(ctx, token); err != nil {
		return err
	}
	if _, err := i.createMachineUser(ctx, UserTypeOrgOwner); err != nil {
		return err
	}
	if _, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{
			Resource: &internal_permission_v2.ResourceType_OrganizationId{OrganizationId: i.DefaultOrg.GetId()},
		},
		UserId: i.Users.Get(UserTypeOrgOwner).ID,
		Roles:  []string{"ORG_OWNER"},
	}); err != nil {
		return err
	}
	if _, err := i.createMachineUser(ctx, UserTypeLogin); err != nil {
		return err
	}
	if _, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{
			Resource: &internal_permission_v2.ResourceType_Instance{Instance: true},
		},
		UserId: i.Users.Get(UserTypeLogin).ID,
		Roles:  []string{"IAM_LOGIN_CLIENT"},
	}); err != nil {
		return err
	}
	if _, err := i.createMachineUser(ctx, UserTypeNoPermission); err != nil {
		return err
	}
	i.createWebAuthNClient()
	return nil
}

// Host returns the primary Domain of the instance with the port.
func (i *Instance) Host() string {
	return fmt.Sprintf("%s:%d", i.Domain, i.Config.Port)
}

func (i *Instance) createMachineUserInstanceOwner(ctx context.Context, token string) error {
	return await(func() error {
		user, err := i.Client.Auth.GetMyUser(WithAuthorizationToken(ctx, token), &auth.GetMyUserRequest{})
		if err != nil {
			return err
		}
		i.Users.Set(UserTypeIAMOwner, &User{
			ID:       user.GetUser().GetId(),
			Username: user.GetUser().GetUserName(),
			Token:    token,
		})
		return nil
	})
}

// awaitReadProjections waits until the key read-model projections have caught
// up to at least the last event written during setupInstance.  Under parallel
// load many instances are set up concurrently and projection workers can fall
// behind, causing "not found" errors in the first few test calls.
//
// We probe two independent projections:
//  1. projections.users (via GetUserByID for UserTypeNoPermission — the last
//     user written in setupInstance).
//  2. projections.secret_generators (via ListSecretGenerators) — populated
//     from instance-creation events and required by ImportHumanUser /
//     AddHumanUser notification flows.
func (i *Instance) awaitReadProjections(ctx context.Context) error {
	ownerCtx := i.WithAuthorizationToken(ctx, UserTypeIAMOwner)
	userID := i.Users.Get(UserTypeNoPermission).ID

	// Wait for projections.users
	if err := await(func() error {
		_, err := i.Client.UserV2.GetUserByID(ownerCtx, &user_v2.GetUserByIDRequest{
			UserId: userID,
		})
		if err != nil {
			logging.WithError(err).Debug("await users projection")
		}
		return err
	}); err != nil {
		return err
	}

	// Wait for projections.secret_generators
	return await(func() error {
		resp, err := i.Client.Admin.ListSecretGenerators(ownerCtx, &admin.ListSecretGeneratorsRequest{})
		if err != nil {
			logging.WithError(err).Debug("await secret_generators projection")
			return err
		}
		if len(resp.GetResult()) == 0 {
			return errors.New("secret_generators projection not yet ready")
		}
		return nil
	})
}

func (i *Instance) setClient(ctx context.Context) error {
	client, err := newClient(ctx, i.Host())
	if err != nil {
		return err
	}
	i.Client = client
	return nil
}

func (i *Instance) setInstance(ctx context.Context) error {
	return await(func() error {
		instance, err := i.Client.Admin.GetMyInstance(ctx, &admin.GetMyInstanceRequest{})
		if err != nil {
			return err
		}
		i.Instance = instance.GetInstance()
		return nil
	})
}

func (i *Instance) setOrganization(ctx context.Context) error {
	return await(func() error {
		resp, err := i.Client.Mgmt.GetMyOrg(ctx, &management.GetMyOrgRequest{})
		if err != nil {
			return err
		}
		i.DefaultOrg = resp.GetOrg()
		return nil
	})
}

func (i *Instance) createMachineUser(ctx context.Context, userType UserType) (userID string, err error) {
	err = await(func() error {
		username := Username()
		userResp, err := i.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
			UserName:        username,
			Name:            username,
			Description:     userType.String(),
			AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
		})
		if err != nil {
			return err
		}
		userID = userResp.GetUserId()
		patResp, err := i.Client.Mgmt.AddPersonalAccessToken(ctx, &management.AddPersonalAccessTokenRequest{
			UserId: userID,
		})
		if err != nil {
			return err
		}
		i.Users.Set(userType, &User{
			ID:       userID,
			Username: username,
			Token:    patResp.GetToken(),
		})
		return nil
	})
	return userID, err
}

func (i *Instance) createWebAuthNClient() {
	i.WebAuthN = webauthn.NewClient(i.Config.WebAuthNName, i.Domain, http_util.BuildOrigin(i.Host(), i.Config.Secure))
}

// Deprecated: WithAuthorization is misleading, as we have Zitadel resources called authorization now.
// It is aliased to WithAuthorizationToken, which sets the Authorization header with a Bearer token.
// Use WithAuthorizationToken directly instead.
func (i *Instance) WithAuthorization(ctx context.Context, u UserType) context.Context {
	return i.WithAuthorizationToken(ctx, u)
}

func (i *Instance) WithAuthorizationToken(ctx context.Context, u UserType) context.Context {
	return WithAuthorizationToken(ctx, i.Users.Get(u).Token)
}

func (i *Instance) GetUserID(u UserType) string {
	return i.Users.Get(u).ID
}

func WithAuthorizationToken(ctx context.Context, token string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = make(metadata.MD)
	}
	md.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return metadata.NewOutgoingContext(ctx, md)
}

func (i *Instance) BearerToken(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	return md.Get("Authorization")[0]
}

func (i *Instance) WithSystemAuthorizationHTTP(u UserType) map[string]string {
	return map[string]string{"Authorization": fmt.Sprintf("Bearer %s", i.Users.Get(u).Token)}
}

func await(af func() error) error {
	maxTimer := time.NewTimer(15 * time.Minute)
	for {
		err := af()
		if err == nil {
			return nil
		}
		select {
		case <-maxTimer.C:
			return err
		case <-time.After(time.Second):
			continue
		}
	}
}

func mustAwait(af func() error) {
	if err := await(af); err != nil {
		panic(err)
	}
}
