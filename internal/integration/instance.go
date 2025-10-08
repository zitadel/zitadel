// Package integration provides helpers for integration testing.
package integration

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zitadel/logging"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/webauthn"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/instance"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/org"
	"github.com/zitadel/zitadel/pkg/grpc/system"
	"github.com/zitadel/zitadel/pkg/grpc/user"
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

// NotEmpty can be used as placeholder, when the returned values is unknown.
// It can be used in tests to assert whether a value should be empty or not.
const NotEmpty = "not empty"

const (
	adminPATFile = "admin-pat.txt"
)

// UserType provides constants that give
// a short explanation with the purpose
// a service user.
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
// http://localhost:8080 by default
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
	AdminUserID string // First human user for password login

	Client   *Client
	WebAuthN *webauthn.Client
}

// NewInstance returns a new instance that can be used for integration tests.
// The instance contains a gRPC client connected to the domain of this instance.
// The included users are the IAM_OWNER, ORG_OWNER of the default org and
// a Login client user.
//
// The instance is isolated and is safe for parallel testing.
func NewInstance(ctx context.Context) *Instance {
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
		panic(err)
	}
	i := &Instance{
		Config: loadedConfig,
		Domain: primaryDomain,
	}
	i.setClient(ctx)
	i.awaitFirstUser(WithAuthorizationToken(ctx, resp.GetPat()))
	i.setupInstance(ctx, resp.GetPat())
	return i
}

func (i *Instance) ID() string {
	return i.Instance.GetId()
}

func (i *Instance) awaitFirstUser(ctx context.Context) {
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
			return
		}
		logging.WithError(err).Debug("await first instance user")
		allErrs = append(allErrs, err)
		select {
		case <-ctx.Done():
			panic(errors.Join(append(allErrs, ctx.Err())...))
		case <-time.After(time.Second):
			continue
		}
	}
}

func (i *Instance) setupInstance(ctx context.Context, token string) {
	i.Users = make(UserMap)
	ctx = WithAuthorizationToken(ctx, token)
	i.setInstance(ctx)
	i.setOrganization(ctx)
	i.createMachineUserInstanceOwner(ctx, token)
	i.createMachineUserOrgOwner(ctx)
	i.createLoginClient(ctx)
	i.createMachineUserNoPermission(ctx)
	i.createWebAuthNClient()
}

// Host returns the primary Domain of the instance with the port.
func (i *Instance) Host() string {
	return fmt.Sprintf("%s:%d", i.Domain, i.Config.Port)
}

func loadInstanceOwnerPAT() string {
	data, err := os.ReadFile(filepath.Join(tmpDir, adminPATFile))
	if err != nil {
		panic(err)
	}
	return string(bytes.TrimSpace(data))
}

func (i *Instance) createMachineUserInstanceOwner(ctx context.Context, token string) {
	mustAwait(func() error {
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

func (i *Instance) createMachineUserOrgOwner(ctx context.Context) {
	_, err := i.Client.Mgmt.AddOrgMember(ctx, &management.AddOrgMemberRequest{
		UserId: i.createMachineUser(ctx, UserTypeOrgOwner),
		Roles:  []string{"ORG_OWNER"},
	})
	if err != nil {
		panic(err)
	}
}

func (i *Instance) createLoginClient(ctx context.Context) {
	_, err := i.Client.Admin.AddIAMMember(ctx, &admin.AddIAMMemberRequest{
		UserId: i.createMachineUser(ctx, UserTypeLogin),
		Roles:  []string{"IAM_LOGIN_CLIENT"},
	})
	if err != nil {
		panic(err)
	}
}

func (i *Instance) createMachineUserNoPermission(ctx context.Context) {
	i.createMachineUser(ctx, UserTypeNoPermission)
}

func (i *Instance) setClient(ctx context.Context) {
	client, err := newClient(ctx, i.Host())
	if err != nil {
		panic(err)
	}
	i.Client = client
}

func (i *Instance) setInstance(ctx context.Context) {
	mustAwait(func() error {
		instance, err := i.Client.Admin.GetMyInstance(ctx, &admin.GetMyInstanceRequest{})
		i.Instance = instance.GetInstance()
		return err
	})
}

func (i *Instance) setOrganization(ctx context.Context) {
	mustAwait(func() error {
		resp, err := i.Client.Mgmt.GetMyOrg(ctx, &management.GetMyOrgRequest{})
		i.DefaultOrg = resp.GetOrg()
		return err
	})
}

func (i *Instance) createMachineUser(ctx context.Context, userType UserType) (userID string) {
	mustAwait(func() error {
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
	return userID
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
